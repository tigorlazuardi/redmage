package events

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/views/components/progress"
)

type ImageDownloadEvent string

const (
	ImageDownloadStart    ImageDownloadEvent = "image.download.start"
	ImageDownloadEnd      ImageDownloadEvent = "image.download.end"
	ImageDownloadError    ImageDownloadEvent = "image.download.error"
	ImageDownloadProgress ImageDownloadEvent = "image.download.progress"
)

type ImageDownload struct {
	EventKind             ImageDownloadEvent `json:"event,omitempty"`
	ImageURL              string             `json:"image_url,omitempty"`
	ImageHeight           int32              `json:"image_height,omitempty"`
	ImageWidth            int32              `json:"image_width,omitempty"`
	ContentLength         int64              `json:"content_length,omitempty"`
	Downloaded            int64              `json:"downloaded,omitempty"`
	Subreddit             string             `json:"subreddit,omitempty"`
	PostURL               string             `json:"post_url,omitempty"`
	PostName              string             `json:"post_name,omitempty"`
	PostTitle             string             `json:"post_title,omitempty"`
	PostCreated           int64              `json:"post_created,omitempty"`
	PostAuthor            string             `json:"post_author,omitempty"`
	PostAuthorURL         string             `json:"post_author_url,omitempty"`
	ImageRelativePath     string             `json:"image_relative_path,omitempty"`
	ImageOriginalURL      string             `json:"image_original_url,omitempty"`
	ImageSize             int64              `json:"image_size,omitempty"`
	ThumbnailRelativePath string             `json:"thumbnail_relative_path,omitempty"`
	NSFW                  int32              `json:"nsfw,omitempty"`
	Error                 error              `json:"error,omitempty"`
	Device                string             `json:"device,omitempty"`
}

// Render the template.
func (im ImageDownload) Render(ctx context.Context, w io.Writer) error {
	switch im.EventKind {
	case ImageDownloadStart:
		return progress.ImageDownloadStartNotification(progress.ImageDownloadStartNotificationData{
			ID:                 fmt.Sprintf("notif-image-download-%s-%s", im.Subreddit, im.PostName),
			Subreddit:          im.Subreddit,
			PostName:           im.PostName,
			PostTitle:          im.PostTitle,
			PostURL:            im.PostURL,
			AutoRemoveDuration: time.Second * 5,
		}).Render(ctx, w)
	case ImageDownloadEnd:
		return progress.ImageDownloadEndNotification(progress.ImageDownloadEndNotificationData{
			ID:                 fmt.Sprintf("notif-image-download-%s-%s", im.Subreddit, im.PostName),
			Subreddit:          im.Subreddit,
			PostURL:            im.PostName,
			PostName:           im.PostTitle,
			PostTitle:          im.PostURL,
			AutoRemoveDuration: time.Second * 5,
		}).Render(ctx, w)
	case ImageDownloadError:
		return progress.ImageDownloadErrorNotification(progress.ImageDownloadErrorNotificationData{
			ID:                 fmt.Sprintf("notif-image-download-%s-%s", im.Subreddit, im.PostName),
			Subreddit:          im.Subreddit,
			PostURL:            im.PostName,
			PostName:           im.PostTitle,
			PostTitle:          im.PostURL,
			Error:              im.Error,
			AutoRemoveDuration: time.Second * 5,
		}).Render(ctx, w)
	case ImageDownloadProgress:
		return progress.ImageDownloadProgressNotification(progress.ImageDownloadProgressNotificationData{
			ID:                 fmt.Sprintf("notif-image-download-%s-%s", im.Subreddit, im.PostName),
			Subreddit:          im.Subreddit,
			PostURL:            im.PostName,
			PostName:           im.PostTitle,
			PostTitle:          im.PostURL,
			ContentLength:      im.ContentLength,
			Downloaded:         im.Downloaded,
			AutoRemoveDuration: time.Second * 5,
		}).Render(ctx, w)
	default:
		return errs.Fail("events.ImageDownload: unknown event kind", "event", im)
	}
}

// Event returns the event name
func (im ImageDownload) Event() string {
	return string(im.EventKind)
}

// SerializeTo writes the event data to the writer.
//
// SerializeTo must not write multiple linebreaks (single linebreak is fine)
// in succession to the writer since it will mess up SSE events.
func (im ImageDownload) SerializeTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(im)
}

func PublishImageDownloadEvent(bc *broadcast.Relay[Event], event ImageDownload) {
	bc.Broadcast(event)
	bc.Broadcast(ImageDownloadSubreddit{event})
	if event.EventKind == ImageDownloadEnd {
		bc.Broadcast(ImageDownloadSubredditCard{event})
	}
}
