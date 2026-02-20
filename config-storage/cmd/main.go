package main

import (
	"config-storage/internal/handlers"
	"config-storage/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize storage (file-based for persistence)
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	store, err := storage.NewFileStorage(dataDir)
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}
	log.Printf("Using file storage at: %s\n", dataDir)

	// Initialize handler
	handler := handlers.NewHandler(store)

	// Set up router
	router := mux.NewRouter()

	// Configure routes
	router.HandleFunc("/configs", handler.CreateConfig).Methods("POST")
	router.HandleFunc("/configs", handler.ListConfigs).Methods("GET")
	router.HandleFunc("/configs/{id}", handler.GetConfig).Methods("GET")
	router.HandleFunc("/configs/{id}", handler.UpdateConfig).Methods("PUT")
	router.HandleFunc("/configs/{id}", handler.DeleteConfig).Methods("DELETE")

	// Start server
	log.Println("Config Storage server starting on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
