package main

import (
	"config-storage/internal/handlers"
	"config-storage/internal/indexing"
	"config-storage/internal/storage"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize storage backend (MinIO is the default)
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "minio"
	}

	var store storage.Storage

	switch storageType {
	case "minio":
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
		bucket := os.Getenv("MINIO_BUCKET")
		if bucket == "" {
			bucket = "configs"
		}

		useSSL := false
		if v := os.Getenv("MINIO_USE_SSL"); v != "" {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				log.Fatalf("invalid MINIO_USE_SSL value %q: %v", v, err)
			}
			useSSL = parsed
		}

		if endpoint == "" || accessKey == "" || secretKey == "" {
			log.Fatal("MINIO_ENDPOINT, MINIO_ACCESS_KEY and MINIO_SECRET_KEY must be set for MinIO storage")
		}

		minioStore, err := storage.NewMinIOStorage(endpoint, accessKey, secretKey, bucket, useSSL)
		if err != nil {
			log.Fatalf("failed to initialize MinIO storage: %v", err)
		}
		store = minioStore

	default:
		log.Fatalf("unsupported STORAGE_TYPE: %s", storageType)
	}

	indexingURL := os.Getenv("INDEXING_URL")
	if indexingURL == "" {
		indexingURL = "http://localhost:8082"
	}
	indexer := indexing.NewClient(indexingURL)

	// Initialize handler
	handler := handlers.NewHandler(store, indexer)

	// Set up router
	router := mux.NewRouter()

	// Configure routes
	router.HandleFunc("/configs", handler.UploadConfig).Methods("POST")
	router.HandleFunc("/configs", handler.ListConfigs).Methods("GET")
	router.HandleFunc("/configs/{id}", handler.GetConfig).Methods("GET")
	router.HandleFunc("/configs/{id}", handler.UpdateConfig).Methods("PUT")
	router.HandleFunc("/configs/{id}", handler.DeleteConfig).Methods("DELETE")
	router.HandleFunc("/configs/download/{id}", handler.DownloadConfig).Methods("GET")

	// Start server
	log.Println("Config Storage server starting on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
