package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/whitenhiemer/crochetbot/internal/models"
)

func TestHandleParsePattern(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    ParseRequest
		expectedStatus int
		expectSuccess  bool
		expectError    bool
	}{
		{
			name: "valid pattern text",
			requestBody: ParseRequest{
				Text: `HEAD & BODY
Rnd 1. 6 sc in magic ring (6)
Rnd 2. 6 inc (12)
Rnd 3. [sc, inc] x 6 (18)`,
			},
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
			expectError:    false,
		},
		{
			name: "empty text",
			requestBody: ParseRequest{
				Text: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
			expectError:    true,
		},
		{
			name: "malformed pattern",
			requestBody: ParseRequest{
				Text: "invalid pattern text without proper format",
			},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/pattern/parse", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handleParsePattern(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			var response ParseResponse
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("expected success=%v, got %v", tt.expectSuccess, response.Success)
			}

			if tt.expectError && response.Error == "" {
				t.Error("expected error message but got none")
			}

			if !tt.expectError && response.Error != "" {
				t.Errorf("unexpected error: %s", response.Error)
			}

			if tt.expectSuccess && response.Pattern == nil {
				t.Error("expected pattern but got nil")
			}
		})
	}
}

func TestHandleFormatPattern(t *testing.T) {
	samplePattern := &models.Pattern{
		Name: "Test Pattern",
		Parts: []models.Part{
			{
				Name:  "Body",
				Type:  "sphere",
				Color: "main color",
				Rounds: []models.Round{
					{Number: 1, Instructions: "6 sc in magic ring", StitchCount: 6, StitchType: "sc"},
					{Number: 2, Instructions: "6 inc", StitchCount: 12, StitchType: "inc", Repeats: 6},
					{Number: 3, Instructions: "[sc, inc] x 6", StitchCount: 18, StitchType: "inc", Repeats: 6},
				},
			},
		},
		Materials: models.Materials{
			YarnWeight:  "worsted",
			HookSize:    "3.5mm",
			Colors:      []models.Color{{Name: "blue", Amount: 50}},
		},
	}

	tests := []struct {
		name           string
		requestBody    FormatRequest
		expectedStatus int
		expectSuccess  bool
	}{
		{
			name: "valid pattern",
			requestBody: FormatRequest{
				Pattern:     samplePattern,
				CompactMode: false,
			},
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name: "nil pattern",
			requestBody: FormatRequest{
				Pattern: nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/pattern/format", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handleFormatPattern(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			var response FormatResponse
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("expected success=%v, got %v", tt.expectSuccess, response.Success)
			}

			if tt.expectSuccess && response.Text == "" {
				t.Error("expected formatted text but got empty string")
			}

			if tt.expectSuccess && response.CompactText == "" {
				t.Error("expected compact text but got empty string")
			}
		})
	}
}

func TestHandleValidatePattern(t *testing.T) {
	validPattern := &models.Pattern{
		Parts: []models.Part{
			{
				Name: "Body",
				Type: "sphere",
				Rounds: []models.Round{
					{Number: 1, StitchCount: 6, StitchType: "sc"},
					{Number: 2, StitchCount: 12, StitchType: "inc"},
					{Number: 3, StitchCount: 18, StitchType: "inc"},
					{Number: 4, StitchCount: 24, StitchType: "inc"},
					{Number: 5, StitchCount: 24, StitchType: "sc"},
					{Number: 6, StitchCount: 18, StitchType: "dec"},
				},
			},
		},
	}

	invalidPattern := &models.Pattern{
		Parts: []models.Part{
			{
				Name: "Invalid",
				Rounds: []models.Round{
					{Number: 1, StitchCount: 6, StitchType: "sc"},
					{Number: 2, StitchCount: 100, StitchType: "inc"}, // Unrealistic jump
				},
			},
		},
	}

	tests := []struct {
		name              string
		requestBody       ValidateRequest
		expectedStatus    int
		expectSuccess     bool
		expectValidResult bool
	}{
		{
			name: "valid pattern",
			requestBody: ValidateRequest{
				Pattern: validPattern,
			},
			expectedStatus:    http.StatusOK,
			expectSuccess:     true,
			expectValidResult: true,
		},
		{
			name: "pattern with unrealistic changes (still valid, but with warnings)",
			requestBody: ValidateRequest{
				Pattern: invalidPattern,
			},
			expectedStatus:    http.StatusOK,
			expectSuccess:     true,
			expectValidResult: true, // IsValid is true, but will have issues/warnings
		},
		{
			name: "nil pattern",
			requestBody: ValidateRequest{
				Pattern: nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
		{
			name: "with reference pattern",
			requestBody: ValidateRequest{
				Pattern:          validPattern,
				ReferencePattern: validPattern,
			},
			expectedStatus:    http.StatusOK,
			expectSuccess:     true,
			expectValidResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/pattern/validate", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handleValidatePattern(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			var response ValidateResponse
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("expected success=%v, got %v", tt.expectSuccess, response.Success)
			}

			if tt.expectSuccess && response.ValidationResult == nil {
				t.Error("expected validation result but got nil")
			}

			if tt.expectSuccess && response.ValidationResult != nil {
				if response.ValidationResult.IsValid != tt.expectValidResult {
					t.Errorf("expected isValid=%v, got %v", tt.expectValidResult, response.ValidationResult.IsValid)
				}

				if response.ValidationResult.Score < 0 || response.ValidationResult.Score > 100 {
					t.Errorf("invalid score: %f (must be 0-100)", response.ValidationResult.Score)
				}
			}

			// Check comparison metrics if reference provided
			if tt.requestBody.ReferencePattern != nil && response.ComparisonMetrics == nil {
				t.Error("expected comparison metrics but got nil")
			}
		})
	}
}

