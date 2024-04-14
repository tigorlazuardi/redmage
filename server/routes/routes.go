package routes

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/server/routes/middleware"
)

type Routes struct {
	API       *api.API
	Config    *config.Config
	PublicDir fs.FS
}

func (routes *Routes) Register(router chi.Router) {
	router.HandleFunc("/ping", routes.HealthCheck)
	router.HandleFunc("/health", routes.HealthCheck)
	if routes.Config.Bool("http.hotreload") {
		router.Get("/hot_reload", routes.CreateHotReloadRoute())
	}

	router.Group(routes.registerWWWRoutes)
	router.Route("/api/v1", routes.registerV1APIRoutes)
}

func (routes *Routes) registerV1APIRoutes(router chi.Router) {
	router.Use(otelchi.Middleware("redmage"))
	router.Use(chimiddleware.RequestLogger(middleware.ChiLogger{}))
	router.Use(chimiddleware.SetHeader("Content-Type", "application/json"))

	router.Get("/subreddits", routes.SubredditsListAPI)
	router.Get("/devices", routes.APIDeviceList)
}

func (routes *Routes) registerWWWRoutes(router chi.Router) {
	router.Mount("/public", http.StripPrefix("/public", http.FileServer(http.FS(routes.PublicDir))))

	router.Group(func(r chi.Router) {
		r.Use(otelchi.Middleware("redmage"))
		r.Use(chimiddleware.RequestID)
		r.Use(chimiddleware.RequestLogger(middleware.ChiLogger{}))
		r.Use(chimiddleware.SetHeader("Content-Type", "text/html; charset=utf-8"))
		r.Get("/", routes.PageHome)
	})
}
