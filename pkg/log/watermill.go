package log

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
)

type WatermillLogger struct {
	with []any
}

func (wa *WatermillLogger) Error(msg string, err error, fields watermill.LogFields) {
	New(context.Background()).With(wa.with...).With(intoAttrs(fields)...).Err(err).Error(msg)
}

func (wa *WatermillLogger) Info(msg string, fields watermill.LogFields) {
	New(context.Background()).With(wa.with...).With(intoAttrs(fields)...).Info(msg)
}

func (wa *WatermillLogger) Debug(msg string, fields watermill.LogFields) {
	New(context.Background()).With(wa.with...).With(intoAttrs(fields)...).Debug(msg)
}

func (wa *WatermillLogger) Trace(msg string, fields watermill.LogFields) {
	New(context.Background()).With(wa.with...).With(intoAttrs(fields)...).Debug(msg)
}

func (wa *WatermillLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &WatermillLogger{with: append(wa.with, intoAttrs(fields)...)}
}

func intoAttrs(fields watermill.LogFields) (attrs []any) {
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	return attrs
}
