package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// HealthResponse defines the structure for the health check JSON response.
type HealthResponse struct {
	Status string `json:"status"`
}

// DataResponse defines the structure for the data endpoint JSON response.
type DataResponse struct {
	Message string `json:"message"`
	Source  string `json:"source"`
}

func callHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Log the received data
	slog.Info("Received POST request on /call", "body", string(body))

	// Send a response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "POST request processed successfully for endpoint /call")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond)
	response := HealthResponse{Status: "healthy"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	slog.Info("Health check requested", "endpoint", "/health", "status", "200 OK")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	response := DataResponse{Message: "This is sample JSON data", Source: "Cloud Foundry Go App"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	slog.Info("Data requested", "endpoint", "/data", "method", r.Method)
}

func main() {
	// Set up slog to use JSON format for better integration with log analysis tools
	// (e.g., Splunk, ELK stack used in CF environments).
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/call", callHandler)
	http.HandleFunc("/call/", callHandler)

	slog.Info("Starting server", "port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
