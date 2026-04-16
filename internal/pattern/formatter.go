package pattern

import (
	"fmt"
	"strings"

	"github.com/whitenhiemer/crochetbot/internal/models"
)

// Formatter converts Pattern structs to human-readable text
type Formatter struct {
	// Configuration
	IncludeNotes      bool
	IncludeTips       bool
	CompactMode       bool
	AbbreviationStyle string // "US" or "UK"
}

// NewFormatter creates a new pattern formatter
func NewFormatter() *Formatter {
	return &Formatter{
		IncludeNotes:      true,
		IncludeTips:       true,
		CompactMode:       false,
		AbbreviationStyle: "US",
	}
}

// FormatPattern converts a Pattern to Woobles-style text
func (f *Formatter) FormatPattern(pattern *models.Pattern) string {
	var sb strings.Builder

	// Header
	sb.WriteString(strings.ToUpper(pattern.Name))
	sb.WriteString("\n\n")

	if pattern.Description != "" {
		sb.WriteString(pattern.Description)
		sb.WriteString("\n\n")
	}

	// Materials
	if len(pattern.Materials.Colors) > 0 {
		sb.WriteString("MATERIALS\n")
		sb.WriteString(fmt.Sprintf("Hook size: %s\n", pattern.Materials.HookSize))
		sb.WriteString(fmt.Sprintf("Yarn weight: %s\n", pattern.Materials.YarnWeight))
		for _, color := range pattern.Materials.Colors {
			sb.WriteString(fmt.Sprintf("- %s (%d yds)\n", color.Name, color.Amount))
		}
		if len(pattern.Materials.OtherSupplies) > 0 {
			sb.WriteString("Other supplies: ")
			sb.WriteString(strings.Join(pattern.Materials.OtherSupplies, ", "))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Finished size
	if pattern.FinishedSize.HeightInches > 0 {
		sb.WriteString("FINISHED SIZE\n")
		sb.WriteString(fmt.Sprintf("Height: %.1f inches (%.1f cm)\n",
			pattern.FinishedSize.HeightInches, pattern.FinishedSize.HeightCm))
		sb.WriteString(fmt.Sprintf("Width: %.1f inches (%.1f cm)\n\n",
			pattern.FinishedSize.WidthInches, pattern.FinishedSize.WidthCm))
	}

	// Abbreviations
	sb.WriteString("ABBREVIATIONS\n")
	sb.WriteString("ch   = chain\n")
	sb.WriteString("sc   = single crochet\n")
	sb.WriteString("hdc  = half double crochet\n")
	sb.WriteString("dc   = double crochet\n")
	sb.WriteString("inc  = increase (2 sc in same stitch)\n")
	sb.WriteString("dec  = (invisible) decrease\n")
	sb.WriteString("sl st = slip stitch\n")
	sb.WriteString("rnd  = round\n\n")

	// Pattern parts
	sb.WriteString("PATTERN\n\n")

	for partIdx, part := range pattern.Parts {
		if partIdx > 0 {
			sb.WriteString("\n")
		}

		// Part header
		sb.WriteString(strings.ToUpper(part.Name))
		sb.WriteString("\n")
		if part.Color != "" && part.Color != "main color" {
			sb.WriteString(fmt.Sprintf("With %s.\n", part.Color))
		} else {
			sb.WriteString("With main color.\n")
		}
		sb.WriteString("\n")

		// Part notes
		if f.IncludeNotes && len(part.Notes) > 0 {
			for _, note := range part.Notes {
				sb.WriteString(fmt.Sprintf("NOTE: %s\n", note))
			}
			sb.WriteString("\n")
		}

		// Rounds
		consecutiveCount := 0
		var consecutiveRounds []models.Round

		for i, round := range part.Rounds {
			// Check if this is part of consecutive identical rounds
			if i > 0 && f.canGroupRounds(part.Rounds[i-1], round) {
				if consecutiveCount == 0 {
					consecutiveRounds = []models.Round{part.Rounds[i-1]}
					consecutiveCount = 1
				}
				consecutiveRounds = append(consecutiveRounds, round)
				consecutiveCount++

				// If this is the last round, flush the group
				if i == len(part.Rounds)-1 {
					sb.WriteString(f.formatRoundGroup(consecutiveRounds))
				}
				continue
			}

			// Flush any accumulated consecutive rounds
			if consecutiveCount > 0 {
				sb.WriteString(f.formatRoundGroup(consecutiveRounds))
				consecutiveRounds = nil
				consecutiveCount = 0
			}

			// Format single round
			sb.WriteString(f.FormatRound(round))
		}

		// Add stuffing note if applicable
		if f.shouldAddStuffingNote(part) {
			sb.WriteString("\nStuff piece firmly, shaping as you go.\n")
		}
	}

	// Assembly instructions
	if len(pattern.Assembly) > 0 {
		sb.WriteString("\n")
		sb.WriteString("ASSEMBLY\n")
		for i, instruction := range pattern.Assembly {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, instruction))
		}
	}

	// Accuracy metrics
	if pattern.AccuracyMetrics.ShapeMatchPercent > 0 {
		sb.WriteString("\n")
		sb.WriteString("PATTERN ACCURACY\n")
		sb.WriteString(fmt.Sprintf("Shape match: %.1f%%\n", pattern.AccuracyMetrics.ShapeMatchPercent))
		if pattern.AccuracyMetrics.Notes != "" {
			sb.WriteString(fmt.Sprintf("%s\n", pattern.AccuracyMetrics.Notes))
		}
	}

	return sb.String()
}

