package api

import (
	"context"
	"errors"

	"github.com/tigorlazuardi/redmage/db/queries"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type DownloadSubredditParams struct {
	Countback int
	NSFW      bool
	Devices   []queries.Device
}

var (
	ErrNoDevices         = errors.New("api: downloading subreddit images requires at least one device")
	ErrDownloadDirNotSet = errors.New("api: downloading subreddit images require download directory to be set")
)

func (api *API) DownloadSubredditImages(ctx context.Context, subredditName string, params DownloadSubredditParams) error {
	if len(params.Devices) == 0 {
		return errs.Wrap(ErrNoDevices)
	}
	return nil
}
