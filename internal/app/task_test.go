package app

import (
	"strings"
	"testing"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTask *models.Task
	}{
		{
			name:  "CreateTask_Success_ValidName",
			input: "test-task",
			wantTask: &models.Task{
				Name:    "test-task",
				URLs:    []string{},
				Status:  "pending",
				Archive: "",
				Errors:  nil,
			},
		},
		{
			name:  "CreateTask_Success_EmptyName",
			input: "",
			wantTask: &models.Task{
				Name:    "",
				URLs:    []string{},
				Status:  "pending",
				Archive: "",
				Errors:  nil,
			},
		},
		{
			name:  "CreateTask_Success_LongName",
			input: strings.Repeat("a", 100),
			wantTask: &models.Task{
				Name:    strings.Repeat("a", 100),
				URLs:    []string{},
				Status:  "pending",
				Archive: "",
				Errors:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createTask(tt.input)
			_, err := uuid.Parse(got.ID)
			assert.NoError(t, err, "ID should be a valid UUID")
			assert.Equal(t, tt.wantTask.Name, got.Name, "Name should match")
			assert.Equal(t, tt.wantTask.URLs, got.URLs, "URLs should match")
			assert.Equal(t, tt.wantTask.Status, got.Status, "Status should match")
			assert.Equal(t, tt.wantTask.Archive, got.Archive, "Archive should match")
			assert.Equal(t, tt.wantTask.Errors, got.Errors, "Errors should match")
		})
	}
}
