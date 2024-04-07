package api

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/tigorlazuardi/redmage/db/queries/subreddits"
	"github.com/tigorlazuardi/redmage/server/routes/middleware"
)

type API struct {
	Subreddits *subreddits.Queries
}

func (api *API) Register(router chi.Router) {
	router.Use(chimiddleware.RequestID)
	router.Get("/", HealthCheck)
	router.Get("/health", HealthCheck)
	router.Route("/subreddits", func(r chi.Router) {
		r.Use(chimiddleware.RequestLogger(middleware.ChiLogger{}))
		r.Get("/", api.ListSubreddits)
	})
}
