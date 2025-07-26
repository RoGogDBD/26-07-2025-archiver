package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreate(t *testing.T) {
	tests := []struct {
		name           string
		storageSetup   func(s *TaskStorage)
		body           string
		wantStatusCode int
	}{
		{
			name: "valid request",
			storageSetup: func(s *TaskStorage) {
			},
			body:           `{"name":"task1"}`,
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "empty name",
			storageSetup:   func(s *TaskStorage) {},
			body:           `{"name":""}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			storageSetup:   func(s *TaskStorage) {},
			body:           `not json`,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewTaskStorage()
			tt.storageSetup(storage)
			handler := &Handler{
				Storage: storage,
				FS:      nil,
				Addr:    &config.NetAddress{Host: "x", Port: 0},
			}

			req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler.HandleCreate(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
		})
	}
}