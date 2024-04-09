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
	if err := homeview.Home(vc).Render(ctx, rw); err != nil {
		log.New(ctx).Err(err).Error("failed to render home view")
	}
}
