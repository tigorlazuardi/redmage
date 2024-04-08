package cli

import (
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/redmage/db"
	"github.com/tigorlazuardi/redmage/db/queries/subreddits"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/server"
	"github.com/tigorlazuardi/redmage/server/routes/api"
	"github.com/tigorlazuardi/redmage/server/routes/www"
)

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

		subreddits := subreddits.New(db)

		api := &api.API{
			Subreddits: subreddits,
		}

		www := &www.WWW{
			Subreddits: subreddits,
		}
		server := server.New(cfg, api, www)

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
