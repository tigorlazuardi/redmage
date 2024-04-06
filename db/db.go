package db

import (
	"database/sql"
	"io/fs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/tigorlazuardi/redmage/config"
)

var Migrations fs.FS

func Open(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.String("db.driver"), cfg.String("db.string"))
	if err != nil {
		return db, err
	}

	if cfg.Bool("db.automigrate") {
		goose.SetBaseFS(Migrations)

		if err := goose.SetDialect(cfg.String("db.driver")); err != nil {
			return db, err
		}

		if err := goose.Up(db, "db/migrations"); err != nil {
			return db, err
		}
	}

	return db, err
}
