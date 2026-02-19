package handlers

import (
	"config-storage/internal/models"
	"config-storage/internal/storage"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateConfig handles POST requests to create a new configuration.
func (h *Handler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	var config models.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate the configuration (ID is optional for creation)
	if err := config.ValidateForCreate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Create(&config); err != nil {
		if err == storage.ErrAlreadyExists {
			respondError(w, http.StatusConflict, "config already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, config)
}

// GetConfig handles GET requests to retrieve a configuration by ID.
func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	config, err := h.storage.Get(id)
	if err != nil {
		if err == storage.ErrNotFound {
			respondError(w, http.StatusNotFound, "config not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, config)
}
