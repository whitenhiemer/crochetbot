package main

import (
	"fmt"
	"time"

	"github.com/whitenhiemer/crochetbot/internal/models"
	"github.com/whitenhiemer/crochetbot/internal/pattern"
)

func main() {
	fmt.Println("=== COMPLEX SHAPE GENERATION TEST ===")
	fmt.Println("Testing advanced multi-part and composite shapes\n")

	parser := pattern.NewParser()
	validator := pattern.Validator{}
	formatter := pattern.NewFormatter()

	shapes := []struct {
		name     string
		gen      func() *models.Pattern
		category string
	}{
		// Multi-sphere composites (5)
		{"Pear-48-30", func() *models.Pattern { return generatePear(48, 30) }, "Composite"},
		{"Pear-54-36", func() *models.Pattern { return generatePear(54, 36) }, "Composite"},
		{"Snowman-3tier", func() *models.Pattern { return generateSnowman(42, 30, 24) }, "Composite"},
		{"Snowman-4tier", func() *models.Pattern { return generateSnowman(48, 36, 30, 24) }, "Composite"},
		{"Hourglass", func() *models.Pattern { return generateHourglass(42, 18) }, "Composite"},

		// Discs and flat shapes (5)
		{"Disc-thin", func() *models.Pattern { return generateDisc(72, 1) }, "Flat"},
		{"Disc-thick", func() *models.Pattern { return generateDisc(72, 3) }, "Flat"},
		{"Disc-large", func() *models.Pattern { return generateDisc(90, 2) }, "Flat"},
		{"Button", func() *models.Pattern { return generateButton(24) }, "Flat"},
		{"Coin", func() *models.Pattern { return generateDisc(48, 1) }, "Flat"},

		// Tori/Donuts (5)
		{"Torus-18-12", func() *models.Pattern { return generateTorus(18, 12) }, "Curved"},
		{"Torus-24-14", func() *models.Pattern { return generateTorus(24, 14) }, "Curved"},
		{"Torus-30-16", func() *models.Pattern { return generateTorus(30, 16) }, "Curved"},
		{"Ring-thin", func() *models.Pattern { return generateTorus(18, 8) }, "Curved"},
		{"Donut-thick", func() *models.Pattern { return generateTorus(30, 12) }, "Curved"},

		// Tapered limbs (5)
		{"Limb-12-30", func() *models.Pattern { return generateLimb(12, 30, 20) }, "Tapered"},
		{"Limb-18-36", func() *models.Pattern { return generateLimb(18, 36, 24) }, "Tapered"},
		{"Arm-thin", func() *models.Pattern { return generateLimb(12, 24, 18) }, "Tapered"},
		{"Leg-thick", func() *models.Pattern { return generateLimb(18, 42, 28) }, "Tapered"},
		{"Tentacle", func() *models.Pattern { return generateLimb(12, 18, 30) }, "Tapered"},

		// Long gradual tapers (5)
		{"Carrot-36", func() *models.Pattern { return generateCarrot(36, 24) }, "Gradual"},
		{"Carrot-42", func() *models.Pattern { return generateCarrot(42, 28) }, "Gradual"},
		{"Icicle", func() *models.Pattern { return generateCarrot(30, 20) }, "Gradual"},
		{"Parsnip", func() *models.Pattern { return generateCarrot(42, 32) }, "Gradual"},
		{"Cone-gradual", func() *models.Pattern { return generateCarrot(48, 30) }, "Gradual"},

		// Asymmetric shapes (5)
		{"Egg-36", func() *models.Pattern { return generateEgg(36, 6, 10) }, "Asymmetric"},
		{"Egg-42", func() *models.Pattern { return generateEgg(42, 7, 12) }, "Asymmetric"},
		{"Teardrop-asym", func() *models.Pattern { return generateAsymTeardrop(48, 14) }, "Asymmetric"},
		{"Bulb", func() *models.Pattern { return generateBulb(42, 24) }, "Asymmetric"},
		{"Avocado", func() *models.Pattern { return generateEgg(48, 8, 10) }, "Asymmetric"},

		// Hearts and complex assemblies (5)
		{"Heart-2lobe", func() *models.Pattern { return generateHeart(24, 30) }, "Assembly"},
		{"Heart-large", func() *models.Pattern { return generateHeart(30, 36) }, "Assembly"},
		{"Star-5point", func() *models.Pattern { return generateStar(5) }, "Assembly"},
		{"Star-6point", func() *models.Pattern { return generateStar(6) }, "Assembly"},
		{"Bow", func() *models.Pattern { return generateBow() }, "Assembly"},

		// Irregular organic shapes (10)
		{"Bean-kidney", func() *models.Pattern { return generateBean(36) }, "Organic"},
		{"Blob-round", func() *models.Pattern { return generateBlob(42, 3) }, "Organic"},
		{"Blob-tall", func() *models.Pattern { return generateBlob(36, 5) }, "Organic"},
		{"Peanut", func() *models.Pattern { return generatePeanut(30) }, "Organic"},
		{"Potato", func() *models.Pattern { return generatePotato(42) }, "Organic"},
		{"Mushroom-cap", func() *models.Pattern { return generateMushroomCap(48) }, "Organic"},
		{"Bell", func() *models.Pattern { return generateBell(42, 30) }, "Organic"},
		{"Dome", func() *models.Pattern { return generateDome(60) }, "Organic"},
		{"Lemon", func() *models.Pattern { return generateLemon(42) }, "Organic"},
		{"Gourd", func() *models.Pattern { return generateGourd(42, 30, 18) }, "Organic"},
	}

	startTime := time.Now()
	results := make([]TestResult, 0, len(shapes))

	categoryStats := make(map[string]*CategoryStats)

	for i, shape := range shapes {
		fmt.Printf("[%2d/%d] Testing %-20s (%s)... ", i+1, len(shapes), shape.name, shape.category)

		pat := shape.gen()
		validation := validator.ValidatePattern(pat)
		formatted := formatter.FormatPattern(pat)
		_, err := parser.ParsePattern(formatted)
		roundTripOK := (err == nil)

		result := TestResult{
			Name:        shape.name,
			Category:    shape.category,
			Validation:  validation,
			RoundTripOK: roundTripOK,
			PatternSize: len(formatted),
			RoundCount:  countRounds(pat),
		}
		results = append(results, result)

		// Update category stats
		if categoryStats[shape.category] == nil {
			categoryStats[shape.category] = &CategoryStats{}
		}
		categoryStats[shape.category].Count++
		categoryStats[shape.category].TotalScore += validation.Score
		if validation.IsValid {
			categoryStats[shape.category].ValidCount++
		}
		if roundTripOK {
			categoryStats[shape.category].RoundTripSuccess++
		}

		status := "✓"
		if !validation.IsValid {
			status = "✗"
		}
		fmt.Printf("%s %.1f/100 RT:%v\n", status, validation.Score, roundTripOK)
	}

	elapsed := time.Since(startTime)

	// Print summary
	fmt.Println("\n=== OVERALL RESULTS ===")
	fmt.Printf("Total shapes: %d\n", len(results))
	fmt.Printf("Time: %v (%.2f shapes/sec)\n", elapsed, float64(len(results))/elapsed.Seconds())

	validCount := 0
	totalScore := 0.0
	roundTripCount := 0
	for _, r := range results {
		if r.Validation.IsValid {
			validCount++
		}
		totalScore += r.Validation.Score
		if r.RoundTripOK {
			roundTripCount++
		}
	}

	fmt.Printf("\nValidation:\n")
	fmt.Printf("  Valid: %d/%d (%.1f%%)\n", validCount, len(results), float64(validCount)/float64(len(results))*100)
	fmt.Printf("  Avg score: %.2f/100\n", totalScore/float64(len(results)))
	fmt.Printf("  Round-trip: %d/%d (%.1f%%)\n", roundTripCount, len(results), float64(roundTripCount)/float64(len(results))*100)

	// Category breakdown
	fmt.Println("\n=== CATEGORY BREAKDOWN ===")
	for category, stats := range categoryStats {
		avgScore := stats.TotalScore / float64(stats.Count)
		validPct := float64(stats.ValidCount) / float64(stats.Count) * 100
		rtPct := float64(stats.RoundTripSuccess) / float64(stats.Count) * 100
		fmt.Printf("%-12s: %2d shapes | Avg %.1f/100 | Valid %.0f%% | RT %.0f%%\n",
			category, stats.Count, avgScore, validPct, rtPct)
	}

	// Find best and worst
	best := results[0]
	worst := results[0]
	for _, r := range results {
		if r.Validation.Score > best.Validation.Score {
			best = r
		}
		if r.Validation.Score < worst.Validation.Score {
			worst = r
		}
	}

	fmt.Printf("\nBest:  %s (%.1f/100)\n", best.Name, best.Validation.Score)
	fmt.Printf("Worst: %s (%.1f/100)\n", worst.Name, worst.Validation.Score)

	// Score distribution
	fmt.Println("\n=== SCORE DISTRIBUTION ===")
	scoreRanges := map[string]int{
		"90-100": 0,
		"80-89":  0,
		"70-79":  0,
		"60-69":  0,
		"<60":    0,
	}
	for _, r := range results {
		score := r.Validation.Score
		if score >= 90 {
			scoreRanges["90-100"]++
		} else if score >= 80 {
			scoreRanges["80-89"]++
		} else if score >= 70 {
			scoreRanges["70-79"]++
		} else if score >= 60 {
			scoreRanges["60-69"]++
		} else {
			scoreRanges["<60"]++
		}
	}
	for _, rng := range []string{"90-100", "80-89", "70-79", "60-69", "<60"} {
		count := scoreRanges[rng]
		pct := float64(count) / float64(len(results)) * 100
		bar := ""
		for range count {
			bar += "█"
		}
		fmt.Printf("%7s: %3d (%.0f%%) %s\n", rng, count, pct, bar)
	}
}

