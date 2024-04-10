package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/tigorlazuardi/redmage/db/queries"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type DownloadSubredditParams struct {
	Countback int
	NSFW      bool
	Devices   []queries.Device
	Type      int
}

var (
	ErrNoDevices         = errors.New("api: no devices set")
	ErrDownloadDirNotSet = errors.New("api: download directory not set")
)

func (api *API) DownloadSubredditImages(ctx context.Context, subredditName string, params DownloadSubredditParams) error {
	downloadDir := api.config.String("download.directory")
	if downloadDir == "" {
		return errs.Wrapw(ErrDownloadDirNotSet, "download directory must be set before images can be downloaded").Code(http.StatusBadRequest)
	}

	if len(params.Devices) == 0 {
		return errs.Wrapw(ErrNoDevices, "downloading images requires at least one device configured").Code(http.StatusBadRequest)
	}

	return nil
}
