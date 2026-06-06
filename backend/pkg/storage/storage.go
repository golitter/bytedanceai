package storage

import (
	"context"
	"io"
)

// Provider is the storage abstraction used by services to upload files.
type Provider interface {
	// UploadBytes uploads in-memory data and returns the public URL.
	UploadBytes(ctx context.Context, key string, data []byte) (string, error)
	// UploadReader uploads from a reader and returns the public URL.
	UploadReader(ctx context.Context, key string, reader io.Reader, size int64) (string, error)
}
