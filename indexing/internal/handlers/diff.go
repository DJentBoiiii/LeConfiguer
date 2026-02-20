package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmezard/go-difflib/difflib"
)

type diffResponse struct {
	ConfigID string `json:"config_id"`
	FromID   uint   `json:"from_id"`
	ToID     uint   `json:"to_id"`
	Diff     string `json:"diff"`
}

func (h *Handler) Diff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	configID := vars["id"]
	if configID == "" {
		respondError(w, http.StatusBadRequest, errors.New("config id is required"))
		return
	}

	changes, err := h.latestChanges(configID, 2)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if len(changes) < 2 {
		respondError(w, http.StatusNotFound, errors.New("not enough versions to diff"))
		return
	}

	toChange := changes[0]
	fromChange := changes[1]

	fromContent := fromChange.Content
	toContent := toChange.Content

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(fromContent),
		B:        difflib.SplitLines(toContent),
		FromFile: fmt.Sprintf("version-%d", fromChange.ID),
		ToFile:   fmt.Sprintf("version-%d", toChange.ID),
		Context:  3,
	}

	text, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(text))
}
