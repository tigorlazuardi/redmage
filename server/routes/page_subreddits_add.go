package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/subredditsview/addview"
)

func (routes *Routes) PageSubredditsAdd(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "Routes.PageSubredditsAdd")
	defer span.End()

	c := views.NewContext(routes.Config, r)

	if err := addview.Addview(c).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render subreddits add page")
	}
}
