package storage

import (
	"config-storage/internal/models"
	"encoding/json"
	"os"
	"path/filepath"
)

// Create adds a new configuration to file storage.
func (f *FileStorage) Create(config *models.Config) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := filepath.Join(f.dataDir, config.ID+".json")

	if _, err := os.Stat(filePath); err == nil {
		return ErrAlreadyExists
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// Get retrieves a configuration by ID from file storage.
func (f *FileStorage) Get(id string) (*models.Config, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	filePath := filepath.Join(f.dataDir, id+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
