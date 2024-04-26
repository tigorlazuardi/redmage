package api

import (
	"context"

	"github.com/tigorlazuardi/redmage/api/reddit"
)

type SubredditCheckParam = reddit.CheckSubredditParams

func (api *API) SubredditCheck(ctx context.Context, params SubredditCheckParam) (actual string, err error) {
	ctx, span := tracer.Start(ctx, "*API.SubredditCheck")
	defer span.End()

	return api.reddit.CheckSubreddit(ctx, params)
}
