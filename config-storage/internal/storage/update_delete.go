package storage

import "config-storage/internal/models"

// Update modifies an existing configuration.
func (m *MemoryStorage) Update(id string, config *models.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.configs[id]; !exists {
		return ErrNotFound
	}

	config.ID = id
	m.configs[id] = config
	return nil
}

// Delete removes a configuration from storage.
func (m *MemoryStorage) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.configs[id]; !exists {
		return ErrNotFound
	}

	delete(m.configs, id)
	return nil
}
