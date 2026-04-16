package api

import (
	"encoding/json"
	"net/http"

	"github.com/whitenhiemer/crochetbot/internal/models"
	"github.com/whitenhiemer/crochetbot/internal/pattern"
)

// ParseRequest represents pattern text to parse
type ParseRequest struct {
	Text string `json:"text"`
}

// ParseResponse returns parsed pattern
type ParseResponse struct {
	Success bool            `json:"success"`
	Pattern *models.Pattern `json:"pattern,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// FormatRequest represents pattern to format
type FormatRequest struct {
	Pattern     *models.Pattern `json:"pattern"`
	CompactMode bool            `json:"compact_mode"`
}

// FormatResponse returns formatted text
type FormatResponse struct {
	Success       bool   `json:"success"`
	Text          string `json:"text,omitempty"`
	CompactText   string `json:"compact_text,omitempty"`
	Error         string `json:"error,omitempty"`
}

// ValidateRequest represents pattern to validate
type ValidateRequest struct {
	Pattern         *models.Pattern `json:"pattern"`
	ReferencePattern *models.Pattern `json:"reference_pattern,omitempty"`
}

// ValidateResponse returns validation results
type ValidateResponse struct {
	Success            bool                        `json:"success"`
	ValidationResult   *pattern.ValidationResult   `json:"validation_result,omitempty"`
	ComparisonMetrics  *pattern.ComparisonMetrics  `json:"comparison_metrics,omitempty"`
	Error              string                      `json:"error,omitempty"`
}

// handleParsePattern parses text pattern
func handleParsePattern(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ParseResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if req.Text == "" {
		respondJSON(w, http.StatusBadRequest, ParseResponse{
			Success: false,
			Error:   "text is required",
		})
		return
	}

	// Parse pattern
	parser := pattern.NewParser()
	pat, err := parser.ParsePattern(req.Text)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ParseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate parsed pattern
	issues := parser.ValidatePattern(pat)
	if len(issues) > 0 {
		// Still return pattern but include warnings
		pat.Description = "Parsed with warnings: " + issues[0]
	}

	respondJSON(w, http.StatusOK, ParseResponse{
		Success: true,
		Pattern: pat,
	})
}

// handleFormatPattern formats pattern to text
func handleFormatPattern(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FormatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, FormatResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if req.Pattern == nil {
		respondJSON(w, http.StatusBadRequest, FormatResponse{
			Success: false,
			Error:   "pattern is required",
		})
		return
	}

	// Format pattern
	formatter := pattern.NewFormatter()
	text := formatter.FormatPattern(req.Pattern)
	compact := formatter.FormatCompact(req.Pattern)

	respondJSON(w, http.StatusOK, FormatResponse{
		Success:     true,
		Text:        text,
		CompactText: compact,
	})
}

// handleValidatePattern validates a pattern
func handleValidatePattern(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ValidateResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if req.Pattern == nil {
		respondJSON(w, http.StatusBadRequest, ValidateResponse{
			Success: false,
			Error:   "pattern is required",
		})
		return
	}

	// Validate pattern
	validator := pattern.NewValidator()
	result := validator.ValidatePattern(req.Pattern)

	response := ValidateResponse{
		Success:          true,
		ValidationResult: &result,
	}

	// Compare to reference if provided
	if req.ReferencePattern != nil {
		comparison := validator.CompareToReference(req.Pattern, req.ReferencePattern)
		response.ComparisonMetrics = &comparison
	}

	respondJSON(w, http.StatusOK, response)
}

// handleComparePatterns compares two patterns
func handleComparePatterns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Generated *models.Pattern `json:"generated"`
		Reference *models.Pattern `json:"reference"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if req.Generated == nil || req.Reference == nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "both generated and reference patterns are required",
		})
		return
	}

	validator := pattern.NewValidator()
	comparison := validator.CompareToReference(req.Generated, req.Reference)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"comparison": comparison,
	})
}

// respondJSON writes JSON response
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
