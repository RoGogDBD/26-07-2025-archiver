package repository

import (
	"strings"
	"testing"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddAndGetTask(t *testing.T) {
	tests := []struct {
		name      string
		addTasks  []*models.Task
		getTaskID string
		wantFound bool
		wantName  string
	}{
		{
			name:      "Task exists",
			addTasks:  []*models.Task{{ID: "1", Name: "task1"}},
			getTaskID: "1",
			wantFound: true,
			wantName:  "task1",
		},
		{
			name:      "Task does not exist",
			addTasks:  []*models.Task{{ID: "1", Name: "task1"}},
			getTaskID: "2",
			wantFound: false,
			wantName:  "",
		},
		{
			name:      "Empty storage",
			addTasks:  []*models.Task{},
			getTaskID: "1",
			wantFound: false,
			wantName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewTaskStorage()
			for _, task := range tt.addTasks {
				storage.AddTask(task)
			}

			gotTask, found := storage.GetTask(tt.getTaskID)
			assert.Equal(t, tt.wantFound, found)
			if found {
				assert.Equal(t, tt.wantName, gotTask.Name)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name        string
		initialTask *models.Task
		updateTask  *models.Task
		wantName    string
	}{
		{
			name:        "Update existing task",
			initialTask: &models.Task{ID: "1", Name: "task1"},
			updateTask:  &models.Task{ID: "1", Name: "updated"},
			wantName:    "updated",
		},
		{
			name:        "Update non-existing task (should add)",
			initialTask: nil,
			updateTask:  &models.Task{ID: "2", Name: "new"},
			wantName:    "new",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewTaskStorage()
			if tt.initialTask != nil {
				storage.AddTask(tt.initialTask)
			}

			storage.UpdateTask(tt.updateTask)
			gotTask, ok := storage.GetTask(tt.updateTask.ID)
			assert.True(t, ok)
			assert.Equal(t, tt.wantName, gotTask.Name)
		})
	}
}

func TestCanCreateTask(t *testing.T) {
	tests := []struct {
		name     string
		addTasks []*models.Task
		want     bool
	}{
		{
			name:     "No tasks - can create",
			addTasks: []*models.Task{},
			want:     true,
		},
		{
			name:     "Less than 3 tasks - can create",
			addTasks: []*models.Task{{ID: "1"}, {ID: "2"}},
			want:     true,
		},
		{
			name:     "3 tasks - cannot create",
			addTasks: []*models.Task{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			want:     false,
		},
		{
			name:     "More than 3 tasks - cannot create",
			addTasks: []*models.Task{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewTaskStorage()
			for _, task := range tt.addTasks {
				storage.AddTask(task)
			}

			got := storage.CanCreateTask()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name         string
		initialTasks []*models.Task
		deleteID     string
		checkID      string
		wantFound    bool
	}{
		{
			name:         "Delete existing task",
			initialTasks: []*models.Task{{ID: "1", Name: "task1"}},
			deleteID:     "1",
			checkID:      "1",
			wantFound:    false,
		},
		{
			name:         "Delete non-existing task",
			initialTasks: []*models.Task{{ID: "1", Name: "task1"}},
			deleteID:     "2",
			checkID:      "1",
			wantFound:    true,
		},
		{
			name:         "Delete from empty storage",
			initialTasks: []*models.Task{},
			deleteID:     "1",
			checkID:      "1",
			wantFound:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewTaskStorage()
			for _, task := range tt.initialTasks {
				storage.AddTask(task)
			}

			storage.DeleteTask(tt.deleteID)
			_, found := storage.GetTask(tt.checkID)
			if found != tt.wantFound {
				t.Errorf("unexpected task presence after DeleteTask: got found=%v, want found=%v", found, tt.wantFound)
			}
		})
	}
}

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
			got := CreateTask(tt.input)
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
