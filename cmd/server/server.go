package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/app"
	"github.com/RoGogDBD/25-07-2025-archiver/internal/config"
	"github.com/viant/afs"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func run() error {
	storage := app.NewTaskStorage()
	fs := afs.New()
	addr := config.ParseFlags()

	router := NewRouter(storage, fs, addr)

	fmt.Println("Server started at", addr.String())
	return http.ListenAndServe(addr.String(), router)
}
