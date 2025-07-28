package repository

import (
	"sync"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/models"
	"github.com/google/uuid"
)

type TaskRepository interface {
	AddTask(*models.Task)
	GetTask(string) (*models.Task, bool)
	UpdateTask(*models.Task)
	DeleteTask(string)
	CanCreateTask() bool
}

type TaskStorage struct {
	mu    sync.RWMutex
	tasks map[string]*models.Task
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks: make(map[string]*models.Task),
	}
}

func (s *TaskStorage) AddTask(t *models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.ID] = t
}

func (s *TaskStorage) GetTask(id string) (*models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *TaskStorage) UpdateTask(t *models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.ID] = t
}

func (s *TaskStorage) CanCreateTask() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tasks) < 3
}

func (s *TaskStorage) DeleteTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id)
}

func CreateTask(name string) *models.Task {
	return &models.Task{
		ID:     uuid.New().String(),
		Name:   name,
		URLs:   []string{},
		Status: "pending",
		Errors: nil,
	}
}
