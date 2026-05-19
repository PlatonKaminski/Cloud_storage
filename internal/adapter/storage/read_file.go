package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud_storage/internal/domain"
)

func (s *LocalStorage) Read(_ context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)

	f, err := os.Open(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, domain.ErrFileNotFound
		}
		return nil, fmt.Errorf("open file: %w", err)
	}

	return f, nil
}
