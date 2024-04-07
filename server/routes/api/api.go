package api

import "github.com/go-chi/chi/v5"

func Register(router chi.Router) {
	router.Get("/", HealthCheck)
	router.Get("/health", HealthCheck)
}
