package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.New(ctx).Err(err).Error("failed to parse id")
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": fmt.Sprintf("bad id: %s", err)})
		return
	}

	var body deviceUpdate

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.New(ctx).Err(err).Error("failed to decode json body")
		rw.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": fmt.Sprintf("cannot decode json body: %s", err)})
		return
	}

	err = routes.API.DevicesUpdate(ctx, id, &models.DeviceSetter{
		Slug:                 omit.FromCond(body.Slug, body.Slug != ""),
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
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": message})
		return
	}

	_ = json.NewEncoder(rw).Encode(map[string]string{"message": "ok"})
}
