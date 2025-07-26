package app

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/models"
)

func (h *Handler) CreateArchive(ctx context.Context, task *models.Task) (string, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for i, fileURL := range task.URLs {
		resource, err := h.FS.OpenURL(ctx, fileURL)
		if err != nil {
			task.Errors = append(task.Errors, fmt.Sprintf("failed to open URL %q: %v", fileURL, err))
			continue
		}

		data, err := io.ReadAll(resource)
		resource.Close()
		if err != nil {
			task.Errors = append(task.Errors, fmt.Sprintf("failed to read data from %q: %v", fileURL, err))
			continue
		}

		fileName := fmt.Sprintf("file%d%s", i+1, filepath.Ext(fileURL))
		f, err := zipWriter.Create(fileName)
		if err != nil {
			task.Errors = append(task.Errors, fmt.Sprintf("failed to create entry in zip for %q: %v", fileName, err))
			continue
		}

		_, err = f.Write(data)
		if err != nil {
			task.Errors = append(task.Errors, fmt.Sprintf("failed to write data to zip for %q: %v", fileName, err))
			continue
		}
	}

	if err := zipWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to close zip writer: %w", err)
	}

	archivePath := fmt.Sprintf("/tmp/%s.zip", task.ID)
	if err := os.WriteFile(archivePath, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("failed to save archive: %w", err)
	}

	archiveURL := fmt.Sprintf("http://%s/archives/%s", h.Addr.String(), task.ID)
	return archiveURL, nil
}
