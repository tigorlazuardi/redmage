package api

import (
	"net/http"

	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/tigorlazuardi/redmage/models"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"golang.org/x/net/context"
)

func (api *API) ScheduleHistoryLatest(ctx context.Context) (result *models.ScheduleHistory, err error) {
	ctx, span := tracer.Start(ctx, "*API.ScheduleHistoryLatest")
	defer span.End()

	result, err = models.ScheduleHistories.Query(ctx, api.db, sm.OrderBy(models.ScheduleHistoryColumns.CreatedAt).Desc()).One()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return result, errs.Wrapw(err, "last schedule history not found").Code(http.StatusNotFound)
		}
		return result, errs.Wrapw(err, "failed to find last schedule history")
	}
	return result, nil
}
