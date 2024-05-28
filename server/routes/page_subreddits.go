package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/subreddits"
)

func (routes *Routes) PageSubreddits(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.PageSubreddits")
	defer span.End()

	c := views.NewContext(routes.Config, r)
	var data subreddits.Data
	data.Params.FillFromQuery(r.URL.Query())

	var err error
	data.Subreddits, err = routes.API.ListSubredditsWithCover(ctx, data.Params)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to list subreddits")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := subreddits.View(c, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render subreddits")
		}
		return
	}

	if err := subreddits.View(c, data).Render(r.Context(), rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render subreddits view")
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
