package mocks

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyReadCloser struct {
	io.Reader
}

func (d dummyReadCloser) Close() error {
	return nil
}

func TestFSMock_OpenURL(t *testing.T) {
	ctx := context.Background()
	fsMock := &FSMock{}

	tests := []struct {
		name       string
		url        string
		mockReturn io.ReadCloser
		mockError  error
		wantErr    bool
		wantData   string
	}{
		{
			name:       "Success returns reader without error",
			url:        "http://example.com/file.txt",
			mockReturn: dummyReadCloser{strings.NewReader("file content")},
			mockError:  nil,
			wantErr:    false,
			wantData:   "file content",
		},
		{
			name:       "Returns error and nil reader",
			url:        "http://example.com/badurl",
			mockReturn: nil,
			mockError:  errors.New("failed to open URL"),
			wantErr:    true,
			wantData:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsMock.On("OpenURL", ctx, tt.url).Return(tt.mockReturn, tt.mockError)

			rc, err := fsMock.OpenURL(ctx, tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rc)
				buf := new(strings.Builder)
				_, readErr := io.Copy(buf, rc)
				assert.NoError(t, readErr)
				assert.Equal(t, tt.wantData, buf.String())
				rc.Close()
			}

			fsMock.AssertCalled(t, "OpenURL", ctx, tt.url)
			fsMock.ExpectedCalls = nil
		})
	}
}
