package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/tigorlazuardi/redmage/cli"
	"github.com/tigorlazuardi/redmage/db"

	_ "github.com/tigorlazuardi/redmage/tools"
)

//go:embed db/migrations/*.sql
var Migrations embed.FS

//go:embed public/*
var PublicFS embed.FS

func main() {
	_ = godotenv.Load()
	db.Migrations = Migrations
	var err error
	cli.PublicDir, err = fs.Sub(PublicFS, "public")
	if err != nil {
		panic(errors.New("failed to create sub filesystem"))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cli.RootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
