package api

import (
	"context"
	"errors"
	"image/jpeg"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/disintegration/imaging"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DownloadSubredditParams struct {
	Countback     int
	Devices       models.DeviceSlice
	SubredditType reddit.SubredditType
}

var (
	ErrNoDevices         = errors.New("api: no devices set")
	ErrDownloadDirNotSet = errors.New("api: download directory not set")
)

func (api *API) DownloadSubredditImages(ctx context.Context, subreddit *models.Subreddit, devices models.DeviceSlice) error {
	downloadDir := api.config.String("download.directory")
	if downloadDir == "" {
		return errs.Wrapw(ErrDownloadDirNotSet, "download directory must be set before images can be downloaded").Code(http.StatusBadRequest)
	}

	if len(devices) == 0 {
		return errs.Wrapw(ErrNoDevices, "downloading images requires at least one device configured").Code(http.StatusBadRequest)
	}

	ctx, span := tracer.Start(ctx, "*API.DownloadSubredditImages", trace.WithAttributes(attribute.String("subreddit", subreddit.Name)))
	defer span.End()

	wg := sync.WaitGroup{}

	countback := int(subreddit.Countback)

	var (
		list reddit.Listing
		err  error
	)
	for countback > 0 {
		limit := 100
		if limit > countback {
			limit = countback
		}
		log.New(ctx).Debug("getting posts", "subreddit", subreddit, "current_countback", countback, "current_limit", limit)
		list, err = api.reddit.GetPosts(ctx, reddit.GetPostsParam{
			Subreddit:     subreddit.Name,
			Limit:         limit,
			After:         list.GetLastAfter(),
			SubredditType: reddit.SubredditType(subreddit.Subtype),
		})
		if err != nil {
			return errs.Wrapw(err, "failed to get posts", "subreddit", subreddit)
		}
		wg.Add(1)
		go func(ctx context.Context, posts reddit.Listing) {
			defer wg.Done()
			err := api.downloadSubredditListImage(ctx, list, subreddit, devices)
			if err != nil {
				log.New(ctx).Err(err).Error("failed to download image")
			}
		}(ctx, list)
		if len(list.GetPosts()) == 0 {
			break
		}
		countback -= len(list.GetPosts())
	}

	wg.Wait()

	return nil
}

func (api *API) downloadSubredditListImage(ctx context.Context, list reddit.Listing, subreddit *models.Subreddit, devices models.DeviceSlice) error {
	ctx, span := tracer.Start(ctx, "*API.downloadSubredditListImage")
	defer span.End()

	wg := sync.WaitGroup{}

	for _, post := range list.GetPosts() {
		if !post.IsImagePost() {
			continue
		}
		devices := api.getDevicesThatAcceptPost(ctx, post, devices)
		if len(devices) == 0 {
			continue
		}
		log.New(ctx).Debug("downloading image", "post_id", post.GetID(), "post_url", post.GetImageURL(), "devices", devices)
		wg.Add(1)
		api.imageSemaphore <- struct{}{}
		go func(ctx context.Context, post reddit.Post) {
			defer func() {
				<-api.imageSemaphore
				wg.Done()
			}()

			if err := api.downloadSubredditImage(ctx, post, subreddit, devices); err != nil {
				log.New(ctx).Err(err).Error("failed to download subreddit image")
			}
		}(ctx, post)
	}

	wg.Wait()

	return nil
}

