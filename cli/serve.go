package cli

import (
	"io/fs"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/db"
	"github.com/tigorlazuardi/redmage/db/queries"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/server"
)

var PublicDir fs.FS = os.DirFS("public")

var serveCmd = &cobra.Command{
	Use:          "serve",
	Short:        "Starts the HTTP Server",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.Open(cfg)
		if err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to connect database")
			os.Exit(1)
		}

		queries := queries.New(db)

		api := api.New(queries, db)

		server := server.New(cfg, api, PublicDir)

		exit := make(chan struct{}, 1)

		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt)
			<-sig
			exit <- struct{}{}
		}()

		if err := server.Start(exit); err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to start server")
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
