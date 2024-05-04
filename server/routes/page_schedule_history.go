package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	scheduleshistoryview "github.com/tigorlazuardi/redmage/views/schedulehistoriesview"
)

func (routes *Routes) PageScheduleHistory(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.PageScheduleHistory")
	defer span.End()

	c := views.NewContext(routes.Config, req)

	var data scheduleshistoryview.Data

	if err := scheduleshistoryview.ScheduleHistoriesview(c, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("Failed to render schedule histories view")
	}
}
