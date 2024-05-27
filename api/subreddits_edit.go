package api

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type SubredditEditParams struct {
	Name           string
	EnableSchedule *int32
	Schedule       *string
	Countback      *int32
}

func (api *API) SubredditsEdit(ctx context.Context, params SubredditEditParams) (subreddit *models.Subreddit, err error) {
	ctx, span := tracer.Start(ctx, "*API.SubredditsEdit")
	defer span.End()

	now := time.Now()

	subreddit = &models.Subreddit{
		Name: params.Name,
	}

	set := &models.SubredditSetter{
		EnableSchedule: omit.FromPtr(params.EnableSchedule),
		Schedule:       omit.FromPtr(params.Schedule),
		Countback:      omit.FromPtr(params.Countback),
		UpdatedAt:      omit.From(now.Unix()),
	}

	api.lockf(func() {
		err = models.Subreddits.Update(ctx, api.db, set, subreddit)
	})

	if err != nil {
		return subreddit, errs.Wrapw(err, "failed to update subreddit", "set", set)
	}

	if err := subreddit.Reload(ctx, api.db); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return subreddit, errs.Wrapw(err, "subreddit not found", "subreddit", subreddit.Name).Code(404)
		}
		return subreddit, errs.Wrapw(err, "failed to reload subreddit")
	}

	if params.Schedule != nil {
		_, _ = api.scheduler.Put(params.Name, *params.Schedule)
	}

	return subreddit, nil
}
