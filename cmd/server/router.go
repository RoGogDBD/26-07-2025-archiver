package main

import (
	"net/http"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/app"
	"github.com/RoGogDBD/25-07-2025-archiver/internal/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(storage *app.TaskStorage, fs app.FSFetcher, addr *config.NetAddress) http.Handler {
	r := chi.NewRouter()

	handler := &app.Handler{
		Storage: storage,
		FS:      fs,
		Addr:    addr,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/task", func(r chi.Router) {
		r.Post("/", handler.HandleCreate)
		r.Route("/{id}", func(r chi.Router) {
			r.Post("/add", handler.HandleAdd)
			r.Get("/status", handler.HandleStatus)
		})
	})

	r.Get("/archives/{id}", handler.HandleDownloadArchive)

	return r
}
