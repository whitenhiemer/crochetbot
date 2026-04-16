package main

import (
	"fmt"
	"os"
	"time"

	"github.com/whitenhiemer/crochetbot/internal/models"
	"github.com/whitenhiemer/crochetbot/internal/pattern"
)

func main() {
	fmt.Println("=== PROCEDURAL PATTERN GENERATION TEST ===")
	fmt.Println("Generating and validating 300 pattern variations\n")

	parser := pattern.NewParser()
	validator := pattern.NewValidator()
	formatter := pattern.NewFormatter()

	results := make([]TestResult, 0, 300)
	startTime := time.Now()

	shapeID := 1

	// Category 1: Spheres (50 variations)
	fmt.Println("Generating spheres...")
	for maxStitches := 24; maxStitches <= 72; maxStitches += 6 {
		for evenRounds := 0; evenRounds <= 4; evenRounds++ {
			if shapeID > 300 {
				break
			}
			result := testSphere(shapeID, maxStitches, evenRounds, parser, validator, formatter)
			results = append(results, result)
			shapeID++
		}
	}

	// Category 2: Cylinders (50 variations)
	fmt.Println("Generating cylinders...")
	for maxStitches := 24; maxStitches <= 48; maxStitches += 6 {
		for cylinderRounds := 6; cylinderRounds <= 20; cylinderRounds += 2 {
			if shapeID > 300 {
				break
			}
			result := testCylinder(shapeID, maxStitches, cylinderRounds, parser, validator, formatter)
			results = append(results, result)
			shapeID++
		}
	}

	// Category 3: Cones (40 variations)
	fmt.Println("Generating cones...")
	for maxStitches := 30; maxStitches <= 72; maxStitches += 6 {
		for taperRate := 1; taperRate <= 3; taperRate++ {
			if shapeID > 300 {
				break
			}
			result := testCone(shapeID, maxStitches, taperRate, parser, validator, formatter)
			results = append(results, result)
			shapeID++
		}
	}

	// Category 4: Ovals (40 variations)
	fmt.Println("Generating ovals...")
	for maxStitches := 36; maxStitches <= 60; maxStitches += 6 {
		for elongation := 6; elongation <= 18; elongation += 3 {
			if shapeID > 300 {
				break
			}
			result := testOval(shapeID, maxStitches, elongation, parser, validator, formatter)
			results = append(results, result)
			shapeID++
		}
	}

	// Category 5: Teardrops (30 variations)
	fmt.Println("Generating teardrops...")
	for maxStitches := 36; maxStitches <= 54; maxStitches += 6 {
		for taperRounds := 8; taperRounds <= 16; taperRounds += 2 {
			if shapeID > 300 {
				break
			}
			result := testTeardrop(shapeID, maxStitches, taperRounds, parser, validator, formatter)
			results = append(results, result)
			shapeID++
		}
	}

	// Category 6: Capsules (30 variations)
	fmt.Println("Generating capsules...")
	for maxStitches := 30; maxStitches <= 42; maxStitches += 6 {
		for middleRounds := 8; middleRounds <= 20; middleRounds += 3 {
			if shapeID > 300 {
				break
			}
			result := testCapsule(shapeID, maxStitches, middleRounds, parser, validator, formatter)
			results = append(results, result)
			shapeID++
		}
	}

	// Category 7: Discs (20 variations)
	fmt.Println("Generating discs...")
	for maxStitches := 54; maxStitches <= 84; maxStitches += 6 {
		for thickness := 1; thickness <= 3; thickness++ {
			if shapeID > 300 {
				break
			}
			result := testDisc(shapeID, maxStitches, thickness, parser, validator, formatter)
			results = append(results, result)
			shapeID++
		}
	}

	// Category 8: Tori (20 variations)
	fmt.Println("Generating tori...")
	for tubeStitches := 18; tubeStitches <= 30; tubeStitches += 3 {
		for tubeLength := 10; tubeLength <= 16; tubeLength += 2 {
			if shapeID > 300 {
				break
			}
			result := testTorus(shapeID, tubeStitches, tubeLength, parser, validator, formatter)
			results = append(results, result)
			shapeID++
		}
	}

	// Category 9: Limbs (20 variations)
	fmt.Println("Generating limbs...")
	for startStitches := 12; startStitches <= 24; startStitches += 3 {
		for endStitches := 24; endStitches <= 36; endStitches += 3 {
			if shapeID > 300 || startStitches >= endStitches {
				continue
			}
			result := testLimb(shapeID, startStitches, endStitches, parser, validator, formatter)
			results = append(results, result)
			shapeID++
			if shapeID > 300 {
				break
			}
		}
	}

	elapsed := time.Since(startTime)

	// Print summary statistics
	fmt.Println("\n=== RESULTS SUMMARY ===")
	fmt.Printf("Total patterns generated: %d\n", len(results))
	fmt.Printf("Total time: %v (%.2f patterns/sec)\n", elapsed, float64(len(results))/elapsed.Seconds())

	successCount := 0
	totalScore := 0.0
	totalTerminology := 0.0
	totalStructural := 0.0
	totalRealism := 0.0
	roundTripSuccess := 0

	for _, r := range results {
		if r.ValidationResult.IsValid {
			successCount++
		}
		totalScore += r.ValidationResult.Score
		totalTerminology += r.ValidationResult.TerminologyScore
		totalStructural += r.ValidationResult.StructuralScore
		totalRealism += r.ValidationResult.RealismScore
		if r.RoundTripSuccess {
			roundTripSuccess++
		}
	}

	fmt.Printf("\nValidation Results:\n")
	fmt.Printf("  Valid patterns: %d/%d (%.1f%%)\n", successCount, len(results), float64(successCount)/float64(len(results))*100)
	fmt.Printf("  Average score: %.2f/100\n", totalScore/float64(len(results)))
	fmt.Printf("  Average terminology: %.2f/100\n", totalTerminology/float64(len(results)))
	fmt.Printf("  Average structural: %.2f/100\n", totalStructural/float64(len(results)))
	fmt.Printf("  Average realism: %.2f/100\n", totalRealism/float64(len(results)))
	fmt.Printf("  Round-trip success: %d/%d (%.1f%%)\n", roundTripSuccess, len(results), float64(roundTripSuccess)/float64(len(results))*100)

	// Category breakdown
	fmt.Println("\nCategory Breakdown:")
	categories := map[string][]TestResult{
		"Spheres":   results[0:min(50, len(results))],
		"Cylinders": results[min(50, len(results)):min(100, len(results))],
		"Cones":     results[min(100, len(results)):min(140, len(results))],
		"Ovals":     results[min(140, len(results)):min(180, len(results))],
		"Teardrops": results[min(180, len(results)):min(210, len(results))],
		"Capsules":  results[min(210, len(results)):min(240, len(results))],
		"Discs":     results[min(240, len(results)):min(260, len(results))],
		"Tori":      results[min(260, len(results)):min(280, len(results))],
		"Limbs":     results[min(280, len(results)):min(300, len(results))],
	}

	for name, catResults := range categories {
		if len(catResults) == 0 {
			continue
		}
		avgScore := 0.0
		for _, r := range catResults {
			avgScore += r.ValidationResult.Score
		}
		avgScore /= float64(len(catResults))
		fmt.Printf("  %s: %.2f/100 (%d patterns)\n", name, avgScore, len(catResults))
	}

	// Find best and worst
	best := results[0]
	worst := results[0]
	for _, r := range results {
		if r.ValidationResult.Score > best.ValidationResult.Score {
			best = r
		}
		if r.ValidationResult.Score < worst.ValidationResult.Score {
			worst = r
		}
	}

	fmt.Printf("\nBest pattern: #%d %s (%.2f/100)\n", best.ShapeID, best.ShapeName, best.ValidationResult.Score)
	fmt.Printf("Worst pattern: #%d %s (%.2f/100)\n", worst.ShapeID, worst.ShapeName, worst.ValidationResult.Score)

	// Save detailed results to file
	saveResults(results)

	fmt.Println("\n=== TEST COMPLETE ===")
	fmt.Println("Detailed results saved to test_300_shapes_results.txt")
}

