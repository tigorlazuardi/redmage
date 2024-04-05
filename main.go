package main

import (
	"context"
	"os"

	"github.com/tigorlazuardi/redmage/cli"
)

func main() {
	if err := cli.RootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
