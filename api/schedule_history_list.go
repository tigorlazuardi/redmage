package api

import (
	"context"
	"strconv"
	"time"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/tigorlazuardi/redmage/api/utils"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type ScheduleHistoryListParams struct {
	Subreddit string
	Time      time.Time
	Direction string

	Limit  int64
	Offset int64
}

func (params *ScheduleHistoryListParams) FillFromQuery(query Queryable) {
	params.Subreddit = query.Get("subreddit")
	params.Direction = query.Get("direction")
	params.Limit, _ = strconv.ParseInt(query.Get("limit"), 10, 64)
	if params.Limit < 1 {
		params.Limit = 100
	}
	if params.Limit > 1000 {
		params.Limit = 1000
	}

	params.Offset, _ = strconv.ParseInt(query.Get("offset"), 10, 64)
	if params.Offset < 0 {
		params.Offset = 0
	}

	now := time.Now()

	timeInt, _ := strconv.ParseInt(query.Get("time"), 10, 64)
	if timeInt > 0 {
		params.Time = time.Unix(timeInt, 0)
	} else if timeInt < 0 {
		params.Time = now.Add(time.Duration(timeInt) * time.Second)
	}
	if params.Time.After(now) {
		params.Time = time.Time{}
	}
}

func (params ScheduleHistoryListParams) CountQuery() (expr []bob.Mod[*dialect.SelectQuery]) {
	if params.Subreddit != "" {
		expr = append(expr, models.SelectWhere.ScheduleHistories.Subreddit.EQ(params.Subreddit))
	}
	if !params.Time.IsZero() {
		if params.Direction == "before" {
			expr = append(expr,
				models.SelectWhere.ScheduleHistories.CreatedAt.GTE(params.Time.Unix()),
			)
		} else {
			expr = append(expr, models.SelectWhere.ScheduleHistories.CreatedAt.LT(params.Time.Unix()))
		}
	}

	return expr
}

func (params ScheduleHistoryListParams) Query() (expr []bob.Mod[*dialect.SelectQuery]) {
	expr = append(expr, params.CountQuery()...)
	if params.Limit > 0 {
		expr = append(expr, sm.Limit(params.Limit))
	}
	if params.Offset > 0 {
		expr = append(expr, sm.Offset(params.Offset))
	}
	expr = append(expr, sm.OrderBy(models.ScheduleHistoryColumns.CreatedAt).Desc())

	return expr
}

type ScheduleHistoryListResult struct {
	Schedules models.ScheduleHistorySlice `json:"schedules"`
	Total     int64                       `json:"count"`
}

func (result ScheduleHistoryListResult) GetLast() *models.ScheduleHistory {
	if len(result.Schedules) > 0 {
		return result.Schedules[len(result.Schedules)-1]
	}
	return nil
}

func (result ScheduleHistoryListResult) GetLastTime() time.Time {
	if schedule := result.GetLast(); schedule != nil {
		return time.Unix(schedule.CreatedAt, 0)
	}
	return time.Now()
}

func (result ScheduleHistoryListResult) GetFirstTime() time.Time {
	if schedule := result.GetFirst(); schedule != nil {
		return time.Unix(schedule.CreatedAt, 0)
	}
	return time.Now()
}

func (result ScheduleHistoryListResult) GetFirst() *models.ScheduleHistory {
	if len(result.Schedules) > 0 {
		return result.Schedules[0]
	}
	return nil
}

func (result ScheduleHistoryListResult) SplitByDay() (out []ScheduleHistoryListResultDay) {
	out = make([]ScheduleHistoryListResultDay, 0, 4)

	var lastDay time.Time
	var lastIdx int
	for _, schedule := range result.Schedules {
		t := utils.StartOfDay(time.Unix(schedule.CreatedAt, 0).In(time.Local))
		if !t.Equal(lastDay) {
			out = append(out, ScheduleHistoryListResultDay{
				Date: t,
			})
			lastDay = t
			lastIdx = len(out) - 1

			out[lastIdx].Schedules = append(out[lastIdx].Schedules, schedule)
			out[lastIdx].Total += 1
		} else {
			out[lastIdx].Schedules = append(out[lastIdx].Schedules, schedule)
			out[lastIdx].Total += 1
		}
	}

	return
}

type ScheduleHistoryListResultDay struct {
	Date time.Time `json:"date"`
	ScheduleHistoryListResult
}

func (resultDay ScheduleHistoryListResultDay) GetLast() *models.ScheduleHistory {
	if len(resultDay.Schedules) > 0 {
		return resultDay.Schedules[len(resultDay.Schedules)-1]
	}
	return nil
}

func (api *API) ScheduleHistoryList(ctx context.Context, params ScheduleHistoryListParams) (result ScheduleHistoryListResult, err error) {
	ctx, span := tracer.Start(ctx, "*API.ScheduleHistoryList")
	defer span.End()

	result.Schedules, err = models.ScheduleHistories.Query(ctx, api.db, params.Query()...).All()
	if err != nil {
		return result, errs.Wrapw(err, "failed to list schedule histories", "query", params)
	}

	result.Total, err = models.ScheduleHistories.Query(ctx, api.db, params.CountQuery()...).Count()
	if err != nil {
		return result, errs.Wrapw(err, "failed to count schedule histories", "query", params)
	}

	return result, nil
}
