package storage

import (
	"io"
	"strings"

	"config-storage/internal/models"

	"github.com/minio/minio-go/v7"
)

// Download retrieves the file content of a configuration from MinIO by ID.
func (s *MinIOStorage) Download(id string) (*models.Config, io.ReadCloser, error) {
	key, err := s.findObjectKeyByID(id)
	if err != nil {
		return nil, nil, err
	}

	obj, err := s.client.GetObject(s.ctx, s.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}

	info, err := s.client.StatObject(s.ctx, s.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		obj.Close()
		return nil, nil, err
	}

	parts := strings.Split(key, "/")
	environment := ""
	configType := ""
	if len(parts) >= 2 {
		environment = parts[0]
		configType = parts[1]
	}

	meta := info.UserMetadata
	configID := meta["id"]
	if configID == "" {
		configID = id
	}
	name := meta["name"]
	if name == "" {
		name = configID
	}
	if metaType := meta["type"]; metaType != "" {
		configType = metaType
	}
	if metaEnv := meta["environment"]; metaEnv != "" {
		environment = metaEnv
	}

	config := &models.Config{
		ID:          configID,
		Name:        name,
		Type:        configType,
		Environment: environment,
	}

	return config, obj, nil
}
