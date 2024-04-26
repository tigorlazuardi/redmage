package reddit

import (
	"context"
	"io"
	"net/http"

	"github.com/alecthomas/units"
	"github.com/tigorlazuardi/redmage/api/bmessage"
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
func (reddit *Reddit) DownloadImage(ctx context.Context, post Post, broadcaster DownloadStatusBroadcaster) (image PostImage, err error) {
	ctx, span := tracer.Start(ctx, "*Reddit.DownloadImage")
	defer span.End()
	imageUrl := post.GetImageURL()
	image.URL = imageUrl

	image.File, err = reddit.downloadImage(ctx, post, bmessage.KindImage, broadcaster)
	return image, err
}

func (reddit *Reddit) DownloadThumbnail(ctx context.Context, post Post, broadcaster DownloadStatusBroadcaster) (image PostImage, err error) {
	ctx, span := tracer.Start(ctx, "*Reddit.DownloadThumbnail")
	defer span.End()
	imageUrl := post.GetThumbnailURL()
	image.URL = imageUrl

	image.File, err = reddit.downloadImage(ctx, post, bmessage.KindThumbnail, broadcaster)
	return image, err
}

func (reddit *Reddit) downloadImage(ctx context.Context, post Post, kind bmessage.ImageKind, broadcaster DownloadStatusBroadcaster) (io.ReadCloser, error) {
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
	ctx, cancel := context.WithTimeout(ctx, reddit.Config.Duration("download.timeout.headers"))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errs.Wrapw(err, "reddit: failed to create request", "url", url)
	}
	resp, err := reddit.Client.Do(req)
	if err != nil {
		return nil, errs.Wrapw(err, "reddit: failed to execute request", "url", url)
	}
	idleSpeedStr := reddit.Config.String("download.timeout.idlespeed")
	metricSpeed, _ := units.ParseMetricBytes(idleSpeedStr)
	if metricSpeed == 0 {
		metricSpeed = 10 * units.KB
	}
	metadata := bmessage.ImageMetadata{
		URL:    url,
		Height: height,
		Width:  width,
		Kind:   kind,
	}
	idr := &ImageDownloadReader{
		OnProgress: func(downloaded int64, contentLength int64, err error) {
			var event bmessage.DownloadEvent
			if err != nil {
				event = bmessage.DownloadError
			} else {
				event = bmessage.DownloadProgress
			}
			broadcaster.Broadcast(bmessage.ImageDownloadMessage{
				Event:         event,
				Metadata:      metadata,
				ContentLength: units.MetricBytes(resp.ContentLength),
				Downloaded:    units.MetricBytes(downloaded),
				Subreddit:     post.GetSubreddit(),
				PostURL:       post.GetPermalink(),
				PostID:        post.GetID(),
				Error:         err,
			})
		},
		OnClose: func(downloaded, contentLength int64, closeErr error) {
			broadcaster.Broadcast(bmessage.ImageDownloadMessage{
				Event:         bmessage.DownloadEnd,
				Metadata:      metadata,
				ContentLength: units.MetricBytes(resp.ContentLength),
				Downloaded:    units.MetricBytes(downloaded),
				Subreddit:     post.GetSubreddit(),
				PostURL:       post.GetPermalink(),
				PostID:        post.GetID(),
				Error:         closeErr,
			})
		},
		IdleTimeout:        reddit.Config.Duration("download.timeout.idle"),
		IdleSpeedThreshold: metricSpeed,
	}

	resp = idr.WrapHTTPResponse(resp)
	reader, writer := io.Pipe()
	go func() {
		defer resp.Body.Close()
		broadcaster.Broadcast(bmessage.ImageDownloadMessage{
			Event: bmessage.DownloadStart,
			Metadata: bmessage.ImageMetadata{
				URL:    url,
				Height: height,
				Width:  width,
				Kind:   kind,
			},
			Subreddit: post.GetSubreddit(),
			PostURL:   post.GetPermalink(),
			PostID:    post.GetID(),
		})
		_, err := io.Copy(writer, resp.Body)
		_ = writer.CloseWithError(err)
	}()
	return reader, nil
}
