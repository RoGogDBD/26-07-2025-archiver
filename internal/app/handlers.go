package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/config"
	"github.com/RoGogDBD/25-07-2025-archiver/internal/models"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Storage *TaskStorage
	FS      FSFetcher
	Addr    *config.NetAddress
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if !h.Storage.CanCreateTask() {
		http.Error(w, "server is busy, maximum 3 tasks allowed", http.StatusServiceUnavailable)
		return
	}

	var req models.CreateTaskRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name is required and cannot be empty", http.StatusBadRequest)
		return
	}

	if len(req.Name) > 100 {
		http.Error(w, "name is too long (max 100 characters)", http.StatusBadRequest)
		return
	}

	task := createTask(req.Name)
	h.Storage.AddTask(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		ID string `json:"id"`
	}{ID: task.ID}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "task ID is required", http.StatusBadRequest)
		return
	}

	task, ok := h.Storage.GetTask(taskID)
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if len(task.URLs) >= 3 {
		http.Error(w, "maximum 3 URLs per task allowed", http.StatusBadRequest)
		return
	}

	var req models.AddURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.URL) == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	resource, err := h.FS.OpenURL(r.Context(), req.URL)
	if err != nil {
		http.Error(w, "unable to access file at provided URL", http.StatusBadRequest)
		return
	}
	defer resource.Close()

	header := make([]byte, 512)
	n, err := resource.Read(header)
	if err != nil && err != io.EOF {
		http.Error(w, "unable to read file header", http.StatusBadRequest)
		return
	}

	contentType := http.DetectContentType(header[:n])
	if !(contentType == "application/pdf" || contentType == "image/jpeg") {
		http.Error(w, "only .pdf and .jpeg files are allowed", http.StatusBadRequest)
		return
	}

	fullReader := io.MultiReader(bytes.NewReader(header[:n]), resource)

	_, err = io.ReadAll(fullReader)
	if err != nil {
		http.Error(w, "failed to read file content", http.StatusInternalServerError)
		return
	}

	task.URLs = append(task.URLs, req.URL)
	h.Storage.UpdateTask(task)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: "URL added successfully"})
}

func (h *Handler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "task ID is required", http.StatusBadRequest)
		return
	}

	task, ok := h.Storage.GetTask(taskID)
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(task.URLs) == 3 {
		if task.Archive == "" {
			archiveURL, err := h.CreateArchive(r.Context(), task)
			if err != nil {
				http.Error(w, "failed to create archive: "+err.Error(), http.StatusInternalServerError)
				return
			}
			task.Archive = archiveURL
			task.Status = "completed"
			h.Storage.UpdateTask(task)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			Status   string   `json:"status"`
			URLCount int      `json:"url_count"`
			Archive  string   `json:"archive_url,omitempty"`
			Errors   []string `json:"errors,omitempty"`
		}{
			Status:   task.Status,
			URLCount: len(task.URLs),
			Archive:  task.Archive,
			Errors:   task.Errors,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status   string   `json:"status"`
		URLCount int      `json:"url_count"`
		Errors   []string `json:"errors,omitempty"`
	}{
		Status:   task.Status,
		URLCount: len(task.URLs),
		Errors:   task.Errors,
	})
}

func (h *Handler) HandleDownloadArchive(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "task ID is required", http.StatusBadRequest)
		return
	}

	task, ok := h.Storage.GetTask(taskID)
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if task.Archive == "" {
		http.Error(w, "archive not available", http.StatusBadRequest)
		return
	}

	archivePath := fmt.Sprintf("/tmp/%s.zip", taskID)
	file, err := os.Open(archivePath)
	if err != nil {
		http.Error(w, "failed to open archive: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	defer os.Remove(archivePath)
	defer h.Storage.DeleteTask(taskID)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.zip"`, taskID))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}
