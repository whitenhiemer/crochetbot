//go:build integration
// +build integration

package pattern

import (
	"os"
	"testing"

	"github.com/whitenhiemer/crochetbot/internal/mesh"
)

// Run with: go test ./internal/pattern/ -tags=integration -v

func TestRealSTLFiles(t *testing.T) {
	tests := []struct {
		name         string
		file         string
		expectedType string
		minAccuracy  float64
		minRounds    int
		maxRounds    int
	}{
		{
			name:         "Groundhog",
			file:         "~/Downloads/1groundhog-standing-alert.stl",
			expectedType: "cylinder",
			minAccuracy:  80.0,
			minRounds:    60,
			maxRounds:    85,
		},
		{
			name:         "Narwhal",
			file:         "~/Downloads/Narwhal.stl",
			expectedType: "cylinder",
			minAccuracy:  75.0,
			minRounds:    60,
			maxRounds:    85,
		},
	}

	g := NewGenerator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expand home directory
			filePath := os.ExpandEnv(tt.file)
			if filePath[0] == '~' {
				home, _ := os.UserHomeDir()
				filePath = home + filePath[1:]
			}

			// Check if file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Skipf("Test file not found: %s", filePath)
			}

			// Load STL
			m, err := mesh.LoadSTL(filePath)
			if err != nil {
				t.Fatalf("Failed to load STL: %v", err)
			}

			m.CalculateBounds()

			// Generate pattern
			pattern, err := g.Generate(m)
			if err != nil {
				t.Fatalf("Failed to generate pattern: %v", err)
			}

			// Verify pattern type
			if len(pattern.Parts) == 0 {
				t.Fatal("No parts generated")
			}

			part := pattern.Parts[0]
			if part.Type != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, part.Type)
			}

			// Verify round count
			roundCount := len(part.Rounds)
			if roundCount < tt.minRounds {
				t.Errorf("Too few rounds: %d (expected >= %d)", roundCount, tt.minRounds)
			}
			if roundCount > tt.maxRounds {
				t.Errorf("Too many rounds: %d (expected <= %d)", roundCount, tt.maxRounds)
			}

			// Verify accuracy
			accuracy := pattern.AccuracyMetrics.ShapeMatchPercent
			if accuracy < tt.minAccuracy {
				t.Errorf("Accuracy too low: %.1f%% (expected >= %.1f%%)", accuracy, tt.minAccuracy)
			}

			// Verify finished size
			if pattern.FinishedSize.HeightInches != 6.0 {
				t.Errorf("Expected 6\" height, got %.1f\"", pattern.FinishedSize.HeightInches)
			}

			// Log results
			t.Logf("✓ Generated %s pattern:", tt.name)
			t.Logf("  - Type: %s", part.Type)
			t.Logf("  - Rounds: %d", roundCount)
			t.Logf("  - Accuracy: %.1f%%", accuracy)
			t.Logf("  - Size: %.1f\" x %.1f\" (%.1f cm x %.1f cm)",
				pattern.FinishedSize.HeightInches,
				pattern.FinishedSize.WidthInches,
				pattern.FinishedSize.HeightCm,
				pattern.FinishedSize.WidthCm)
		})
	}
}
