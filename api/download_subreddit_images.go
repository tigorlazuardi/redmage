package api

import (
	"context"
	"errors"
	"image/jpeg"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/db/queries"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DownloadSubredditParams struct {
	Countback     int
	Devices       []queries.Device
	SubredditType reddit.SubredditType
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

	ctx, span := tracer.Start(ctx, "*API.DownloadSubredditImages", trace.WithAttributes(attribute.String("subreddit", subredditName)))
	defer span.End()

	wg := sync.WaitGroup{}

	countback := params.Countback

	for page := 1; countback > 0; page += 1 {
		limit := countback
		if limit > 100 {
			limit = 100
		}
		list, err := api.reddit.GetPosts(ctx, reddit.GetPostsParam{
			Subreddit:     subredditName,
			Limit:         limit,
			Page:          page,
			SubredditType: params.SubredditType,
		})
		if err != nil {
			return errs.Wrapw(err, "failed to get posts", "subreddit_name", subredditName, "params", params)
		}
		wg.Add(1)
		go func(ctx context.Context, posts reddit.Listing) {
			defer wg.Done()
			err := api.downloadSubredditListImage(ctx, list, params)
			if err != nil {
				log.New(ctx).Err(err).Error("failed to download image")
			}
		}(ctx, list)
		countback -= len(list.GetPosts())
	}

	wg.Wait()

	return nil
}

func (api *API) downloadSubredditListImage(ctx context.Context, list reddit.Listing, params DownloadSubredditParams) error {
	ctx, span := tracer.Start(ctx, "*API.downloadSubredditListImage")
	defer span.End()

	wg := sync.WaitGroup{}

	for _, post := range list.GetPosts() {
		if !post.IsImagePost() {
			continue
		}
		devices := getDevicesThatAcceptPost(post, params.Devices)
		if len(devices) == 0 {
			continue
		}
		wg.Add(1)
		api.imageSemaphore <- struct{}{}
		go func(ctx context.Context, post reddit.Post) {
			defer func() {
				<-api.imageSemaphore
				wg.Done()
			}()

			if err := api.downloadSubredditImage(ctx, post, devices); err != nil {
				log.New(ctx).Err(err).Error("failed to download subreddit image")
			}
		}(ctx, post)
	}

	wg.Wait()

	return nil
}

func (api *API) downloadSubredditImage(ctx context.Context, post reddit.Post, devices []queries.Device) error {
	ctx, span := tracer.Start(ctx, "*API.downloadSubredditImage")
	defer span.End()

	imageHandler, err := api.reddit.DownloadImage(ctx, post, api.downloadBroadcast)
	if err != nil {
		return errs.Wrapw(err, "failed to download image")
	}
	defer imageHandler.Close()

	// copy to temp dir first to avoid copying incomplete files.
	tmpImageFile, err := api.copyImageToTempDir(ctx, imageHandler)
	if err != nil {
		return errs.Wrapw(err, "failed to download image to temp file")
	}
	defer tmpImageFile.Close()

	w, close, err := api.createDeviceImageWriters(post, devices)
	if err != nil {
		return errs.Wrapw(err, "failed to create image files")
	}
	defer close()
	_, err = io.Copy(w, tmpImageFile)
	if err != nil {
		return errs.Wrapw(err, "failed to save image files")
	}
	thumbnailPath := post.GetThumbnailTargetPath(api.config)
	_, errStat := os.Stat(thumbnailPath)
	if errStat == nil {
		// file exist
		return nil
	}
	if !errors.Is(errStat, os.ErrNotExist) {
		return errs.Wrapw(err, "failed to check thumbail existence", "path", thumbnailPath)
	}

	thumbnailSource, err := imaging.Open(tmpImageFile.filename)
	if err != nil {
		return errs.Wrapw(err, "failed to open temp thumbnail file", "filename", tmpImageFile.filename)
	}

	thumbnail := imaging.Resize(thumbnailSource, 256, 0, imaging.Lanczos)
	thumbnailFile, err := os.Create(thumbnailPath)
	if err != nil {
		return errs.Wrapw(err, "failed to create thumbnail file", "filename", thumbnailPath)
	}
	defer thumbnailFile.Close()

	err = jpeg.Encode(thumbnailFile, thumbnail, nil)
	if err != nil {
		return errs.Wrapw(err, "failed to encode thumbnail file to jpeg", "filename", thumbnailPath)
	}

	return nil
}

