package cli

import (
	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

var serveCmd = &cobra.Command{
	Use:          "serve",
	Short:        "Starts the HTTP Server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostPort := cfg.String("http.host") + ":" + cfg.String("http.port")

		log.Log(cmd.Context()).Info("starting http server", "host", hostPort)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
