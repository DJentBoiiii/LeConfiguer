package storage

import (
	"config-storage/internal/models"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// List returns all configurations from file storage.
func (f *FileStorage) List() ([]*models.Config, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	entries, err := os.ReadDir(f.dataDir)
	if err != nil {
		return nil, err
	}

	configs := make([]*models.Config, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(f.dataDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var config models.Config
		if err := json.Unmarshal(data, &config); err != nil {
			continue
		}

		configs = append(configs, &config)
	}

	return configs, nil
}
