package cli

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/db"
	"github.com/tigorlazuardi/redmage/db/queries"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/telemetry"
	"github.com/tigorlazuardi/redmage/server"
)

var PublicDir fs.FS = os.DirFS("public")

var serveCmd = &cobra.Command{
	Use:          "serve",
	Short:        "Starts the HTTP Server",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		tele, err := telemetry.New(cmd.Context(), cfg)
		if err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to start telemetry")
			os.Exit(1)
		}
		defer tele.Close()

		db, err := db.Open(cfg)
		if err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to connect database")
			os.Exit(1)
		}

		queries := queries.New(db)

		red := &reddit.Reddit{
			Client: http.DefaultClient,
			Config: cfg,
		}

		api := api.New(api.Dependencies{
			Queries: queries,
			DB:      db,
			Config:  cfg,
			Reddit:  red,
		})

		server := server.New(cfg, api, PublicDir)

		if err := server.Start(cmd.Context().Done()); err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to start server")
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
