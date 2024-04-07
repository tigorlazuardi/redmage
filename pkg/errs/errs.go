package errs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/tigorlazuardi/redmage/pkg/caller"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

type Error interface {
	error
	Message(msg string, args ...any) Error
	GetMessage() string
	Code(status int) Error
	GetCode() int
	Caller(pc caller.Caller) Error
	GetCaller() caller.Caller
	Details(...any) Error
	GetDetails() []any
	Log(ctx context.Context) Error
}

var _ Error = (*Err)(nil)

type Err struct {
	message string
	code    int
	caller  caller.Caller
	details []any
	origin  error
}

func (er *Err) LogValue() slog.Value {
	values := make([]slog.Attr, 0, 5)

	if er.message != "" {
		values = append(values, slog.String("message", er.message))
	}

	if er.code != 0 {
		values = append(values, slog.Int("code", er.code))
	}
	if er.caller.PC != 0 {
		values = append(values, slog.Any("origin", er.caller))
	}

	if len(er.details) > 0 {
		values = append(values, slog.Group("details", er.details...))
	}

	if er.origin == nil {
		values = append(values, slog.Any("error", er.origin))
	} else {
		errGroupValues := make([]slog.Attr, 0, 3)
		errGroupValues = append(errGroupValues, slog.String("type", reflect.TypeOf(er.origin).String()))
		if lv, ok := er.origin.(slog.LogValuer); ok {
			errGroupValues = append(errGroupValues, slog.Attr{Key: "data", Value: lv.LogValue()})
		} else {
			errGroupValues = append(errGroupValues, slog.String("message", er.origin.Error()), slog.Any("data", er.origin))
		}

		values = append(values, slog.Attr{Key: "error", Value: slog.GroupValue(errGroupValues...)})
	}

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
		if er.message != "" {
			s.WriteString(er.message)
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
	er.message = fmt.Sprintf(msg, args...)
	return er
}

func (er *Err) GetMessage() string {
	return er.message
}

func (er *Err) Code(status int) Error {
	er.code = status
	return er
}

func (er *Err) GetCode() int {
	return er.code
}

func (er *Err) Caller(pc caller.Caller) Error {
	er.caller = pc
	return er
}

func (er *Err) GetCaller() caller.Caller {
	return er.caller
}

func (er *Err) Details(ctx ...any) Error {
	er.details = ctx
	return er
}

func (er *Err) GetDetails() []any {
	return er.details
}

func (er *Err) Log(ctx context.Context) Error {
	log.Log(ctx).Caller(er.caller).Err(er).Error(er.message)
	return er
}

func Wrap(err error) Error {
	return &Err{
		origin: err,
		caller: caller.New(3),
	}
}

func Wrapw(err error, message string, details ...any) Error {
	return &Err{
		origin:  err,
		details: details,
		message: message,
		caller:  caller.New(3),
	}
}

func Wrapf(err error, message string, args ...any) Error {
	message = fmt.Sprintf(message, args...)
	return &Err{
		origin:  err,
		message: message,
		caller:  caller.New(3),
	}
}

func Fail(message string, details ...any) Error {
	return &Err{
		origin:  errors.New(message),
		details: details,
		caller:  caller.New(3),
	}
}

func Failf(message string, args ...any) Error {
	return &Err{
		origin: fmt.Errorf(message, args...),
		caller: caller.New(3),
	}
}