type TestResult struct {
	ShapeID          int
	ShapeName        string
	ValidationResult pattern.ValidationResult
	FormattedLength  int
	RoundTripSuccess bool
	ParseTime        time.Duration
	ValidateTime     time.Duration
	FormatTime       time.Duration
}

func testSphere(id, maxStitches, evenRounds int, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	pat := generateSphere(maxStitches, evenRounds)
	return runTest(id, fmt.Sprintf("Sphere-%dst-%dev", maxStitches, evenRounds), pat, parser, validator, formatter)
}

func testCylinder(id, maxStitches, cylinderRounds int, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	pat := generateCylinder(maxStitches, cylinderRounds)
	return runTest(id, fmt.Sprintf("Cylinder-%dst-%dr", maxStitches, cylinderRounds), pat, parser, validator, formatter)
}

func testCone(id, maxStitches, taperRate int, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	pat := generateCone(maxStitches, taperRate)
	return runTest(id, fmt.Sprintf("Cone-%dst-taper%d", maxStitches, taperRate), pat, parser, validator, formatter)
}

func testOval(id, maxStitches, elongation int, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	pat := generateOval(maxStitches, elongation)
	return runTest(id, fmt.Sprintf("Oval-%dst-%delong", maxStitches, elongation), pat, parser, validator, formatter)
}

