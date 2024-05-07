package api

import (
	"context"
	"strconv"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type DevicesListParams struct {
	Q      string
	Status int

	Limit   int64
	Offset  int64
	OrderBy string
	Sort    string
}

func (dlp *DevicesListParams) FillFromQuery(query Queryable) {
	statusStr := query.Get("status")
	switch statusStr {
	case "0":
		dlp.Status = 0
	case "1":
		dlp.Status = 1
	default:
		dlp.Status = -1
	}
	dlp.Q = query.Get("q")

	dlp.Limit, _ = strconv.ParseInt(query.Get("limit"), 10, 64)
	if dlp.Limit < 1 {
		dlp.Limit = 20
	}
	dlp.Offset, _ = strconv.ParseInt(query.Get("offset"), 10, 64)
	if dlp.Offset < 0 {
		dlp.Offset = 0
	}

	dlp.OrderBy = query.Get("order_by")
	dlp.Sort = query.Get("sort")
}

func (dlp DevicesListParams) Query() (expr []bob.Mod[*dialect.SelectQuery]) {
	expr = append(expr, dlp.CountQuery()...)

	if dlp.Limit > 0 {
		expr = append(expr, sm.Limit(dlp.Limit))
	}

	if dlp.Offset > 0 {
		expr = append(expr, sm.Offset(dlp.Offset))
	}

	if dlp.OrderBy != "" {
		order := sm.OrderBy(sqlite.Quote(dlp.OrderBy))
		if dlp.Sort == "desc" {
			expr = append(expr, order.Desc())
		} else {
			expr = append(expr, order.Asc())
		}
	} else {
		expr = append(expr, sm.OrderBy(models.DeviceColumns.Name).Asc())
	}

	return expr
}

func (dlp DevicesListParams) CountQuery() (expr []bob.Mod[*dialect.SelectQuery]) {
	if dlp.Status >= 0 {
		expr = append(expr, models.SelectWhere.Devices.Enable.EQ(int32(dlp.Status)))
	}

	if dlp.Q != "" {
		arg := sqlite.Arg("%" + dlp.Q + "%")
		expr = append(expr,
			sm.Where(
				models.DeviceColumns.Name.Like(arg).Or(models.DeviceColumns.Slug.Like(arg)),
			),
		)
	}

	return expr
}

type DevicesListResult struct {
	Devices models.DeviceSlice `json:"devices"`
	Total   int64              `json:"total"`
}

func (api *API) DevicesList(ctx context.Context, params DevicesListParams) (result DevicesListResult, err error) {
	ctx, span := tracer.Start(ctx, "*API.DevicesList")
	defer span.End()

	result.Devices, err = api.GetDevices(ctx, params)
	if err != nil {
		return result, errs.Wrapw(err, "failed to query devices", "params", params)
	}

	result.Total, err = models.Devices.Query(ctx, api.db, params.CountQuery()...).Count()
	if err != nil {
		return result, errs.Wrapw(err, "failed to count devices", "params", params)
	}

	return result, nil
}

func (api *API) GetDevices(ctx context.Context, params DevicesListParams) (result models.DeviceSlice, err error) {
	ctx, span := tracer.Start(ctx, "*API.GetDevices")
	defer span.End()

	result, err = models.Devices.Query(ctx, api.db, params.Query()...).All()
	if err != nil {
		return result, errs.Wrapw(err, "failed to query devices", "params", params)
	}

	return result, nil
}