func (api *API) downloadSubredditImage(ctx context.Context, post reddit.Post, subreddit *models.Subreddit, devices models.DeviceSlice) error {
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

	thumbnailPath := post.GetThumbnailTargetPath(api.config)
	_, errStat := os.Stat(thumbnailPath)
	if errStat == nil {
		// file exist
		return nil
	}
	if !errors.Is(errStat, os.ErrNotExist) {
		return errs.Wrapw(err, "failed to check thumbnail existence", "path", thumbnailPath)
	}

	_ = os.MkdirAll(post.GetThumbnailTargetDir(api.config), 0o777)

	thumbnailSource, err := imaging.Open(tmpImageFile.filename)
	if err != nil {
		return errs.Wrapw(err, "failed to open temp thumbnail file",
			"filename", tmpImageFile.filename,
			"post_url", post.GetPermalink(),
			"image_url", post.GetImageURL(),
		)
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

	w, close, err := api.createDeviceImageWriters(post, devices)
	if err != nil {
		return errs.Wrapw(err, "failed to create image files")
	}
	log.New(ctx).Debug("saving image files", "post_id", post.GetID(), "post_url", post.GetImageURL(), "devices", devices)
	defer close()
	_, err = io.Copy(w, tmpImageFile)
	if err != nil {
		return errs.Wrapw(err, "failed to save image files")
	}

	var many []*models.ImageSetter
	for _, device := range devices {
		var nsfw int32
		if post.IsNSFW() {
			nsfw = 1
		}
		many = append(many, &models.ImageSetter{
			SubredditID:           omit.From(subreddit.ID),
			DeviceID:              omit.From(device.ID),
			Title:                 omit.From(post.GetTitle()),
			PostID:                omit.From(post.GetID()),
			PostURL:               omit.From(post.GetImageURL()),
			PostCreated:           omit.From(post.GetCreated().Format(time.RFC3339)),
			PostName:              omit.From(post.GetName()),
			Poster:                omit.From(post.GetAuthor()),
			PosterURL:             omit.From(post.GetAuthorURL()),
			ImageRelativePath:     omit.From(post.GetImageRelativePath(device)),
			ThumbnailRelativePath: omit.From(post.GetThumbnailRelativePath()),
			ImageOriginalURL:      omit.From(post.GetImageURL()),
			ThumbnailOriginalURL:  omit.From(post.GetThumbnailURL()),
			NSFW:                  omit.From(nsfw),
		})
	}

	log.New(ctx).Debug("inserting images to database", "images", many)
	_, err = models.Images.InsertMany(ctx, api.db, many...)
	if err != nil {
		return errs.Wrapw(err, "failed to insert images to database", "params", many)
	}

	return nil
}

func (api *API) createDeviceImageWriters(post reddit.Post, devices models.DeviceSlice) (writer io.Writer, close func(), err error) {
	// open file for each device
	var files []*os.File
	var writers []io.Writer
	for _, device := range devices {
		var filename string
		if device.WindowsWallpaperMode == 1 {
			filename = post.GetWindowsWallpaperImageTargetPath(api.config, device)
			dir := post.GetWindowsWallpaperImageTargetDir(api.config, device)
			_ = os.MkdirAll(dir, 0o777)
		} else {
			filename = post.GetImageTargetPath(api.config, device)
			dir := post.GetImageTargetDir(api.config, device)
			if err := os.MkdirAll(dir, 0o777); err != nil {
				for _, f := range files {
					_ = f.Close()
				}
				return nil, nil, errs.Wrapw(err, "failed to create target image dir")
			}
		}
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			for _, f := range files {
				_ = f.Close()
			}
			return nil, nil, errs.Wrapw(err, "failed to open target image file",
				"device_name", device.Name,
				"device_slug", device.Slug,
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

func (api *API) getDevicesThatAcceptPost(ctx context.Context, post reddit.Post, devices models.DeviceSlice) (devs models.DeviceSlice) {
	for _, device := range devices {
		if shouldDownloadPostForDevice(post, device) && !api.isImageExists(ctx, post, device) {
			devs = append(devs, device)
		}
	}
	return devs
}

func (api *API) isImageExists(ctx context.Context, post reddit.Post, device *models.Device) (found bool) {
	ctx, span := tracer.Start(ctx, "*API.IsImageExists")
	defer span.End()

	_, err := models.Images.Query(ctx, api.db,
		models.SelectWhere.Images.DeviceID.EQ(device.ID),
		models.SelectWhere.Images.PostID.EQ(post.GetID()),
	).One()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false
		}
	}

	// Image does not exist in target path.
	if _, err := os.Stat(post.GetImageTargetPath(api.config, device)); err != nil {
		return false
	}

	return true
}

func shouldDownloadPostForDevice(post reddit.Post, device *models.Device) bool {
	if post.IsNSFW() && device.NSFW == 0 {
		return false
	}
	devAspectRatio := deviceAspectRatio(device)
	rangeStart := devAspectRatio - device.AspectRatioTolerance
	rangeEnd := devAspectRatio + device.AspectRatioTolerance

	imgAspectRatio := post.GetImageAspectRatio()

	width, height := post.GetImageSize()
	log.New(context.Background()).Debug("checking image aspect ratio",
		"device", device.Slug,
		"device_height", device.ResolutionY,
		"device_width", device.ResolutionX,
		"device_aspect_ratio", devAspectRatio,
		"image_aspect_ratio", imgAspectRatio,
		"range_start", rangeStart,
		"range_end", rangeEnd,
		"success_fulfill_download_range_start", (imgAspectRatio > rangeStart),
		"success_fulfill_download_range_end", (imgAspectRatio < rangeEnd),
		"url", post.GetImageURL(),
		"image.width", width,
		"image.height", height,
	)

	if imgAspectRatio < rangeStart {
		return false
	}

	if imgAspectRatio > rangeEnd {
		return false
	}

	if device.MaxX > 0 && width > int64(device.MaxX) {
		return false
	}
	if device.MaxY > 0 && height > int64(device.MaxY) {
		return false
	}
	if device.MinX > 0 && width < int64(device.MinX) {
		return false
	}
	if device.MinY > 0 && height < int64(device.MinY) {
		return false
	}
	return true
}

func deviceAspectRatio(device *models.Device) float64 {
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
	err = os.MkdirAll(tmpDirname, 0777)
	if err != nil {
		return nil, errs.Wrapw(err, "failed to create temporary dir", "dir_name", tmpDirname)
	}
	tmpFilename := path.Join(tmpDirname, imageFilename)

	file, err := os.OpenFile(tmpFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o777)
	if err != nil {
		return nil, errs.Wrapw(err, "failed to open temp image file",
			"temp_file_path", tmpFilename,
			"image_url", img.URL,
		)
	}

	// File must be closed by end of function because kernel stuffs.
	//
	// A fresh fd must be used to properly get the new data.
	defer file.Close()

	_, err = io.Copy(file, img.File)
	if err != nil {
		return nil, errs.Wrapw(err, "failed to download image to temp file",
			"temp_file_path", tmpFilename,
			"image_url", img.URL,
		)
	}

	filew, err := os.OpenFile(tmpFilename, os.O_RDONLY, 0o777)
	if err != nil {
		return nil, errs.Wrapw(err, "failed to download image to temp file",
			"temp_file_path", tmpFilename,
			"image_url", img.URL,
		)
	}

	return &tempFile{
		file:     filew,
		filename: tmpFilename,
	}, err
}
