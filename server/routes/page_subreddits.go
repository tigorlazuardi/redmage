package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/subredditsview"
)

func (routes *Routes) PageSubreddits(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.PageSubreddits")
	defer span.End()

	data := subredditsview.Data{
		// Subreddits: routes.API.SubredditsList(ctx),
	}

	c := views.NewContext(routes.Config, r)

	if err := subredditsview.Subreddit(c, data).Render(r.Context(), rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render subreddits view")
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
