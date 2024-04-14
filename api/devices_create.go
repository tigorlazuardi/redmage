package api

import (
	"context"

	"github.com/tigorlazuardi/redmage/db/queries"
)

type DeviceCreateParams = queries.DeviceCreateParams

func (api *API) DevicesCreate(ctx context.Context, params DeviceCreateParams) (queries.Device, error) {
	ctx, span := tracer.Start(ctx, "*API.DevicesCreate")
	defer span.End()
	return api.queries.DeviceCreate(ctx, params)
}
