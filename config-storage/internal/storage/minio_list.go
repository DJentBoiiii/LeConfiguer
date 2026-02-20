package storage

import (
	"config-storage/internal/models"
	"encoding/json"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
)

// List returns all configurations from MinIO bucket.
func (m *MinIOStorage) List() ([]*models.Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	configs := make([]*models.Config, 0)

	for object := range m.client.ListObjects(m.ctx, m.bucketName, minio.ListObjectsOptions{}) {
		if object.Err != nil {
			continue
		}

		if !strings.HasSuffix(object.Key, ".json") {
			continue
		}

		objectReader, err := m.client.GetObject(m.ctx, m.bucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			continue
		}

		data, err := io.ReadAll(objectReader)
		objectReader.Close()
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
