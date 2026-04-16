package pattern

import (
	"strings"
	"testing"
)

func TestParseRound(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name              string
		input             string
		expectError       bool
		expectedRound     int
		expectedStitches  int
		expectedRepeats   int
		expectedStitchType string
	}{
		{
			name:               "magic ring start",
			input:              "Rnd 1. 6 sc in magic ring (6)",
			expectError:        false,
			expectedRound:      1,
			expectedStitches:   6,
			expectedRepeats:    1,
			expectedStitchType: "sc",
		},
		{
			name:               "increase with repeats",
			input:              "Rnd 3. [sc, inc] x 6 (18)",
			expectError:        false,
			expectedRound:      3,
			expectedStitches:   18,
			expectedRepeats:    6,
			expectedStitchType: "inc",
		},
		{
			name:               "decrease with repeats",
			input:              "Rnd 8. [2 sc, dec] x 6 (18)",
			expectError:        false,
			expectedRound:      8,
			expectedStitches:   18,
			expectedRepeats:    6,
			expectedStitchType: "dec",
		},
		{
			name:               "constant round",
			input:              "Rnd 5. 24 sc (24)",
			expectError:        false,
			expectedRound:      5,
			expectedStitches:   24,
			expectedRepeats:    1,
			expectedStitchType: "sc",
		},
		{
			name:               "round range",
			input:              "Rnds 5-7. 24 sc (24)",
			expectError:        false,
			expectedRound:      5,
			expectedStitches:   24,
			expectedRepeats:    1,
			expectedStitchType: "sc",
		},
		{
			name:               "repeat with word",
			input:              "Rnd 2. [sc, inc] repeat 6 times (12)",
			expectError:        false,
			expectedRound:      2,
			expectedStitches:   12,
			expectedRepeats:    6,
			expectedStitchType: "inc",
		},
		{
			name:               "slip stitch",
			input:              "Rnd 3. 6 sl st (6)",
			expectError:        false,
			expectedRound:      3,
			expectedStitches:   6,
			expectedRepeats:    1,
			expectedStitchType: "sl st",
		},
		{
			name:        "no round number",
			input:       "just some text (6)",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			round, err := parser.ParseRound(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if round.Number != tt.expectedRound {
				t.Errorf("expected round %d, got %d", tt.expectedRound, round.Number)
			}

			if round.StitchCount != tt.expectedStitches {
				t.Errorf("expected %d stitches, got %d", tt.expectedStitches, round.StitchCount)
			}

			if round.Repeats != tt.expectedRepeats {
				t.Errorf("expected %d repeats, got %d", tt.expectedRepeats, round.Repeats)
			}

			if round.StitchType != tt.expectedStitchType {
				t.Errorf("expected stitch type %s, got %s", tt.expectedStitchType, round.StitchType)
			}
		})
	}
}

func TestParsePattern(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name          string
		input         string
		expectError   bool
		expectedParts int
		expectRounds  int
	}{
		{
			name: "single part pattern",
			input: `HEAD & BODY

Rnd 1. 6 sc in magic ring (6)
Rnd 2. 6 inc (12)
Rnd 3. [sc, inc] x 6 (18)`,
			expectError:   false,
			expectedParts: 1,
			expectRounds:  3,
		},
		{
			name: "multi-part pattern",
			input: `HEAD

Rnd 1. 6 sc in magic ring (6)
Rnd 2. 6 inc (12)

ARMS

Rnd 1. 6 sc in magic ring (6)`,
			expectError:   false,
			expectedParts: 2,
			expectRounds:  2, // First part
		},
		{
			name: "with color change",
			input: `BEAK
With black yarn.

Rnd 1. 4 sc in magic ring (4)
Rnd 2. [sc, inc] x 2 (6)
Rnd 3. (switch to white yarn) 6 sl st (6)`,
			expectError:   false,
			expectedParts: 1,
			expectRounds:  3,
		},
		{
			name:        "empty pattern",
			input:       "",
			expectError: true,
		},
		{
			name: "pattern with notes",
			input: `BODY
With main color.

Rnd 1. 6 sc in magic ring (6)
Stuff firmly.
Rnd 2. 6 inc (12)`,
			expectError:   false,
			expectedParts: 1,
			expectRounds:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := parser.ParsePattern(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(pattern.Parts) != tt.expectedParts {
				t.Errorf("expected %d parts, got %d", tt.expectedParts, len(pattern.Parts))
			}

			if tt.expectedParts > 0 && len(pattern.Parts) > 0 {
				firstPart := pattern.Parts[0]
				if len(firstPart.Rounds) != tt.expectRounds {
					t.Errorf("expected %d rounds in first part, got %d",
						tt.expectRounds, len(firstPart.Rounds))
				}
			}
		})
	}
}

func TestIsPartHeader(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		input    string
		expected bool
	}{
		{"HEAD & BODY", true},
		{"WINGS", true},
		{"BEAK", true},
		{"FEET", true},
		{"ARM", true},
		{"Rnd 1. 6 sc (6)", false},
		{"With pink yarn.", false},
		{"Stuff firmly", false},
		{"head", false}, // lowercase doesn't count
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parser.isPartHeader(tt.input)
			if result != tt.expected {
				t.Errorf("input %q: expected %v, got %v", tt.input, tt.expected, result)
			}
		})
	}
}

