package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/app"
	"github.com/RoGogDBD/25-07-2025-archiver/internal/config"
	"github.com/RoGogDBD/25-07-2025-archiver/internal/models"
	"github.com/viant/afs"
)

func TestNewRouter(t *testing.T) {
	storage := app.NewTaskStorage()
	fs := afs.New()
	addr := &config.NetAddress{Host: "localhost", Port: 8080}

	storage.AddTask(&models.Task{
		ID:      "taskWithNoArchive",
		URLs:    []string{"http://example.com/file1.pdf"},
		Archive: "",
	})

	router := NewRouter(storage, fs, addr)

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		wantStatusCode int
	}{
		{
			name:           "Create task - valid",
			method:         "POST",
			url:            "/task/",
			body:           `{"name":"test task"}`,
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "Create task - empty name",
			method:         "POST",
			url:            "/task/",
			body:           `{"name":""}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Add URL - nonexistent task",
			method:         "POST",
			url:            "/task/nonexistent/add",
			body:           `{"url":"http://example.com/file.pdf"}`,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "Status - nonexistent task",
			method:         "GET",
			url:            "/task/nonexistent/status",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "Download archive - nonexistent task",
			method:         "GET",
			url:            "/archives/nonexistent",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "Download archive - no archive",
			method:         "GET",
			url:            "/archives/taskWithNoArchive",
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, tt.url, nil)
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d, response body: %s", tt.wantStatusCode, rec.Code, rec.Body.String())
			}
		})
	}
}
