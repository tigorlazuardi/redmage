package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/mattn/go-sqlite3"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) SubredditsCreate(ctx context.Context, params *models.Subreddit) (subreddit *models.Subreddit, err error) {
	ctx, span := tracer.Start(ctx, "*API.SubredditsCreate")
	defer span.End()

	now := time.Now()
	set := &models.SubredditSetter{
		Name:           omit.From(params.Name),
		EnableSchedule: omit.From(params.EnableSchedule),
		Subtype:        omit.From(params.Subtype),
		Schedule:       omit.From(params.Schedule),
		Countback:      omit.From(params.Countback),
		CreatedAt:      omit.From(now.Unix()),
		UpdatedAt:      omit.From(now.Unix()),
	}

	api.lockf(func() {
		subreddit, err = models.Subreddits.Insert(ctx, api.db, set)
	})
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return subreddit, errs.Wrapw(err, "subreddit with that name already exists", "set", set).Code(http.StatusConflict)
			}
			return subreddit, errs.Wrapw(err, "failed to create subreddit", "set", set)
		}
	}

	_, err = api.ScheduleSet(ctx, ScheduleSetParams{
		Subreddit: subreddit.Name,
		Status:    ScheduleStatus(params.EnableSchedule), // Possible value should only be 0 or 1
	})
	if err != nil {
		return subreddit, errs.Wrapw(err, "failed to set schedule status")
	}

	if params.EnableSchedule == 1 {
		_, err = api.scheduler.Put(subreddit.Name, subreddit.Schedule)
		if err != nil {
			return subreddit, errs.Wrapw(err, "failed to put job to scheduler")
		}
	}

	return subreddit, nil
}
