package handlers

import (
	"net/http"
)

// ListConfigs handles GET /configs
func (h *Handler) ListConfigs(w http.ResponseWriter, r *http.Request) {
	configs, err := h.store.List()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusOK, configs)
}
