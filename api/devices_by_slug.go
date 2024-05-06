package api

import (
	"context"
	"net/http"

	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) DeviceBySlug(ctx context.Context, slug string) (device *models.Device, err error) {
	ctx, span := tracer.Start(ctx, "*API.DeviceByName")
	defer span.End()

	device, err = models.FindDevice(ctx, api.db, slug)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return device, errs.Wrapw(err, "device not found", "device", device).Code(http.StatusNotFound)
		}

		return device, errs.Wrapw(err, "failed to find device", "device", device)
	}

	return device, nil
}
