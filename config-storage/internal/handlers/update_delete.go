package handlers

import (
	"config-storage/internal/models"
	"config-storage/internal/storage"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// UpdateConfig handles PUT requests to update a configuration.
func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var config models.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate the configuration (ID will be set from URL)
	config.ID = id
	if err := config.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Update(id, &config); err != nil {
		if err == storage.ErrNotFound {
			respondError(w, http.StatusNotFound, "config not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, config)
}

// DeleteConfig handles DELETE requests to remove a configuration.
func (h *Handler) DeleteConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.storage.Delete(id); err != nil {
		if err == storage.ErrNotFound {
			respondError(w, http.StatusNotFound, "config not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
