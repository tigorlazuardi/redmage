package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tigorlazuardi/redmage/api/reddit"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/components"
	"github.com/tigorlazuardi/redmage/views/subreddits/put"
)

func (routes *Routes) PageSubredditsEdit(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.PageSubredditsEdit")
	defer span.End()

	c := views.NewContext(routes.Config, r)

	name := chi.URLParam(r, "name")

	sub, err := routes.API.SubredditsGetByName(ctx, name)
	if err != nil {
		code, message := errs.HTTPMessage(err)
		if code >= 500 {
			log.New(ctx).Err(err).Error("failed to get device by slug")
		}
		rw.WriteHeader(code)
		msg := fmt.Sprintf("%d: %s", code, message)
		if err := components.PageError(c, msg).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render subreddit edit page")
		}
	}

	data := put.Data{
		Title:          fmt.Sprintf("Edit %s", sub.Name),
		EditMode:       true,
		PostAction:     fmt.Sprintf("/subreddits/edit/%s", sub.Name),
		NameInput:      put.NameInputData{Value: sub.Name},
		TypeInput:      put.TypeInputData{Value: reddit.SubredditType(sub.Subtype)},
		ScheduleInput:  put.ScheduleInputData{Value: sub.Schedule},
		CountbackInput: put.CountbackInputData{Value: int64(sub.Countback)},
	}

	if err := put.View(c, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render subreddit edit page")
	}
}
