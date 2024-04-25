package api

import (
	"context"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) DevicesUpdate(ctx context.Context, id int, update *models.DeviceSetter) (err error) {
	ctx, span := tracer.Start(ctx, "*API.DevicesUpdate")
	defer span.End()

	err = models.Devices.Update(ctx, api.db, update, &models.Device{ID: int32(id)})
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return errs.Wrapw(err, "a device with the same slug id already exists").Code(409)
			}
		}
		return errs.Wrapw(err, "failed to update device", "id", id, "values", update)
	}

	return
}
