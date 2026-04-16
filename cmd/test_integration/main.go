package main

import (
	"fmt"
	"os"

	"github.com/whitenhiemer/crochetbot/internal/pattern"
)

func main() {
	fmt.Println("=== CROCHET PATTERN INTEGRATION TEST ===\n")

	// Zoe the Flamingo sample pattern
	flamingoText := `HEAD & BODY
With pink yarn.

Rnd 1. start 6 sc in a magic loop (6)
Rnd 2. 6 inc (12)
Rnd 3. [sc, inc] x 6 (18)
Rnd 4. [2 sc, inc] x 6 (24)
Rnds 5-7. 24 sc (24)
Rnd 8. [2 sc, dec] x 6 (18)
Rnd 9. [4 sc, dec] x 3 (15)
Rnds 10-11. 15 sc (15)
Rnd 12. [4 sc, inc] x 3 (18)
Rnd 13. [2 sc, inc] x 6 (24)
Rnd 14. [5 sc, inc] x 4 (28)
Rnd 15. [sc, inc] x 2, 20 sc, [sc, inc] x 2 (32)
Rnd 16. 32 sc (32)
Rnd 17. [2 sc, dec] x 8 (24)
Rnd 18. 24 sc (24)
Rnd 19. [sc, dec] x 8 (16)
Rnd 20. 8 dec (8)

BEAK
With black and white yarn.

Rnd 1. (black yarn) start 4 sc in a magic loop (4)
Rnd 2. [sc, inc] x 2 (6)
Rnd 3. (switch to white yarn) 6 sl st (6)`

	// Step 1: Parse
	fmt.Println("Step 1: Parse text pattern")
	parser := pattern.NewParser()
	parsedPattern, err := parser.ParsePattern(flamingoText)
	if err != nil {
		fmt.Printf("❌ Parse failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Parsed %d parts\n", len(parsedPattern.Parts))
	for i, part := range parsedPattern.Parts {
		fmt.Printf("  Part %d: %s (%d rounds)\n", i+1, part.Name, len(part.Rounds))
	}
	fmt.Println()

	// Step 2: Validate
	fmt.Println("Step 2: Validate pattern quality")
	validator := pattern.NewValidator()
	validation := validator.ValidatePattern(parsedPattern)
	fmt.Printf("✓ Overall Score: %.1f/100\n", validation.Score)
	fmt.Printf("  - Terminology: %.1f/100\n", validation.TerminologyScore)
	fmt.Printf("  - Structure: %.1f/100\n", validation.StructuralScore)
	fmt.Printf("  - Realism: %.1f/100\n", validation.RealismScore)
	fmt.Printf("  - Issues: %d\n", len(validation.Issues))
	fmt.Printf("  - Warnings: %d\n", len(validation.Warnings))

	if len(validation.Issues) > 0 {
		fmt.Println("\n  Top issues:")
		for i, issue := range validation.Issues {
			if i >= 3 {
				break
			}
			fmt.Printf("    [%s] %s: %s\n", issue.Severity, issue.Location, issue.Message)
		}
	}
	fmt.Println()

	// Step 3: Format
	fmt.Println("Step 3: Format pattern to text")
	formatter := pattern.NewFormatter()
	formatted := formatter.FormatPattern(parsedPattern)
	fmt.Println("✓ Formatted pattern:")
	fmt.Println("---")
	// Print first 30 lines
	lines := 0
	for _, ch := range formatted {
		if ch == '\n' {
			lines++
			if lines >= 30 {
				fmt.Println("... (truncated)")
				break
			}
		}
		fmt.Printf("%c", ch)
	}
	fmt.Println("---\n")

	// Step 4: Progression Analysis
	fmt.Println("Step 4: Analyze stitch progression")
	prog := validation.StitchProgression
	fmt.Printf("  Avg increase rate: %.1f%%\n", prog.AverageIncreaseRate)
	fmt.Printf("  Avg decrease rate: %.1f%%\n", prog.AverageDecreaseRate)
	fmt.Printf("  Max jump: %d stitches\n", prog.MaxJump)
	fmt.Printf("  Max drop: %d stitches\n", prog.MaxDrop)
	fmt.Printf("  Smooth transitions: %d/%d (%.1f%%)\n",
		prog.SmoothTransitions, prog.TotalTransitions,
		float64(prog.SmoothTransitions)/float64(prog.TotalTransitions)*100)
	fmt.Println()

	// Step 5: Compare to itself (baseline)
	fmt.Println("Step 5: Compare pattern to itself (baseline)")
	comparison := validator.CompareToReference(parsedPattern, parsedPattern)
	fmt.Printf("  Structural similarity: %.2f (should be 1.0)\n", comparison.StructuralSimilarity)
	fmt.Printf("  Length ratio: %.2f (should be 1.0)\n", comparison.LengthRatio)
	fmt.Printf("  Stitch count drift: %.4f (should be ~0.0)\n", comparison.StitchCountDrift)
	fmt.Printf("  Progression match: %.2f (should be 1.0)\n", comparison.ProgressionMatch)
	fmt.Printf("  Terminology match: %.2f (should be 1.0)\n", comparison.TerminologyMatch)
	fmt.Println()

	// Step 6: Round-trip test
	fmt.Println("Step 6: Round-trip test (parse → format → parse)")
	reparsed, err := parser.ParsePattern(formatted)
	if err != nil {
		fmt.Printf("❌ Re-parse failed: %v\n", err)
	} else {
		fmt.Printf("✓ Successfully re-parsed\n")
		fmt.Printf("  Original parts: %d, Re-parsed parts: %d\n",
			len(parsedPattern.Parts), len(reparsed.Parts))

		if len(parsedPattern.Parts) > 0 && len(reparsed.Parts) > 0 {
			origRounds := len(parsedPattern.Parts[0].Rounds)
			reparsedRounds := len(reparsed.Parts[0].Rounds)
			fmt.Printf("  Original rounds (part 1): %d, Re-parsed: %d\n",
				origRounds, reparsedRounds)

			if origRounds == reparsedRounds {
				fmt.Println("  ✓ Round count preserved")
			} else {
				fmt.Println("  ⚠ Round count changed")
			}
		}

		// Compare round-trip similarity
		rtComparison := validator.CompareToReference(reparsed, parsedPattern)
		fmt.Printf("  Round-trip similarity: %.2f\n", rtComparison.StructuralSimilarity)
		if rtComparison.StructuralSimilarity > 0.95 {
			fmt.Println("  ✓ High fidelity round-trip")
		} else {
			fmt.Println("  ⚠ Some information lost in round-trip")
		}
	}

	fmt.Println("\n=== ALL TESTS COMPLETE ===")
	fmt.Println("\nSummary:")
	fmt.Printf("  Pattern: Zoe the Flamingo\n")
	fmt.Printf("  Parts: %d\n", len(parsedPattern.Parts))
	fmt.Printf("  Total rounds: %d\n",
		len(parsedPattern.Parts[0].Rounds)+len(parsedPattern.Parts[1].Rounds))
	fmt.Printf("  Quality score: %.1f/100\n", validation.Score)
	fmt.Printf("  Validation: %s\n", func() string {
		if validation.IsValid {
			return "✓ PASS"
		}
		return "✗ FAIL"
	}())
}
