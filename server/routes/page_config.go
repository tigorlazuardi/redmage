package routes

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/configview"
)

func (routes *Routes) PageConfig(rw http.ResponseWriter, req *http.Request) {
	vc := views.NewContext(routes.Config, req)

	if err := configview.Config(vc).Render(req.Context(), rw); err != nil {
		log.New(req.Context()).Err(err).Error("Failed to render config page")
	}
}
