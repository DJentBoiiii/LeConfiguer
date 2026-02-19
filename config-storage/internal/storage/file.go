package storage

import (
	"os"
	"sync"
)

// FileStorage is a file-based implementation of the Storage interface.
type FileStorage struct {
	mu      sync.RWMutex
	dataDir string
}

// NewFileStorage creates a new file-based storage instance.
func NewFileStorage(dataDir string) (*FileStorage, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	return &FileStorage{
		dataDir: dataDir,
	}, nil
}
