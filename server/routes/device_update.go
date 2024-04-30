package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/go-chi/chi/v5"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

type deviceUpdate struct {
	models.Device

	MinX                 *int32   `json:"min_x"`
	MinY                 *int32   `json:"min_y"`
	MaxX                 *int32   `json:"max_x"`
	MaxY                 *int32   `json:"max_y"`
	NSFW                 *int32   `json:"nsfw"`
	WindowsWallpaperMode *int32   `json:"windows_wallpaper_mode"`
	AspectRatioTolerance *float64 `json:"aspect_ratio_tolerance"`
}

func (routes *Routes) APIDeviceUpdate(rw http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := tracer.Start(r.Context(), "*Routes.APIDeviceUpdate")
	defer span.End()

	var (
		enc = json.NewEncoder(rw)
		dec = json.NewDecoder(r.Body)
	)

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": "missing name"})
		return
	}

	var body deviceUpdate

	if err = dec.Decode(&body); err != nil {
		log.New(ctx).Err(err).Error("failed to decode json body")
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": fmt.Sprintf("cannot decode json body: %s", err)})
		return
	}

	device, err := routes.API.DevicesUpdate(ctx, slug, &models.DeviceSetter{
		Name:                 omit.FromCond(body.Name, body.Name != ""),
		ResolutionX:          omit.FromCond(body.ResolutionX, body.ResolutionX != 0),
		ResolutionY:          omit.FromCond(body.ResolutionY, body.ResolutionY != 0),
		AspectRatioTolerance: omit.FromPtr(body.AspectRatioTolerance),
		MinX:                 omit.FromPtr(body.MinX),
		MinY:                 omit.FromPtr(body.MinY),
		MaxX:                 omit.FromPtr(body.MaxX),
		MaxY:                 omit.FromPtr(body.MaxY),
		NSFW:                 omit.FromPtr(body.NSFW),
		WindowsWallpaperMode: omit.FromPtr(body.WindowsWallpaperMode),
		UpdatedAt:            omit.From(time.Now().Unix()),
	})
	if err != nil {
		log.New(ctx).Err(err).Error("failed to update device")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = enc.Encode(map[string]string{"error": message})
		return
	}

	_ = enc.Encode(device)
}
