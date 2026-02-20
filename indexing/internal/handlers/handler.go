package handlers

import (
	"indexing/internal/storage"

	"gorm.io/gorm"
)

type Handler struct {
	db            *gorm.DB
	storageClient *storage.Client
}

func New(db *gorm.DB, storageClient *storage.Client) *Handler {
	return &Handler{
		db:            db,
		storageClient: storageClient,
	}
}
