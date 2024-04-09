package api

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/redmage/db/queries"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

type API struct {
	queries *queries.Queries
	db      *sql.DB

	scheduler   *cron.Cron
	scheduleMap map[cron.EntryID]queries.Subreddit

	downloadBroadcast *broadcast.Relay[DownloadStatusMessage]
}

func New(q *queries.Queries, db *sql.DB) *API {
	return &API{
		queries:           q,
		db:                db,
		scheduler:         cron.New(),
		scheduleMap:       make(map[cron.EntryID]queries.Subreddit, 8),
		downloadBroadcast: broadcast.NewRelay[DownloadStatusMessage](),
	}
}

func (api *API) StartScheduler(ctx context.Context) error {
	subreddits, err := api.queries.SubredditsGetAll(ctx)
	if err != nil {
		return errs.Wrapw(err, "failed to get all subreddits")
	}

	for _, subreddit := range subreddits {
		id, err := api.scheduler.AddFunc(subreddit.Schedule, func() {
			//  TODO: Add download
		})
		if err != nil {
			log.
				New(ctx).
				Err(err).
				Error(
					fmt.Sprintf("failed to start scheduler for subreddit '%s'", subreddit.Name),
					"subreddit", subreddit,
				)
			continue
		}
		api.scheduleMap[id] = subreddit
	}

	return nil
}
