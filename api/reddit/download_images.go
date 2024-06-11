package reddit

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/alecthomas/units"
	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/redmage/api/bmessage"
	"github.com/tigorlazuardi/redmage/api/events"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type DownloadStatusBroadcaster interface {
	Broadcast(bmessage.ImageDownloadMessage)
}

type NullDownloadStatusBroadcaster struct{}

func (NullDownloadStatusBroadcaster) Broadcast(bmessage.ImageDownloadMessage) {}

type PostImage struct {
	URL  string
	File io.ReadCloser
}

func (po *PostImage) Read(p []byte) (n int, err error) {
	return po.File.Read(p)
}

func (po *PostImage) Close() error {
	return po.File.Close()
}

// DownloadImage downloades the image.
//
// If downloading image or thumbnail fails
func (reddit *Reddit) DownloadImage(ctx context.Context, post Post, broadcaster *broadcast.Relay[events.Event]) (image PostImage, err error) {
	ctx, span := tracer.Start(ctx, "*Reddit.DownloadImage")
	defer span.End()
	imageUrl := post.GetImageURL()
	image.URL = imageUrl

	image.File, err = reddit.downloadImage(ctx, post, bmessage.KindImage, broadcaster)
	return image, err
}

func (reddit *Reddit) DownloadThumbnail(ctx context.Context, post Post, broadcaster *broadcast.Relay[events.Event]) (image PostImage, err error) {
	ctx, span := tracer.Start(ctx, "*Reddit.DownloadThumbnail")
	defer span.End()
	imageUrl := post.GetThumbnailURL()
	image.URL = imageUrl

	image.File, err = reddit.downloadImage(ctx, post, bmessage.KindThumbnail, broadcaster)
	return image, err
}

func (reddit *Reddit) downloadImage(ctx context.Context, post Post, kind bmessage.ImageKind, broadcaster *broadcast.Relay[events.Event]) (io.ReadCloser, error) {
	var (
		url    string
		height int64
		width  int64
	)
	if kind == bmessage.KindImage {
		url = post.GetImageURL()
		width, height = post.GetImageSize()
	} else {
		url = post.GetThumbnailURL()
		width, height = post.GetThumbnailSize()
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errs.Wrapw(err, "reddit: failed to create request", "url", url)
	}
	resp, err := reddit.Client.Do(req)
	if err != nil {
		return nil, errs.Wrapw(err, "reddit: failed to execute request", "url", url)
	}
	if resp.StatusCode >= 400 {
		return nil, errs.Fail("unexpected status code when trying to download images",
			"url", url,
			"status_code", resp.StatusCode,
		)
	}
	idleSpeedStr := reddit.Config.String("download.timeout.idlespeed")
	metricSpeed, _ := units.ParseMetricBytes(idleSpeedStr)
	if metricSpeed == 0 {
		metricSpeed = 10 * units.KB
	}
	eventData := events.ImageDownload{
		ImageURL:         post.GetImageURL(),
		ImageHeight:      int32(height),
		ImageWidth:       int32(width),
		Subreddit:        post.GetSubreddit(),
		PostURL:          post.GetPostURL(),
		PostName:         post.GetPostName(),
		PostTitle:        post.GetPostTitle(),
		PostCreated:      post.GetPostCreated(),
		PostAuthor:       post.GetAuthor(),
		PostAuthorURL:    post.GetAuthorURL(),
		ImageOriginalURL: post.GetImageURL(),
		NSFW:             int32(post.IsNSFWInt()),
	}
	idr := &ImageDownloadReader{
		OnProgress: func(downloaded int64, contentLength int64, err error) {
			var kind events.ImageDownloadEvent
			if errors.Is(err, io.EOF) {
				err = nil
			}
			if err != nil {
				kind = events.ImageDownloadError
			} else {
				kind = events.ImageDownloadProgress
			}
			ev := eventData.Clone()
			ev.EventKind = kind
			ev.Downloaded = downloaded
			ev.ImageSize = contentLength
			ev.ContentLength = contentLength
			events.PublishImageDownloadEvent(broadcaster, ev)
		},
		OnClose: func(downloaded, contentLength int64, closeErr error) {
			ev := eventData.Clone()
			ev.EventKind = events.ImageDownloadEnd
			ev.Downloaded = downloaded
			ev.ImageSize = contentLength
			ev.ContentLength = contentLength
			ev.Error = closeErr
			events.PublishImageDownloadEvent(broadcaster, ev)
		},
		IdleTimeout:        reddit.Config.Duration("download.timeout.idle"),
		IdleSpeedThreshold: metricSpeed,
	}

	resp = idr.WrapHTTPResponse(resp)
	reader, writer := io.Pipe()
	go func() {
		defer resp.Body.Close()
		ev := eventData.Clone()
		ev.EventKind = events.ImageDownloadStart
		ev.ImageSize = resp.ContentLength
		ev.ContentLength = resp.ContentLength
		events.PublishImageDownloadEvent(broadcaster, ev)
		_, err := io.Copy(writer, resp.Body)
		_ = writer.CloseWithError(err)
	}()
	return reader, nil
}
