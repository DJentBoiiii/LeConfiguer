package storage

import (
	"bytes"
	"errors"
	"path"
	"strings"

	"config-storage/internal/models"

	"github.com/minio/minio-go/v7"
)

// objectName builds the MinIO object key for a given config.
// Layout: environment/type/id
func objectName(config *models.Config) string {
	return path.Join(config.Environment, config.Type, config.ID)
}

// Create stores a new configuration file in MinIO.
// The file content is expected in config.JSONContent as []byte.
func (s *MinIOStorage) Create(config *models.Config) error {
	if config == nil {
		return errors.New("config is nil")
	}

	data, ok := config.JSONContent.([]byte)
	if !ok {
		return errors.New("config JSONContent must be []byte for MinIO storage")
	}

	key := objectName(config)

	// Check if object already exists
	_, err := s.client.StatObject(s.ctx, s.bucketName, key, minio.StatObjectOptions{})
	if err == nil {
		return ErrAlreadyExists
	}
	if resp := minio.ToErrorResponse(err); resp.Code != "NoSuchKey" && resp.Code != "NotFound" {
		return err
	}

	reader := bytes.NewReader(data)
	_, err = s.client.PutObject(s.ctx, s.bucketName, key, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
		UserMetadata: map[string]string{
			"id":          config.ID,
			"name":        config.Name,
			"type":        config.Type,
			"environment": config.Environment,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves a configuration and its file content from MinIO by ID.
func (s *MinIOStorage) Get(id string) (*models.Config, error) {
	key, err := s.findObjectKeyByID(id)
	if err != nil {
		return nil, err
	}

	obj, err := s.client.GetObject(s.ctx, s.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	info, err := s.client.StatObject(s.ctx, s.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
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

	return &models.Config{
		ID:          configID,
		Name:        name,
		Type:        configType,
		Environment: environment,
	}, nil
}
