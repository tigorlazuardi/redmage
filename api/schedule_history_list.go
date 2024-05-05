package api

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type ScheduleHistoryListParams struct {
	Subreddit string
	After     time.Time
	Before    time.Time

	Limit   int64
	Offset  int64
	OrderBy string
	Sort    string
}

func (params *ScheduleHistoryListParams) FillFromQuery(query Queryable) {
	params.Subreddit = query.Get("subreddit")
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

	afterInt, _ := strconv.ParseInt(query.Get("after"), 10, 64)
	if afterInt > 0 {
		params.After = time.Unix(afterInt, 0)
	} else if afterInt < 0 {
		params.After = now.Add(time.Duration(afterInt) * time.Second)
	}

	beforeInt, _ := strconv.ParseInt(query.Get("before"), 10, 64)
	if beforeInt > 0 {
		params.Before = time.Unix(beforeInt, 0)
	} else if beforeInt < 0 {
		params.Before = now.Add(time.Duration(beforeInt) * time.Second)
	}

	params.OrderBy = query.Get("order_by")
	params.Sort = query.Get("sort")
}

func (params ScheduleHistoryListParams) CountQuery() (expr []bob.Mod[*dialect.SelectQuery]) {
	if params.Subreddit != "" {
		expr = append(expr, models.SelectWhere.ScheduleHistories.Subreddit.EQ(params.Subreddit))
	}
	if !params.After.IsZero() {
		expr = append(expr, models.SelectWhere.ScheduleHistories.CreatedAt.GTE(params.After.Unix()))
	}
	if !params.Before.IsZero() {
		expr = append(expr, models.SelectWhere.ScheduleHistories.CreatedAt.LTE(params.Before.Unix()))
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
	if params.OrderBy != "" {
		if strings.ToLower(params.Sort) == "desc" {
			expr = append(expr, sm.OrderBy(sqlite.Quote(params.OrderBy)).Desc())
		} else {
			expr = append(expr, sm.OrderBy(sqlite.Quote(params.OrderBy)).Asc())
		}
	} else {
		expr = append(expr, sm.OrderBy(models.ScheduleHistoryColumns.CreatedAt).Desc(), sm.OrderBy(models.ScheduleHistoryColumns.Status).Desc())
	}

	return expr
}

type ScheduleHistoryListResult struct {
	Schedules models.ScheduleHistorySlice `json:"schedules"`
	Total     int64                       `json:"count"`
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
