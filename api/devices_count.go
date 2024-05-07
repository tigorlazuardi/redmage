package api

import (
	"context"

	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) DevicesCountEnabled(ctx context.Context) (int64, error) {
	ctx, span := tracer.Start(ctx, "*API.DevicesCountEnabled")
	defer span.End()

	count, err := models.Devices.Query(ctx, api.db, models.SelectWhere.Devices.Enable.EQ(1)).Count()
	if err != nil {
		return 0, errs.Wrapw(err, "failed to count enabled devices")
	}
	return count, nil
}
