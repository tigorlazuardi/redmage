package main

import (
	"context"
	"embed"
	"os"

	"github.com/tigorlazuardi/redmage/cli"
	"github.com/tigorlazuardi/redmage/db"
)

//go:embed db/migrations/*.sql
var Migrations embed.FS

func main() {
	db.Migrations = Migrations
	if err := cli.RootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
