package main

import (
	"fmt"
	"os"

	"github.com/whitenhiemer/crochetbot/internal/pattern"
)

func main() {
	// Test data: Zoe the Flamingo Head & Body section
	flamingoPattern := `
HEAD & BODY
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
`

	parser := pattern.NewParser()
	formatter := pattern.NewFormatter()
	validator := pattern.NewValidator()

	fmt.Println("=== TESTING PATTERN PARSER ===\n")

	// Test 1: Parse individual rounds
	fmt.Println("Test 1: Parse individual rounds")
	testRounds := []string{
		"Rnd 1. 6 sc in magic ring (6)",
		"Rnd 3. [sc, inc] x 6 (18)",
		"Rnds 5-7. 24 sc (24)",
		"Rnd 8. [2 sc, dec] x 6 (18)",
	}

	for _, line := range testRounds {
		round, err := parser.ParseRound(line)
		if err != nil {
			fmt.Printf("  ✗ Failed to parse: %s\n    Error: %v\n", line, err)
		} else {
			fmt.Printf("  ✓ Rnd %d: %d stitches, type=%s, repeats=%d\n",
				round.Number, round.StitchCount, round.StitchType, round.Repeats)
		}
	}

	// Test 2: Parse complete pattern
	fmt.Println("\nTest 2: Parse complete pattern")
	parsedPattern, err := parser.ParsePattern(flamingoPattern)
	if err != nil {
		fmt.Printf("  ✗ Parse error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("  ✓ Parsed %d part(s)\n", len(parsedPattern.Parts))
	if len(parsedPattern.Parts) > 0 {
		part := parsedPattern.Parts[0]
		fmt.Printf("    Part: %s\n", part.Name)
		fmt.Printf("    Rounds: %d\n", len(part.Rounds))
		fmt.Printf("    First: Rnd %d (%d st)\n", part.Rounds[0].Number, part.Rounds[0].StitchCount)
		fmt.Printf("    Last: Rnd %d (%d st)\n",
			part.Rounds[len(part.Rounds)-1].Number,
			part.Rounds[len(part.Rounds)-1].StitchCount)
	}

	// Test 3: Format pattern back to text
	fmt.Println("\nTest 3: Format pattern to text")
	formatted := formatter.FormatPattern(parsedPattern)
	fmt.Println("  ✓ Formatted output:")
	fmt.Println("---")
	fmt.Println(formatted)
	fmt.Println("---")

	// Test 4: Validate pattern
	fmt.Println("\nTest 4: Validate pattern")
	result := validator.ValidatePattern(parsedPattern)
	fmt.Printf("  Valid: %v\n", result.IsValid)
	fmt.Printf("  Score: %.1f/100\n", result.Score)
	fmt.Printf("  Terminology: %.1f/100\n", result.TerminologyScore)
	fmt.Printf("  Structure: %.1f/100\n", result.StructuralScore)
	fmt.Printf("  Realism: %.1f/100\n", result.RealismScore)

	if len(result.Issues) > 0 {
		fmt.Printf("\n  Issues found: %d\n", len(result.Issues))
		for _, issue := range result.Issues {
			fmt.Printf("    [%s] %s: %s\n", issue.Severity, issue.Location, issue.Message)
		}
	} else {
		fmt.Println("  ✓ No issues found")
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("\n  Warnings: %d\n", len(result.Warnings))
		for _, warning := range result.Warnings {
			fmt.Printf("    - %s\n", warning)
		}
	}

	// Test 5: Progression metrics
	fmt.Println("\nTest 5: Stitch progression analysis")
	prog := result.StitchProgression
	fmt.Printf("  Avg increase rate: %.1f%%\n", prog.AverageIncreaseRate)
	fmt.Printf("  Avg decrease rate: %.1f%%\n", prog.AverageDecreaseRate)
	fmt.Printf("  Max jump: %d stitches\n", prog.MaxJump)
	fmt.Printf("  Max drop: %d stitches\n", prog.MaxDrop)
	fmt.Printf("  Smooth transitions: %d/%d\n", prog.SmoothTransitions, prog.TotalTransitions)
	fmt.Printf("  Unrealistic changes: %d\n", prog.UnrealisticChanges)

	fmt.Println("\n=== ALL TESTS COMPLETE ===")
}
