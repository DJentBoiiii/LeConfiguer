package handlers

import (
	"config-storage/internal/storage"
)

// Handler manages HTTP requests for configuration operations.
type Handler struct {
	storage storage.Storage
}

// NewHandler creates a new Handler instance.
func NewHandler(s storage.Storage) *Handler {
	return &Handler{storage: s}
}