// FormatRound converts a single round to text
func (f *Formatter) FormatRound(round models.Round) string {
	var sb strings.Builder

	// Round number
	sb.WriteString(fmt.Sprintf("Rnd %d. ", round.Number))

	// Format instruction in Woobles style
	instruction := f.formatInstruction(round)
	sb.WriteString(instruction)

	// Add stitch count
	if round.StitchCount > 0 {
		sb.WriteString(fmt.Sprintf(" (%d)", round.StitchCount))
	}

	sb.WriteString("\n")

	// Add notes if present
	if f.IncludeNotes && round.Notes != "" && !strings.Contains(round.Notes, "repeated") {
		sb.WriteString(fmt.Sprintf("   (%s)\n", round.Notes))
	}

	return sb.String()
}

// formatRoundGroup formats consecutive identical rounds
func (f *Formatter) formatRoundGroup(rounds []models.Round) string {
	if len(rounds) == 0 {
		return ""
	}

	if len(rounds) == 1 {
		return f.FormatRound(rounds[0])
	}

	// Format as range: "Rnds 5-7. 24 sc (24)"
	first := rounds[0]
	last := rounds[len(rounds)-1]

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Rnds %d-%d. ", first.Number, last.Number))
	sb.WriteString(f.formatInstruction(first))

	if first.StitchCount > 0 {
		sb.WriteString(fmt.Sprintf(" (%d)", first.StitchCount))
	}

	sb.WriteString("\n")

	return sb.String()
}

// formatInstruction formats the instruction part of a round
func (f *Formatter) formatInstruction(round models.Round) string {
	// If instruction is already well-formatted, use it
	if round.Instructions != "" && !strings.Contains(round.Instructions, "increase evenly") && !strings.Contains(round.Instructions, "decrease evenly") {
		return round.Instructions
	}

	// Generate instruction from data
	switch round.StitchType {
	case "finish":
		return "Fasten off, leaving long tail. Close opening with yarn needle."

	case "inc":
		if round.Repeats > 1 {
			// Calculate sc between increases
			scBetween := (round.StitchCount / round.Repeats) - 2
			if scBetween <= 0 {
				return "2 sc in each st around"
			}
			return fmt.Sprintf("[inc, %d sc] repeat %d times", scBetween, round.Repeats)
		}
		return fmt.Sprintf("%d inc", round.StitchCount)

	case "dec":
		if round.Repeats > 1 {
			scBetween := (round.StitchCount / round.Repeats)
			if scBetween <= 0 {
				return "dec around"
			}
			return fmt.Sprintf("[dec, %d sc] repeat %d times", scBetween, round.Repeats)
		}
		return fmt.Sprintf("%d dec", round.StitchCount)

	case "sc":
		if round.StitchCount > 0 {
			return fmt.Sprintf("sc in each st around (%d sc)", round.StitchCount)
		}
		return "sc around"

	default:
		// Use the raw instruction if available
		if round.Instructions != "" {
			return round.Instructions
		}
		return fmt.Sprintf("%d %s", round.StitchCount, round.StitchType)
	}
}

// canGroupRounds checks if two rounds can be grouped together
func (f *Formatter) canGroupRounds(r1, r2 models.Round) bool {
	// Must be consecutive
	if r2.Number != r1.Number+1 {
		return false
	}

	// Must have same stitch count
	if r1.StitchCount != r2.StitchCount {
		return false
	}

	// Must have same instruction pattern
	if r1.Instructions != r2.Instructions {
		return false
	}

	// Don't group special rounds
	if r1.StitchType == "finish" || r2.StitchType == "finish" {
		return false
	}

	return true
}

// shouldAddStuffingNote determines if we should add a stuffing reminder
func (f *Formatter) shouldAddStuffingNote(part models.Part) bool {
	if !f.IncludeNotes {
		return false
	}

	// Check if part has decreases (meaning it closes)
	hasDecreases := false
	for _, round := range part.Rounds {
		if round.StitchType == "dec" || strings.Contains(strings.ToLower(round.Instructions), "dec") {
			hasDecreases = true
			break
		}
	}

	// Only add stuffing note for 3D parts that close
	return hasDecreases && (part.Type == "sphere" || part.Type == "cylinder")
}

// FormatCompact creates a compact, single-line representation
func (f *Formatter) FormatCompact(pattern *models.Pattern) string {
	var parts []string
	for _, part := range pattern.Parts {
		parts = append(parts, fmt.Sprintf("%s (%d rnds)", part.Name, len(part.Rounds)))
	}
	return fmt.Sprintf("%s: %s", pattern.Name, strings.Join(parts, ", "))
}

// FormatPartSummary creates a summary of a part
func (f *Formatter) FormatPartSummary(part models.Part) string {
	if len(part.Rounds) == 0 {
		return fmt.Sprintf("%s: (empty)", part.Name)
	}

	first := part.Rounds[0]
	last := part.Rounds[len(part.Rounds)-1]

	return fmt.Sprintf("%s: %d rounds, %d → %d stitches",
		part.Name, len(part.Rounds), first.StitchCount, last.StitchCount)
}

// FormatRoundRange creates a summary of a round range
func (f *Formatter) FormatRoundRange(rounds []models.Round) string {
	if len(rounds) == 0 {
		return "(no rounds)"
	}

	first := rounds[0]
	last := rounds[len(rounds)-1]

	stitchProgression := fmt.Sprintf("%d → %d st", first.StitchCount, last.StitchCount)

	return fmt.Sprintf("Rnds %d-%d: %s", first.Number, last.Number, stitchProgression)
}
