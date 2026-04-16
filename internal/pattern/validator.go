package pattern

import (
	"fmt"
	"math"
	"strings"

	"github.com/whitenhiemer/crochetbot/internal/models"
)

// Validator checks pattern quality and realism
type Validator struct {
	// Thresholds for validation
	MaxIncreasePercent float64 // Max stitch increase per round (%)
	MaxDecreasePercent float64 // Max stitch decrease per round (%)
	MinStitches        int     // Minimum stitches per round
	MaxStitches        int     // Maximum stitches per round
}

// NewValidator creates a validator with standard crochet constraints
func NewValidator() *Validator {
	return &Validator{
		MaxIncreasePercent: 30.0, // 30% max increase per round is aggressive but possible
		MaxDecreasePercent: 30.0, // 30% max decrease per round
		MinStitches:        6,    // Magic ring standard
		MaxStitches:        200,  // Very large amigurumi
	}
}

// ValidationResult contains all validation metrics
type ValidationResult struct {
	IsValid            bool
	Score              float64 // 0-100
	Issues             []ValidationIssue
	Warnings           []string
	StitchProgression  ProgressionMetrics
	TerminologyScore   float64
	StructuralScore    float64
	RealismScore       float64
	ComparisonToRef    *ComparisonMetrics // Only if reference provided
}

// ValidationIssue represents a specific problem
type ValidationIssue struct {
	Severity string // "error", "warning", "info"
	Location string // e.g., "Part Body, Round 5"
	Message  string
	Details  string
}

// ProgressionMetrics analyzes stitch count progression
type ProgressionMetrics struct {
	AverageIncreaseRate float64
	AverageDecreaseRate float64
	MaxJump             int
	MaxDrop             int
	UnrealisticChanges  int
	SmoothTransitions   int
	TotalTransitions    int
}

// ComparisonMetrics compares generated to reference pattern
type ComparisonMetrics struct {
	StructuralSimilarity float64 // 0-1
	LengthRatio          float64 // generated/reference
	StitchCountDrift     float64 // average difference
	ProgressionMatch     float64 // how similar is the shape progression
	TerminologyMatch     float64 // instruction format similarity
}

// ValidatePattern performs comprehensive validation
func (v *Validator) ValidatePattern(pattern *models.Pattern) ValidationResult {
	result := ValidationResult{
		IsValid:  true,
		Issues:   []ValidationIssue{},
		Warnings: []string{},
	}

	if len(pattern.Parts) == 0 {
		result.IsValid = false
		result.Issues = append(result.Issues, ValidationIssue{
			Severity: "error",
			Location: "Pattern",
			Message:  "No parts found",
			Details:  "Pattern must contain at least one part",
		})
		return result
	}

	// Validate each part
	for _, part := range pattern.Parts {
		v.validatePart(part, &result)
	}

	// Calculate progression metrics
	result.StitchProgression = v.calculateProgression(pattern)

	// Score terminology
	result.TerminologyScore = v.scoreTerminology(pattern)

	// Score structure
	result.StructuralScore = v.scoreStructure(pattern)

	// Score realism
	result.RealismScore = v.scoreRealism(pattern, result.StitchProgression)

	// Calculate overall score
	result.Score = (result.TerminologyScore*0.3 + result.StructuralScore*0.3 + result.RealismScore*0.4)

	return result
}

