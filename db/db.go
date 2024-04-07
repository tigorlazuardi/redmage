package db

import (
	"database/sql"
	"io/fs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

var Migrations fs.FS

func Open(cfg *config.Config) (*sql.DB, error) {
	driver := cfg.String("db.driver")
	db, err := sql.Open(driver, cfg.String("db.string"))
	if err != nil {
		return db, errs.Wrapw(err, "failed to open database", "driver", driver)
	}

	if cfg.Bool("db.automigrate") {
		goose.SetLogger(&gooseLogger{})
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
