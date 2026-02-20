package storage

import (
	"bytes"
	"config-storage/internal/models"
	"encoding/json"

	"github.com/minio/minio-go/v7"
)

// Create adds a new configuration to MinIO bucket.
func (m *MinIOStorage) Create(config *models.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if object already exists
	_, err := m.client.StatObject(m.ctx, m.bucketName, config.ID+".json", minio.StatObjectOptions{})
	if err == nil {
		return ErrAlreadyExists
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	_, err = m.client.PutObject(
		m.ctx,
		m.bucketName,
		config.ID+".json",
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/json"},
	)

	return err
}

// Get retrieves a configuration by ID from MinIO bucket.
func (m *MinIOStorage) Get(id string) (*models.Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	object, err := m.client.GetObject(m.ctx, m.bucketName, id+".json", minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	// Read the object content
	var data []byte
	buffer := make([]byte, 1024)
	for {
		n, err := object.Read(buffer)
		if n > 0 {
			data = append(data, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
