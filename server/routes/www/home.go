package www

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/views/homeview"
)

func (www *WWW) Home(rw http.ResponseWriter, r *http.Request) {
	_ = homeview.Home().Render(r.Context(), rw)
}
