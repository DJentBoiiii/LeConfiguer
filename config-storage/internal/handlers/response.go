package handlers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, err error) {
	if err == nil {
		respondJSON(w, status, errorResponse{Error: http.StatusText(status)})
		return
	}
	respondJSON(w, status, errorResponse{Error: err.Error()})
}
