package api

import (
	"context"

	"github.com/tigorlazuardi/redmage/db/queries"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type ListSubredditsParams struct {
	Name   string `json:"name"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

type ListSubredditsResult struct {
	Total int64
	Data  []queries.Subreddit
}

func (api *API) ListSubreddits(ctx context.Context, arg ListSubredditsParams) (result ListSubredditsResult, err error) {
	if arg.Name != "" {
		result.Data, err = api.Queries.SubredditsSearch(ctx, queries.SubredditsSearchParams{
			Name:   "%" + arg.Name + "%",
			Limit:  arg.Limit,
			Offset: arg.Offset,
		})
		if err != nil {
			return result, errs.Wrapw(err, "failed to search subreddit", "query", arg)
		}
		result.Total, err = api.Queries.SubredditsSearchCount(ctx, queries.SubredditsSearchCountParams{
			Name:   "%" + arg.Name + "%",
			Limit:  arg.Limit,
			Offset: arg.Offset,
		})
		if err != nil {
			return result, errs.Wrapw(err, "failed to count subreddit search", "query", arg)
		}
		return result, err
	}

	result.Data, err = api.Queries.SubredditsList(ctx, queries.SubredditsListParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return result, errs.Wrapw(err, "failed to list subreddit", "query", arg)
	}
	result.Total, err = api.Queries.SubredditsListCount(ctx)
	if err != nil {
		return result, errs.Wrapw(err, "failed to count subreddit list", "query", arg)
	}
	return result, err
}
