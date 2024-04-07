package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
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
	message := fmt.Sprintf("%s %s %d %dms", ch.request.Method, ch.request.URL.Path, status, elapsed.Milliseconds())

	requestLog := slog.Attr{Key: "request", Value: ch.extractRequestLog()}
	responseLog := slog.Group("response", "status", status, "header", header, "bytes", bytes)
	roundtripLog := slog.String("roundtrip", fmt.Sprintf("%dms", elapsed.Milliseconds()))

	if status >= 400 {
		log.Log(ch.request.Context()).Error(message, requestLog, responseLog, roundtripLog)
		return
	}

	log.Log(ch.request.Context()).Info(message, requestLog, responseLog, roundtripLog)
}

func (ch *ChiEntry) Panic(v interface{}, stack []byte) {
	entry := log.Log(ch.request.Context())
	if err, ok := v.(error); ok {
		entry.Err(err).Error("panic occurred", "stack", string(stack))
	} else {
		entry.Error("panic occurred", "panic_data", v, "stack", string(stack))
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
		values = append(values, slog.Any("query", queries))
	}
	values = append(values, slog.Any("headers", ch.request.Header))
	return slog.GroupValue(values...)
}
