package handlers

import (
	"net/http"
)

// ListConfigs handles GET requests to retrieve all configurations.
func (h *Handler) ListConfigs(w http.ResponseWriter, r *http.Request) {
	configs, err := h.storage.List()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, configs)
}
