package api

import (
	"context"

	"github.com/tigorlazuardi/redmage/db/queries"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type DevicesListParams struct {
	All    bool
	Query  string
	Limit  int64
	Offset int64
}

type DevicesListResult struct {
	Devices []queries.Device `json:"devices"`
	Total   int64            `json:"total"`
}

func (api *API) DevicesList(ctx context.Context, params DevicesListParams) (result DevicesListResult, err error) {
	ctx, span := tracer.Start(ctx, "*API.DevicesList")
	defer span.End()

	q := params.Query

	if params.All {
		result.Devices, err = api.queries.DeviceGetAll(ctx)
		if err != nil {
			return result, errs.Wrapw(err, "failed to get all devices", "params", params)
		}
		result.Total, err = api.queries.DeviceCount(ctx)
		if err != nil {
			return result, errs.Wrapw(err, "failed to count all devices", "params", params)
		}
		return result, nil
	}

	if q != "" {
		like := "%" + q + "%"
		result.Devices, err = api.queries.DeviceSearch(ctx, queries.DeviceSearchParams{
			Name:   like,
			Slug:   like,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
		if err != nil {
			return result, errs.Wrapw(err, "failed to search device", "params", params)
		}
		result.Total, err = api.queries.DeviceSearchCount(ctx, queries.DeviceSearchCountParams{
			Name:   like,
			Slug:   like,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
		if err != nil {
			return result, errs.Wrapw(err, "failed to count device search", "params", params)
		}
		return result, nil
	}

	result.Devices, err = api.queries.DeviceList(ctx, queries.DeviceListParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return result, errs.Wrapw(err, "failed to list device", "params", params)
	}

	result.Total, err = api.queries.DeviceCount(ctx)
	if err != nil {
		return result, errs.Wrapw(err, "failed to count all devices", "params", params)
	}

	return result, nil
}
