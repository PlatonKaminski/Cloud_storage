package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

func (s *LocalStorage) Delete(_ context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	err := os.Remove(fullPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}
