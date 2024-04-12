package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type OpenObserveHandler struct {
	opts             OpenObserveHandlerOptions
	semaphore        chan struct{}
	withAttrs        []slog.Attr
	withGroup        []string
	buffer           *bytes.Buffer
	mu               sync.Mutex
	sendDebounceFunc *time.Timer
}

func (sl *OpenObserveHandler) clone() *OpenObserveHandler {
	return &OpenObserveHandler{
		opts:      sl.opts,
		semaphore: sl.semaphore,
		withAttrs: sl.withAttrs,
		withGroup: sl.withGroup,
		buffer:    &bytes.Buffer{},
	}
}

func (sl *OpenObserveHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return sl.opts.HandlerOptions.Level.Level() <= lvl
}

func (sl *OpenObserveHandler) Handle(ctx context.Context, record slog.Record) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	if sl.sendDebounceFunc == nil {
		sl.sendDebounceFunc = time.AfterFunc(sl.opts.BufferTimeout, func() {
			sl.mu.Lock()
			defer sl.mu.Unlock()
			if sl.buffer.Len() < 1 {
				return
			}
			b := sl.extractBuffer().Bytes()
			if b[len(b)-1] == ',' {
				b = b[:len(b)-1]
			}
			b = append(b, ']')
			sl.semaphore <- struct{}{}
			go func() {
				defer func() { <-sl.semaphore }()
				sl.postLog(bytes.NewReader(b))
			}()
		})
	}

	if sl.buffer.Len() < 1 {
		sl.buffer.WriteRune('[')
	}

	jsonHandler := sl.jsonHandler(sl.buffer)
	if err := jsonHandler.Handle(ctx, record); err != nil {
		return err
	}

	if sl.buffer.Len() < sl.opts.BufferSize-1 {
		sl.buffer.WriteRune(',')
	} else {
		sl.sendDebounceFunc.Stop()
		sl.sendDebounceFunc = nil
		sl.buffer.WriteRune(']')
		buf := sl.extractBuffer()
		sl.semaphore <- struct{}{}
		go func() {
			defer func() { <-sl.semaphore }()
			sl.postLog(buf)
		}()
	}

	return nil
}

func (sl *OpenObserveHandler) extractBuffer() *bytes.Buffer {
	b := sl.buffer.Bytes()
	newb := make([]byte, len(b))
	copy(newb, b)
	if sl.buffer.Cap() > 512*1024 {
		sl.buffer = &bytes.Buffer{}
	} else {
		sl.buffer.Reset()
	}
	return bytes.NewBuffer(newb)
}

func noopReplaceAttr(_ []string, attr slog.Attr) slog.Attr { return attr }

func (sl *OpenObserveHandler) jsonHandler(w io.Writer) slog.Handler {
	handler := slog.NewJSONHandler(w, wrapHandlerOptions(sl.opts.HandlerOptions)).
		WithAttrs(sl.opts.WithAttrs).
		WithAttrs(sl.withAttrs)
	for _, name := range sl.withGroup {
		handler = handler.WithGroup(name)
	}
	return handler
}

func (sl *OpenObserveHandler) postLog(buf io.Reader) {
	req, err := http.NewRequest(http.MethodPost, sl.opts.Endpoint, buf)
	if err != nil {
		fmt.Printf("openobserve: failed to create request: %s\n", err)
		return
	}
	req.SetBasicAuth(sl.opts.Username, sl.opts.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := sl.opts.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("openobserve: failed to execute request: %s\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		fmt.Printf("openobserve: unexpected %d status code from openobserve instance when sending logs\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}

func (sl *OpenObserveHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	sl2 := sl.clone()
	sl2.withAttrs = append(sl2.withAttrs, attrs...)
	return sl2
}

func (sl *OpenObserveHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return sl
	}
	sl2 := sl.clone()
	sl2.withGroup = append(sl2.withGroup, name)
	return sl2
}

func wrapHandlerOptions(in *slog.HandlerOptions) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource:   in.AddSource,
		Level:       in.Level,
		ReplaceAttr: wrapReplaceAttr(in.ReplaceAttr),
	}
}

type replaceAttrFunc = func(group []string, attr slog.Attr) slog.Attr

func wrapReplaceAttr(replaceAttr replaceAttrFunc) replaceAttrFunc {
	return func(group []string, attr slog.Attr) slog.Attr {
		if len(group) > 0 {
			return replaceAttr(group, attr)
		}
		if attr.Key == slog.TimeKey {
			return slog.Attr{
				Key:   "_timestamp",
				Value: attr.Value,
			}
		}
		return replaceAttr(group, attr)
	}
}

type OpenObserveHandlerOptions struct {
	HandlerOptions *slog.HandlerOptions

	// Maximum size for the buffer to store log messages before flushing.
	BufferSize int

	// Maximum time to wait before flushing the buffer.
	BufferTimeout time.Duration

	// Maximum number of concurrent requests to send logs.
	Concurrency int

	// Endpoint to send logs to.
	Endpoint string

	// HTTPClient to use for sending logs.
	HTTPClient HTTPClient

	// Attributes to include in every log message.
	WithAttrs []slog.Attr

	Username string
	Password string
}

func NewOpenObserveHandler(opts OpenObserveHandlerOptions) *OpenObserveHandler {
	if opts.HandlerOptions == nil {
		opts.HandlerOptions = &slog.HandlerOptions{}
	}
	if opts.HandlerOptions.ReplaceAttr == nil {
		opts.HandlerOptions.ReplaceAttr = noopReplaceAttr
	}
	return &OpenObserveHandler{
		opts:      opts,
		semaphore: make(chan struct{}, opts.Concurrency),
		buffer:    &bytes.Buffer{},
	}
}
