package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views/subredditsview/addview"
)

func (routes *Routes) SubredditValidateScheduleHTMX(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.SubredditValidateScheduleHTMX")
	defer span.End()

	var data addview.ScheduleInputData

	enabled, _ := strconv.Atoi(r.FormValue("enable_schedule"))
	data.Disabled = enabled == 0
	data.Value = r.FormValue("schedule")

	if data.Value == "" {
		if err := addview.ScheduleInput(data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render schedule input")
		}
		return
	}

	if data.Disabled {
		if err := addview.ScheduleInput(data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render schedule input")
		}
		return
	}

	scheduler, err := cronParser.Parse(data.Value)
	if err != nil {
		data.Error = fmt.Sprintf("Invalid schedule format: %s", err.Error())
		if err := addview.ScheduleInput(data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render schedule input")
		}
		return
	}

	next := scheduler.Next(time.Now())

	data.Valid = fmt.Sprintf("Syntax is valid. Next run at: %s", next.Format("Monday, _2 January 2006 15:04 MST"))

	if err := addview.ScheduleInput(data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render schedule input")
	}
}
