package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/subreddits/put"
)

func (routes *Routes) PageSubredditsAdd(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "Routes.PageSubredditsAdd")
	defer span.End()

	c := views.NewContext(routes.Config, r)

	data := put.Data{
		Title:      "Add Subreddit",
		PostAction: "/htmx/subreddits/add",
	}

	if err := put.View(c, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render subreddits add page")
	}
}