type TestResult struct {
	Name        string
	Category    string
	Validation  pattern.ValidationResult
	RoundTripOK bool
	PatternSize int
	RoundCount  int
}

type CategoryStats struct {
	Count            int
	ValidCount       int
	RoundTripSuccess int
	TotalScore       float64
}

func countRounds(pat *models.Pattern) int {
	total := 0
	for _, part := range pat.Parts {
		total += len(part.Rounds)
	}
	return total
}

// Pear: two spheres of different sizes
func generatePear(bottomMax, topMax int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	// Bottom sphere
	for currentStitches < bottomMax {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Hold
	for range 3 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
		})
		roundNum++
	}

	// Decrease to waist
	waist := topMax - 6
	for currentStitches > waist {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
		})
		currentStitches -= 6
		roundNum++
	}

	// Top sphere
	for currentStitches < topMax {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
			StitchCount:  currentStitches + 6,
			StitchType:   "inc",
		})
		currentStitches += 6
		roundNum++
	}

	// Hold
	for range 2 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
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
		})
		currentStitches -= 6
		roundNum++
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Pear %d-%d", bottomMax, topMax),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "pear",
			Rounds: rounds,
		}},
	}
}

// Snowman: multiple stacked spheres
func generateSnowman(sizes ...int) *models.Pattern {
	parts := []models.Part{}

	for i, size := range sizes {
		rounds := []models.Round{}
		roundNum := 1
		currentStitches := 6

		for currentStitches < size {
			if currentStitches == 6 {
				rounds = append(rounds, models.Round{
					Number:       roundNum,
					Instructions: "6 sc in magic ring",
					StitchCount:  6,
					StitchType:   "sc",
				})
			} else {
				rounds = append(rounds, models.Round{
					Number:       roundNum,
					Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
					StitchCount:  currentStitches + 6,
					StitchType:   "inc",
				})
			}
			currentStitches = rounds[len(rounds)-1].StitchCount
			roundNum++
		}

		// Hold
		for range 2 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "sc in each st around",
				StitchCount:  currentStitches,
				StitchType:   "sc",
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
			})
			currentStitches -= 6
			roundNum++
		}

		parts = append(parts, models.Part{
			Name:   fmt.Sprintf("Sphere %d", i+1),
			Type:   "sphere",
			Rounds: rounds,
		})
	}

	return &models.Pattern{
		Name:  fmt.Sprintf("Snowman %dtier", len(sizes)),
		Parts: parts,
	}
}