// validatePart checks a single part
func (v *Validator) validatePart(part models.Part, result *ValidationResult) {
	if len(part.Rounds) == 0 {
		result.Issues = append(result.Issues, ValidationIssue{
			Severity: "error",
			Location: fmt.Sprintf("Part %s", part.Name),
			Message:  "Part has no rounds",
			Details:  "Each part must have at least one round",
		})
		result.IsValid = false
		return
	}

	prevStitches := 0
	for i, round := range part.Rounds {
		location := fmt.Sprintf("Part %s, Round %d", part.Name, round.Number)

		// Check stitch count bounds
		if round.StitchCount < 0 {
			result.Issues = append(result.Issues, ValidationIssue{
				Severity: "error",
				Location: location,
				Message:  "Negative stitch count",
				Details:  fmt.Sprintf("Stitch count cannot be negative: %d", round.StitchCount),
			})
			result.IsValid = false
		}

		if round.StitchCount > v.MaxStitches {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("%s: Very large stitch count (%d) - may be difficult to work",
					location, round.StitchCount))
		}

		// Check progression from previous round
		if i > 0 && prevStitches > 0 && round.StitchCount > 0 {
			diff := round.StitchCount - prevStitches
			changePercent := (float64(diff) / float64(prevStitches)) * 100

			if changePercent > v.MaxIncreasePercent {
				result.Issues = append(result.Issues, ValidationIssue{
					Severity: "warning",
					Location: location,
					Message:  "Unrealistic increase rate",
					Details: fmt.Sprintf("%.1f%% increase (%d → %d) exceeds typical crochet constraints",
						changePercent, prevStitches, round.StitchCount),
				})
			}

			if changePercent < -v.MaxDecreasePercent {
				result.Issues = append(result.Issues, ValidationIssue{
					Severity: "warning",
					Location: location,
					Message:  "Unrealistic decrease rate",
					Details: fmt.Sprintf("%.1f%% decrease (%d → %d) may be difficult to execute",
						changePercent, prevStitches, round.StitchCount),
				})
			}

			// Check for stitch count not divisible by 6 on increase/decrease rounds
			if diff != 0 && round.StitchCount%6 != 0 && round.StitchType != "finish" {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("%s: Stitch count %d not divisible by 6 - may be irregular",
						location, round.StitchCount))
			}
		}

		// Check for missing instructions
		if round.Instructions == "" && round.StitchType == "" {
			result.Issues = append(result.Issues, ValidationIssue{
				Severity: "warning",
				Location: location,
				Message:  "Missing instructions",
				Details:  "Round has no instructions or stitch type",
			})
		}

		prevStitches = round.StitchCount
	}
}

// calculateProgression analyzes stitch count changes
func (v *Validator) calculateProgression(pattern *models.Pattern) ProgressionMetrics {
	metrics := ProgressionMetrics{}

	totalIncrease := 0.0
	totalDecrease := 0.0
	increaseCount := 0
	decreaseCount := 0

	for _, part := range pattern.Parts {
		for i := 1; i < len(part.Rounds); i++ {
			prev := part.Rounds[i-1].StitchCount
			curr := part.Rounds[i].StitchCount

			if prev == 0 || curr == 0 {
				continue
			}

			diff := curr - prev
			metrics.TotalTransitions++

			if diff > 0 {
				rate := (float64(diff) / float64(prev)) * 100
				totalIncrease += rate
				increaseCount++

				if diff > metrics.MaxJump {
					metrics.MaxJump = diff
				}

				if rate > v.MaxIncreasePercent {
					metrics.UnrealisticChanges++
				} else {
					metrics.SmoothTransitions++
				}
			} else if diff < 0 {
				rate := (float64(-diff) / float64(prev)) * 100
				totalDecrease += rate
				decreaseCount++

				if -diff > metrics.MaxDrop {
					metrics.MaxDrop = -diff
				}

				if rate > v.MaxDecreasePercent {
					metrics.UnrealisticChanges++
				} else {
					metrics.SmoothTransitions++
				}
			} else {
				metrics.SmoothTransitions++
			}
		}
	}

	if increaseCount > 0 {
		metrics.AverageIncreaseRate = totalIncrease / float64(increaseCount)
	}
	if decreaseCount > 0 {
		metrics.AverageDecreaseRate = totalDecrease / float64(decreaseCount)
	}

	return metrics
}

// scoreTerminology checks if instructions use standard terms
func (v *Validator) scoreTerminology(pattern *models.Pattern) float64 {
	standardTerms := []string{"sc", "hdc", "dc", "inc", "dec", "sl st", "ch", "magic ring", "repeat"}
	score := 100.0
	totalRounds := 0
	issuesFound := 0

	for _, part := range pattern.Parts {
		for _, round := range part.Rounds {
			totalRounds++
			instruction := strings.ToLower(round.Instructions)

			// Check if instruction contains at least one standard term
			hasStandardTerm := false
			for _, term := range standardTerms {
				if strings.Contains(instruction, term) {
					hasStandardTerm = true
					break
				}
			}

			if !hasStandardTerm && round.StitchType == "" && instruction != "" {
				issuesFound++
			}

			// Check for common formatting issues
			if strings.Contains(instruction, "increase evenly") ||
				strings.Contains(instruction, "decrease evenly") {
				issuesFound++
			}
		}
	}

	if totalRounds > 0 {
		errorRate := float64(issuesFound) / float64(totalRounds)
		score -= errorRate * 50 // Deduct up to 50 points
	}

	if score < 0 {
		score = 0
	}

	return score
}

