package handlers

import "config-storage/internal/storage"

// Handler aggregates all HTTP handlers for the config storage service.
type Handler struct {
	store storage.Storage
}

// NewHandler creates a new Handler with the given storage backend.
func NewHandler(store storage.Storage) *Handler {
	return &Handler{store: store}
}
