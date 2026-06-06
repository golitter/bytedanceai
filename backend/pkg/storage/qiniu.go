package storage

import (
	"context"
	"io"

	"agenthub/backend/pkg/qiniu"
)

// QiniuStorage adapts the existing qiniu.Uploader to the Provider interface.
type QiniuStorage struct {
	uploader *qiniu.Uploader
}

// NewQiniuStorage wraps a qiniu.Uploader as a Provider.
func NewQiniuStorage(uploader *qiniu.Uploader) *QiniuStorage {
	return &QiniuStorage{uploader: uploader}
}

func (s *QiniuStorage) UploadBytes(ctx context.Context, key string, data []byte) (string, error) {
	return s.uploader.UploadBytes(ctx, key, data)
}

func (s *QiniuStorage) UploadReader(ctx context.Context, key string, reader io.Reader, size int64) (string, error) {
	return s.uploader.UploadReader(ctx, key, reader, size)
}
