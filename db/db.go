package db

import (
	"database/sql"
	"io/fs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/XSAM/otelsql"
)

var Migrations fs.FS

func Open(cfg *config.Config) (*sql.DB, error) {
	driver := cfg.String("db.driver")
	dsn := cfg.String("db.string")
	db, err := otelsql.Open(driver, dsn, otelsql.WithAttributes(
		semconv.DBSystemSqlite,
	))
	if err != nil {
		return db, errs.Wrapw(err, "failed to open database", "driver", driver)
	}
	if cfg.Bool("db.automigrate") {
		goose.SetLogger(goose.NopLogger())
		goose.SetBaseFS(Migrations)

		if err := goose.SetDialect(driver); err != nil {
			return db, errs.Wrapw(err, "failed to set goose dialect", "dialect", driver)
		}

		if err := goose.Up(db, "db/migrations"); err != nil {
			return db, errs.Wrapw(err, "failed to migrate database", "dialect", driver)
		}
	}
	return db, err
}

func OpenPubsub(cfg *config.Config) (*sql.DB, error) {
	driver := cfg.String("pubsub.db.driver")
	dsn := cfg.String("pubsub.db.string")
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return db, errs.Wrapw(err, "failed to open database", "driver", driver)
	}
	return db, err
}

func ApplyLogger(cfg *config.Config, db *sql.DB) *sql.DB {
	dsn := cfg.String("db.string")
	return sqldblogger.OpenDriver(dsn, db.Driver(), sqlLogger{},
		sqldblogger.WithSQLQueryAsMessage(true),
	)
}
