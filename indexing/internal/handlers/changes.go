package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"indexing/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type createChangeRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Environment string `json:"environment"`
	Action      string `json:"action"`
	Content     string `json:"content"`
}

var allowedActions = map[string]struct{}{
	"create": {},
	"update": {},
	"delete": {},
}

func (h *Handler) CreateChange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	configID := vars["id"]
	if configID == "" {
		respondError(w, http.StatusBadRequest, errors.New("config id is required"))
		return
	}

	var req createChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if _, ok := allowedActions[req.Action]; !ok {
		respondError(w, http.StatusBadRequest, errors.New("action must be create, update, or delete"))
		return
	}

	if req.Action != "delete" && req.Content == "" {
		respondError(w, http.StatusBadRequest, errors.New("content is required for create/update"))
		return
	}

	change := models.ConfigChange{
		ConfigID:    configID,
		Name:        req.Name,
		Type:        req.Type,
		Environment: req.Environment,
		Action:      req.Action,
		Content:     req.Content,
	}

	if err := h.db.Create(&change).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusCreated, change)
}

func (h *Handler) ListVersions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	configID := vars["id"]
	if configID == "" {
		respondError(w, http.StatusBadRequest, errors.New("config id is required"))
		return
	}

	var changes []models.ConfigChange
	if err := h.db.Where("config_id = ?", configID).
		Order("created_at desc").
		Find(&changes).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	for i := range changes {
		changes[i].Content = ""
	}

	respondJSON(w, http.StatusOK, changes)
}

func (h *Handler) GetVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	configID := vars["id"]
	if configID == "" {
		respondError(w, http.StatusBadRequest, errors.New("config id is required"))
		return
	}

	versionID, err := parseUint(vars["versionId"])
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	var change models.ConfigChange
	if err := h.db.Where("id = ? AND config_id = ?", versionID, configID).First(&change).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, err)
			return
		}
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusOK, change)
}

func parseUint(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}
