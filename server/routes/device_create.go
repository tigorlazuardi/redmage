package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/telemetry"
	"github.com/tigorlazuardi/redmage/views/components"
)

func (routes *Routes) APIDeviceCreate(rw http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := tracer.Start(r.Context(), "*Routes.APIDeviceCreate")
	defer func() { telemetry.EndWithStatus(span, err) }()

	var body *models.Device

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.New(ctx).Err(err).Error("failed to decode json body")
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": fmt.Sprintf("cannot decode json body: %s", err)})
		return
	}

	if err = validateCreateDevice(body); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": err.Error()})
		return
	}

	device, err := routes.API.DevicesCreate(ctx, body)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to create device", "body", body)
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": message})
		return
	}

	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(device); err != nil {
		log.New(ctx).Err(err).Error("failed to marshal json api devices")
	}
}

var slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func validateCreateDevice(params *models.Device) error {
	if params.Name == "" {
		return errors.New("name is required")
	}
	if params.Slug == "" {
		return errors.New("slug is required")
	}
	if !slugRegex.MatchString(params.Slug) {
		return errors.New("slug must be lowercase alphanumeric with dash. eg: my-awesome-laptop")
	}
	if params.ResolutionX < 1 {
		return errors.New("device width resolution is required")
	}
	if params.ResolutionY < 1 {
		return errors.New("device height resolution is required")
	}
	if params.AspectRatioTolerance < 0 {
		return errors.New("aspect ratio tolerance cannot be negative value")
	}
	if params.MaxX < 0 {
		params.MaxX = 0
	}
	if params.MaxY < 0 {
		params.MaxY = 0
	}
	if params.MinX < 0 {
		params.MinX = 0
	}
	if params.MinY < 0 {
		params.MinY = 0
	}
	if params.NSFW < 0 {
		params.NSFW = 0
	}
	if params.NSFW > 1 {
		params.NSFW = 1
	}
	if params.WindowsWallpaperMode < 0 {
		params.WindowsWallpaperMode = 0
	}
	if params.WindowsWallpaperMode > 1 {
		params.WindowsWallpaperMode = 1
	}
	return nil
}

func (routes *Routes) DevicesCreateHTMX(rw http.ResponseWriter, req *http.Request) {
	var err error
	ctx, span := tracer.Start(req.Context(), "*Routes.DevicesCreateHTMX")
	defer func() { telemetry.EndWithStatus(span, err) }()

	device, err := createDeviceFromParams(req)
	if err != nil {
		rw.WriteHeader(400)
		if err := components.ErrorToast(err.Error()).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	_, err = routes.API.DevicesCreate(ctx, device)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to create device", "device", device)
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		if err := components.ErrorToast(message).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	rw.Header().Set("HX-Redirect", "/devices")
	rw.WriteHeader(http.StatusCreated)

	if err := components.SuccessToast("device created").Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render success notification")
	}
}

func createDeviceFromParams(req *http.Request) (*models.Device, error) {
	device := new(models.Device)

	device.Enable = 1
	device.Name = req.FormValue("name")
	device.Slug = req.FormValue("slug")
	device.ResolutionX, _ = strconv.ParseFloat(req.FormValue("resolution_x"), 32)
	device.ResolutionY, _ = strconv.ParseFloat(req.FormValue("resolution_y"), 32)
	device.AspectRatioTolerance, _ = strconv.ParseFloat(req.FormValue("aspect_ratio_tolerance"), 32)

	maxX, _ := strconv.ParseInt(req.FormValue("max_x"), 10, 32)
	device.MaxX = int32(maxX)

	maxY, _ := strconv.ParseInt(req.FormValue("max_y"), 10, 32)
	device.MaxY = int32(maxY)

	minX, _ := strconv.ParseInt(req.FormValue("min_x"), 10, 32)
	device.MinX = int32(minX)

	minY, _ := strconv.ParseInt(req.FormValue("min_y"), 10, 32)
	device.MinY = int32(minY)

	nsfw, _ := strconv.ParseInt(req.FormValue("nsfw"), 10, 32)
	device.NSFW = int32(nsfw)

	windowsWallpaperMode, _ := strconv.ParseInt(req.FormValue("windows_wallpaper_mode"), 10, 32)
	device.WindowsWallpaperMode = int32(windowsWallpaperMode)

	return device, validateCreateDevice(device)
}
