package api

import (
	"context"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type ListSubredditsParams struct {
	Name    string `json:"name"`
	Limit   int64  `json:"limit"`
	Offset  int64  `json:"offset"`
	OrderBy string `json:"order_by"`
	Sort    string `json:"sort"`
}

func (l ListSubredditsParams) Query() (expr []bob.Mod[*dialect.SelectQuery]) {
	if l.Name != "" {
		expr = append(expr, models.SelectWhere.Subreddits.Name.Like("%"+l.Name+"%"))
	}
	if l.Limit > 0 {
		expr = append(expr, sm.Limit(l.Limit))
	}
	if l.Offset > 0 {
		expr = append(expr, sm.Offset(l.Offset))
	}
	if l.OrderBy != "" {
		if l.Sort == "desc" {
			expr = append(expr, sm.OrderBy(l.OrderBy).Desc())
		} else {
			expr = append(expr, sm.OrderBy(l.OrderBy).Asc())
		}
	}

	return expr
}

func (l ListSubredditsParams) CountQuery() (expr []bob.Mod[*dialect.SelectQuery]) {
	if l.Name != "" {
		expr = append(expr, models.SelectWhere.Subreddits.Name.Like("%"+l.Name+"%"))
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
