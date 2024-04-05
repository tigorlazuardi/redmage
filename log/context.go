package log

import (
	"context"
	"log/slog"
)

type loggerKey struct{}

func FromContext(ctx context.Context) slog.Handler {
	h, _ := ctx.Value(loggerKey{}).(slog.Handler)
	return h
}

func WithContext(ctx context.Context, l slog.Handler) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

func NullHandlerContext(ctx context.Context) context.Context {
	return WithContext(ctx, NullHandler{})
}
