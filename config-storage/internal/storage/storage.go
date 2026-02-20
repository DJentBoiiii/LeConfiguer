package storage

import (
	"config-storage/internal/models"
	"errors"
	"io"
)

var (
	ErrNotFound      = errors.New("config not found")
	ErrAlreadyExists = errors.New("config already exists")
)

// Storage defines the interface for configuration storage operations.
type Storage interface {
	Create(config *models.Config) error
	Get(id string) (*models.Config, error)
	Update(id string, config *models.Config) error
	Delete(id string) error
	List() ([]*models.Config, error)
	Download(id string) (*models.Config, io.ReadCloser, error)
}
