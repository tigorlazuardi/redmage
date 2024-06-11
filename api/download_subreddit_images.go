package api

import (
	"context"
	"errors"
	"image/jpeg"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
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
		return errs.Wrapw(ErrNoDevices, "downloading images requires at least one device configured and enabled").Code(http.StatusBadRequest)
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
		acceptedDevices := api.getDevicesThatAcceptPost(ctx, post, devices)
		if len(acceptedDevices) == 0 {
			continue
		}
		log.New(ctx).Debug("downloading image", "post_id", post.GetID(), "post_url", post.GetImageURL(), "devices", acceptedDevices)
		wg.Add(1)
		api.imageSemaphore <- struct{}{}
		go func(ctx context.Context, post reddit.Post) {
			defer func() {
				<-api.imageSemaphore
				wg.Done()
			}()

			if imageFile := api.findImageFileForDevices(ctx, post, devices); imageFile != nil {
				defer cleanup(imageFile)
				err := api.saveImageToFSAndDatabase(ctx, imageFile, subreddit, post, acceptedDevices)
				if err != nil {
					log.New(ctx).Err(err).Error("failed to download subreddit image")
				}
				return
			}

			if err := api.downloadSubredditImage(ctx, post, subreddit, acceptedDevices); err != nil {
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

	imageHandler, err := api.reddit.DownloadImage(ctx, post, api.eventBroadcast)
	if err != nil {
		return errs.Wrapw(err, "failed to download image")
	}
	defer imageHandler.Close()

	// copy to temp dir first to avoid copying incomplete files.
	tmpImageFile, err := api.copyImageToTempDir(ctx, imageHandler)
	if err != nil {
		return errs.Wrapw(err, "failed to download image to temp file")
	}
	defer cleanup(tmpImageFile)

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

	return api.saveImageToFSAndDatabase(ctx, tmpImageFile, subreddit, post, devices)
}

type stat interface {
	Stat() (fs.FileInfo, error)
}

func (api *API) saveImageToFSAndDatabase(ctx context.Context, image io.ReadCloser, subreddit *models.Subreddit, post reddit.Post, devices models.DeviceSlice) (err error) {
	ctx, span := tracer.Start(ctx, "*API.saveImageToFSAndDatabase")
	defer span.End()
	defer image.Close()

	w, close, err := api.createDeviceImageWriters(post, devices)
	defer close()
	if err != nil {
		return errs.Wrapw(err, "failed to create image files")
	}
	log.New(ctx).Debug("saving image files", "post_id", post.GetID(), "post_url", post.GetImageURL(), "devices", devices)
	size, err := io.Copy(w, image)
	if err != nil {
		return errs.Wrapw(err, "failed to save image files")
	}
	if size == 0 {
		if s, ok := image.(stat); ok {
			if fi, err := s.Stat(); err == nil {
				size = fi.Size()
			}
		}
	}

	var many []*models.ImageSetter
	now := time.Now()
	for _, device := range devices {
		var nsfw int32
		if post.IsNSFW() {
			nsfw = 1
		}
		width, height := post.GetImageSize()

		many = append(many, &models.ImageSetter{
			Subreddit:             omit.From(subreddit.Name),
			Device:                omit.From(device.Slug),
			PostTitle:             omit.From(post.GetTitle()),
			PostURL:               omit.From(post.GetPostURL()),
			PostCreated:           omit.From(post.GetCreated().Unix()),
			PostName:              omit.From(post.GetName()),
			PostAuthor:            omit.From(post.GetAuthor()),
			PostAuthorURL:         omit.From(post.GetAuthorURL()),
			ImageWidth:            omit.From(int32(width)),
			ImageHeight:           omit.From(int32(height)),
			ImageSize:             omit.From(size),
			ImageRelativePath:     omit.From(post.GetImageRelativePath(device)),
			ThumbnailRelativePath: omit.From(post.GetThumbnailRelativePath()),
			ImageOriginalURL:      omit.From(post.GetImageURL()),
			NSFW:                  omit.From(nsfw),
			CreatedAt:             omit.From(now.Unix()),
			UpdatedAt:             omit.From(now.Unix()),
		})
	}

	log.New(ctx).Debug("inserting images to database", "images", many)
	api.lockf(func() {
		_, err = models.Images.InsertMany(ctx, api.db, many...)
	})
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
		if shouldDownloadPostForDevice(post, device) && !api.isImageEntryExists(ctx, post, device) {
			devs = append(devs, device)
		}
	}
	return devs
}

// isImageEntryExists checks if the image entry already exists in the database and
// the image file actually exists in the filesystem.
func (api *API) isImageEntryExists(ctx context.Context, post reddit.Post, device *models.Device) (found bool) {
	ctx, span := tracer.Start(ctx, "*API.IsImageExists")
	defer span.End()

	exist, errQuery := models.Images.Query(ctx, api.db,
		models.SelectWhere.Images.Device.EQ(device.Slug),
		models.SelectWhere.Images.PostName.EQ(post.GetName()),
	).Exists()
	if errQuery != nil {
		return false
	}
	if !exist {
		return false
	}
	stat, errStat := os.Stat(post.GetImageTargetPath(api.config, device))
	if errStat != nil {
		return false
	}
	return stat.Size() > 0
}

// findImageFileForDevice finds if any of the image file exists for given devices.
//
// This helps to avoid downloading the same image for different devices.
//
// Return nil if no image file exists for the devices.
//
// Ensure to close the file after use.
func (api *API) findImageFileForDevices(ctx context.Context, post reddit.Post, devices models.DeviceSlice) *os.File {
	for _, device := range devices {
		stat, err := os.Stat(post.GetImageTargetPath(api.config, device))
		if err == nil {
			var err error
			oldImageFile, err := os.Open(post.GetImageTargetPath(api.config, device))
			if err != nil {
				log.New(ctx).Err(err).Error("failed to open image file", "filename", post.GetImageTargetPath(api.config, device))
				return nil
			}
			defer oldImageFile.Close()

			tempFilename := filepath.Join(os.TempDir(), "redmage", stat.Name())

			tempFileWrite, err := os.Create(tempFilename)
			if err != nil {
				log.New(ctx).Err(err).Error("failed to create temp file", "filename", post.GetImageTargetPath(api.config, device))
				return nil
			}
			defer tempFileWrite.Close()

			_, err = io.Copy(tempFileWrite, oldImageFile)
			if err != nil {
				log.New(ctx).Err(err).Error("failed to copy image file", "filename", post.GetImageTargetPath(api.config, device))
				return nil
			}

			rf, err := os.Open(tempFilename)
			if err != nil {
				log.New(ctx).Err(err).Error("failed to open temp file", "filename", tempFileWrite.Name())
				return nil
			}

			return rf
		}
	}

	return nil
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

func (te *tempFile) Name() string {
	return te.filename
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
	err = os.MkdirAll(tmpDirname, 0o777)
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

type removeableFile interface {
	io.ReadCloser
	Name() string
}

func cleanup(file removeableFile) {
	_ = file.Close()
	_ = os.Remove(file.Name())
}
