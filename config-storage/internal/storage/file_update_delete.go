package storage

import (
	"config-storage/internal/models"
	"encoding/json"
	"os"
	"path/filepath"
)

// Update modifies an existing configuration in file storage.
func (f *FileStorage) Update(id string, config *models.Config) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := filepath.Join(f.dataDir, id+".json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ErrNotFound
	}

	config.ID = id
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// Delete removes a configuration from file storage.
func (f *FileStorage) Delete(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := filepath.Join(f.dataDir, id+".json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ErrNotFound
	}

	return os.Remove(filePath)
}
