package www

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/views"
	"github.com/tigorlazuardi/redmage/views/homeview"
)

func (www *WWW) Home(rw http.ResponseWriter, r *http.Request) {
	_ = homeview.Home(views.NewContext(www.Config, r)).Render(r.Context(), rw)
}
