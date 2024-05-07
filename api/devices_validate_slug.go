package api

import (
	"context"

	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func (api *API) DevicesValidateSlug(ctx context.Context, slug string) (exist bool, err error) {
	ctx, span := tracer.Start(ctx, "*API.DevicesValidateSlug")
	defer span.End()

	exist, err = models.Devices.Query(ctx, api.db, models.SelectWhere.Devices.Slug.EQ(slug)).Exists()
	if err != nil {
		return exist, errs.Wrapw(err, "failed to check device slug existence", "slug", slug)
	}
	return exist, err
}
