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
	"github.com/tigorlazuardi/redmage/views/components"
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
		Enable:               omit.FromCond(body.Enable, body.Enable == 1 || body.Enable == 0),
		ResolutionX:          omit.FromCond(body.ResolutionX, body.ResolutionX != 0),
		ResolutionY:          omit.FromCond(body.ResolutionY, body.ResolutionY != 0),
		AspectRatioTolerance: omit.FromPtr(body.AspectRatioTolerance),
		MinX:                 omit.FromPtr(body.MinX),
		MinY:                 omit.FromPtr(body.MinY),
		MaxX:                 omit.FromPtr(body.MaxX),
		MaxY:                 omit.FromPtr(body.MaxY),
		NSFW:                 omit.FromPtr(body.NSFW),
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

func (routes *Routes) DevicesUpdateHTMX(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.DevicesUpdateHTMX")
	defer span.End()

	slug := chi.URLParam(req, "slug")
	exist, err := routes.API.DevicesExist(ctx, slug)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to check device slug existence")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		if err := components.ErrorToast(message).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	if !exist {
		rw.WriteHeader(http.StatusNotFound)
		if err := components.ErrorToast("Device with slug identifier '%s' does not exist", slug).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	device, err := routes.API.DevicesUpdate(ctx, slug, deviceSetterFromRequest(req))
	if err != nil {
		log.New(ctx).Err(err).Error("failed to update device")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		if err := components.ErrorToast(message).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	rw.Header().Set("HX-Redirect", fmt.Sprintf("/devices/details/%s", slug))
	if err := components.SuccessToast("Device %q has been updated", device.Name).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render success notification")
	}
}

func deviceSetterFromRequest(req *http.Request) *models.DeviceSetter {
	setter := &models.DeviceSetter{
		UpdatedAt: omit.From(time.Now().Unix()),
	}

	name := req.FormValue("name")
	setter.Name = omit.FromCond(name, name != "")

	enable, err := strconv.Atoi(req.FormValue("enable"))
	if err == nil {
		if enable > 1 {
			enable = 1
		} else if enable < 0 {
			enable = 0
		}
		setter.Enable = omit.From(int32(enable))
	}

	resx, _ := strconv.Atoi(req.FormValue("resolution_x"))
	setter.ResolutionX = omit.FromCond(float64(resx), resx != 0)

	resy, _ := strconv.Atoi(req.FormValue("resolution_y"))
	setter.ResolutionY = omit.FromCond(float64(resy), resy != 0)

	art, err := strconv.ParseFloat(req.FormValue("aspect_ratio_tolerance"), 64)
	if err == nil {
		setter.AspectRatioTolerance = omit.FromCond(art, art >= 0)
	}

	minX, err := strconv.Atoi(req.FormValue("min_x"))
	if err == nil {
		setter.MinX = omit.FromCond(int32(minX), minX >= 0)
	}

	minY, err := strconv.Atoi(req.FormValue("min_y"))
	if err == nil {
		setter.MinY = omit.FromCond(int32(minY), minY >= 0)
	}

	maxX, err := strconv.Atoi(req.FormValue("max_x"))
	if err == nil {
		setter.MaxX = omit.FromCond(int32(maxX), maxX >= 0)
	}

	maxY, err := strconv.Atoi(req.FormValue("max_y"))
	if err == nil {
		setter.MaxY = omit.FromCond(int32(maxY), maxY >= 0)
	}

	nsfw, err := strconv.Atoi(req.FormValue("nsfw"))
	if err == nil {
		if nsfw > 1 {
			nsfw = 1
		} else if nsfw < 0 {
			nsfw = 0
		}
		setter.NSFW = omit.From(int32(nsfw))
	}

	return setter
}
