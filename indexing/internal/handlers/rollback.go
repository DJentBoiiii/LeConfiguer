package handlers

import (
	"errors"
	"net/http"

	"indexing/internal/models"

	"github.com/gorilla/mux"
)

func (h *Handler) Rollback(w http.ResponseWriter, r *http.Request) {
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
	if len(changes) == 0 {
		respondError(w, http.StatusNotFound, errors.New("no versions found"))
		return
	}

	latestChange := changes[0]

	if err := h.db.Delete(&latestChange).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if len(changes) > 1 && h.storageClient != nil {
		previousChange := changes[1]

		if previousChange.Content != "" {
			err := h.storageClient.UpdateConfig(
				r.Context(),
				configID,
				previousChange.Name,
				previousChange.Type,
				previousChange.Environment,
				previousChange.Content,
			)
			if err != nil {
				respondError(w, http.StatusBadGateway, err)
				return
			}
		}
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message":   "rollback successful",
		"config_id": configID,
	})
}

func (h *Handler) latestChanges(configID string, limit int) ([]models.ConfigChange, error) {
	var changes []models.ConfigChange
	if err := h.db.Where("config_id = ?", configID).
		Order("created_at desc").
		Limit(limit).
		Find(&changes).Error; err != nil {
		return nil, err
	}
	return changes, nil
}