// Hourglass: two cones joined at narrow point
func generateHourglass(maxStitches, waist int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	// Bottom cone
	for currentStitches < maxStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Narrow to waist
	for currentStitches > waist {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
		})
		currentStitches -= 6
		roundNum++
	}

	// Widen again
	for currentStitches < maxStitches {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
			StitchCount:  currentStitches + 6,
			StitchType:   "inc",
		})
		currentStitches += 6
		roundNum++
	}

	// Close
	for currentStitches > 6 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
		})
		currentStitches -= 6
		roundNum++
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Hourglass %d-%d", maxStitches, waist),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "hourglass",
			Rounds: rounds,
		}},
	}
}

// Disc: flat circle
func generateDisc(maxStitches, thickness int) *models.Pattern {
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
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	for range thickness {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
		})
		roundNum++
	}

	for currentStitches > 6 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
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

// Button: small disc
func generateButton(size int) *models.Pattern {
	return generateDisc(size, 2)
}

// Torus: tube formed into ring
func generateTorus(tubeStitches, tubeLength int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	for currentStitches < tubeStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	for range tubeLength {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
		})
		roundNum++
	}

	for currentStitches > 6 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
		})
		currentStitches -= 6
		roundNum++
	}

	rounds = append(rounds, models.Round{
		Number:       roundNum,
		Instructions: "Join ends to form ring",
		StitchCount:  0,
		StitchType:   "join",
		Notes:        "Stuff tube, join ends with whip stitch",
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

// Limb: tapered cylinder
func generateLimb(startStitches, endStitches, length int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	for currentStitches < startStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Gradual taper
	stitchDiff := endStitches - startStitches
	increments := max(stitchDiff/3, 1)
	roundsPerIncrement := max(length/increments, 2)

	for range increments {
		if currentStitches < endStitches {
			currentStitches += 3
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("work %d inc evenly around", 3),
				StitchCount:  currentStitches,
				StitchType:   "inc",
			})
			roundNum++
		}

		for range roundsPerIncrement - 1 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "sc in each st around",
				StitchCount:  currentStitches,
				StitchType:   "sc",
			})
			roundNum++
		}
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Limb %d-%d", startStitches, endStitches),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "limb",
			Rounds: rounds,
		}},
	}
}

