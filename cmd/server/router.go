package main

import (
	"net/http"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/task", func(r chi.Router) {
		r.Post("/", h.HandleCreate)
		r.Route("/{id}", func(r chi.Router) {
			r.Post("/add", h.HandleAdd)
			r.Get("/status", h.HandleStatus)
		})
	})

	r.Get("/archives/{id}", h.HandleDownloadArchive)

	return r
}
