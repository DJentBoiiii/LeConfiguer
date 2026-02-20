package storage

import (
	"bytes"
	"config-storage/internal/models"
	"encoding/json"

	"github.com/minio/minio-go/v7"
)

// Update modifies an existing configuration in MinIO bucket.
func (m *MinIOStorage) Update(id string, config *models.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if object exists
	_, err := m.client.StatObject(m.ctx, m.bucketName, id+".json", minio.StatObjectOptions{})
	if err != nil {
		return ErrNotFound
	}

	config.ID = id
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	_, err = m.client.PutObject(
		m.ctx,
		m.bucketName,
		id+".json",
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/json"},
	)

	return err
}

// Delete removes a configuration from MinIO bucket.
func (m *MinIOStorage) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if object exists
	_, err := m.client.StatObject(m.ctx, m.bucketName, id+".json", minio.StatObjectOptions{})
	if err != nil {
		return ErrNotFound
	}

	err = m.client.RemoveObject(m.ctx, m.bucketName, id+".json", minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
