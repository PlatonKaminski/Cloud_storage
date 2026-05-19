package storage

import (
	"fmt"
	"os"

	"cloud_storage/internal/usecase"
)

// compile-time check
var _ usecase.FileStorage = (*LocalStorage)(nil)

type LocalStorage struct {
	basePath string
}

func New(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}
