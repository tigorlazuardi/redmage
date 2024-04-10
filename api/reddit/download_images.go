package reddit

import (
	"context"
	"io"

	"github.com/tigorlazuardi/redmage/api/bmessage"
)

type DownloadStatusBroadcaster interface {
	Broadcast(bmessage.DownloadStatusMessage)
}

type NullDownloadStatusBroadcaster struct{}

func (NullDownloadStatusBroadcaster) Broadcast(bmessage.DownloadStatusMessage) {}

type PostImage struct {
	ImageURL      string
	ImageFile     io.Reader
	ThumbnailURL  string
	ThumbnailFile io.Reader
}

func (reddit *Reddit) DownloadImage(ctx context.Context, post Post, broadcaster DownloadStatusBroadcaster) (images PostImage, err error) {
	return images, err
}
