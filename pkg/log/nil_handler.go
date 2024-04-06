package log

import (
	"context"
	"log/slog"
)

type NullHandler struct{}

func (NullHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (NullHandler) Handle(context.Context, slog.Record) error { return nil }
func (nu NullHandler) WithAttrs([]slog.Attr) slog.Handler     { return nu }
func (nu NullHandler) WithGroup(string) slog.Handler          { return nu }
