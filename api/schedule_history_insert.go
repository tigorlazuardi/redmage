package api

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) ScheduleHistoryInsert(ctx context.Context, params ScheduleSetParams) (history *models.ScheduleHistory, err error) {
	ctx, span := tracer.Start(ctx, "*API.ScheduleHistoryInsert")
	defer span.End()

	return api.scheduleHistoryInsert(ctx, api.db, params)
}

func (api *API) scheduleHistoryInsert(ctx context.Context, exec bob.Executor, params ScheduleSetParams) (history *models.ScheduleHistory, err error) {
	ctx, span := tracer.Start(ctx, "*API.scheduleHistoryInsert")
	defer span.End()

	now := time.Now()

	api.lockf(func() {
		history, err = models.ScheduleHistories.Insert(ctx, exec, &models.ScheduleHistorySetter{
			Subreddit:    omit.FromCond(params.Subreddit, params.Subreddit != ""),
			Status:       omit.From(params.Status.Int8()),
			ErrorMessage: omit.FromCond(params.ErrorMessage, params.Status == ScheduleStatusError),
			CreatedAt:    omit.From(now.Unix()),
		})
	})
	if err != nil {
		return history, errs.Wrapw(err, "failed to insert schedule history", "params", params)
	}
	return history, err
}
