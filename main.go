package main

import (
	"context"
	"embed"
	"os"

	"github.com/tigorlazuardi/redmage/cli"
	"github.com/tigorlazuardi/redmage/db"
)

// go:embed db/migrations/*.sql
var migrations embed.FS

func main() {
	db.Migrations = migrations

	if err := cli.RootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
