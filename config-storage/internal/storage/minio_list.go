package storage

import (
	"path"
	"strings"

	"config-storage/internal/models"

	"github.com/minio/minio-go/v7"
)

// List returns metadata for all stored configurations.
// File contents are not loaded to keep the response lightweight.
func (s *MinIOStorage) List() ([]*models.Config, error) {
	var configs []*models.Config

	objectCh := s.client.ListObjects(s.ctx, s.bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		info, err := s.client.StatObject(s.ctx, s.bucketName, object.Key, minio.StatObjectOptions{})
		if err != nil {
			return nil, err
		}

		parts := strings.Split(object.Key, "/")
		environment := ""
		configType := ""
		if len(parts) >= 2 {
			environment = parts[0]
			configType = parts[1]
		}

		meta := info.UserMetadata
		id := meta["id"]
		if id == "" {
			id = path.Base(object.Key)
		}
		name := meta["name"]
		if name == "" {
			name = id
		}
		if metaType := meta["type"]; metaType != "" {
			configType = metaType
		}
		if metaEnv := meta["environment"]; metaEnv != "" {
			environment = metaEnv
		}

		configs = append(configs, &models.Config{
			ID:          id,
			Name:        name,
			Type:        configType,
			Environment: environment,
		})
	}

	return configs, nil
}