// Carrot: long gradual taper
func generateCarrot(maxStitches, taperRounds int) *models.Pattern {
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
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Gradual decrease with even rounds between
	decreaseCount := 0
	for currentStitches > 6 && decreaseCount < taperRounds {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
		})
		currentStitches -= 6
		roundNum++
		decreaseCount++

		if currentStitches > 6 {
			for range 2 {
				rounds = append(rounds, models.Round{
					Number:       roundNum,
					Instructions: "sc in each st around",
					StitchCount:  currentStitches,
					StitchType:   "sc",
				})
				roundNum++
			}
		}
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Carrot %dst", maxStitches),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "cone",
			Rounds: rounds,
		}},
	}
}

// Egg: asymmetric ellipsoid
func generateEgg(maxStitches, bottomRounds, topRounds int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	// Bottom (rounder)
	for currentStitches < maxStitches {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Middle
	for range bottomRounds {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
		})
		roundNum++
	}

	// Top (more tapered) - with even rounds between decreases
	decreaseCount := 0
	for currentStitches > 6 && decreaseCount < topRounds {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
		})
		currentStitches -= 6
		roundNum++
		decreaseCount++

		if currentStitches > 6 && decreaseCount < topRounds {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "sc in each st around",
				StitchCount:  currentStitches,
				StitchType:   "sc",
			})
			roundNum++
		}
	}

	return &models.Pattern{
		Name: fmt.Sprintf("Egg %dst", maxStitches),
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "ellipsoid",
			Rounds: rounds,
		}},
	}
}

