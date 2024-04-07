package caller

import (
	"log/slog"
	"os"
	"runtime"
	"strings"
)

type Caller struct {
	PC    uintptr
	Frame runtime.Frame
}

func (ca Caller) File() string {
	return ca.Frame.File
}

func (ca Caller) ShortFile() string {
	wd, err := os.Getwd()
	if err != nil {
		return ca.Frame.File
	}
	if after, found := strings.CutPrefix(ca.Frame.File, wd); found {
		return strings.TrimPrefix(after, string(os.PathSeparator))
	}
	return ca.Frame.File
}

func (ca Caller) Line() int {
	return ca.Frame.Line
}

func (ca Caller) Function() string {
	return ca.Frame.Function
}

func (ca Caller) ShortFunction() string {
	split := strings.Split(ca.Frame.Function, string(os.PathSeparator))
	return split[len(split)-1]
}

func (ca Caller) LogValue() slog.Value {
	if ca.PC == 0 {
		return slog.AnyValue(nil)
	}

	return slog.GroupValue(
		slog.String("file", ca.ShortFile()),
		slog.Int("line", ca.Line()),
		slog.String("function", ca.ShortFunction()),
	)
}

func New(skip int) Caller {
	var c Caller
	pcs := make([]uintptr, 1)
	n := runtime.Callers(skip, pcs)
	if n == 0 {
		return c
	}
	c.PC = pcs[0]
	c.Frame, _ = runtime.CallersFrames(pcs).Next()
	return c
}

func From(pc uintptr) Caller {
	var c Caller
	c.PC = pc
	c.Frame, _ = runtime.CallersFrames([]uintptr{pc}).Next()
	return c
}
