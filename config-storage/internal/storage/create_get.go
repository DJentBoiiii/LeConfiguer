package storage

import "config-storage/internal/models"

// Create adds a new configuration to the storage.
func (m *MemoryStorage) Create(config *models.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.configs[config.ID]; exists {
		return ErrAlreadyExists
	}

	m.configs[config.ID] = config
	return nil
}

// Get retrieves a configuration by ID.
func (m *MemoryStorage) Get(id string) (*models.Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, exists := m.configs[id]
	if !exists {
		return nil, ErrNotFound
	}

	return config, nil
}
