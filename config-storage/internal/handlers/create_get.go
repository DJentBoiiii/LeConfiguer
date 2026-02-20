package handlers

import (
	"config-storage/internal/models"
	"config-storage/internal/storage"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// UploadConfig handles POST /configs
// It expects multipart/form-data with fields: id (optional), name, type, environment, file.
func (h *Handler) UploadConfig(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB
		respondError(w, http.StatusBadRequest, err)
		return
	}

	id := r.FormValue("id")
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

	if id == "" {
		id = uuid.New().String()
	}

	config := &models.Config{
		ID:          id,
		Name:        name,
		Type:        configType,
		Environment: environment,
		JSONContent: content,
	}

	if err := config.ValidateForCreate(); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.Create(config); err != nil {
		if err == storage.ErrAlreadyExists {
			respondError(w, http.StatusConflict, err)
			return
		}
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Do not return file content in the response.
	config.JSONContent = nil
	respondJSON(w, http.StatusCreated, config)
}

// GetConfig handles GET /configs/{id}
func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	config, err := h.store.Get(id)
	if err != nil {
		if err == storage.ErrNotFound {
			respondError(w, http.StatusNotFound, err)
			return
		}
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusOK, config)
}
