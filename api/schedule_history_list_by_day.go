package api

import (
	"context"
	"time"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type ScheduleHistoryListByDateParams struct {
	Date time.Time
}

func (params *ScheduleHistoryListByDateParams) FillFromQuery(query Queryable) {
	var err error
	params.Date, err = time.Parse("2006-01-02", query.Get("date"))
	if err != nil {
		params.Date = time.Now()
	}
}

func (params *ScheduleHistoryListByDateParams) CountQuery() (expr []bob.Mod[*dialect.SelectQuery]) {
	unixTopTime := time.Date(params.Date.Year(), params.Date.Month(), params.Date.Day(), 23, 59, 59, 0, params.Date.Location()).Unix()
	unixLowTime := time.Date(params.Date.Year(), params.Date.Month(), params.Date.Day(), 0, 0, 0, 0, params.Date.Location()).Unix()
	expr = append(expr, models.SelectWhere.ScheduleHistories.CreatedAt.GTE(unixLowTime))
	expr = append(expr, models.SelectWhere.ScheduleHistories.CreatedAt.LTE(unixTopTime))
	return
}

func (params *ScheduleHistoryListByDateParams) Query() (expr []bob.Mod[*dialect.SelectQuery]) {
	expr = params.CountQuery()
	return
}

func (api *API) ScheduleHistoryListByDate(ctx context.Context, params ScheduleHistoryListByDateParams) (result ScheduleHistoryListResult, err error) {
	ctx, span := tracer.Start(ctx, "*API.ScheduleHistoryListByDate")
	defer span.End()

	result.Schedules, err = models.ScheduleHistories.Query(ctx, api.db, params.Query()...).All()
	if err != nil {
		return result, errs.Wrapw(err, "failed to list schedule histories", "query", params)
	}

	result.Total, err = models.ScheduleHistories.Query(ctx, api.db, params.CountQuery()...).Count()
	if err != nil {
		return result, errs.Wrapw(err, "failed to count schedule histories", "query", params)
	}

	return
}
