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

	// Calculate target dimensions (assume 6 inches height)
	targetHeightInches := 6.0
	targetHeightCm := targetHeightInches * 2.54

	// Calculate width based on mesh proportions
	width, height, depth := m.GetDimensions()
	var targetWidthInches float64
	if shapeType == "cylinder" {
		// Use average of width and depth
		avgDiameter := (width + depth) / 2
		if height > 0 {
			targetWidthInches = targetHeightInches * (avgDiameter / height)
		} else {
			targetWidthInches = 2.0 // fallback
		}
	} else {
		// Sphere-like: use average radius
		avgRadius := m.GetAverageRadius()
		if height > 0 {
			targetWidthInches = targetHeightInches * (avgRadius * 2 / height)
		} else {
			targetWidthInches = targetHeightInches // spherical default
		}
	}
	targetWidthCm := targetWidthInches * 2.54

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
		FinishedSize: models.FinishedSize{
			HeightInches: targetHeightInches,
			HeightCm:     targetHeightCm,
			WidthInches:  targetWidthInches,
			WidthCm:      targetWidthCm,
		},
	}

	// Generate parts based on shape type
	switch shapeType {
	case "sphere":
		part := g.generateSpherePart(m)
		pattern.Parts = append(pattern.Parts, part)
	case "cylinder":
		part := g.generateCylinderPart(m)
		pattern.Parts = append(pattern.Parts, part)
	default:
		return nil, fmt.Errorf("unsupported shape type: %s", shapeType)
	}

	// Calculate accuracy metrics
	pattern.AccuracyMetrics = g.calculateAccuracy(m, pattern)

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

	// Check if it's a vertical/standing shape
	width, height, depth := m.GetDimensions()

	// If height is significantly larger than both width and depth, treat as cylinder
	// (vertical amigurumi like standing animals)
	heightToWidth := height / width
	heightToDepth := height / depth

	if (heightToWidth > 1.3 || heightToDepth > 1.3) {
		return "cylinder"
	}

	// If depth is significantly larger (forward-facing), also cylinder
	depthToWidth := depth / width
	depthToHeight := depth / height

	if (depthToWidth > 1.3 || depthToHeight > 1.3) {
		return "cylinder"
	}

	// Default to sphere for compact shapes
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

// generateCylinderPart creates rounds for a cylindrical/vertical shape
func (g *Generator) generateCylinderPart(m *mesh.Mesh) models.Part {
	// Calculate dimensions
	width, height, depth := m.GetDimensions()

	// Determine which dimension is the "length" (height for standing, depth for horizontal)
	var length float64
	if height > width && height > depth {
		// Standing/vertical
		length = height
	} else if depth > width && depth > height {
		// Forward-facing/horizontal
		length = depth
	} else {
		// Default to height
		length = height
	}

	// Calculate target number of rounds (higher resolution for better detail)
	lengthRounds := int(length * 5) // 5 rounds per unit of length for more detail
	if lengthRounds < 20 {
		lengthRounds = 20
	}
	if lengthRounds > 80 {
		lengthRounds = 80
	}

	// Get radius profile at different heights
	radiusProfile := m.GetRadiusProfile(lengthRounds)
	if len(radiusProfile) == 0 {
		// Fallback to constant diameter
		diameter := (width + depth) / 2
		radiusProfile = make([]float64, lengthRounds)
		for i := range radiusProfile {
			radiusProfile[i] = diameter / 2
		}
	}

	// Find max radius to normalize stitch counts
	maxRadius := 0.0
	for _, r := range radiusProfile {
		if r > maxRadius {
			maxRadius = r
		}
	}

	// Convert radii to stitch counts
	rawStitchProfile := make([]int, len(radiusProfile))
	maxStitches := int(maxRadius * 10) // 10 stitches per unit radius
	if maxStitches < 12 {
		maxStitches = 12
	}
	if maxStitches > 72 {
		maxStitches = 72
	}
	maxStitches = ((maxStitches + 5) / 6) * 6 // Round to multiple of 6

	for i, radius := range radiusProfile {
		stitches := int((radius / maxRadius) * float64(maxStitches))
		// Round to multiple of 6 and enforce minimum
		stitches = ((stitches + 5) / 6) * 6
		if stitches < 6 {
			stitches = 6
		}
		rawStitchProfile[i] = stitches
	}

	// Smooth the profile to ensure gradual increases/decreases
	// With more rounds, we can use smaller increments for better detail
	stitchProfile := make([]int, len(rawStitchProfile))
	stitchProfile[0] = 6 // Start with magic ring

	for i := 1; i < len(rawStitchProfile); i++ {
		targetStitches := rawStitchProfile[i]
		prevStitches := stitchProfile[i-1]

		diff := targetStitches - prevStitches
		maxChange := 6 // Smaller increments for finer detail

		if diff > maxChange {
			// Limit increase
			stitchProfile[i] = prevStitches + maxChange
		} else if diff < -maxChange {
			// Limit decrease
			stitchProfile[i] = prevStitches - maxChange
		} else {
			stitchProfile[i] = targetStitches
		}

		// Ensure multiple of 6
		stitchProfile[i] = ((stitchProfile[i] + 5) / 6) * 6
		if stitchProfile[i] < 6 {
			stitchProfile[i] = 6
		}
	}

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

	// Generate rounds following the stitch profile
	currentStitches := 6
	roundNum := 2
	stuffingStarted := false

	for i, targetStitches := range stitchProfile {
		if targetStitches == currentStitches {
			// Constant round
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("sc in each st around (%d sc)", currentStitches),
				StitchCount:  currentStitches,
				StitchType:   "sc",
				Repeats:      1,
				Notes:        "",
			})
		} else if targetStitches > currentStitches {
			// Increase round
			stitchDiff := targetStitches - currentStitches
			if stitchDiff%6 == 0 && stitchDiff <= currentStitches {
				// Can do even increases
				increments := stitchDiff / 6
				scBetween := (currentStitches / 6) - 1
				var instruction string
				if scBetween == 0 {
					instruction = "2 sc in each st around"
				} else {
					instruction = fmt.Sprintf("[%d inc, %d sc] repeat 6 times", increments, scBetween)
				}
				rounds = append(rounds, models.Round{
					Number:       roundNum,
					Instructions: instruction,
					StitchCount:  targetStitches,
					StitchType:   "sc",
					Repeats:      6,
					Notes:        "",
				})
			} else {
				// Irregular increase - just add stitches evenly
				rounds = append(rounds, models.Round{
					Number:       roundNum,
					Instructions: fmt.Sprintf("increase evenly to %d sts", targetStitches),
					StitchCount:  targetStitches,
					StitchType:   "sc",
					Repeats:      1,
					Notes:        "",
				})
			}
		} else {
			// Decrease round
			stitchDiff := currentStitches - targetStitches
			stuffNote := ""
			if !stuffingStarted {
				stuffNote = "Begin stuffing firmly"
				stuffingStarted = true
			} else if targetStitches <= 18 {
				stuffNote = "Finish stuffing"
			}

			if stitchDiff%6 == 0 && stitchDiff <= currentStitches/2 {
				// Can do even decreases
				decrements := stitchDiff / 6
				scBetween := (currentStitches / 6) - decrements - 1
				var instruction string
				if scBetween <= 0 {
					instruction = "dec around"
				} else {
					instruction = fmt.Sprintf("[%d dec, %d sc] repeat 6 times", decrements, scBetween)
				}
				rounds = append(rounds, models.Round{
					Number:       roundNum,
					Instructions: instruction,
					StitchCount:  targetStitches,
					StitchType:   "sc",
					Repeats:      6,
					Notes:        stuffNote,
				})
			} else {
				// Irregular decrease
				rounds = append(rounds, models.Round{
					Number:       roundNum,
					Instructions: fmt.Sprintf("decrease evenly to %d sts", targetStitches),
					StitchCount:  targetStitches,
					StitchType:   "sc",
					Repeats:      1,
					Notes:        stuffNote,
				})
			}
		}

		currentStitches = targetStitches
		roundNum++

		// Safety limit
		if i >= lengthRounds-1 {
			break
		}
	}

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
		Type:         "cylinder",
		Rounds:       rounds,
		Color:        "main color",
		StartingType: "magic ring",
		Notes:        []string{"Use stitch marker to track rounds", "Stuff firmly as you go", "Keep tension consistent for even shape"},
	}
}

