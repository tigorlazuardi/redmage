package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/tigorlazuardi/redmage/pkg/caller"
)

type PrettyHandler struct {
	mu          sync.Mutex
	opts        *slog.HandlerOptions
	output      io.Writer
	replaceAttr func(groups []string, attr slog.Attr) slog.Attr
	withAttrs   []slog.Attr
	withGroup   []string
}

// NewPrettyHandler creates a human friendly readable logs.
func NewPrettyHandler(writer io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	if opts.ReplaceAttr == nil {
		opts.ReplaceAttr = func(groups []string, attr slog.Attr) slog.Attr { return attr }
	}
	return &PrettyHandler{
		opts:   opts,
		output: writer,
		replaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if len(groups) > 0 {
				return opts.ReplaceAttr(groups, attr)
			}
			switch attr.Key {
			case slog.TimeKey, slog.LevelKey, slog.SourceKey, slog.MessageKey:
				return slog.Attr{}
			default:
				return opts.ReplaceAttr(groups, attr)
			}
		},
	}
}

// Enabled implements slog.Handler interface.
func (pr *PrettyHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return pr.opts.Level.Level() <= lvl
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		buf := &bytes.Buffer{}
		buf.Grow(1024)
		return buf
	},
}

func putBuffer(buf *bytes.Buffer) {
	const limit = 1024 * 512 // 512KB
	if buf.Cap() < limit {
		buf.Reset()
		bufferPool.Put(buf)
	}
}

// Handle implements slog.Handler interface.
func (pr *PrettyHandler) Handle(ctx context.Context, record slog.Record) error {
	var levelColor *color.Color
	switch {
	case record.Level >= slog.LevelError:
		levelColor = color.New(color.FgRed)
	case record.Level >= slog.LevelWarn:
		levelColor = color.New(color.FgYellow)
	case record.Level >= slog.LevelInfo:
		levelColor = color.New(color.FgGreen)
	default:
		levelColor = color.New(color.FgWhite)
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	jsonBuf := bufferPool.Get().(*bytes.Buffer)
	defer putBuffer(buf)
	defer putBuffer(jsonBuf)

	if record.PC != 0 && pr.opts.AddSource {
		frame := caller.From(record.PC).Frame
		levelColor.Fprint(buf, frame.File)
		levelColor.Fprint(buf, ":")
		levelColor.Fprint(buf, frame.Line)
		levelColor.Fprint(buf, " -- ")
		split := strings.Split(frame.Function, string(os.PathSeparator))
		fnName := split[len(split)-1]
		levelColor.Fprint(buf, fnName)
		buf.WriteByte('\n')
	}

	if !record.Time.IsZero() {
		buf.WriteString(record.Time.Format("[2006-01-02 15:04:05] "))
	}

	buf.WriteByte('[')
	levelColor.Add(color.Bold).Fprint(buf, record.Level.String())
	buf.WriteString("] ")

	if record.Message != "" {
		buf.WriteString(record.Message)
	}

	buf.WriteByte('\n')

	serializer := pr.createSerializer(jsonBuf)
	_ = serializer.Handle(ctx, record)
	if jsonBuf.Len() > 3 { // Ignore empty json like "{}\n"
		_ = json.Indent(buf, jsonBuf.Bytes(), "", "  ")
		// json indent includes new line, no need to add extra new line.
	} else {
		buf.WriteByte('\n')
	}

	pr.mu.Lock()
	defer pr.mu.Unlock()
	_, err := buf.WriteTo(pr.output)
	return err
}

func (pr *PrettyHandler) createSerializer(w io.Writer) slog.Handler {
	var jsonHandler slog.Handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: pr.replaceAttr,
	})

	if len(pr.withAttrs) > 0 {
		jsonHandler = jsonHandler.WithAttrs(pr.withAttrs)
	}

	if len(pr.withGroup) > 0 {
		for _, group := range pr.withGroup {
			jsonHandler = jsonHandler.WithGroup(group)
		}
	}

	return jsonHandler
}

func (pr *PrettyHandler) clone() *PrettyHandler {
	return &PrettyHandler{
		opts:        pr.opts,
		output:      pr.output,
		replaceAttr: pr.replaceAttr,
		withAttrs:   pr.withAttrs,
		withGroup:   pr.withGroup,
	}
}

// WithAttrs implements slog.Handler interface.
func (pr *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	p := pr.clone()
	p.withAttrs = append(p.withAttrs, attrs...)
	return p
}

// WithGroup implements slog.Handler interface.
func (pr *PrettyHandler) WithGroup(name string) slog.Handler {
	p := pr.clone()
	p.withGroup = append(p.withGroup, name)
	return p
}
