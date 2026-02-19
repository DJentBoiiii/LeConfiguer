package main

import (
	"config-storage/internal/handlers"
	"config-storage/internal/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize storage
	store := storage.NewMemoryStorage()

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
	log.Println("Config Storage server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
