package www

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/db/queries/subreddits"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/server/routes/middleware"
)

var PublicFS = os.DirFS("public")

type WWW struct {
	Subreddits *subreddits.Queries
	Config     *config.Config
}

func (www *WWW) Register(router chi.Router) {
	router.Use(chimiddleware.RealIP)
	router.
		With(chimiddleware.RequestLogger(middleware.ChiSimpleLogger{})).
		Mount("/public", http.StripPrefix("/public", http.FileServer(http.FS(PublicFS))))

	if www.Config.Bool("http.hotreload") {
		log.New(context.Background()).Debug("enabled hot reload")
		router.
			With(chimiddleware.RequestLogger(middleware.ChiSimpleLogger{})).
			Get("/hot_reload", www.CreateHotReloadRoute())
	}

	router.Group(func(r chi.Router) {
		r.Use(chimiddleware.RequestID)
		r.Use(chimiddleware.RequestLogger(middleware.ChiLogger{}))
		r.Use(chimiddleware.SetHeader("Content-Type", "text/html; utf-8"))
		r.Get("/", www.Home)
	})
}
