package handlers

import (
	"context"
	"errors"
	"net/http"

	"config-storage/internal/indexing"

	"github.com/gorilla/mux"
)

func getConfigID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return "", errors.New("config id is required")
	}
	return id, nil
}

func (h *Handler) sendIndexChange(ctx context.Context, id, name, configType, environment, action, content string) error {
	if h.indexer == nil {
		return nil
	}

	return h.indexer.SendChange(ctx, id, indexing.ChangeRequest{
		Name:        name,
		Type:        configType,
		Environment: environment,
		Action:      action,
		Content:     content,
	})
}
