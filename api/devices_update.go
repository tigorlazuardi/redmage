package api

import (
	"context"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) DevicesUpdate(ctx context.Context, slug string, update *models.DeviceSetter) (device *models.Device, err error) {
	ctx, span := tracer.Start(ctx, "*API.DevicesUpdate")
	defer span.End()

	device = &models.Device{Slug: slug}

	api.lockf(func() {
		err = models.Devices.Update(ctx, api.db, update, device)
	})
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return device, errs.Wrapw(err, "a device with the same slug id already exists").Code(409)
			}
		}
		return device, errs.Wrapw(err, "failed to update device", "slug", slug, "values", update)
	}

	if err := device.Reload(ctx, api.db); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return device, errs.Wrapw(err, "device not found", "slug", slug).Code(404)
		}
		return device, errs.Wrapw(err, "failed to reload device", "slug", slug)
	}

	return device, nil
}
