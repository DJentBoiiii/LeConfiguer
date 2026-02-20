package storage

import (
	"bytes"
	"errors"
	"path"

	"config-storage/internal/models"

	"github.com/minio/minio-go/v7"
)

// findObjectKeyByID searches for the MinIO object key whose basename matches the given ID.
func (s *MinIOStorage) findObjectKeyByID(id string) (string, error) {
	objectCh := s.client.ListObjects(s.ctx, s.bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return "", object.Err
		}
		if path.Base(object.Key) == id {
			return object.Key, nil
		}
	}

	return "", ErrNotFound
}

// Update replaces an existing configuration's file and metadata.
func (s *MinIOStorage) Update(id string, config *models.Config) error {
	if config == nil {
		return errors.New("config is nil")
	}
	if config.ID == "" {
		config.ID = id
	}

	// Ensure the object exists
	oldKey, err := s.findObjectKeyByID(id)
	if err != nil {
		return err
	}

	data, ok := config.JSONContent.([]byte)
	if !ok {
		return errors.New("config JSONContent must be []byte for MinIO storage")
	}

	// Remove old object (in case environment/type changed)
	if err := s.client.RemoveObject(s.ctx, s.bucketName, oldKey, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	newKey := objectName(config)
	reader := bytes.NewReader(data)
	_, err = s.client.PutObject(s.ctx, s.bucketName, newKey, reader, int64(len(data)), minio.PutObjectOptions{
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

// Delete removes a configuration by ID.
func (s *MinIOStorage) Delete(id string) error {
	key, err := s.findObjectKeyByID(id)
	if err != nil {
		return err
	}

	if err := s.client.RemoveObject(s.ctx, s.bucketName, key, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}
