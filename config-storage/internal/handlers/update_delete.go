package handlers

import (
	"config-storage/internal/models"
	"config-storage/internal/storage"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// UpdateConfig handles PUT /configs/{id}
// It expects multipart/form-data with fields: name, type, environment, file.
func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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
		Name:        name + "." + configType,
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

	// Do not return file content in the response.
	config.JSONContent = nil
	respondJSON(w, http.StatusOK, config)
}

// DeleteConfig handles DELETE /configs/{id}
func (h *Handler) DeleteConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.store.Delete(id); err != nil {
		if err == storage.ErrNotFound {
			respondError(w, http.StatusNotFound, err)
			return
		}
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
