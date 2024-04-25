package api

import (
	"context"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type DeviceCreateParams = models.DeviceSetter

func (api *API) DevicesCreate(ctx context.Context, params *DeviceCreateParams) (*models.Device, error) {
	ctx, span := tracer.Start(ctx, "*API.DevicesCreate")
	defer span.End()

	device, err := models.Devices.Insert(ctx, api.db, params)
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
