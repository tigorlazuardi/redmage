package api

import (
	"context"

	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) DevicesExist(ctx context.Context, slug string) (exist bool, err error) {
	ctx, span := tracer.Start(ctx, "API.DevicesExist")
	defer span.End()

	exist, err = models.Devices.Query(ctx, api.db, models.SelectWhere.Devices.Slug.EQ(slug)).Exists()
	if err != nil {
		return false, errs.Wrapw(err, "failed to check device existence", "slug", slug)
	}
	return exist, nil
}
