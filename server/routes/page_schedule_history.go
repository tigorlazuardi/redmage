package routes

import (
	"net/http"
	"time"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/schedulehistoriesview"
)

func (routes *Routes) PageScheduleHistory(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.PageScheduleHistory")
	defer span.End()

	c := views.NewContext(routes.Config, req)

	var data schedulehistoriesview.Data
	if tz := req.URL.Query().Get("tz"); tz == "" {
		data.Timezone = time.Local
	} else {
		var err error
		data.Timezone, err = time.LoadLocation(tz)
		if err != nil {
			data.Timezone = time.Local
		}
	}

	var params api.ScheduleHistoryListParams
	params.FillFromQuery(req.URL.Query())

	result, err := routes.API.ScheduleHistoryList(ctx, params)
	if err != nil {
		log.New(ctx).Err(err).Error("Failed to list schedule histories")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := schedulehistoriesview.ScheduleHistoriesview(c, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("Failed to render schedule histories view")
		}
		return
	}

	data.Schedules = result.Schedules
	data.Total = result.Total

	if err := schedulehistoriesview.ScheduleHistoriesview(c, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("Failed to render schedule histories view")
	}
}
