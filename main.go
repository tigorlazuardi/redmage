package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"os"

	"github.com/joho/godotenv"
	"github.com/tigorlazuardi/redmage/cli"
	"github.com/tigorlazuardi/redmage/db"
	"github.com/tigorlazuardi/redmage/server/routes/www"
)

//go:embed db/migrations/*.sql
var Migrations embed.FS

//go:embed public/*
var PublicFS embed.FS

func main() {
	_ = godotenv.Load()
	db.Migrations = Migrations
	var err error
	www.PublicFS, err = fs.Sub(PublicFS, "public")
	if err != nil {
		panic(errors.New("failed to create sub filesystem"))
	}
	if err := cli.RootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
