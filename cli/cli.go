package cli

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "redmage",
	Short: "Redmage is an HTTP server to download images from Reddit.",
}
