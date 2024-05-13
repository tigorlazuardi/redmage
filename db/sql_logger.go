package db

import (
	"context"
	"log/slog"
	"strings"

	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

type sqlLogger struct{}

func (sqlLogger) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	var lvl slog.Level

	msg = strings.TrimSpace(msg)

	switch level {
	case sqldblogger.LevelDebug, sqldblogger.LevelTrace, sqldblogger.LevelInfo:
		lvl = slog.LevelDebug
	case sqldblogger.LevelError:
		lvl = slog.LevelError
	}

	log.New(ctx).With("sql", data).Level(lvl).Log(msg)
}
