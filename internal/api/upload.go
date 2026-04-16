package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/whitenhiemer/crochetbot/internal/models"
)

const (
	maxUploadSize = 50 << 20 // 50 MB
	uploadFormKey = "file"
)

var allowedExtensions = map[string]bool{
	".obj": true,
	// ".stl": true, // TODO: Implement STL parser
}

type UploadResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	File    *models.Mesh3D `json:"file,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// handleUploadRequest processes file uploads
func handleUploadRequest(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with size limit
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		respondError(w, http.StatusBadRequest, "File too large or invalid form data")
		return
	}

	// Get file from form
	file, header, err := r.FormFile(uploadFormKey)
	if err != nil {
		respondError(w, http.StatusBadRequest, "No file uploaded or invalid form field")
		return
	}
	defer file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid file type. Allowed: %v", getAllowedExtensions()))
		return
	}

	// Validate file size
	if header.Size > maxUploadSize {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("File size exceeds maximum of %d MB", maxUploadSize>>20))
		return
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	safeFilename := sanitizeFilename(header.Filename)
	uniqueFilename := fmt.Sprintf("%d_%s", timestamp, safeFilename)

	// Get upload directory from env or use default
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Create destination file
	destPath := filepath.Join(uploadDir, uniqueFilename)
	destFile, err := os.Create(destPath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create destination file")
		return
	}
	defer destFile.Close()

	// Copy uploaded file to destination
	written, err := io.Copy(destFile, file)
	if err != nil {
		os.Remove(destPath) // Clean up on error
		respondError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Create mesh metadata
	mesh := &models.Mesh3D{
		ID:         fmt.Sprintf("mesh-%d", timestamp),
		Filename:   uniqueFilename,
		UploadedAt: time.Now(),
		Format:     strings.TrimPrefix(ext, "."),
	}

	// TODO: Parse file to get vertex/face counts
	// For now, just return basic info

	response := UploadResponse{
		Success: true,
		Message: fmt.Sprintf("File uploaded successfully: %s (%d bytes)", header.Filename, written),
		File:    mesh,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sanitizeFilename removes potentially dangerous characters
func sanitizeFilename(filename string) string {
	// Remove path separators
	filename = filepath.Base(filename)

	// Replace spaces and special characters
	replacer := strings.NewReplacer(
		" ", "_",
		"(", "",
		")", "",
		"[", "",
		"]", "",
		"{", "",
		"}", "",
		"&", "",
		"$", "",
		"!", "",
		"@", "",
		"#", "",
		"%", "",
		"^", "",
		"*", "",
		"+", "",
		"=", "",
		"|", "",
		"\\", "",
		"/", "",
		":", "",
		";", "",
		"'", "",
		"\"", "",
		"<", "",
		">", "",
		"?", "",
	)

	filename = replacer.Replace(filename)

	// Limit length
	if len(filename) > 200 {
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		filename = name[:200-len(ext)] + ext
	}

	return filename
}

// getAllowedExtensions returns list of allowed file extensions
func getAllowedExtensions() []string {
	exts := make([]string, 0, len(allowedExtensions))
	for ext := range allowedExtensions {
		exts = append(exts, ext)
	}
	return exts
}

// respondError sends error response
func respondError(w http.ResponseWriter, statusCode int, message string) {
	response := UploadResponse{
		Success: false,
		Error:   message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
