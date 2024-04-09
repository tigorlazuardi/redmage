package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/homeview"
)

func (routes *Routes) PageHome(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vc := views.NewContext(routes.Config, r)

	params := parseSubredditListQuery(r)

	list, err := routes.API.ListSubreddits(ctx, params)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	data := homeview.Data{
		SubredditsList: list,
		Error: err,
	}

	if err := homeview.Home(vc, data).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render home view")
	}
}
