package app

import (
	"context"
	"io"

	"github.com/viant/afs/storage"
)

type FSFetcher interface {
	OpenURL(ctx context.Context, URL string, options ...storage.Option) (io.ReadCloser, error)
}
