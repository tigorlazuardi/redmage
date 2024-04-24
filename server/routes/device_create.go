package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/aarondl/opt/omit"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/pkg/telemetry"
)

func (routes *Routes) APIDeviceCreate(rw http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := tracer.Start(r.Context(), "*Routes.APIDeviceCreate")
	defer func() { telemetry.EndWithStatus(span, err) }()

	var body models.Device

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

	device, err := routes.API.DevicesCreate(ctx, &models.DeviceSetter{
		Slug:                 omit.From(body.Slug),
		Name:                 omit.From(body.Name),
		ResolutionX:          omit.From(body.ResolutionX),
		ResolutionY:          omit.From(body.ResolutionY),
		AspectRatioTolerance: omit.From(body.AspectRatioTolerance),
		MinX:                 omit.From(body.MinX),
		MinY:                 omit.From(body.MinY),
		MaxX:                 omit.From(body.MaxX),
		MaxY:                 omit.From(body.MaxY),
		NSFW:                 omit.From(body.NSFW),
		WindowsWallpaperMode: omit.From(body.WindowsWallpaperMode),
	})
	if err != nil {
		rw.WriteHeader(errs.FindCodeOrDefault(err, http.StatusInternalServerError))
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": errs.FindMessage(err)})
		return
	}

	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(device); err != nil {
		log.New(ctx).Err(err).Error("failed to marshal json api devices")
	}
}

var slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func validateCreateDevice(params models.Device) error {
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
