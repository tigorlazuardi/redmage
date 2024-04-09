package api

import "context"

type DownloadSubredditParams struct {
	Countback int
	NSFW      bool
}

func (api *API) DownloadSubredditImages(ctx context.Context, subredditName string, params DownloadSubredditParams) error {
	return nil
}
