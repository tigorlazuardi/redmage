package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/caller"
)

var handler slog.Handler = NullHandler{}

func NewHandler(cfg *config.Config) slog.Handler {
	if !cfg.Bool("log.enable") {
		return NullHandler{}
	}
	var output io.Writer
	if strings.ToLower(cfg.String("log.output")) == "stdout" {
		output = colorable.NewColorableStdout()
	} else {
		output = colorable.NewColorableStderr()
	}

	var lvl slog.Level
	_ = lvl.UnmarshalText(cfg.Bytes("log.level"))
	opts := &slog.HandlerOptions{
		AddSource: cfg.Bool("log.source"),
		Level:     lvl,
	}

	format := strings.ToLower(cfg.String("log.format"))
	if isatty.IsTerminal(os.Stdout.Fd()) && format == "pretty" {
		return NewPrettyHandler(output, opts)
	} else {
		return slog.NewJSONHandler(output, opts)
	}
}

type Entry struct {
	ctx     context.Context
	handler slog.Handler
	caller  caller.Caller
	time    time.Time
}

// Log prepares a new entry to write logs.
func Log(ctx context.Context) *Entry {
	h := FromContext(ctx)
	if h == nil {
		h = handler
	}
	return &Entry{ctx: ctx, handler: h, time: time.Now()}
}

func (entry *Entry) Caller(caller caller.Caller) *Entry {
	entry.caller = caller
	return entry
}

func (entry *Entry) Info(message string, fields ...any) {
	record := slog.NewRecord(entry.time, slog.LevelInfo, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	record.AddAttrs(slog.Group("context", fields...))
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Infof(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	record := slog.NewRecord(entry.time, slog.LevelInfo, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Error(message string, fields ...any) {
	record := slog.NewRecord(entry.time, slog.LevelError, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	record.AddAttrs(slog.Group("context", fields...))
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Errorf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	record := slog.NewRecord(entry.time, slog.LevelError, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Debug(message string, fields ...any) {
	record := slog.NewRecord(entry.time, slog.LevelDebug, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	record.AddAttrs(slog.Group("context", fields...))
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Debugf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	record := slog.NewRecord(entry.time, slog.LevelDebug, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Warn(message string, fields ...any) {
	record := slog.NewRecord(entry.time, slog.LevelWarn, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	record.AddAttrs(slog.Group("context", fields...))
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Warnf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	record := slog.NewRecord(entry.time, slog.LevelWarn, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) getCaller() caller.Caller {
	if entry.caller.PC != 0 {
		return entry.caller
	}
	return caller.New(4)
}

func (entry *Entry) getExtra() []slog.Attr {
	out := make([]slog.Attr, 0, 1)
	if reqid := middleware.GetReqID(entry.ctx); reqid != "" {
		out = append(out, slog.String("request.id", reqid))
	}

	return out
}

func SetDefault(h slog.Handler) {
	handler = h
}