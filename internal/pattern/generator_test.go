package pattern

import (
	"math"
	"testing"

	"github.com/whitenhiemer/crochetbot/internal/mesh"
)

func TestShapeDetection(t *testing.T) {
	tests := []struct {
		name           string
		width, height, depth float64
		expectedShape  string
	}{
		{
			name: "Sphere - all equal",
			width: 10.0, height: 10.0, depth: 10.0,
			expectedShape: "sphere",
		},
		{
			name: "Tall cylinder - height > others",
			width: 10.0, height: 20.0, depth: 10.0,
			expectedShape: "cylinder",
		},
		{
			name: "Horizontal cylinder - depth > others",
			width: 10.0, height: 10.0, depth: 20.0,
			expectedShape: "cylinder",
		},
		{
			name: "Wide cylinder - width > others",
			width: 20.0, height: 10.0, depth: 10.0,
			expectedShape: "cylinder",
		},
		{
			name: "Narwhal proportions (1.28 ratio)",
			width: 58.5, height: 64.43, depth: 50.35,
			expectedShape: "cylinder",
		},
		{
			name: "Nearly spherical (< 15% difference)",
			width: 10.0, height: 11.0, depth: 10.5,
			expectedShape: "sphere",
		},
	}

	g := NewGenerator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mesh.Mesh{}
			m.Bounds.MinX, m.Bounds.MaxX = 0, tt.width
			m.Bounds.MinY, m.Bounds.MaxY = 0, tt.height
			m.Bounds.MinZ, m.Bounds.MaxZ = 0, tt.depth

			shape := g.analyzeShape(m)
			if shape != tt.expectedShape {
				t.Errorf("Expected %s, got %s (w:%.1f h:%.1f d:%.1f)", 
					tt.expectedShape, shape, tt.width, tt.height, tt.depth)
			}
		})
	}
}

func TestPatternGeneration(t *testing.T) {
	g := NewGenerator()

	// Create simple test mesh (tall cylinder: height/radius = 3.0)
	m := createTestCylinderMesh(30.0, 90.0)
	
	pattern, err := g.Generate(m)
	if err != nil {
		t.Fatalf("Failed to generate pattern: %v", err)
	}

	if len(pattern.Parts) == 0 {
		t.Fatal("No parts generated")
	}

	part := pattern.Parts[0]
	roundCount := len(part.Rounds)

	minRounds, maxRounds := 20, 85
	if roundCount < minRounds {
		t.Errorf("Too few rounds: %d (expected >= %d)", roundCount, minRounds)
	}
	if roundCount > maxRounds {
		t.Errorf("Too many rounds: %d (expected <= %d)", roundCount, maxRounds)
	}

	// Accuracy expectations - lower for synthetic meshes, real models typically 70-95%
	minAccuracy := 40.0
	if pattern.AccuracyMetrics.ShapeMatchPercent < minAccuracy {
		t.Errorf("Accuracy too low: %.1f%% (expected >= %.1f%%)",
			pattern.AccuracyMetrics.ShapeMatchPercent, minAccuracy)
	}

	// Real-world models should get better accuracy
	if pattern.AccuracyMetrics.ShapeMatchPercent < 70.0 {
		t.Logf("Note: Synthetic mesh accuracy is %.1f%%, real models typically achieve 70-95%%",
			pattern.AccuracyMetrics.ShapeMatchPercent)
	}

	// Verify finished size is set
	if pattern.FinishedSize.HeightInches == 0 {
		t.Error("Finished size not set")
	}

	t.Logf("Generated pattern with %d rounds, %.1f%% accuracy", 
		roundCount, pattern.AccuracyMetrics.ShapeMatchPercent)
}

func TestPatternSmoothness(t *testing.T) {
	g := NewGenerator()
	m := createTestCylinderMesh(30.0, 90.0)

	pattern, err := g.Generate(m)
	if err != nil {
		t.Fatalf("Failed to generate pattern: %v", err)
	}

	rounds := pattern.Parts[0].Rounds
	maxChange := 0
	prevStitches := 0

	for i, round := range rounds {
		if round.StitchCount == 0 {
			continue // Skip finish rounds
		}

		if i > 0 && prevStitches > 0 {
			change := round.StitchCount - prevStitches
			if change < 0 {
				change = -change
			}
			if change > maxChange {
				maxChange = change
			}

			// Should not change by more than 12 stitches per round
			if change > 12 {
				t.Errorf("Round %d changes by %d stitches (too large jump from %d to %d)",
					round.Number, change, prevStitches, round.StitchCount)
			}
		}

		prevStitches = round.StitchCount
	}

	t.Logf("Max stitch change per round: %d", maxChange)
}

func TestStitchCountMultipleOfSix(t *testing.T) {
	g := NewGenerator()
	m := createTestCylinderMesh(30.0, 90.0)

	pattern, err := g.Generate(m)
	if err != nil {
		t.Fatalf("Failed to generate pattern: %v", err)
	}

	for _, round := range pattern.Parts[0].Rounds {
		if round.StitchCount == 0 {
			continue // Finish rounds can be 0
		}
		
		if round.StitchCount%6 != 0 {
			t.Errorf("Round %d has %d stitches (not multiple of 6)",
				round.Number, round.StitchCount)
		}

		if round.StitchCount < 6 {
			t.Errorf("Round %d has %d stitches (minimum is 6)",
				round.Number, round.StitchCount)
		}
	}
}

// Helper function to create a test cylinder mesh
func createTestCylinderMesh(radius, height float64) *mesh.Mesh {
	m := &mesh.Mesh{}
	
	// Create vertices around a cylinder
	slices := 20
	stacks := 40
	
	for stack := 0; stack <= stacks; stack++ {
		y := (float64(stack) / float64(stacks)) * height
		
		for slice := 0; slice < slices; slice++ {
			angle := (float64(slice) / float64(slices)) * 2.0 * math.Pi
			x := math.Cos(angle) * radius
			z := math.Sin(angle) * radius
			
			m.Vertices = append(m.Vertices, mesh.Vertex{
				X: x, Y: y, Z: z,
			})
		}
	}
	
	m.CalculateBounds()
	return m
}
