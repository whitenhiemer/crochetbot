package pattern

import (
	"fmt"
	"time"

	"github.com/whitenhiemer/crochetbot/internal/mesh"
	"github.com/whitenhiemer/crochetbot/internal/models"
)

// Generator creates crochet patterns from 3D meshes
type Generator struct {
	// Configuration options
	DefaultYarnWeight string
	DefaultHookSize   string
}

// NewGenerator creates a new pattern generator with defaults
func NewGenerator() *Generator {
	return &Generator{
		DefaultYarnWeight: "worsted",
		DefaultHookSize:   "3.5mm",
	}
}

// Generate creates a crochet pattern from a mesh
func (g *Generator) Generate(m *mesh.Mesh) (*models.Pattern, error) {
	if len(m.Vertices) == 0 {
		return nil, fmt.Errorf("mesh has no vertices")
	}

	// Calculate bounds if not already done
	if m.Bounds.MaxX == 0 {
		m.CalculateBounds()
	}

	// Analyze shape type
	shapeType := g.analyzeShape(m)

	// Generate pattern based on shape
	pattern := &models.Pattern{
		ID:          generateID(),
		Name:        fmt.Sprintf("Generated %s Pattern", shapeType),
		CreatedAt:   time.Now(),
		Description: fmt.Sprintf("Auto-generated pattern from 3D model"),
		Difficulty:  "beginner",
		Parts:       []models.Part{},
		Materials:   g.generateMaterials(),
		Assembly:    []string{},
	}

	// Generate parts based on shape type
	switch shapeType {
	case "sphere":
		part := g.generateSpherePart(m)
		pattern.Parts = append(pattern.Parts, part)
	default:
		return nil, fmt.Errorf("unsupported shape type: %s", shapeType)
	}

	return pattern, nil
}

// analyzeShape determines the basic shape type of the mesh
func (g *Generator) analyzeShape(m *mesh.Mesh) string {
	if m.IsApproximatelySphere() {
		return "sphere"
	}
	if m.IsApproximatelyCylinder() {
		return "cylinder"
	}
	// Default to sphere for now
	return "sphere"
}

// generateSpherePart creates rounds for a spherical shape
func (g *Generator) generateSpherePart(m *mesh.Mesh) models.Part {
	// Calculate approximate diameter in stitches
	// Rough estimate: 1 unit = 5 stitches for worsted weight
	avgRadius := m.GetAverageRadius()
	maxStitches := int(avgRadius * 2 * 5) // diameter * stitches per unit

	// Ensure reasonable stitch count
	if maxStitches < 12 {
		maxStitches = 12
	}
	if maxStitches > 72 {
		maxStitches = 72
	}

	// Round to multiple of 6 for clean increases
	maxStitches = ((maxStitches + 5) / 6) * 6

	rounds := []models.Round{}

	// Starting magic ring
	rounds = append(rounds, models.Round{
		Number:       1,
		Instructions: "6 sc in magic ring",
		StitchCount:  6,
		StitchType:   "sc",
		Repeats:      1,
		Notes:        "Pull tight to close",
	})

	// Increase rounds (until max width)
	currentStitches := 6
	roundNum := 2
	increaseRounds := []models.Round{}

	for currentStitches < maxStitches {
		nextStitches := currentStitches + 6
		if nextStitches > maxStitches {
			nextStitches = maxStitches
		}

		// Calculate sc between increases
		scBetween := (currentStitches / 6) - 1
		var instruction string
		if scBetween == 0 {
			instruction = "2 sc in each st around"
		} else {
			instruction = fmt.Sprintf("[inc, %d sc] repeat 6 times", scBetween)
		}

		increaseRounds = append(increaseRounds, models.Round{
			Number:       roundNum,
			Instructions: instruction,
			StitchCount:  nextStitches,
			StitchType:   "sc",
			Repeats:      6,
			Notes:        "",
		})

		currentStitches = nextStitches
		roundNum++
	}
	rounds = append(rounds, increaseRounds...)

	// Constant rounds (equator) - about 1/3 of total rounds
	numConstantRounds := len(increaseRounds) / 3
	if numConstantRounds < 2 {
		numConstantRounds = 2
	}

	for i := 0; i < numConstantRounds; i++ {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("sc in each st around (%d sc)", currentStitches),
			StitchCount:  currentStitches,
			StitchType:   "sc",
			Repeats:      1,
			Notes:        "",
		})
		roundNum++
	}

	// Decrease rounds (symmetrical to increases)
	decreaseRounds := []models.Round{}
	stuffNote := ""

	for currentStitches > 6 {
		nextStitches := currentStitches - 6
		if nextStitches < 6 {
			nextStitches = 6
		}

		scBetween := (currentStitches / 6) - 2
		var instruction string
		if scBetween <= 0 {
			instruction = "dec around"
		} else {
			instruction = fmt.Sprintf("[dec, %d sc] repeat 6 times", scBetween)
		}

		// Add stuffing note on first decrease
		if len(decreaseRounds) == 0 {
			stuffNote = "Begin stuffing firmly"
		} else if currentStitches <= 18 {
			stuffNote = "Finish stuffing"
		} else {
			stuffNote = ""
		}

		decreaseRounds = append(decreaseRounds, models.Round{
			Number:       roundNum,
			Instructions: instruction,
			StitchCount:  nextStitches,
			StitchType:   "sc",
			Repeats:      6,
			Notes:        stuffNote,
		})

		currentStitches = nextStitches
		roundNum++
	}
	rounds = append(rounds, decreaseRounds...)

	// Final close
	rounds = append(rounds, models.Round{
		Number:       roundNum,
		Instructions: "Fasten off, leaving long tail. Close opening with yarn needle.",
		StitchCount:  0,
		StitchType:   "finish",
		Repeats:      1,
		Notes:        "Weave in ends",
	})

	return models.Part{
		Name:         "Body",
		Type:         "sphere",
		Rounds:       rounds,
		Color:        "main color",
		StartingType: "magic ring",
		Notes:        []string{"Use stitch marker to track rounds", "Stuff firmly for best shape"},
	}
}

// generateMaterials creates default materials list
func (g *Generator) generateMaterials() models.Materials {
	return models.Materials{
		YarnWeight:  g.DefaultYarnWeight,
		YarnYardage: 50, // TODO: Calculate based on pattern
		HookSize:    g.DefaultHookSize,
		Colors: []models.Color{
			{Name: "main color", Amount: 50},
		},
		OtherSupplies: []string{"stuffing", "yarn needle", "stitch marker"},
	}
}

// generateID creates a simple unique ID
func generateID() string {
	return fmt.Sprintf("pattern-%d", time.Now().Unix())
}
