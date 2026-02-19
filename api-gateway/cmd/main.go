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
	r.HandleFunc("/configurations", createConfiguration(proxy)).Methods("POST")
	r.HandleFunc("/configurations/{id}", readConfiguration(proxy)).Methods("GET")
	r.HandleFunc("/configurations/{id}", updateConfiguration(proxy)).Methods("PUT")
	r.HandleFunc("/configurations/{id}", deleteConfiguration(proxy)).Methods("DELETE")
	r.HandleFunc("/configurations", filterConfigurations(proxy)).Methods("GET")
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

func createConfiguration(proxy *storage.Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.Forward(w, r)
	}
}

func readConfiguration(proxy *storage.Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.Forward(w, r)
	}
}

func updateConfiguration(proxy *storage.Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.Forward(w, r)
	}
}

func deleteConfiguration(proxy *storage.Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.Forward(w, r)
	}
}

func filterConfigurations(proxy *storage.Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.Forward(w, r)
	}
}
