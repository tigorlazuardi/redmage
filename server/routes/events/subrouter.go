package events

import (
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/tigorlazuardi/redmage/server/routes/middleware"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func (handler *Handler) Route(router chi.Router) {
	router.Use(otelchi.Middleware("redmage"))
	router.Use(chimiddleware.RequestLogger(middleware.ChiLogger{}))
	router.Get("/", handler.HTMXEvents)
	router.Get("/simple", handler.SimpleEvents)
	router.Get("/json", handler.JSONEvents)
}
