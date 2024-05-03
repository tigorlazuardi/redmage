package api

import (
	"context"

	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type SubredditCheckParam = reddit.CheckSubredditParams

func (api *API) SubredditCheck(ctx context.Context, params SubredditCheckParam) (actual string, err error) {
	ctx, span := tracer.Start(ctx, "*API.SubredditCheck")
	defer span.End()

	return api.reddit.CheckSubreddit(ctx, params)
}

func (api *API) SubredditRegistered(ctx context.Context, params SubredditCheckParam) (registered bool, err error) {
	ctx, span := tracer.Start(ctx, "*API.SubredditRegistered")
	defer span.End()

	registered, err = models.Subreddits.Query(ctx, api.db,
		models.SelectWhere.Subreddits.Name.EQ(params.Subreddit),
		models.SelectWhere.Subreddits.Subtype.EQ(int32(params.SubredditType)),
	).Exists()
	if err != nil {
		return registered, errs.Wrapw(err, "failed to check subreddit registered state", "params", params)
	}

	return registered, err
}