func testTeardrop(id, maxStitches, taperRounds int, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	pat := generateTeardrop(maxStitches, taperRounds)
	return runTest(id, fmt.Sprintf("Teardrop-%dst-%dtr", maxStitches, taperRounds), pat, parser, validator, formatter)
}

func testCapsule(id, maxStitches, middleRounds int, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	pat := generateCapsule(maxStitches, middleRounds)
	return runTest(id, fmt.Sprintf("Capsule-%dst-%dmr", maxStitches, middleRounds), pat, parser, validator, formatter)
}

func testDisc(id, maxStitches, thickness int, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	pat := generateDisc(maxStitches, thickness)
	return runTest(id, fmt.Sprintf("Disc-%dst-%dthick", maxStitches, thickness), pat, parser, validator, formatter)
}

func testTorus(id, tubeStitches, tubeLength int, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	pat := generateTorus(tubeStitches, tubeLength)
	return runTest(id, fmt.Sprintf("Torus-%dtube-%dlen", tubeStitches, tubeLength), pat, parser, validator, formatter)
}

func testLimb(id, startStitches, endStitches int, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	pat := generateLimb(startStitches, endStitches)
	return runTest(id, fmt.Sprintf("Limb-%dto%d", startStitches, endStitches), pat, parser, validator, formatter)
}

func runTest(id int, name string, pat *models.Pattern, parser *pattern.Parser, validator *pattern.Validator, formatter *pattern.Formatter) TestResult {
	result := TestResult{
		ShapeID:   id,
		ShapeName: name,
	}

	// Validate
	start := time.Now()
	result.ValidationResult = validator.ValidatePattern(pat)
	result.ValidateTime = time.Since(start)

	// Format
	start = time.Now()
	formatted := formatter.FormatPattern(pat)
	result.FormatTime = time.Since(start)
	result.FormattedLength = len(formatted)

	// Parse (round-trip)
	start = time.Now()
	reparsed, err := parser.ParsePattern(formatted)
	result.ParseTime = time.Since(start)
	result.RoundTripSuccess = (err == nil && len(reparsed.Parts) > 0)

	return result
}

