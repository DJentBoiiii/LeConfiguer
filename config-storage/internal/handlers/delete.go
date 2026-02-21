package handlers

import (
	"config-storage/internal/storage"
	"net/http"
)

func (h *Handler) DeleteConfig(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.Delete(id); err != nil {
		if err == storage.ErrNotFound {
			respondError(w, http.StatusNotFound, err)
			return
		}
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.sendIndexChange(r.Context(), id, config.Name, config.Type, config.Environment, "delete", ""); err != nil {
		respondError(w, http.StatusBadGateway, err)
		return
	}

	if h.indexer != nil {
		if err := h.indexer.DeleteConfig(r.Context(), id); err != nil {
			respondError(w, http.StatusBadGateway, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
