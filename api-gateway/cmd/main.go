package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"api-gateway/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	baseURL := os.Getenv("CONFIG_STORAGE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}

	proxy := storage.NewProxy(baseURL)

	r := mux.NewRouter()

	// Define routes
	proxyHandler := proxy.Forward
	r.HandleFunc("/configs", proxyHandler).Methods("POST")
	r.HandleFunc("/configs/{id}", proxyHandler).Methods("GET")
	r.HandleFunc("/configs/{id}", proxyHandler).Methods("PUT")
	r.HandleFunc("/configs/{id}", proxyHandler).Methods("DELETE")
	r.HandleFunc("/configs", proxyHandler).Methods("GET")

	// Start server
	fmt.Println("API Gateway running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
