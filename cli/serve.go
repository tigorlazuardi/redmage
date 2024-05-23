package cli

import (
	"io/fs"
	"os"

	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/db"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/pubsub"
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

		if err := db.RunMigrations(cfg); err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to run migrations")
			os.Exit(1)
		}

		database, err := db.NewComfy(cfg)
		if err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to open database")
			os.Exit(1)
		}

		// database, err := db.Open(cfg)
		// if err != nil {
		// 	log.New(cmd.Context()).Err(err).Error("failed to open connection to database")
		// 	os.Exit(1)
		// }

		pubsubDB, err := pubsub.New(cfg)
		if err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to open connection to pubsub database")
			os.Exit(1)
		}

		publisher, err := pubsub.NewPublisher(pubsubDB)
		if err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to create publisher")
			os.Exit(1)
		}
		subscriber, err := pubsub.NewSubscriber(cfg, pubsubDB)
		if err != nil {
			log.New(cmd.Context()).Err(err).Error("failed to create subscriber")
			os.Exit(1)
		}

		// loggedDb := db.ApplyLogger(cfg, database)
		// if err != nil {
		// 	log.New(cmd.Context()).Err(err).Error("failed to connect database")
		// 	os.Exit(1)
		// }
		red := &reddit.Reddit{
			Client: reddit.NewRedditHTTPClient(cfg),
			Config: cfg,
		}

		api := api.New(api.Dependencies{
			DB:         database,
			Config:     cfg,
			Reddit:     red,
			Publisher:  publisher,
			Subscriber: subscriber,
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
