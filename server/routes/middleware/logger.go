package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

type ChiLogger struct{}

func (ChiLogger) NewLogEntry(r *http.Request) chimiddleware.LogEntry {
	return &ChiEntry{request: r}
}

type ChiEntry struct {
	request *http.Request
}

func (ch *ChiEntry) Write(status int, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	elasedStr := formatDuration(elapsed)
	message := fmt.Sprintf("%s %s %d %s", ch.request.Method, ch.request.URL, status, elasedStr)

	requestLog := slog.Attr{Key: "request", Value: ch.extractRequestLog()}
	responseLog := slog.Group("response", "status", status, "headers", flat(header), "bytes", bytes)
	roundtripLog := slog.String("elapsed", elasedStr)

	group := slog.Group("http", requestLog, responseLog, roundtripLog)
	if status >= 400 {
		log.New(ch.request.Context()).With(group).Error(message)
		return
	}

	log.New(ch.request.Context()).With(group).Info(message)
}

func (ch *ChiEntry) Panic(v interface{}, stack []byte) {
	group := slog.Group("http", slog.Attr{Key: "request", Value: ch.extractRequestLog()})
	entry := log.New(ch.request.Context())
	message := fmt.Sprintf("[PANIC] %s %s", ch.request.Method, ch.request.URL)
	if err, ok := v.(error); ok {
		entry.Err(err).With(group).Error(message, "stack", string(stack))
	} else {
		entry.With(group).Error(message, "panic_data", v, "stack", string(stack))
	}
}

func (ch *ChiEntry) extractRequestLog() slog.Value {
	values := make([]slog.Attr, 0, 4)
	values = append(values,
		slog.String("method", ch.request.Method),
		slog.String("path", ch.request.URL.Path),
	)
	queries := ch.request.URL.Query()
	if len(queries) > 0 {
		values = append(values, slog.Any("query", flat(queries)))
	}
	values = append(values, slog.Any("headers", flat(ch.request.Header)))
	return slog.GroupValue(values...)
}

func flat(header map[string][]string) map[string]string {
	m := make(map[string]string, len(header))

	for k := range header {
		m[k] = strings.Join(header[k], ", ")
	}
	return m
}

func formatDuration(dur time.Duration) string {
	nanosecs := float64(dur)

	return fmt.Sprintf("%.3fms", nanosecs/float64(time.Millisecond))
}

type ChiSimpleLogger struct{}

func (ChiSimpleLogger) NewLogEntry(r *http.Request) chimiddleware.LogEntry {
	return &ChiSimpleEntry{request: r}
}

type ChiSimpleEntry struct {
	request *http.Request
}

func (ch *ChiSimpleEntry) Write(status int, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	elapsedStr := formatDuration(elapsed)
	message := fmt.Sprintf("%s %s %d %s", ch.request.Method, ch.request.URL, status, elapsedStr)

	level := slog.LevelInfo
	if status >= 400 {
		level = slog.LevelError
	}
	log.New(ch.request.Context()).Level(level).Log(message)
}

func (ch *ChiSimpleEntry) Panic(v interface{}, stack []byte) {
	(&ChiEntry{ch.request}).Panic(v, stack)
}