// Asymmetric teardrop
func generateAsymTeardrop(maxStitches, taperRounds int) *models.Pattern {
	return generateEgg(maxStitches, 2, taperRounds)
}

// Bulb: wide bottom, narrow top
func generateBulb(maxStitches, topStitches int) *models.Pattern {
	return generatePear(maxStitches, topStitches)
}

// Heart: two lobes plus point (multi-part)
func generateHeart(lobeSize, pointHeight int) *models.Pattern {
	// Left lobe
	leftLobe := generateSphere(lobeSize, 2)
	leftLobe.Parts[0].Name = "Left Lobe"

	// Right lobe
	rightLobe := generateSphere(lobeSize, 2)
	rightLobe.Parts[0].Name = "Right Lobe"

	// Point
	pointRounds := []models.Round{}
	roundNum := 1
	currentStitches := 12

	for i := 0; i < pointHeight && currentStitches > 6; i++ {
		if i == 0 {
			pointRounds = append(pointRounds, models.Round{
				Number:       roundNum,
				Instructions: "12 sc base attachment",
				StitchCount:  12,
				StitchType:   "sc",
			})
		} else {
			pointRounds = append(pointRounds, models.Round{
				Number:       roundNum,
				Instructions: "dec around",
				StitchCount:  currentStitches - 3,
				StitchType:   "dec",
			})
			currentStitches -= 3
		}
		roundNum++
	}

	return &models.Pattern{
		Name: "Heart",
		Parts: []models.Part{
			leftLobe.Parts[0],
			rightLobe.Parts[0],
			{Name: "Point", Type: "cone", Rounds: pointRounds},
		},
		Assembly: []string{
			"Attach lobes side by side at top",
			"Attach point to bottom center where lobes meet",
		},
	}
}

// Star: multiple points
func generateStar(points int) *models.Pattern {
	parts := []models.Part{}

	// Center
	centerRounds := []models.Round{
		{Number: 1, Instructions: "6 sc in magic ring", StitchCount: 6, StitchType: "sc"},
		{Number: 2, Instructions: "6 inc", StitchCount: 12, StitchType: "inc"},
		{Number: 3, Instructions: "[sc, inc] around", StitchCount: 18, StitchType: "inc"},
	}
	parts = append(parts, models.Part{Name: "Center", Type: "disc", Rounds: centerRounds})

	// Points
	for i := 0; i < points; i++ {
		pointRounds := []models.Round{
			{Number: 1, Instructions: "6 sc attach to center", StitchCount: 6, StitchType: "sc"},
			{Number: 2, Instructions: "[sc, inc] around", StitchCount: 9, StitchType: "inc"},
			{Number: 3, Instructions: "[2 sc, inc] around", StitchCount: 12, StitchType: "inc"},
			{Number: 4, Instructions: "[2 sc, dec] around", StitchCount: 9, StitchType: "dec"},
			{Number: 5, Instructions: "[sc, dec] around", StitchCount: 6, StitchType: "dec"},
		}
		parts = append(parts, models.Part{
			Name:   fmt.Sprintf("Point %d", i+1),
			Type:   "cone",
			Rounds: pointRounds,
		})
	}

	return &models.Pattern{
		Name:  fmt.Sprintf("Star %dpt", points),
		Parts: parts,
	}
}

