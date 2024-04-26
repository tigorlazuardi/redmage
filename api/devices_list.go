package api

import (
	"context"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type DevicesListParams struct {
	All     bool
	Q       string
	Limit   int64
	Offset  int64
	OrderBy string
	Sort    string
	Active  bool
}

func (dlp DevicesListParams) Query() (expr []bob.Mod[*dialect.SelectQuery]) {
	expr = append(expr, dlp.CountQuery()...)
	if dlp.Active {
		expr = append(expr, models.SelectWhere.Devices.Enable.EQ(1))
	}

	if dlp.All {
		return expr
	}

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
	}

	return expr
}

func (dlp DevicesListParams) CountQuery() (expr []bob.Mod[*dialect.SelectQuery]) {
	if dlp.Active {
		expr = append(expr, models.SelectWhere.Devices.Enable.EQ(1))
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

	result.Devices, err = models.Devices.Query(ctx, api.db, params.Query()...).All()
	if err != nil {
		return result, errs.Wrapw(err, "failed to query devices", "params", params)
	}

	result.Total, err = models.Devices.Query(ctx, api.db, params.CountQuery()...).Count()
	if err != nil {
		return result, errs.Wrapw(err, "failed to count devices", "params", params)
	}

	return result, nil
}
