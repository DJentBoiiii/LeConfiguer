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
	// Initialize MinIO storage
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		bucketName = "configs"
	}
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	store, err := storage.NewMinIOStorage(endpoint, accessKey, secretKey, bucketName, useSSL)
	if err != nil {
		log.Fatal("Failed to initialize MinIO storage:", err)
	}
	log.Printf("Using MinIO storage at: %s (bucket: %s)\n", endpoint, bucketName)

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
