package main

import (
	"fmt"
	"time"

	"github.com/whitenhiemer/crochetbot/internal/models"
	"github.com/whitenhiemer/crochetbot/internal/pattern"
)

func main() {
	fmt.Println("=== QUICK PATTERN GENERATION TEST (30 shapes) ===\n")

	parser := pattern.NewParser()
	validator := pattern.Validator{}
	formatter := pattern.NewFormatter()

	shapes := []struct {
		name string
		gen  func() *models.Pattern
	}{
		// Spheres (5)
		{"Sphere-24st", func() *models.Pattern { return generateSphere(24, 2) }},
		{"Sphere-36st", func() *models.Pattern { return generateSphere(36, 3) }},
		{"Sphere-48st", func() *models.Pattern { return generateSphere(48, 4) }},
		{"Sphere-60st", func() *models.Pattern { return generateSphere(60, 4) }},
		{"Sphere-72st", func() *models.Pattern { return generateSphere(72, 5) }},

		// Cylinders (5)
		{"Cylinder-24st-10r", func() *models.Pattern { return generateCylinder(24, 10) }},
		{"Cylinder-30st-12r", func() *models.Pattern { return generateCylinder(30, 12) }},
		{"Cylinder-36st-14r", func() *models.Pattern { return generateCylinder(36, 14) }},
		{"Cylinder-42st-16r", func() *models.Pattern { return generateCylinder(42, 16) }},
		{"Cylinder-48st-18r", func() *models.Pattern { return generateCylinder(48, 18) }},

		// Cones (5)
		{"Cone-36st-t1", func() *models.Pattern { return generateCone(36, 1) }},
		{"Cone-42st-t2", func() *models.Pattern { return generateCone(42, 2) }},
		{"Cone-48st-t2", func() *models.Pattern { return generateCone(48, 2) }},
		{"Cone-54st-t3", func() *models.Pattern { return generateCone(54, 3) }},
		{"Cone-60st-t3", func() *models.Pattern { return generateCone(60, 3) }},

		// Ovals (5)
		{"Oval-36st-6e", func() *models.Pattern { return generateOval(36, 6) }},
		{"Oval-42st-9e", func() *models.Pattern { return generateOval(42, 9) }},
		{"Oval-48st-12e", func() *models.Pattern { return generateOval(48, 12) }},
		{"Oval-54st-15e", func() *models.Pattern { return generateOval(54, 15) }},
		{"Oval-60st-18e", func() *models.Pattern { return generateOval(60, 18) }},

		// Teardrops (5)
		{"Teardrop-36st-8tr", func() *models.Pattern { return generateTeardrop(36, 8) }},
		{"Teardrop-42st-10tr", func() *models.Pattern { return generateTeardrop(42, 10) }},
		{"Teardrop-48st-12tr", func() *models.Pattern { return generateTeardrop(48, 12) }},
		{"Teardrop-54st-14tr", func() *models.Pattern { return generateTeardrop(54, 14) }},
		{"Teardrop-60st-16tr", func() *models.Pattern { return generateTeardrop(60, 16) }},

		// Capsules (5)
		{"Capsule-30st-8m", func() *models.Pattern { return generateCapsule(30, 8) }},
		{"Capsule-36st-12m", func() *models.Pattern { return generateCapsule(36, 12) }},
		{"Capsule-42st-16m", func() *models.Pattern { return generateCapsule(42, 16) }},
		{"Capsule-48st-18m", func() *models.Pattern { return generateCapsule(48, 18) }},
		{"Capsule-54st-20m", func() *models.Pattern { return generateCapsule(54, 20) }},
	}

	startTime := time.Now()
	successCount := 0
	totalScore := 0.0

	for i, shape := range shapes {
		fmt.Printf("[%2d/%d] Testing %s... ", i+1, len(shapes), shape.name)

		pat := shape.gen()

		// Validate
		validation := validator.ValidatePattern(pat)

		// Format
		formatted := formatter.FormatPattern(pat)

		// Parse (round-trip)
		_, err := parser.ParsePattern(formatted)
		roundTripOK := (err == nil)

		totalScore += validation.Score
		if validation.IsValid {
			successCount++
		}

		status := "✓"
		if !validation.IsValid {
			status = "✗"
		}
		rtStatus := "✓"
		if !roundTripOK {
			rtStatus = "✗"
		}

		fmt.Printf("%s Valid:%s Score:%.1f RT:%s\n",
			status,
			func() string { if validation.IsValid { return "YES" } else { return "NO " } }(),
			validation.Score,
			rtStatus)
	}

	elapsed := time.Since(startTime)

	fmt.Println("\n=== RESULTS ===")
	fmt.Printf("Total: %d patterns\n", len(shapes))
	fmt.Printf("Valid: %d/%d (%.1f%%)\n", successCount, len(shapes), float64(successCount)/float64(len(shapes))*100)
	fmt.Printf("Avg Score: %.2f/100\n", totalScore/float64(len(shapes)))
	fmt.Printf("Time: %v (%.2f patterns/sec)\n", elapsed, float64(len(shapes))/elapsed.Seconds())
}

func generateSphere(maxStitches, evenRounds int) *models.Pattern {
	rounds := []models.Round{}
	currentStitches := 6
	roundNum := 1

	rounds = append(rounds, models.Round{
		Number:       roundNum,
		Instructions: "6 sc in magic ring",
		StitchCount:  6,
		StitchType:   "sc",
		Repeats:      1,
	})
	roundNum++

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

	for range evenRounds {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
			Repeats:      1,
		})
		roundNum++
	}

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
	currentStitches := 6

	for currentStitches <= maxStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
				Repeats:      1,
			})
		} else if currentStitches < maxStitches {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
				Repeats:      6,
			})
		}
		if currentStitches < maxStitches {
			currentStitches += 6
		} else {
			break
		}
		roundNum++
	}

	for range cylinderRounds {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  maxStitches,
			StitchType:   "sc",
			Repeats:      1,
		})
		roundNum++
	}

	currentStitches = maxStitches
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

		for range taperRate - 1 {
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

	for range middleRounds {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
			Repeats:      1,
		})
		roundNum++
	}

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
