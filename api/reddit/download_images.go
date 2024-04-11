package reddit

import (
	"context"
	"io"
	"net/http"

	"github.com/alecthomas/units"
	"github.com/tigorlazuardi/redmage/api/bmessage"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"golang.org/x/sync/errgroup"
)

type DownloadStatusBroadcaster interface {
	Broadcast(bmessage.ImageDownloadMessage)
}

type NullDownloadStatusBroadcaster struct{}

func (NullDownloadStatusBroadcaster) Broadcast(bmessage.ImageDownloadMessage) {}

type PostImage struct {
	ImageURL      string
	ImageFile     io.Reader
	ThumbnailURL  string
	ThumbnailFile io.Reader
}

func (reddit *Reddit) DownloadImage(ctx context.Context, post Post, broadcaster DownloadStatusBroadcaster) (image PostImage, err error) {
	imageUrl, thumbnailUrl := post.GetImageURL(), post.GetThumbnailURL()
	image.ImageURL = imageUrl
	image.ThumbnailURL = thumbnailUrl

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		var err error
		image.ImageFile, err = reddit.downloadImage(groupCtx, post, bmessage.KindImage, broadcaster)
		return err
	})
	group.Go(func() error {
		var err error
		image.ThumbnailFile, err = reddit.downloadImage(groupCtx, post, bmessage.KindThumbnail, broadcaster)
		return err
	})

	err = group.Wait()
	return image, err
}

func (reddit *Reddit) downloadImage(ctx context.Context, post Post, kind bmessage.ImageKind, broadcaster DownloadStatusBroadcaster) (io.Reader, error) {
	var (
		url    string
		height int
		width  int
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
	idr := &ImageDownloadReader{
		OnProgress: func(downloaded int64, contentLength int64, err error) {
			broadcaster.Broadcast(bmessage.ImageDownloadMessage{
				Metadata: bmessage.ImageMetadata{
					URL:    url,
					Height: height,
					Width:  width,
					Kind:   kind,
				},
				ContentLength: units.MetricBytes(resp.ContentLength),
				Downloaded:    units.MetricBytes(downloaded),
				Subreddit:     post.GetSubreddit(),
				PostURL:       post.GetPermalink(),
				PostID:        post.GetID(),
				Error:         err,
			})
		},
		IdleTimeout:        reddit.Config.Duration("download.timeout.idle"),
		IdleSpeedThreshold: metricSpeed,
	}

	resp = idr.WrapHTTPResponse(resp)
	reader, writer := io.Pipe()
	go func() {
		defer resp.Body.Close()
		_, err := io.Copy(writer, resp.Body)
		_ = writer.CloseWithError(err)
	}()
	return reader, nil
}