func generateSphere(maxStitches, evenRounds int) *models.Pattern {
	rounds := []models.Round{}
	currentStitches := 6
	roundNum := 1

	// Magic ring
	rounds = append(rounds, models.Round{
		Number:       roundNum,
		Instructions: "6 sc in magic ring",
		StitchCount:  6,
		StitchType:   "sc",
		Repeats:      1,
	})
	roundNum++

	// Increase to max
	for currentStitches < maxStitches {
		currentStitches += 6
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
			StitchCount:  currentStitches,
			StitchType:   "inc",
			Repeats:      6,
		})
		roundNum++
	}

	// Even rounds
	for i := 0; i < evenRounds; i++ {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
			Repeats:      1,
		})
		roundNum++
	}

	// Decrease
	for currentStitches > 6 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
			Repeats:      6,
		})
		currentStitches -= 6
		roundNum++
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Sphere %dst", maxStitches),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "sphere",
			Rounds: rounds,
		}},
	}
}

func generateCylinder(maxStitches, cylinderRounds int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1

	// Increase to max
	currentStitches := 6
	for currentStitches < maxStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
				Repeats:      1,
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
				Repeats:      6,
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Cylinder rounds
	for i := 0; i < cylinderRounds; i++ {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
			Repeats:      1,
		})
		roundNum++
	}

	// Decrease
	for currentStitches > 6 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
			Repeats:      6,
		})
		currentStitches -= 6
		roundNum++
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Cylinder %dst %dr", maxStitches, cylinderRounds),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "cylinder",
			Rounds: rounds,
		}},
	}
}

func generateCone(maxStitches, taperRate int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	rounds = append(rounds, models.Round{
		Number:       1,
		Instructions: "6 sc in magic ring",
		StitchCount:  6,
		StitchType:   "sc",
		Repeats:      1,
	})
	roundNum++

	// Increase to max with taper
	for currentStitches < maxStitches {
		currentStitches += 6
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
			StitchCount:  currentStitches,
			StitchType:   "inc",
			Repeats:      6,
		})
		roundNum++

		// Add even rounds based on taper rate
		for i := 0; i < taperRate-1; i++ {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "sc in each st around",
				StitchCount:  currentStitches,
				StitchType:   "sc",
				Repeats:      1,
			})
			roundNum++
		}
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Cone %dst taper%d", maxStitches, taperRate),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "cone",
			Rounds: rounds,
		}},
	}
}

func generateOval(maxStitches, elongation int) *models.Pattern {
	pat := generateSphere(maxStitches, elongation)
	pat.Name = fmt.Sprintf("Oval %dst %delong", maxStitches, elongation)
	pat.Parts[0].Type = "ellipsoid"
	return pat
}

func generateTeardrop(maxStitches, taperRounds int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	// Bottom (fast increase)
	for currentStitches < maxStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
				Repeats:      1,
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
				Repeats:      6,
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Top (slow decrease with even rounds)
	decreaseCount := 0
	for currentStitches > 6 && decreaseCount < taperRounds {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
			Repeats:      6,
		})
		currentStitches -= 6
		roundNum++
		decreaseCount++

		if currentStitches > 6 && decreaseCount < taperRounds {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "sc in each st around",
				StitchCount:  currentStitches,
				StitchType:   "sc",
				Repeats:      1,
			})
			roundNum++
		}
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Teardrop %dst %dtr", maxStitches, taperRounds),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "teardrop",
			Rounds: rounds,
		}},
	}
}

func generateCapsule(maxStitches, middleRounds int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	// First hemisphere
	for currentStitches < maxStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
				Repeats:      1,
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
				Repeats:      6,
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Middle
	for i := 0; i < middleRounds; i++ {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
			Repeats:      1,
		})
		roundNum++
	}

	// Second hemisphere
	for currentStitches > 6 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
			Repeats:      6,
		})
		currentStitches -= 6
		roundNum++
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Capsule %dst %dmr", maxStitches, middleRounds),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "capsule",
			Rounds: rounds,
		}},
	}
}

func generateDisc(maxStitches, thickness int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	// Expand
	for currentStitches < maxStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
				Repeats:      1,
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
				Repeats:      6,
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Thickness
	for i := 0; i < thickness; i++ {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
			Repeats:      1,
		})
		roundNum++
	}

	// Contract
	for currentStitches > 6 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
			Repeats:      6,
		})
		currentStitches -= 6
		roundNum++
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Disc %dst %dthick", maxStitches, thickness),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "disc",
			Rounds: rounds,
		}},
	}
}

