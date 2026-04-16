package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/whitenhiemer/crochetbot/internal/mesh"
	"github.com/whitenhiemer/crochetbot/internal/models"
	"github.com/whitenhiemer/crochetbot/internal/pattern"
)

type GenerateRequest struct {
	FileID   string `json:"file_id"`
	Filename string `json:"filename"`
}

type GenerateResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Pattern *models.Pattern `json:"pattern,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// handleGenerateRequest processes pattern generation requests
func handleGenerateRequest(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest

	// Parse JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondGenerateError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Filename == "" {
		respondGenerateError(w, http.StatusBadRequest, "filename is required")
		return
	}

	// Get upload directory
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Build file path
	filePath := filepath.Join(uploadDir, req.Filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		respondGenerateError(w, http.StatusNotFound, "File not found")
		return
	}

	// Load mesh based on file extension
	ext := filepath.Ext(filePath)
	var m *mesh.Mesh
	var err error

	switch ext {
	case ".obj":
		m, err = mesh.LoadOBJ(filePath)
	case ".stl":
		m, err = mesh.LoadSTL(filePath)
	default:
		respondGenerateError(w, http.StatusBadRequest, "Unsupported file format")
		return
	}

	if err != nil {
		respondGenerateError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to load mesh: %v", err))
		return
	}

	// Reorient mesh so longest dimension is height (for optimal pattern generation)
	m.CalculateBounds()
	m.ReorientToLongestAxis()

	// Generate pattern
	gen := pattern.NewGenerator()
	pat, err := gen.Generate(m)
	if err != nil {
		respondGenerateError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate pattern: %v", err))
		return
	}

	// Store pattern
	if err := store.Save(pat); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to save pattern: %v\n", err)
	}

	// Success response
	response := GenerateResponse{
		Success: true,
		Message: fmt.Sprintf("Pattern generated successfully with %d part(s)", len(pat.Parts)),
		Pattern: pat,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// respondGenerateError sends error response for generation
func respondGenerateError(w http.ResponseWriter, statusCode int, message string) {
	response := GenerateResponse{
		Success: false,
		Error:   message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
