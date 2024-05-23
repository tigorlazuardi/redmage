package db

import (
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"

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
	if driver == "sqlite3" {
		path, err := filepath.Abs(dsn)
		if err != nil {
			return nil, errs.Wrapw(err, "failed to get absolute path of sqlite3 database", "path", dsn)
		}
		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return nil, errs.Wrapw(err, "failed to create directory for sqlite3 database", "dir", dir)
		}
	}
	db, err := otelsql.Open(driver, dsn, otelsql.WithAttributes(
		semconv.DBSystemSqlite,
	))
	if err != nil {
		return db, errs.Wrapw(err, "failed to open database", "driver", driver, "db.string", dsn)
	}
	if cfg.Bool("db.automigrate") {
		goose.SetLogger(goose.NopLogger())
		goose.SetBaseFS(Migrations)

		if err := goose.SetDialect(driver); err != nil {
			return db, errs.Wrapw(err, "failed to set goose dialect", "dialect", driver, "dsn", dsn)
		}

		if err := goose.Up(db, "db/migrations"); err != nil {
			return db, errs.Wrapw(err, "failed to migrate database", "dialect", driver, "dsn", dsn)
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

func RunMigrations(cfg *config.Config) error {
	if cfg.Bool("db.automigrate") {
		driver := cfg.String("db.driver")
		dsn := cfg.String("db.string")
		db, err := sql.Open(driver, dsn)
		if err != nil {
			return errs.Wrapw(err, "migration: failed to open database", "driver", driver, "db.string", dsn)
		}
		defer db.Close()
		goose.SetLogger(goose.NopLogger())
		goose.SetBaseFS(Migrations)

		if err := goose.SetDialect(driver); err != nil {
			return errs.Wrapw(err, "failed to set goose dialect", "dialect", driver, "dsn", dsn)
		}

		if err := goose.Up(db, "db/migrations"); err != nil {
			return errs.Wrapw(err, "failed to migrate database", "dialect", driver, "dsn", dsn)
		}
	}
	return nil
}