func generateTorus(tubeStitches, tubeLength int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1

	// Tube start
	rounds = append(rounds, models.Round{
		Number:       roundNum,
		Instructions: "6 sc in magic ring",
		StitchCount:  6,
		StitchType:   "sc",
		Repeats:      1,
	})
	roundNum++

	currentStitches := 6
	for currentStitches < tubeStitches {
		currentStitches += 6
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
			StitchCount:  currentStitches,
			StitchType:   "inc",
			Repeats:      6,
		})
		roundNum++
	}

	// Tube length
	for i := 0; i < tubeLength; i++ {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
			Repeats:      1,
		})
		roundNum++
	}

	// Close
	for currentStitches > 6 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
			Repeats:      6,
		})
		currentStitches -= 6
		roundNum++
	}

	rounds = append(rounds, models.Round{
		Number:       roundNum,
		Instructions: "Join ends with whip stitch",
		StitchCount:  0,
		StitchType:   "join",
		Repeats:      1,
		Notes:        "Form torus by joining tube ends",
	})

	return &models.Pattern{
		Name: fmt.Sprintf("Torus %dtube %dlen", tubeStitches, tubeLength),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "torus",
			Rounds: rounds,
		}},
	}
}

func generateLimb(startStitches, endStitches int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	// Start
	for currentStitches < startStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
				Repeats:      1,
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
				Repeats:      6,
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Gradual taper
	stitchDiff := endStitches - startStitches
	increments := stitchDiff / 3

	for i := 0; i < increments; i++ {
		// Add 3 stitches
		currentStitches += 3
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, inc] around", currentStitches/3-1),
			StitchCount:  currentStitches,
			StitchType:   "inc",
			Repeats:      3,
		})
		roundNum++

		// Even rounds
		for j := 0; j < 3; j++ {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "sc in each st around",
				StitchCount:  currentStitches,
				StitchType:   "sc",
				Repeats:      1,
			})
			roundNum++
		}
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Limb %dto%d", startStitches, endStitches),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "limb",
			Rounds: rounds,
		}},
	}
}

func saveResults(results []TestResult) {
	f, err := os.Create("test_300_shapes_results.txt")
	if err != nil {
		fmt.Println("Error creating results file:", err)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "=== 300 SHAPE GENERATION TEST RESULTS ===\n\n")
	fmt.Fprintf(f, "%-5s %-30s %-8s %-8s %-8s %-8s %-8s %-12s\n",
		"ID", "Name", "Valid", "Score", "Term", "Struct", "Real", "RoundTrip")
	fmt.Fprintf(f, "%s\n", string(make([]byte, 100)))

	for _, r := range results {
		validStr := "NO"
		if r.ValidationResult.IsValid {
			validStr = "YES"
		}
		rtStr := "FAIL"
		if r.RoundTripSuccess {
			rtStr = "PASS"
		}

		fmt.Fprintf(f, "%-5d %-30s %-8s %-8.2f %-8.2f %-8.2f %-8.2f %-12s\n",
			r.ShapeID,
			r.ShapeName,
			validStr,
			r.ValidationResult.Score,
			r.ValidationResult.TerminologyScore,
			r.ValidationResult.StructuralScore,
			r.ValidationResult.RealismScore,
			rtStr,
		)
	}

	fmt.Fprintf(f, "\n=== PERFORMANCE METRICS ===\n")
	totalParse := time.Duration(0)
	totalValidate := time.Duration(0)
	totalFormat := time.Duration(0)

	for _, r := range results {
		totalParse += r.ParseTime
		totalValidate += r.ValidateTime
		totalFormat += r.FormatTime
	}

	fmt.Fprintf(f, "Average parse time: %v\n", totalParse/time.Duration(len(results)))
	fmt.Fprintf(f, "Average validate time: %v\n", totalValidate/time.Duration(len(results)))
	fmt.Fprintf(f, "Average format time: %v\n", totalFormat/time.Duration(len(results)))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
