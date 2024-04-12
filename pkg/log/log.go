package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/caller"
	"go.opentelemetry.io/otel/trace"
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
	ctx       context.Context
	handler   slog.Handler
	caller    caller.Caller
	time      time.Time
	err       error
	level     slog.Level
	withAttrs []slog.Attr
}

// New prepares a new entry to write logs.
func New(ctx context.Context) *Entry {
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

func (entry *Entry) Err(err error) *Entry {
	entry.err = err
	return entry
}

func (entry *Entry) Level(lvl slog.Level) *Entry {
	entry.level = lvl
	return entry
}

// With adds fields to top level of the log entry.
func (entry *Entry) With(fields ...any) *Entry {
	entry.withAttrs = append(entry.withAttrs, argsToAttrSlice(fields)...)
	return entry
}

const badKey = "!BADKEY"

func argsToAttrSlice(args []any) []slog.Attr {
	var (
		attr  slog.Attr
		attrs []slog.Attr
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}

func argsToAttr(args []any) (slog.Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String(badKey, x), nil
		}
		return slog.Any(x, args[1]), args[2:]

	case slog.Attr:
		return x, args[1:]

	default:
		return slog.Any(badKey, x), args[1:]
	}
}

func (entry *Entry) Log(message string, fields ...any) {
	if !entry.handler.Enabled(entry.ctx, entry.level) {
		return
	}
	record := slog.NewRecord(entry.time, entry.level, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	record.AddAttrs(slog.Group("details", fields...))
	if entry.err != nil {
		record.AddAttrs(slog.Any("error", entry.err))
	}
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Info(message string, fields ...any) {
	if !entry.handler.Enabled(entry.ctx, slog.LevelInfo) {
		return
	}
	record := slog.NewRecord(entry.time, slog.LevelInfo, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	record.AddAttrs(slog.Group("details", fields...))
	if entry.err != nil {
		record.AddAttrs(slog.Any("error", entry.err))
	}
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Infof(format string, args ...any) {
	if !entry.handler.Enabled(entry.ctx, slog.LevelInfo) {
		return
	}
	message := fmt.Sprintf(format, args...)
	record := slog.NewRecord(entry.time, slog.LevelInfo, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	if entry.err != nil {
		record.AddAttrs(slog.Any("error", entry.err))
	}
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Error(message string, fields ...any) {
	if !entry.handler.Enabled(entry.ctx, slog.LevelError) {
		return
	}
	record := slog.NewRecord(entry.time, slog.LevelError, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	record.AddAttrs(slog.Group("details", fields...))
	if entry.err != nil {
		record.AddAttrs(slog.Any("error", entry.err))
	}
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Errorf(format string, args ...any) {
	if !entry.handler.Enabled(entry.ctx, slog.LevelError) {
		return
	}
	message := fmt.Sprintf(format, args...)
	record := slog.NewRecord(entry.time, slog.LevelError, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	if entry.err != nil {
		record.AddAttrs(slog.Any("details", entry.err))
	}
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Debug(message string, fields ...any) {
	if !entry.handler.Enabled(entry.ctx, slog.LevelDebug) {
		return
	}
	record := slog.NewRecord(entry.time, slog.LevelDebug, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	record.AddAttrs(slog.Group("details", fields...))
	if entry.err != nil {
		record.AddAttrs(slog.Any("error", entry.err))
	}
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Debugf(format string, args ...any) {
	if !entry.handler.Enabled(entry.ctx, slog.LevelDebug) {
		return
	}
	message := fmt.Sprintf(format, args...)
	record := slog.NewRecord(entry.time, slog.LevelDebug, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	if entry.err != nil {
		record.AddAttrs(slog.Any("error", entry.err))
	}
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Warn(message string, fields ...any) {
	if !entry.handler.Enabled(entry.ctx, slog.LevelWarn) {
		return
	}
	record := slog.NewRecord(entry.time, slog.LevelWarn, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	record.AddAttrs(slog.Group("details", fields...))
	if entry.err != nil {
		record.AddAttrs(slog.Any("error", entry.err))
	}
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) Warnf(format string, args ...any) {
	if !entry.handler.Enabled(entry.ctx, slog.LevelWarn) {
		return
	}
	message := fmt.Sprintf(format, args...)
	record := slog.NewRecord(entry.time, slog.LevelWarn, message, entry.getCaller().PC)
	record.AddAttrs(entry.getExtra()...)
	if entry.err != nil {
		record.AddAttrs(slog.Any("error", entry.err))
	}
	_ = entry.handler.Handle(entry.ctx, record)
}

func (entry *Entry) getCaller() caller.Caller {
	if entry.caller.PC != 0 {
		return entry.caller
	}
	return caller.New(4)
}

func (entry *Entry) getExtra() []slog.Attr {
	out := make([]slog.Attr, 0, 4)
	if span := trace.SpanFromContext(entry.ctx); span.IsRecording() {
		out = append(out,
			slog.String("trace.id", span.SpanContext().TraceID().String()),
			slog.String("span.id", span.SpanContext().SpanID().String()),
		)
	}

	out = append(out, entry.withAttrs...)

	return out
}

func SetDefault(h slog.Handler) {
	handler = h
}