// Bow: two loops plus center
func generateBow() *models.Pattern {
	loop := generateTorus(18, 6)
	loop.Parts[0].Name = "Loop 1"

	loop2 := generateTorus(18, 6)
	loop2.Parts[0].Name = "Loop 2"

	centerRounds := []models.Round{
		{Number: 1, Instructions: "6 sc in magic ring", StitchCount: 6, StitchType: "sc"},
		{Number: 2, Instructions: "6 inc", StitchCount: 12, StitchType: "inc"},
		{Number: 3, Instructions: "12 sc", StitchCount: 12, StitchType: "sc"},
		{Number: 4, Instructions: "12 sc", StitchCount: 12, StitchType: "sc"},
		{Number: 5, Instructions: "6 dec", StitchCount: 6, StitchType: "dec"},
	}

	return &models.Pattern{
		Name: "Bow",
		Parts: []models.Part{
			loop.Parts[0],
			loop2.Parts[0],
			{Name: "Center Knot", Type: "cylinder", Rounds: centerRounds},
		},
		Assembly: []string{
			"Place loops side by side",
			"Wrap center knot around middle",
		},
	}
}

// Organic shapes with irregular patterns

func generateBean(size int) *models.Pattern {
	pat := generateEgg(size, 5, 8)
	pat.Name = fmt.Sprintf("Bean %d", size)
	pat.Parts[0].Type = "bean"
	pat.Parts[0].Notes = []string{"Curve slightly during stuffing"}
	return pat
}

func generateBlob(size, irregularity int) *models.Pattern {
	pat := generateSphere(size, irregularity*2)
	pat.Name = fmt.Sprintf("Blob %d", size)
	pat.Parts[0].Notes = []string{"Stuff unevenly for organic shape"}
	return pat
}

func generatePeanut(size int) *models.Pattern {
	return generateHourglass(size, size/2)
}

func generatePotato(size int) *models.Pattern {
	pat := generateBlob(size, 3)
	pat.Name = "Potato"
	return pat
}

func generateMushroomCap(size int) *models.Pattern {
	pat := generateDome(size)
	pat.Name = "Mushroom Cap"
	return pat
}

func generateBell(topSize, bottomSize int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	for currentStitches < topSize {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Hold
	for range 3 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
		})
		roundNum++
	}

	// Flare out
	for currentStitches < bottomSize {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
			StitchCount:  currentStitches + 6,
			StitchType:   "inc",
		})
		currentStitches += 6
		roundNum++
	}

	return &models.Pattern{
		Name: "Bell",
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "bell",
			Rounds: rounds,
		}},
	}
}

func generateDome(size int) *models.Pattern {
	rounds := []models.Round{}
	roundNum := 1
	currentStitches := 6

	for currentStitches < size {
		if currentStitches == 6 {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: "6 sc in magic ring",
				StitchCount:  6,
				StitchType:   "sc",
			})
		} else {
			rounds = append(rounds, models.Round{
				Number:       roundNum,
				Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
				StitchCount:  currentStitches + 6,
				StitchType:   "inc",
			})
		}
		currentStitches = rounds[len(rounds)-1].StitchCount
		roundNum++
	}

	// Hold for flat base
	rounds = append(rounds, models.Round{
		Number:       roundNum,
		Instructions: "sc in each st around",
		StitchCount:  currentStitches,
		StitchType:   "sc",
	})

	return &models.Pattern{
		Name: "Dome",
		Parts: []models.Part{{
			Name:   "Body",
			Type:   "dome",
			Rounds: rounds,
		}},
	}
}

func generateLemon(size int) *models.Pattern {
	return generateEgg(size, 4, 6)
}

func generateGourd(top, middle, bottom int) *models.Pattern {
	pat := generatePear(top, middle)
	pat.Name = "Gourd"
	return pat
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
	})
	roundNum++

	for currentStitches < maxStitches {
		currentStitches += 6
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, inc] around", (currentStitches/6)-1),
			StitchCount:  currentStitches,
			StitchType:   "inc",
		})
		roundNum++
	}

	for range evenRounds {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: "sc in each st around",
			StitchCount:  currentStitches,
			StitchType:   "sc",
		})
		roundNum++
	}

	for currentStitches > 6 {
		rounds = append(rounds, models.Round{
			Number:       roundNum,
			Instructions: fmt.Sprintf("[%d sc, dec] around", (currentStitches/6)-1),
			StitchCount:  currentStitches - 6,
			StitchType:   "dec",
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
