package cli

import (
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/server"
)

var serveCmd = &cobra.Command{
	Use:          "serve",
	Short:        "Starts the HTTP Server",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		server := server.New(cfg)

		exit := make(chan struct{}, 1)

		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt)
			<-sig
			exit <- struct{}{}
		}()

		if err := server.Start(exit); err != nil {
			log.Log(cmd.Context()).Err(err).Error("failed to start server")
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