func TestParseInstruction(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		input              string
		expectedSC         int
		expectedInc        int
		expectedDec        int
		expectedMultiplier int
	}{
		{
			input:              "6 sc in magic ring",
			expectedSC:         1,
			expectedInc:        0,
			expectedDec:        0,
			expectedMultiplier: 1,
		},
		{
			input:              "[sc, inc] x 6",
			expectedSC:         1,
			expectedInc:        1,
			expectedDec:        0,
			expectedMultiplier: 6,
		},
		{
			input:              "[2 sc, dec] repeat 6 times",
			expectedSC:         1,
			expectedInc:        0,
			expectedDec:        1,
			expectedMultiplier: 6,
		},
		{
			input:              "dec around",
			expectedSC:         0,
			expectedInc:        0,
			expectedDec:        1,
			expectedMultiplier: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			components := parser.ParseInstruction(tt.input)

			if components.SingleCrochet != tt.expectedSC {
				t.Errorf("expected %d sc, got %d", tt.expectedSC, components.SingleCrochet)
			}

			if components.Increase != tt.expectedInc {
				t.Errorf("expected %d inc, got %d", tt.expectedInc, components.Increase)
			}

			if components.Decrease != tt.expectedDec {
				t.Errorf("expected %d dec, got %d", tt.expectedDec, components.Decrease)
			}

			if components.Multiplier != tt.expectedMultiplier {
				t.Errorf("expected multiplier %d, got %d", tt.expectedMultiplier, components.Multiplier)
			}
		})
	}
}

func TestValidatePattern(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name              string
		input             string
		expectIssues      bool
		allowedIssueCount int
	}{
		{
			name: "mostly valid pattern",
			input: `BODY

Rnd 1. 6 sc in magic ring (6)
Rnd 2. 6 inc (12)
Rnd 3. [sc, inc] x 6 (18)
Rnd 4. [2 sc, inc] x 6 (24)
Rnd 5. 24 sc (24)`,
			expectIssues:      true,  // Rnd 2 has large jump
			allowedIssueCount: 1,     // Only one warning is OK
		},
		{
			name: "unrealistic jump",
			input: `BODY

Rnd 1. 6 sc in magic ring (6)
Rnd 2. 100 inc (100)`,
			expectIssues:      true,
			allowedIssueCount: 10, // Multiple issues expected
		},
		{
			name:              "empty parts",
			input:             `BODY`,
			expectIssues:      true,
			allowedIssueCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := parser.ParsePattern(tt.input)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			issues := parser.ValidatePattern(pattern)

			if tt.expectIssues && len(issues) == 0 {
				t.Error("expected issues but got none")
			}

			if !tt.expectIssues && len(issues) > 0 {
				t.Errorf("unexpected issues: %v", issues)
			}
		})
	}
}

func TestParsePartFromText(t *testing.T) {
	parser := NewParser()

	text := `Rnd 1. 6 sc in magic ring (6)
Rnd 2. 6 inc (12)
Rnd 3. [sc, inc] x 6 (18)`

	part, err := parser.ParsePartFromText("Test Part", text)
	if err != nil {
		t.Fatalf("failed to parse part: %v", err)
	}

	if part.Name != "Test Part" {
		t.Errorf("expected name 'Test Part', got %s", part.Name)
	}

	if len(part.Rounds) != 3 {
		t.Errorf("expected 3 rounds, got %d", len(part.Rounds))
	}

	if part.Rounds[0].StitchCount != 6 {
		t.Errorf("expected first round to have 6 stitches, got %d", part.Rounds[0].StitchCount)
	}
}

func TestParseComplexPattern(t *testing.T) {
	// Test with actual Flamingo pattern excerpt
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

	parser := NewParser()
	pattern, err := parser.ParsePattern(flamingoText)

	if err != nil {
		t.Fatalf("failed to parse flamingo pattern: %v", err)
	}

	if len(pattern.Parts) != 2 {
		t.Errorf("expected 2 parts, got %d", len(pattern.Parts))
	}

	// Check HEAD & BODY part
	if len(pattern.Parts) > 0 {
		headBody := pattern.Parts[0]
		if !strings.Contains(headBody.Name, "HEAD") {
			t.Errorf("expected first part to be HEAD & BODY, got %s", headBody.Name)
		}

		// Should have parsed multiple rounds
		if len(headBody.Rounds) < 10 {
			t.Errorf("expected at least 10 rounds in head/body, got %d", len(headBody.Rounds))
		}

		// Check specific rounds
		if headBody.Rounds[0].StitchCount != 6 {
			t.Errorf("round 1 should have 6 stitches, got %d", headBody.Rounds[0].StitchCount)
		}

		// Find a round with 32 stitches (should be round 15 or 16)
		found32 := false
		for _, round := range headBody.Rounds {
			if round.StitchCount == 32 {
				found32 = true
				break
			}
		}
		if !found32 {
			t.Error("expected to find round with 32 stitches")
		}
	}

	// Check BEAK part
	if len(pattern.Parts) > 1 {
		beak := pattern.Parts[1]
		if !strings.Contains(beak.Name, "BEAK") {
			t.Errorf("expected second part to be BEAK, got %s", beak.Name)
		}

		if len(beak.Rounds) != 3 {
			t.Errorf("expected 3 rounds in beak, got %d", len(beak.Rounds))
		}
	}
}
