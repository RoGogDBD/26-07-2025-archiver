package mocks

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
	"github.com/viant/afs/storage"
)

type FSMock struct {
	mock.Mock
}

func (m *FSMock) OpenURL(ctx context.Context, url string, options ...storage.Option) (io.ReadCloser, error) {
	args := m.Called(ctx, url)
	if rc := args.Get(0); rc != nil {
		return rc.(io.ReadCloser), args.Error(1)
	}
	return nil, args.Error(1)
}
