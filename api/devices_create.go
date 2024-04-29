package api

import (
	"context"
	"errors"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/mattn/go-sqlite3"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type DeviceCreateParams = models.DeviceSetter

func (api *API) DevicesCreate(ctx context.Context, params *models.Device) (*models.Device, error) {
	ctx, span := tracer.Start(ctx, "*API.DevicesCreate")
	defer span.End()

	now := time.Now()
	device, err := models.Devices.Insert(ctx, api.db, &models.DeviceSetter{
		Slug:                 omit.From(params.Slug),
		Name:                 omit.From(params.Name),
		ResolutionX:          omit.From(params.ResolutionX),
		ResolutionY:          omit.From(params.ResolutionY),
		AspectRatioTolerance: omit.From(params.AspectRatioTolerance),
		MinX:                 omit.From(params.MinX),
		MinY:                 omit.From(params.MinY),
		MaxX:                 omit.From(params.MaxX),
		MaxY:                 omit.From(params.MaxY),
		NSFW:                 omit.From(params.NSFW),
		WindowsWallpaperMode: omit.From(params.WindowsWallpaperMode),
		Enable:               omit.From(params.Enable),
		CreatedAt:            omit.From(now.Unix()),
		UpdatedAt:            omit.From(now.Unix()),
	})
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return nil, errs.Wrapw(sqliteErr, "device already exists", "params", params).Code(409)
			}
		}
		return nil, errs.Wrapw(err, "failed to create device", "params", params)
	}
	return device, nil
}