func TestHandleComparePatterns(t *testing.T) {
	pattern1 := &models.Pattern{
		Parts: []models.Part{
			{
				Name: "Body",
				Rounds: []models.Round{
					{Number: 1, StitchCount: 6},
					{Number: 2, StitchCount: 12},
					{Number: 3, StitchCount: 18},
				},
			},
		},
	}

	pattern2 := &models.Pattern{
		Parts: []models.Part{
			{
				Name: "Body",
				Rounds: []models.Round{
					{Number: 1, StitchCount: 6},
					{Number: 2, StitchCount: 12},
					{Number: 3, StitchCount: 18},
					{Number: 4, StitchCount: 24},
				},
			},
		},
	}

	tests := []struct {
		name           string
		generated      *models.Pattern
		reference      *models.Pattern
		expectedStatus int
		expectSuccess  bool
	}{
		{
			name:           "compare two valid patterns",
			generated:      pattern1,
			reference:      pattern2,
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name:           "missing generated pattern",
			generated:      nil,
			reference:      pattern2,
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
		{
			name:           "missing reference pattern",
			generated:      pattern1,
			reference:      nil,
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]*models.Pattern{
				"generated": tt.generated,
				"reference": tt.reference,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/pattern/compare", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handleComparePatterns(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			var response map[string]any
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			success, ok := response["success"].(bool)
			if !ok {
				t.Fatal("response missing 'success' field")
			}

			if success != tt.expectSuccess {
				t.Errorf("expected success=%v, got %v", tt.expectSuccess, success)
			}

			if tt.expectSuccess {
				if _, ok := response["comparison"]; !ok {
					t.Error("expected comparison metrics but not found")
				}
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	endpoints := []struct {
		path    string
		handler http.HandlerFunc
	}{
		{"/api/pattern/parse", handleParsePattern},
		{"/api/pattern/format", handleFormatPattern},
		{"/api/pattern/validate", handleValidatePattern},
		{"/api/pattern/compare", handleComparePatterns},
	}

	for _, ep := range endpoints {
		t.Run(ep.path+" GET", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ep.path, nil)
			rec := httptest.NewRecorder()

			ep.handler(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
			}
		})
	}
}

func TestRoundTripParseFormat(t *testing.T) {
	// Test that we can parse text, then format it back
	originalText := `HEAD & BODY
With pink yarn.

Rnd 1. 6 sc in magic ring (6)
Rnd 2. 6 inc (12)
Rnd 3. [sc, inc] x 6 (18)
Rnd 4. [2 sc, inc] x 6 (24)
Rnds 5-7. 24 sc (24)
Rnd 8. [2 sc, dec] x 6 (18)`

	// Parse
	parseReq := ParseRequest{Text: originalText}
	parseBody, _ := json.Marshal(parseReq)
	parseHTTPReq := httptest.NewRequest(http.MethodPost, "/api/pattern/parse", bytes.NewReader(parseBody))
	parseRec := httptest.NewRecorder()

	handleParsePattern(parseRec, parseHTTPReq)

	if parseRec.Code != http.StatusOK {
		t.Fatalf("parse failed with status %d", parseRec.Code)
	}

	var parseResp ParseResponse
	if err := json.NewDecoder(parseRec.Body).Decode(&parseResp); err != nil {
		t.Fatalf("failed to decode parse response: %v", err)
	}

	if !parseResp.Success || parseResp.Pattern == nil {
		t.Fatal("parse unsuccessful or pattern nil")
	}

	// Format
	formatReq := FormatRequest{Pattern: parseResp.Pattern}
	formatBody, _ := json.Marshal(formatReq)
	formatHTTPReq := httptest.NewRequest(http.MethodPost, "/api/pattern/format", bytes.NewReader(formatBody))
	formatRec := httptest.NewRecorder()

	handleFormatPattern(formatRec, formatHTTPReq)

	if formatRec.Code != http.StatusOK {
		t.Fatalf("format failed with status %d", formatRec.Code)
	}

	var formatResp FormatResponse
	if err := json.NewDecoder(formatRec.Body).Decode(&formatResp); err != nil {
		t.Fatalf("failed to decode format response: %v", err)
	}

	if !formatResp.Success || formatResp.Text == "" {
		t.Fatal("format unsuccessful or text empty")
	}

	// Verify formatted text contains key elements
	if !bytes.Contains([]byte(formatResp.Text), []byte("HEAD & BODY")) {
		t.Error("formatted text missing part name")
	}

	if !bytes.Contains([]byte(formatResp.Text), []byte("Rnd 1.")) {
		t.Error("formatted text missing round 1")
	}
}
