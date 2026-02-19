package storage

import "config-storage/internal/models"

// List returns all configurations in storage.
func (m *MemoryStorage) List() ([]*models.Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	configs := make([]*models.Config, 0, len(m.configs))
	for _, config := range m.configs {
		configs = append(configs, config)
	}

	return configs, nil
}
