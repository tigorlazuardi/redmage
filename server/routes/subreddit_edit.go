package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/robfig/cron/v3"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views/components"
)

func (routes *Routes) SubredditsEditHTMX(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.SubredditsEditHTMX")
	defer span.End()

	name := chi.URLParam(r, "name")
	countbackInt, _ := strconv.Atoi(r.FormValue("countback"))
	countback := int32(countbackInt)
	schedule := r.FormValue("schedule")

	if countback < 1 {
		rw.WriteHeader(http.StatusBadRequest)
		const msg = "Countback must be greater than 0"
		if err := components.ErrorNotication(msg).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	if _, err := cron.ParseStandard(schedule); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Invalid schedule format: %s", err)
		if err := components.ErrorNotication(msg).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	_, err := routes.API.SubredditsEdit(ctx, api.SubredditEditParams{
		Name:      name,
		Countback: &countback,
		Schedule:  &schedule,
	})
	if err != nil {
		log.New(ctx).Err(err).Error("failed to update device")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		if err := components.ErrorToast(message).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render error notification")
		}
		return
	}

	rw.Header().Set("HX-Retarget", "#root-content")
	rw.Header().Set("HX-Reselect", "#root-content")
	rw.Header().Set("HX-Push-Url", "/subreddits/details/"+name)
	r.Header.Set("X-Flash-Message-Success", fmt.Sprintf("Subreddit %s updated", name))
	routes.PageSubredditsDetails(rw, r)
}
