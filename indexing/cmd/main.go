package main

import (
	"log"
	"net/http"
	"os"

	"indexing/config"
	"indexing/internal/database"
	"indexing/internal/handlers"
	"indexing/internal/models"
	"indexing/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	db, err := database.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatalf("database connect error: %v", err)
	}

	if err := db.AutoMigrate(&models.ConfigChange{}); err != nil {
		log.Fatalf("database migrate error: %v", err)
	}

	storageURL := os.Getenv("STORAGE_URL")
	if storageURL == "" {
		storageURL = "http://localhost:8081"
	}
	storageClient := storage.NewClient(storageURL)

	h := handlers.New(db, storageClient)

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	r.HandleFunc("/configs/{id}/versions/{versionId}", h.GetVersion).Methods(http.MethodGet)
	r.HandleFunc("/configs/{id}/versions", h.ListVersions).Methods(http.MethodGet)
	r.HandleFunc("/configs/{id}/changes", h.CreateChange).Methods(http.MethodPost)
	r.HandleFunc("/configs/{id}/rollback", h.Rollback).Methods(http.MethodPost)
	r.HandleFunc("/diff/{id}", h.Diff).Methods(http.MethodGet)

	log.Printf("indexing service listening on %s", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
