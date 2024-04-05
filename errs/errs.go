package errs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type Error interface {
	error
	Message(msg string, args ...any) Error
	GetMessage() string
	Code(status int) Error
	GetCode() int
	Caller(pc uintptr) Error
	GetCaller() uintptr
	Details(...any) Error
	GetDetails() []any
	Log(ctx context.Context) Error
}

var _ Error = (*Err)(nil)

type Err struct {
	msg     string
	code    int
	caller  uintptr
	details []any
	origin  error
}

func (er *Err) LogValue() slog.Value {
	values := make([]slog.Attr, 0, 5)

	if er.msg != "" {
		values = append(values, slog.String("message", er.msg))
	}

	if er.code != 0 {
		values = append(values, slog.Int("code", er.code))
	}
	if er.caller != 0 {
		frame, _ := runtime.CallersFrames([]uintptr{er.caller}).Next()
		split := strings.Split(frame.Function, string(os.PathSeparator))
		fnName := split[len(split)-1]

		values = append(values, slog.Group("origin",
			slog.String("file", frame.File),
			slog.Int("line", frame.Line),
			slog.String("function", fnName),
		))
	}

	if len(er.details) > 0 {
		values = append(values, slog.Group("details", er.details...))
	}

	values = append(values, slog.Group("error",
		slog.String("type", reflect.TypeOf(er.origin).String()),
		slog.Any("data", er.origin),
	))

	return slog.GroupValue(values...)
}

func (er *Err) Error() string {
	var (
		s      = strings.Builder{}
		source = er.origin
		msg    = source.Error()
		unwrap = errors.Unwrap(source)
	)
	if unwrap == nil {
		if er.msg != "" {
			s.WriteString(er.msg)
			s.WriteString(": ")
		}
		s.WriteString(msg)
		return s.String()
	}
	for unwrap := errors.Unwrap(source); unwrap != nil; source = unwrap {
		originMsg := unwrap.Error()
		var write string
		if cut, found := strings.CutSuffix(msg, originMsg); found {
			write = cut
		} else {
			write = msg
		}
		msg = originMsg
		if write != "" {
			s.WriteString(write)
			s.WriteString(": ")
		}
	}
	return s.String()
}

func (er *Err) Message(msg string, args ...any) Error {
	er.msg = fmt.Sprintf(msg, args...)
	return er
}

func (er *Err) GetMessage() string {
	panic("not implemented") // TODO: Implement
}

func (er *Err) Code(status int) Error {
	panic("not implemented") // TODO: Implement
}

func (er *Err) GetCode() int {
	panic("not implemented") // TODO: Implement
}

func (er *Err) Caller(pc uintptr) Error {
	panic("not implemented") // TODO: Implement
}

func (er *Err) GetCaller() uintptr {
	panic("not implemented") // TODO: Implement
}

func (er *Err) Details(_ ...any) Error {
	panic("not implemented") // TODO: Implement
}

func (er *Err) GetDetails() []any {
	panic("not implemented") // TODO: Implement
}

func (er *Err) Log(ctx context.Context) Error {
	panic("not implemented") // TODO: Implement
}