// scoreStructure checks if pattern structure is sound
func (v *Validator) scoreStructure(pattern *models.Pattern) float64 {
	score := 100.0

	for _, part := range pattern.Parts {
		if len(part.Rounds) == 0 {
			score -= 20
			continue
		}

		// Check for proper start
		firstRound := part.Rounds[0]
		if firstRound.StitchCount != 6 && !strings.Contains(strings.ToLower(firstRound.Instructions), "magic") {
			score -= 5
		}

		// Check for proper close
		lastRound := part.Rounds[len(part.Rounds)-1]
		if part.Type == "sphere" || part.Type == "cylinder" {
			if lastRound.StitchCount > 12 && lastRound.StitchType != "finish" {
				score -= 5 // Should close smaller or have finish instruction
			}
		}

		// Check for balanced increases/decreases
		increases := 0
		decreases := 0
		for i := 1; i < len(part.Rounds); i++ {
			diff := part.Rounds[i].StitchCount - part.Rounds[i-1].StitchCount
			if diff > 0 {
				increases++
			} else if diff < 0 {
				decreases++
			}
		}

		// Spheres should have balanced inc/dec
		if part.Type == "sphere" {
			imbalance := math.Abs(float64(increases - decreases))
			if imbalance > float64(len(part.Rounds))/4 {
				score -= 10
			}
		}
	}

	if score < 0 {
		score = 0
	}

	return score
}

