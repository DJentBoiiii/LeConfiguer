package handlers

import (
	"config-storage/internal/models"
	"config-storage/internal/storage"
	"io"
	"net/http"
)

func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	id, err := getConfigID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB
		respondError(w, http.StatusBadRequest, err)
		return
	}

	name := r.FormValue("name")
	configType := r.FormValue("type")
	environment := r.FormValue("environment")

	file, _, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	config := &models.Config{
		ID:          id,
		Name:        name,
		Type:        configType,
		Environment: environment,
		JSONContent: content,
	}

	if err := config.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.Update(id, config); err != nil {
		if err == storage.ErrNotFound {
			respondError(w, http.StatusNotFound, err)
			return
		}
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.sendIndexChange(r.Context(), config.ID, config.Name, config.Type, config.Environment, "update", string(content)); err != nil {
		respondError(w, http.StatusBadGateway, err)
		return
	}

	config.JSONContent = nil
	respondJSON(w, http.StatusOK, config)
}
