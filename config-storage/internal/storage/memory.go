package storage

import (
	"config-storage/internal/models"
	"sync"
)

// MemoryStorage is an in-memory implementation of the Storage interface.
type MemoryStorage struct {
	mu      sync.RWMutex
	configs map[string]*models.Config
}
