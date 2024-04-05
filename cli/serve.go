package cli

import "github.com/spf13/cobra"

var serveCmd = &cobra.Command{
	Use:          "serve",
	Short:        "Starts the HTTP Server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
