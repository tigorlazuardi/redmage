package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/server/routes/api"
	"github.com/tigorlazuardi/redmage/server/routes/htmx"
)

type Server struct {
	handler http.Handler
}

func (svr *Server) Serve() {
}

func New(cfg *config.Config) *Server {
	router := chi.NewRouter()

	router.Route("/api", api.Register)

	router.Route("/htmx", htmx.Register)

	return &Server{
		handler: router,
	}
}
