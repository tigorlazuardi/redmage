package events

import (
	"context"
	"io"
	"time"

	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/views/components"
)

type ImageDownloadSubreddit struct {
	ImageDownload
}

func (im ImageDownloadSubreddit) Event() string {
	return string(im.EventKind) + "." + im.Subreddit
}

type ImageDownloadSubredditCard struct {
	ImageDownload
}

func (im ImageDownloadSubredditCard) Event() string {
	return string(im.EventKind) + "." + im.Subreddit + ".card"
}

func (im ImageDownloadSubredditCard) Render(ctx context.Context, w io.Writer) error {
	if im.EventKind == ImageDownloadEnd {
		now := time.Now().Unix()
		data := &models.Image{
			Subreddit:             im.Subreddit,
			Device:                im.Device,
			PostTitle:             im.PostTitle,
			PostName:              im.PostName,
			PostURL:               im.PostURL,
			PostCreated:           im.PostCreated,
			PostAuthor:            im.PostAuthor,
			PostAuthorURL:         im.PostAuthor,
			ImageRelativePath:     im.ImageRelativePath,
			ImageOriginalURL:      im.ImageOriginalURL,
			ImageHeight:           im.ImageHeight,
			ImageWidth:            im.ImageWidth,
			ImageSize:             im.ImageSize,
			ThumbnailRelativePath: im.ThumbnailRelativePath,
			NSFW:                  im.NSFW,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
		return components.ImageCard(data, 0).Render(ctx, w)
	} else {
		return im.ImageDownload.Render(ctx, w)
	}
}
