package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/redmage/config"
)

var RootCmd = &cobra.Command{
	Use:   "redmage",
	Short: "Redmage is an HTTP server to download images from Reddit.",
}

func init() {
	flags := RootCmd.PersistentFlags()
	for key, value := range config.DefaultConfig {
		const usage = ""
		switch v := value.(type) {
		case bool:
			flags.Bool(key, v, usage)
		case string:
			flags.String(key, v, usage)
		case int:
			flags.Int(key, v, usage)
		case float32:
			flags.Float32(key, v, usage)
		case float64:
			flags.Float64(key, v, usage)
		default:
			flags.String(key, fmt.Sprintf("%v", v), usage)
		}
	}

	cobra.OnInitialize(initConfig)
}
