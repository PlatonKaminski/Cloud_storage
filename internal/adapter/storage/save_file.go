package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (s *LocalStorage) Save(_ context.Context, path string, r io.Reader) error {
	fullPath := filepath.Join(s.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("create dirs: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err = io.Copy(f, r); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
