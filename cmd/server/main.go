package main

import (
	"log"
	"net/http"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/config"
	"github.com/RoGogDBD/25-07-2025-archiver/internal/handler"
	"github.com/RoGogDBD/25-07-2025-archiver/internal/repository"
	"github.com/viant/afs"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func run() error {
	storage := repository.NewTaskStorage()
	fs := afs.New()
	addr := config.ParseFlags()

	h := handler.NewHandler(storage, fs, addr)

	router := NewRouter(h)

	log.Printf("Server started at %s\n", addr.String())
	return http.ListenAndServe(addr.String(), router)
}
