package db

import (
	"context"
	"strings"

	"github.com/tigorlazuardi/redmage/pkg/caller"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

type gooseLogger struct{}

func (gl *gooseLogger) Fatalf(format string, v ...interface{}) {
	format = strings.TrimSuffix(format, "\n")
	log.Log(context.Background()).Caller(caller.New(3)).Errorf(format, v...)
}

func (gl *gooseLogger) Printf(format string, v ...interface{}) {
	format = strings.TrimSuffix(format, "\n")
	log.Log(context.Background()).Caller(caller.New(3)).Infof(format, v...)
}