// scoreRealism evaluates how realistic the pattern is to crochet
func (v *Validator) scoreRealism(pattern *models.Pattern, prog ProgressionMetrics) float64 {
	score := 100.0

	// Penalize unrealistic changes
	if prog.TotalTransitions > 0 {
		unrealisticRatio := float64(prog.UnrealisticChanges) / float64(prog.TotalTransitions)
		score -= unrealisticRatio * 40
	}

	// Check average rates
	if prog.AverageIncreaseRate > v.MaxIncreasePercent {
		score -= 15
	}
	if prog.AverageDecreaseRate > v.MaxDecreasePercent {
		score -= 15
	}

	// Check max jumps
	if prog.MaxJump > 30 {
		score -= 10
	}
	if prog.MaxDrop > 30 {
		score -= 10
	}

	// Reward smooth transitions
	if prog.TotalTransitions > 0 {
		smoothRatio := float64(prog.SmoothTransitions) / float64(prog.TotalTransitions)
		if smoothRatio > 0.8 {
			score += 10
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// CompareToReference compares a generated pattern to a reference pattern
func (v *Validator) CompareToReference(generated, reference *models.Pattern) ComparisonMetrics {
	metrics := ComparisonMetrics{}

	if len(generated.Parts) == 0 || len(reference.Parts) == 0 {
		return metrics
	}

	// Compare primary parts (first part of each)
	genPart := generated.Parts[0]
	refPart := reference.Parts[0]

	// Length ratio
	metrics.LengthRatio = float64(len(genPart.Rounds)) / float64(len(refPart.Rounds))

	// Structural similarity (compare stitch progressions)
	metrics.StructuralSimilarity = v.compareStructure(genPart, refPart)

	// Stitch count drift
	metrics.StitchCountDrift = v.calculateDrift(genPart, refPart)

	// Progression match
	metrics.ProgressionMatch = v.compareProgression(genPart, refPart)

	// Terminology match
	metrics.TerminologyMatch = v.compareTerminology(genPart, refPart)

	return metrics
}

// compareStructure compares the overall structure of two parts
func (v *Validator) compareStructure(gen, ref models.Part) float64 {
	// Normalize round counts for comparison
	genProfile := v.normalizeStitchProfile(gen)
	refProfile := v.normalizeStitchProfile(ref)

	// Resample to same length for comparison
	targetLen := 100
	genResampled := v.resampleProfile(genProfile, targetLen)
	refResampled := v.resampleProfile(refProfile, targetLen)

	// Calculate similarity (1 - normalized error)
	totalError := 0.0
	for i := 0; i < targetLen; i++ {
		error := math.Abs(genResampled[i] - refResampled[i])
		totalError += error
	}

	avgError := totalError / float64(targetLen)
	similarity := 1.0 - avgError

	if similarity < 0 {
		similarity = 0
	}

	return similarity
}

// normalizeStitchProfile extracts normalized stitch counts (0-1 range)
func (v *Validator) normalizeStitchProfile(part models.Part) []float64 {
	profile := []float64{}
	maxStitches := 0

	for _, round := range part.Rounds {
		if round.StitchCount > maxStitches {
			maxStitches = round.StitchCount
		}
	}

	if maxStitches == 0 {
		return profile
	}

	for _, round := range part.Rounds {
		normalized := float64(round.StitchCount) / float64(maxStitches)
		profile = append(profile, normalized)
	}

	return profile
}

// resampleProfile resamples a profile to target length
func (v *Validator) resampleProfile(profile []float64, targetLen int) []float64 {
	if len(profile) == 0 {
		return make([]float64, targetLen)
	}

	resampled := make([]float64, targetLen)
	for i := 0; i < targetLen; i++ {
		srcPos := float64(i) * float64(len(profile)-1) / float64(targetLen-1)
		srcIdx := int(srcPos)
		fraction := srcPos - float64(srcIdx)

		if srcIdx >= len(profile)-1 {
			resampled[i] = profile[len(profile)-1]
		} else {
			// Linear interpolation
			resampled[i] = profile[srcIdx]*(1-fraction) + profile[srcIdx+1]*fraction
		}
	}

	return resampled
}

// calculateDrift calculates average stitch count difference
func (v *Validator) calculateDrift(gen, ref models.Part) float64 {
	genProfile := v.normalizeStitchProfile(gen)
	refProfile := v.normalizeStitchProfile(ref)

	if len(genProfile) == 0 || len(refProfile) == 0 {
		return 0
	}

	// Compare at matching proportions
	numSamples := 20
	totalDrift := 0.0

	for i := 0; i < numSamples; i++ {
		proportion := float64(i) / float64(numSamples-1)
		genIdx := int(proportion * float64(len(genProfile)-1))
		refIdx := int(proportion * float64(len(refProfile)-1))

		drift := math.Abs(genProfile[genIdx] - refProfile[refIdx])
		totalDrift += drift
	}

	return totalDrift / float64(numSamples)
}

// compareProgression compares shape progression patterns
func (v *Validator) compareProgression(gen, ref models.Part) float64 {
	genChanges := v.extractChangePattern(gen)
	refChanges := v.extractChangePattern(ref)

	if len(genChanges) == 0 || len(refChanges) == 0 {
		return 0
	}

	// Compare patterns: both increasing, both decreasing, both flat
	matches := 0
	comparisons := 0

	minLen := len(genChanges)
	if len(refChanges) < minLen {
		minLen = len(refChanges)
	}

	for i := 0; i < minLen; i++ {
		if genChanges[i] == refChanges[i] {
			matches++
		}
		comparisons++
	}

	if comparisons == 0 {
		return 0
	}

	return float64(matches) / float64(comparisons)
}

// extractChangePattern creates a pattern of changes (increase/decrease/flat)
func (v *Validator) extractChangePattern(part models.Part) []int {
	pattern := []int{}

	for i := 1; i < len(part.Rounds); i++ {
		diff := part.Rounds[i].StitchCount - part.Rounds[i-1].StitchCount
		if diff > 0 {
			pattern = append(pattern, 1) // increase
		} else if diff < 0 {
			pattern = append(pattern, -1) // decrease
		} else {
			pattern = append(pattern, 0) // flat
		}
	}

	return pattern
}

// compareTerminology checks how similar instruction formatting is
func (v *Validator) compareTerminology(gen, ref models.Part) float64 {
	genTerms := v.extractTerminology(gen)
	refTerms := v.extractTerminology(ref)

	// Count overlap
	overlap := 0
	for term := range genTerms {
		if refTerms[term] {
			overlap++
		}
	}

	totalTerms := len(genTerms)
	if len(refTerms) > totalTerms {
		totalTerms = len(refTerms)
	}

	if totalTerms == 0 {
		return 0
	}

	return float64(overlap) / float64(totalTerms)
}

// extractTerminology extracts terminology used in a part
func (v *Validator) extractTerminology(part models.Part) map[string]bool {
	terms := make(map[string]bool)
	keywords := []string{"sc", "hdc", "dc", "inc", "dec", "sl st", "ch", "repeat", "magic ring"}

	for _, round := range part.Rounds {
		instruction := strings.ToLower(round.Instructions)
		for _, keyword := range keywords {
			if strings.Contains(instruction, keyword) {
				terms[keyword] = true
			}
		}
	}

	return terms
}
