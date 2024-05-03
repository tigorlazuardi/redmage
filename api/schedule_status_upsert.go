package api

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) ScheduleStatusUpsert(ctx context.Context, params ScheduleSetParams) (schedule *models.ScheduleStatus, err error) {
	ctx, span := tracer.Start(ctx, "*API.CreateNewScheduleStatus")
	defer span.End()
	return api.scheduleStatusUpsert(ctx, api.db, params)
}

func (api *API) scheduleStatusUpsert(ctx context.Context, exec bob.Executor, params ScheduleSetParams) (schedule *models.ScheduleStatus, err error) {
	ctx, span := tracer.Start(ctx, "*API.createNewScheduleStatus")
	defer span.End()
	now := time.Now()
	ss, err := models.ScheduleStatuses.Upsert(ctx, exec, true, []string{"subreddit"}, []string{
		"subreddit",
		"status",
		"error_message",
		"updated_at",
	}, &models.ScheduleStatusSetter{
		Subreddit:    omit.FromCond(params.Subreddit, params.Subreddit != ""),
		Status:       omit.From(params.Status.Int8()),
		ErrorMessage: omit.From(params.ErrorMessage),
		CreatedAt:    omit.From(now.Unix()),
		UpdatedAt:    omit.From(now.Unix()),
	})
	if err != nil {
		return ss, errs.Wrapw(err, "failed to upsert schedule status", "params", params)
	}
	return ss, err
}
