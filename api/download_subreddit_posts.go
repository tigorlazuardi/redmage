package api

import (
	"context"

	"github.com/tigorlazuardi/redmage/api/reddit"
)

type DownloadSubredditPostsParams struct {
	Page  int
	Limit int
}

func (api *API) DownloadSubredditPosts(ctx context.Context, subredditName string, params DownloadSubredditPostsParams) (posts reddit.Listing, err error) {
	return posts, err
}
