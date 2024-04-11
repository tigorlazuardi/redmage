package api

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/redmage/api/bmessage"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/db/queries"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

type API struct {
	queries *queries.Queries
	db      *sql.DB

	scheduler   *cron.Cron
	scheduleMap map[cron.EntryID]queries.Subreddit

	downloadBroadcast *broadcast.Relay[bmessage.ImageDownloadMessage]

	config *config.Config
}

func New(q *queries.Queries, db *sql.DB, cfg *config.Config) *API {
	return &API{
		queries:           q,
		db:                db,
		scheduler:         cron.New(),
		scheduleMap:       make(map[cron.EntryID]queries.Subreddit, 8),
		downloadBroadcast: broadcast.NewRelay[bmessage.ImageDownloadMessage](),
		config:            cfg,
	}
}

func (api *API) StartScheduler(ctx context.Context) error {
	subreddits, err := api.queries.SubredditsGetAll(ctx)
	if err != nil {
		return errs.Wrapw(err, "failed to get all subreddits")
	}

	for _, subreddit := range subreddits {
		err := api.scheduleSubreddit(subreddit)
		if err != nil {
			log.New(ctx).Err(err).Error(
				fmt.Sprintf("failed to start scheduler for subreddit '%s'", subreddit.Name),
				"subreddit", subreddit,
			)
			continue
		}
	}

	return nil
}

func (api *API) scheduleSubreddit(subreddit queries.Subreddit) error {
	id, err := api.scheduler.AddFunc(subreddit.Schedule, func() {
	})
	if err != nil {
		return errs.Wrap(err)
	}

	api.scheduleMap[id] = subreddit

	return nil
}
