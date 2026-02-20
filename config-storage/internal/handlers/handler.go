package handlers

import (
	"config-storage/internal/indexing"
	"config-storage/internal/storage"
)

type Handler struct {
	store   storage.Storage
	indexer *indexing.Client
}

func NewHandler(store storage.Storage, indexer *indexing.Client) *Handler {
	return &Handler{store: store, indexer: indexer}
}
