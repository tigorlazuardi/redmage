package routes

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/server/routes/events"
	"github.com/tigorlazuardi/redmage/server/routes/middleware"
)

type Routes struct {
	API       *api.API
	Config    *config.Config
	PublicDir fs.FS
}

func (routes *Routes) Register(router chi.Router) {
	router.Use(chimiddleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "image/svg+xml", "image/x-icon"))

	router.HandleFunc("/ping", routes.HealthCheck)
	router.HandleFunc("/health", routes.HealthCheck)
	if routes.Config.Bool("http.hotreload") {
		router.Get("/hot_reload", routes.CreateHotReloadRoute())
	}

	router.Route("/htmx", routes.registerHTMXRoutes)
	router.Route("/api/v1", routes.registerV1APIRoutes)
	eventHandler := events.NewHandler(routes.Config, routes.API.GetEventBroadcaster())
	router.Route("/events", eventHandler.Route)

	router.Group(routes.registerWWWRoutes)
}

func (routes *Routes) registerV1APIRoutes(router chi.Router) {
	router.Use(otelchi.Middleware("redmage"))
	router.Use(chimiddleware.RequestLogger(middleware.ChiLogger{}))
	router.Use(chimiddleware.SetHeader("Content-Type", "application/json"))

	router.Post("/subreddits/start", routes.SubredditStartDownloadAPI)
	router.Get("/subreddits", routes.SubredditsListAPI)
	router.Post("/subreddits", routes.SubredditsCreateAPI)
	router.Post("/subreddits/check", routes.SubredditsCheckAPI)

	router.Get("/devices", routes.APIDeviceList)
	router.Post("/devices", routes.APIDeviceCreate)
	router.Patch("/devices/{slug}", routes.APIDeviceUpdate)

	router.Get("/images", routes.ImagesListAPI)

	router.Get("/events", routes.EventsAPI)
}

func (routes *Routes) registerHTMXRoutes(router chi.Router) {
	router.Use(otelchi.Middleware("redmage"))
	router.Use(chimiddleware.RequestLogger(middleware.ChiLogger{}))
	router.Use(chimiddleware.SetHeader("Content-Type", "text/html; charset=utf-8"))

	router.Post("/subreddits/add", routes.SubredditsCreateHTMX)
	router.Post("/subreddits/start", routes.SubredditStartDownloadHTMX)
	router.Get("/subreddits/start", routes.SubredditStartDownloadHTMX)
	router.Get("/subreddits/check", routes.SubredditCheckHTMX)
	router.Get("/subreddits/validate/schedule", routes.SubredditValidateScheduleHTMX)

	router.Get("/devices/add/validate/slug", routes.DevicesValidateSlugHTMX)
	router.Get("/devices/add/validate/name", routes.DevicesValidateNameHTMX)
}

func (routes *Routes) registerWWWRoutes(router chi.Router) {
	router.Mount("/public", http.StripPrefix("/public", http.FileServer(http.FS(routes.PublicDir))))

	imagesDir := http.FS(os.DirFS(routes.Config.String("download.directory")))

	router.Mount("/img", http.StripPrefix("/img", http.FileServer(imagesDir)))

	router.Group(func(r chi.Router) {
		r.Use(otelchi.Middleware("redmage"))
		r.Use(chimiddleware.RequestID)
		r.Use(chimiddleware.RequestLogger(middleware.ChiLogger{}))
		r.Use(chimiddleware.SetHeader("Content-Type", "text/html; charset=utf-8"))
		r.Get("/", routes.PageHome)
		r.Get("/subreddits", routes.PageSubreddits)
		r.Get("/subreddits/details/{name}", routes.PageSubredditsDetails)
		r.Get("/subreddits/add", routes.PageSubredditsAdd)
		r.Get("/subreddits/edit/{name}", routes.PageSubredditsEdit)
		r.Post("/subreddits/edit/{name}", routes.SubredditsEditHTMX)
		r.Get("/config", routes.PageConfig)
		r.Get("/devices", routes.PageDevices)
		r.Get("/devices/add", routes.PageDevicesAdd)
		r.Post("/devices/add", routes.DevicesCreateHTMX)
		r.Get("/devices/details/{slug}", routes.PageDeviceDetails)
		r.Get("/devices/edit/{slug}", routes.PageDevicesEdit)
		r.Post("/devices/edit/{slug}", routes.DevicesUpdateHTMX)
		r.Get("/history", routes.PageScheduleHistory)
	})
}
