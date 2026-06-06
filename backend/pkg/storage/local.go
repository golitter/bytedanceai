package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage stores files on the local filesystem and serves them via HTTP.
type LocalStorage struct {
	dir       string
	urlPrefix string
}

// NewLocalStorage creates a LocalStorage backed by dir, with URLs prefixed by urlPrefix.
// It ensures dir exists on disk.
func NewLocalStorage(dir, urlPrefix string) (*LocalStorage, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve storage dir: %w", err)
	}
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}
	return &LocalStorage{dir: absDir, urlPrefix: strings.TrimRight(urlPrefix, "/")}, nil
}

// Dir returns the on-disk directory (used by main.go to register Gin static route).
func (s *LocalStorage) Dir() string { return s.dir }

func (s *LocalStorage) UploadBytes(_ context.Context, key string, data []byte) (string, error) {
	if err := validateKey(key); err != nil {
		return "", err
	}
	fullPath := filepath.Join(s.dir, key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", fmt.Errorf("create directory: %w", err)
	}
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return s.publicURL(key), nil
}

func (s *LocalStorage) UploadReader(_ context.Context, key string, reader io.Reader, _ int64) (string, error) {
	if err := validateKey(key); err != nil {
		return "", err
	}
	fullPath := filepath.Join(s.dir, key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", fmt.Errorf("create directory: %w", err)
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return s.publicURL(key), nil
}

func (s *LocalStorage) publicURL(key string) string {
	return s.urlPrefix + "/" + key
}

func validateKey(key string) error {
	if strings.Contains(key, "..") {
		return fmt.Errorf("invalid key: path traversal")
	}
	return nil
}
