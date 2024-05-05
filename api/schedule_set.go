package api

import (
	"context"

	"github.com/stephenafamo/bob"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type ScheduleStatus int8

const (
	ScheduleStatusDisabled ScheduleStatus = iota
	ScheduleStatusEnabled
	ScheduleStatusStandby
	ScheduleStatusEnqueued
	ScheduleStatusDownloading
	ScheduleStatusError
)

func (ss ScheduleStatus) String() string {
	switch ss {
	case ScheduleStatusDisabled:
		return "Disabled"
	case ScheduleStatusEnabled:
		return "Enabled"
	case ScheduleStatusStandby:
		return "Standby"
	case ScheduleStatusEnqueued:
		return "Enqueued"
	case ScheduleStatusDownloading:
		return "Downloading"
	case ScheduleStatusError:
		return "Error"
	}
	return "Unknown"
}

func (ss ScheduleStatus) Int8() int8 {
	return int8(ss)
}

type ScheduleSetParams struct {
	Subreddit    string
	Status       ScheduleStatus
	ErrorMessage string
}

func (api *API) ScheduleSet(ctx context.Context, params ScheduleSetParams) (schedule *models.ScheduleStatus, err error) {
	ctx, span := tracer.Start(ctx, "*API.ScheduleSet")
	defer span.End()

	errTx := api.withTransaction(ctx, func(exec bob.Executor) error {
		schedule, err = api.scheduleSet(ctx, exec, params)
		return err
	})

	return schedule, errTx
}

func (api *API) scheduleSet(ctx context.Context, exec bob.Executor, params ScheduleSetParams) (schedule *models.ScheduleStatus, err error) {
	ctx, span := tracer.Start(ctx, "*API.scheduleSet")
	defer span.End()

	schedule, err = api.scheduleStatusUpsert(ctx, exec, params)
	if err != nil {
		return schedule, errs.Wrapw(err, "failed to set schedule status", "params", params)
	}

	_, err = api.scheduleHistoryInsert(ctx, exec, params)
	if err != nil {
		return schedule, errs.Wrapw(err, "failed to insert schedule history", "params", params)
	}

	// TODO: Create cron job schedule rebalancing if schedule is set to disabled.
	return schedule, nil
}