func (api *API) createDeviceImageWriters(post reddit.Post, devices []queries.Device) (writer io.Writer, close func(), err error) {
	// open file for each device
	var files []*os.File
	var writers []io.Writer
	for _, device := range devices {
		var filename string
		if device.WindowsWallpaperMode == 1 {
			filename = post.GetWindowsWallpaperImageTargetPath(api.config, device)
		} else {
			filename = post.GetImageTargetPath(api.config, device)
		}
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			for _, f := range files {
				_ = f.Close()
			}
			return nil, nil, errs.Wrapw(err, "failed to open temp image file",
				"device_name", device.Name,
				"filename", filename,
			)
		}
		files = append(files, file)
		writers = append(writers, file)
	}

	return io.MultiWriter(writers...), func() {
		for _, file := range files {
			_ = file.Close()
		}
	}, nil
}

func getDevicesThatAcceptPost(post reddit.Post, devices []queries.Device) []queries.Device {
	var devs []queries.Device
	for _, device := range devices {
		if shouldDownloadPostForDevice(post, device) {
			devs = append(devices, device)
		}
	}
	return devs
}

func shouldDownloadPostForDevice(post reddit.Post, device queries.Device) bool {
	if post.IsNSFW() && device.Nsfw == 0 {
		return false
	}
	if math.Abs(deviceAspectRatio(device)-post.GetImageAspectRatio()) > device.AspectRatioTolerance { // outside of aspect ratio tolerance
		return false
	}
	width, height := post.GetImageSize()
	if device.MaxX > 0 && width > device.MaxX {
		return false
	}
	if device.MaxY > 0 && height > device.MaxY {
		return false
	}
	if device.MinX > 0 && width < device.MinX {
		return false
	}
	if device.MinY > 0 && height < device.MinY {
		return false
	}
	return true
}

func deviceAspectRatio(device queries.Device) float64 {
	return float64(device.ResolutionX) / float64(device.ResolutionY)
}

type tempFile struct {
	filename string
	file     *os.File
}

func (te *tempFile) Read(p []byte) (n int, err error) {
	return te.file.Read(p)
}

func (te *tempFile) Close() error {
	return te.file.Close()
}

// copyImageToTempDir copies the image to a temporary directory and returns the file handle
//
// file must be closed by the caller after use.
//
// file is nil if an error occurred.
func (api *API) copyImageToTempDir(ctx context.Context, img reddit.PostImage) (tmp *tempFile, err error) {
	_, span := tracer.Start(ctx, "*API.copyImageToTempDir")
	defer func() { telemetry.EndWithStatus(span, err) }()

	// ignore error because url is always valid if this
	// function is called
	url, _ := url.Parse(img.URL)

	split := strings.Split(url.Path, "/")
	imageFilename := split[len(split)-1]
	tmpDirname := path.Join(os.TempDir(), "redmage")
	_ = os.MkdirAll(tmpDirname, 0644)
	tmpFilename := path.Join(tmpDirname, imageFilename)

	file, err := os.OpenFile(tmpFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, errs.Wrapw(err, "failed to open temp image file",
			"temp_file_path", tmpFilename,
			"image_url", img.URL,
		)
	}

	_, err = io.Copy(file, img.File)
	if err != nil {
		_ = file.Close()
		return nil, errs.Wrapw(err, "failed to download image to temp file",
			"temp_file_path", tmpFilename,
			"image_url", img.URL,
		)
	}

	return &tempFile{
		file:     file,
		filename: tmpFilename,
	}, err
}
