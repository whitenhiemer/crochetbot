package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewRouter creates and configures the HTTP router
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", handleHealth)

	// API routes
	mux.HandleFunc("/api/upload", handleUpload)
	mux.HandleFunc("/api/generate", handleGenerate)
	mux.HandleFunc("/api/pattern/", handleGetPattern)

	// CORS middleware wrapper
	return enableCORS(mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "crochetbot",
	})
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	handleUploadRequest(w, r)
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	handleGenerateRequest(w, r)
}

func handleGetPattern(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract pattern ID from URL path
	// URL format: /api/pattern/{id}
	pathPrefix := "/api/pattern/"
	if !strings.HasPrefix(r.URL.Path, pathPrefix) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	patternID := strings.TrimPrefix(r.URL.Path, pathPrefix)
	if patternID == "" {
		http.Error(w, "Pattern ID required", http.StatusBadRequest)
		return
	}

	// Retrieve pattern
	pattern, err := store.Get(patternID)
	if err != nil {
		http.Error(w, "Pattern not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pattern)
}

// enableCORS adds CORS headers to responses
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
