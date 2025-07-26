package app

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/RoGogDBD/25-07-2025-archiver/internal/app/mocks"
	"github.com/RoGogDBD/25-07-2025-archiver/internal/config"
	"github.com/RoGogDBD/25-07-2025-archiver/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateArchive(t *testing.T) {
	tests := []struct {
		name         string
		urls         []string
		mockSetup    func(m *mocks.FSMock)
		expectErr    bool
		expectErrors []string
		overrideID   string
	}{
		{
			name: "single valid file",
			urls: []string{"http://example.com/file1.txt"},
			mockSetup: func(m *mocks.FSMock) {
				m.On("OpenURL", mock.Anything, "http://example.com/file1.txt").
					Return(io.NopCloser(bytes.NewReader([]byte("file 1 content"))), nil)
			},
			expectErr:    false,
			expectErrors: nil,
		},
		{
			name: "OpenURL returns error",
			urls: []string{"http://example.com/notfound.txt"},
			mockSetup: func(m *mocks.FSMock) {
				m.On("OpenURL", mock.Anything, "http://example.com/notfound.txt").
					Return(nil, errors.New("404 not found"))
			},
			expectErr:    false,
			expectErrors: []string{"failed to open URL"},
		},
		{
			name: "read error from URL",
			urls: []string{"http://example.com/broken.txt"},
			mockSetup: func(m *mocks.FSMock) {
				m.On("OpenURL", mock.Anything, "http://example.com/broken.txt").
					Return(&errorReader{}, nil)
			},
			expectErr:    false,
			expectErrors: []string{"failed to read data from"},
		},
		{
			name: "write file error",
			urls: []string{"http://example.com/file1.txt"},
	    	mockSetup: func(m *mocks.FSMock) {
        	m.On("OpenURL", mock.Anything, "http://example.com/file1.txt").
	                Return(io.NopCloser(bytes.NewReader([]byte("ok"))), nil)
	        },
	        overrideID:   "no_such_dir/task1",
	        expectErr:    true,
	        expectErrors: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := new(mocks.FSMock)
			tt.mockSetup(mockFS)

			handler := &Handler{
				FS: mockFS,
				Addr: &config.NetAddress{
					Host: "127.0.0.1",
					Port: 8080,
				},
			}

			id := "test-task-" + strings.ReplaceAll(tt.name, " ", "_")
			if tt.overrideID != "" {
				id = tt.overrideID
			}

			task := &models.Task{ID: id, URLs: tt.urls}

			archiveURL, err := handler.CreateArchive(context.Background(), task)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, archiveURL, "/archives/")
				assert.FileExists(t, "/tmp/"+task.ID+".zip")
				defer os.Remove("/tmp/" + task.ID + ".zip")
			}

			for _, expected := range tt.expectErrors {
				found := false
				for _, actual := range task.Errors {
					if strings.Contains(actual, expected) {
						found = true
						break
					}
				}
				assert.True(t, found, "expected error containing: %q", expected)
			}

			mockFS.AssertExpectations(t)
		})
	}
}


type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func (e *errorReader) Close() error {
	return nil
}
