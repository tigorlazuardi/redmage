package api

import (
	"context"
	"strconv"
	"strings"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type ListSubredditsParams struct {
	Q       string
	All     bool
	Limit   int64
	Offset  int64
	OrderBy string
	Sort    string
}

func (l *ListSubredditsParams) FillFromQuery(q Queryable) {
	l.Q = q.Get("q")
	l.All, _ = strconv.ParseBool(q.Get("all"))
	l.Limit, _ = strconv.ParseInt(q.Get("limit"), 10, 64)
	if l.Limit < 0 {
		l.Limit = 25
	} else if l.Limit > 100 {
		l.Limit = 100
	}
	l.Offset, _ = strconv.ParseInt(q.Get("offset"), 10, 64)
	if l.Offset < 0 {
		l.Offset = 0
	}
	l.OrderBy = q.Get("order_by")
	l.Sort = strings.ToLower(q.Get("sort"))
}

func (l ListSubredditsParams) Query() (expr []bob.Mod[*dialect.SelectQuery]) {
	if l.Q != "" {
		expr = append(expr, models.SelectWhere.Subreddits.Name.Like("%"+l.Q+"%"))
	}
	if l.Limit > 0 {
		expr = append(expr, sm.Limit(l.Limit))
	}
	if l.Offset > 0 {
		expr = append(expr, sm.Offset(l.Offset))
	}
	if l.OrderBy != "" {
		if l.Sort == "desc" {
			expr = append(expr, sm.OrderBy(sqlite.Quote(l.OrderBy)).Desc())
		} else {
			expr = append(expr, sm.OrderBy(sqlite.Quote(l.OrderBy)).Asc())
		}
	} else {
		expr = append(expr, sm.OrderBy(models.SubredditColumns.Name).Asc())
	}

	return expr
}

func (l ListSubredditsParams) CountQuery() (expr []bob.Mod[*dialect.SelectQuery]) {
	if l.Q != "" {
		expr = append(expr, models.SelectWhere.Subreddits.Name.Like("%"+l.Q+"%"))
	}

	return expr
}

type ListSubredditsResult struct {
	Total int64
	Data  models.SubredditSlice
}

func (api *API) ListSubreddits(ctx context.Context, arg ListSubredditsParams) (result ListSubredditsResult, err error) {
	ctx, span := tracer.Start(ctx, "api.ListSubreddits")
	defer span.End()

	result.Data, err = models.Subreddits.Query(ctx, api.db, arg.Query()...).All()
	if err != nil {
		return result, errs.Wrapw(err, "failed to list subreddits", "query", arg)
	}

	result.Total, err = models.Subreddits.Query(ctx, api.db, arg.CountQuery()...).Count()
	if err != nil {
		return result, errs.Wrapw(err, "failed to count subreddits", "query", arg)
	}

	return result, nil
}

// ListSubredditsWithCover returns list of subreddits with cover image.
//
// Image Relationship `R` struct will be nil if there is no cover image.
func (api *API) ListSubredditsWithCover(ctx context.Context, arg ListSubredditsParams) (result ListSubredditsResult, err error) {
	ctx, span := tracer.Start(ctx, "api.ListSubredditsWithCover")
	defer span.End()

	result, err = api.ListSubreddits(ctx, arg)
	if err != nil {
		return result, errs.Wrapw(err, "failed to list subreddits with cover")
	}

	// Cannot do batch query because we cannot use GROUP BY with ORDER BY in consistent manner.
	//
	// The problem gets worse when using custom ORDER BY from client.
	//
	// For consistency, we query images one by one.
	//
	// Subreddit list is expected to be small, so this should be fine since SQLITE has no network latency.
	for _, subreddit := range result.Data {
		subreddit.R.Images, err = models.Images.Query(ctx, api.db,
			models.SelectWhere.Images.Subreddit.EQ(subreddit.Name),
			sm.Limit(1),
			sm.OrderBy(models.ImageColumns.CreatedAt).Desc(),
		).All()
	}

	return result, nil
}
