package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/schedulehistories"
)

func (routes *Routes) PageScheduleHistory(rw http.ResponseWriter, req *http.Request) {
	ctx, span := tracer.Start(req.Context(), "*Routes.PageScheduleHistory")
	defer span.End()

	c := views.NewContext(routes.Config, req)

	var data schedulehistories.Data

	data.Params.FillFromQuery(req.URL.Query())
	result, err := routes.API.ScheduleHistoryList(ctx, data.Params)
	if err != nil {
		log.New(ctx).Err(err).Error("Failed to list schedule histories")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := schedulehistories.View(c, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("Failed to render schedule histories view")
		}
		return
	}
	data.ScheduleHistories = result

	latest, err := routes.API.ScheduleHistoryLatest(ctx)
	if err != nil {
		log.New(ctx).Err(err).Error("Failed to list schedule histories")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := schedulehistories.View(c, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("Failed to render schedule histories view")
		}
		return
	}

	if first := data.ScheduleHistories.GetFirst(); first != nil {
		if first.ID == latest.ID {
			data.IsCurrent = true
		}

		if data.IsCurrent && len(data.ScheduleHistories.Schedules) < int(data.ScheduleHistories.Total) {
			data.Params = api.ScheduleHistoryListParams{
				Subreddit: data.Params.Subreddit,
				Limit:     data.Params.Limit,
			}
			data.ScheduleHistories, err = routes.API.ScheduleHistoryList(ctx, data.Params)
			if err != nil {
				log.New(ctx).Err(err).Error("Failed to list schedule histories")
				code, message := errs.HTTPMessage(err)
				rw.WriteHeader(code)
				data.Error = message
				if err := schedulehistories.View(c, data).Render(ctx, rw); err != nil {
					log.New(ctx).Err(err).Error("Failed to render schedule histories view")
				}
				return
			}
		}
	}

	if err := schedulehistories.View(c, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("Failed to render schedule histories view")
	}
}
