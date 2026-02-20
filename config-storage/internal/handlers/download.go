package handlers

import (
	"config-storage/internal/storage"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// DownloadConfig handles GET /configs/{id}/download
// It downloads the configuration file with filename format: name.type
func (h *Handler) DownloadConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	config, reader, err := h.store.Download(id)
	if err != nil {
		if err == storage.ErrNotFound {
			respondError(w, http.StatusNotFound, err)
			return
		}
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	defer reader.Close()

	// Set headers for file download
	filename := fmt.Sprintf("%s.%s", config.Name, config.Type)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Stream the file content to the response
	_, err = io.Copy(w, reader)
	if err != nil {
		// If headers are already sent, we can't return an error
		// Log the error in production systems
		return
	}
}
