package handlers

import (
	"config-storage/internal/storage"
	"net/http"
)

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	id, err := getConfigID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

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
