package handlers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) Diff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	configID := vars["id"]
	if configID == "" {
		respondError(w, http.StatusBadRequest, errors.New("config id is required"))
		return
	}

	changes, err := h.latestChanges(configID, 1)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if len(changes) == 0 {
		respondError(w, http.StatusNotFound, errors.New("no version found"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(changes[0].Content))
}
