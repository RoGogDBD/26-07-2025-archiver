package app

import (
	"github.com/RoGogDBD/25-07-2025-archiver/internal/models"
	"github.com/google/uuid"
)

func createTask(name string) *models.Task {
	return &models.Task{
		ID:     uuid.New().String(),
		Name:   name,
		URLs:   []string{},
		Status: "pending",
		Errors: nil,
	}
}
