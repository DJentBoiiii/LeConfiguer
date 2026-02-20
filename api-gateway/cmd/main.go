package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	proxy "api-gateway/internal/proxy"

	"github.com/gorilla/mux"
)

func main() {
	storageURL := os.Getenv("CONFIG_STORAGE_URL")
	if storageURL == "" {
		storageURL = "http://localhost:8081"
	}

	indexingURL := os.Getenv("INDEXING_URL")
	if indexingURL == "" {
		indexingURL = "http://localhost:8082"
	}

	proxy := proxy.NewProxy(storageURL, indexingURL)

	r := mux.NewRouter()

	// Define routes
	proxyHandler := proxy.Forward
	r.HandleFunc("/configs", proxyHandler).Methods("POST")
	r.HandleFunc("/configs/{id}", proxyHandler).Methods("GET")
	r.HandleFunc("/configs/{id}", proxyHandler).Methods("PUT")
	r.HandleFunc("/configs/{id}", proxyHandler).Methods("DELETE")
	r.HandleFunc("/configs", proxyHandler).Methods("GET")
	r.HandleFunc("/configs/{id}/versions/{versionId}", proxyHandler).Methods("GET")
	r.HandleFunc("/configs/{id}/versions", proxyHandler).Methods("GET")
	r.HandleFunc("/configs/{id}/rollback", proxyHandler).Methods("POST")
	r.HandleFunc("/diff/{id}", proxyHandler).Methods("GET")

	// Start server
	fmt.Println("API Gateway running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
