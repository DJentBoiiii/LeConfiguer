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
	r.HandleFunc("/echo", echoHandler).Methods("POST")
	// Start server
	fmt.Println("API Gateway running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Echo from API Gateway"}`))
}