// calculateAccuracy compares pattern stitch profile against mesh radius profile
func (g *Generator) calculateAccuracy(m *mesh.Mesh, pattern *models.Pattern) models.AccuracyMetrics {
	if len(pattern.Parts) == 0 {
		return models.AccuracyMetrics{
			ShapeMatchPercent: 0,
			Notes:             "No parts generated",
		}
	}

	part := pattern.Parts[0]
	rounds := part.Rounds

	// Get actual mesh radius profile
	numSlices := len(rounds)
	if numSlices == 0 {
		return models.AccuracyMetrics{
			ShapeMatchPercent: 0,
			Notes:             "No rounds generated",
		}
	}

	meshProfile := m.GetRadiusProfile(numSlices)
	if len(meshProfile) == 0 {
		return models.AccuracyMetrics{
			ShapeMatchPercent: 0,
			Notes:             "Could not analyze mesh",
		}
	}

	// Normalize both profiles
	maxMeshRadius := 0.0
	for _, r := range meshProfile {
		if r > maxMeshRadius {
			maxMeshRadius = r
		}
	}

	maxStitches := 0
	for _, round := range rounds {
		if round.StitchCount > maxStitches {
			maxStitches = round.StitchCount
		}
	}

	if maxMeshRadius == 0 || maxStitches == 0 {
		return models.AccuracyMetrics{
			ShapeMatchPercent: 0,
			Notes:             "Invalid profile data",
		}
	}

	// Calculate error metrics
	totalError := 0.0
	maxError := 0.0
	validComparisons := 0

	for i := 0; i < len(rounds) && i < len(meshProfile); i++ {
		if rounds[i].StitchCount == 0 {
			continue // Skip finish rounds
		}

		// Normalize to 0-1 range
		meshNorm := meshProfile[i] / maxMeshRadius
		stitchNorm := float64(rounds[i].StitchCount) / float64(maxStitches)

		error := meshNorm - stitchNorm
		if error < 0 {
			error = -error
		}

		totalError += error
		if error > maxError {
			maxError = error
		}
		validComparisons++
	}

	avgError := 0.0
	if validComparisons > 0 {
		avgError = totalError / float64(validComparisons)
	}

	// Calculate match percentage (100% - average error%)
	matchPercent := (1.0 - avgError) * 100
	if matchPercent < 0 {
		matchPercent = 0
	}

	// Generate notes based on accuracy
	var notes string
	if matchPercent >= 90 {
		notes = "Excellent match - pattern closely replicates the model shape"
	} else if matchPercent >= 75 {
		notes = "Good match - pattern captures the main features of the model"
	} else if matchPercent >= 60 {
		notes = "Fair match - pattern approximates the general shape"
	} else {
		notes = "Basic approximation - complex details simplified"
	}

	return models.AccuracyMetrics{
		ShapeMatchPercent: matchPercent,
		AverageError:      avgError * 100, // Convert to percentage
		MaxError:          maxError * 100,
		Notes:             notes,
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
