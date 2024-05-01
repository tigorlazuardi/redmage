package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/subredditsview"
)

func (routes *Routes) PageSubreddits(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.PageSubreddits")
	defer span.End()

	c := views.NewContext(routes.Config, r)

	var params api.ListSubredditsParams
	params.FillFromQuery(r.URL.Query())

	var data subredditsview.Data
	var err error

	data.Subreddits, err = routes.API.ListSubredditsWithCover(ctx, params)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to list subreddits")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		data.Error = message
		if err := subredditsview.Subreddit(c, data).Render(ctx, rw); err != nil {
			log.New(ctx).Err(err).Error("failed to render subreddits")
		}
		return
	}

	if err := subredditsview.Subreddit(c, data).Render(r.Context(), rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render subreddits view")
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
