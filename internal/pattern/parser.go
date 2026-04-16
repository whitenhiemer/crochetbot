package pattern

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/whitenhiemer/crochetbot/internal/models"
)

// Parser reads crochet patterns from text format
type Parser struct {
	// Regex patterns for parsing
	roundPattern    *regexp.Regexp
	bracketPattern  *regexp.Regexp
	stitchPattern   *regexp.Regexp
	colorPattern    *regexp.Regexp
}

// NewParser creates a new pattern parser
func NewParser() *Parser {
	return &Parser{
		// Matches: "Rnd 3." or "Rnds 5-7." or "Rnd 1:"
		roundPattern: regexp.MustCompile(`(?i)Rnds?\s*(\d+)(?:-(\d+))?[.:]`),
		// Matches: "[sc, inc] x 6" or "[2 sc, inc] repeat 6 times"
		bracketPattern: regexp.MustCompile(`\[([^\]]+)\]\s*(?:x|repeat)\s*(\d+)(?:\s*times?)?`),
		// Matches final stitch count: "(18)" or "(18 sc)"
		stitchPattern: regexp.MustCompile(`\((\d+)(?:\s+\w+)?\)`),
		// Matches: "(black yarn)" or "(switch to white yarn)"
		colorPattern: regexp.MustCompile(`\((?:switch to\s+)?([a-zA-Z]+)\s+yarn\)`),
	}
}

// ParsePattern parses a complete pattern from text
func (p *Parser) ParsePattern(text string) (*models.Pattern, error) {
	lines := strings.Split(text, "\n")

	pattern := &models.Pattern{
		Parts: []models.Part{},
	}

	var currentPart *models.Part
	var currentColor string = "main color"

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this is a new part header
		if p.isPartHeader(line) {
			if currentPart != nil {
				pattern.Parts = append(pattern.Parts, *currentPart)
			}
			currentPart = &models.Part{
				Name:         line,
				Rounds:       []models.Round{},
				Color:        currentColor,
				StartingType: "magic ring",
			}
			continue
		}

		// Skip metadata lines (With X yarn, color notes, etc)
		lowerLine := strings.ToLower(line)
		if strings.HasPrefix(lowerLine, "with ") {
			if strings.Contains(lowerLine, "yarn") {
				// Extract color if present
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					currentColor = parts[1] + " yarn"
				}
			}
			continue
		}

		// Try to parse as a round first (it might have inline color notes)
		if round, err := p.ParseRound(line); err == nil && currentPart != nil {
			// Check for inline color change
			if colorMatch := p.colorPattern.FindStringSubmatch(line); colorMatch != nil {
				currentColor = colorMatch[1] + " yarn"
			}
			// Apply current color if round doesn't specify
			if round.Notes == "" {
				// Color might be in the instruction
				if strings.Contains(strings.ToLower(line), "switch to") {
					// Already captured above
				}
			}
			currentPart.Rounds = append(currentPart.Rounds, round)
		}
	}

	// Add final part
	if currentPart != nil {
		pattern.Parts = append(pattern.Parts, *currentPart)
	}

	if len(pattern.Parts) == 0 {
		return nil, fmt.Errorf("no parts found in pattern")
	}

	return pattern, nil
}

// ParseRound parses a single round from text
// Examples:
//   "Rnd 1. 6 sc in magic ring (6)"
//   "Rnd 3. [sc, inc] x 6 (18)"
//   "Rnds 5-7. 24 sc (24)"
func (p *Parser) ParseRound(text string) (models.Round, error) {
	text = strings.TrimSpace(text)

	// Extract round number(s)
	roundMatch := p.roundPattern.FindStringSubmatch(text)
	if roundMatch == nil {
		return models.Round{}, fmt.Errorf("no round number found")
	}

	roundNum, _ := strconv.Atoi(roundMatch[1])

	// Check for round range (e.g., "Rnds 5-7")
	isRange := roundMatch[2] != ""

	// Remove round prefix to get instructions
	instruction := p.roundPattern.ReplaceAllString(text, "")
	instruction = strings.TrimSpace(instruction)

	// Extract final stitch count
	stitchCount := 0
	if stitchMatch := p.stitchPattern.FindStringSubmatch(instruction); stitchMatch != nil {
		stitchCount, _ = strconv.Atoi(stitchMatch[1])
	}

	// Parse instruction details
	repeats := 1
	stitchType := "sc" // default
	notes := ""

	// Check for bracketed pattern with repeats
	if bracketMatch := p.bracketPattern.FindStringSubmatch(instruction); bracketMatch != nil {
		repeats, _ = strconv.Atoi(bracketMatch[2])
		// Keep full instruction for display
	}

	// Detect stitch type
	lowerInst := strings.ToLower(instruction)
	if strings.Contains(lowerInst, "inc") {
		stitchType = "inc"
	} else if strings.Contains(lowerInst, "dec") {
		stitchType = "dec"
	} else if strings.Contains(lowerInst, "sl st") {
		stitchType = "sl st"
	} else if strings.Contains(lowerInst, "hdc") {
		stitchType = "hdc"
	} else if strings.Contains(lowerInst, "dc") {
		stitchType = "dc"
	}

	// Extract notes (common patterns)
	if strings.Contains(lowerInst, "stuff") {
		notes = "stuffing"
	} else if strings.Contains(lowerInst, "magic ring") || strings.Contains(lowerInst, "magic loop") {
		notes = "magic ring start"
	} else if strings.Contains(lowerInst, "fasten off") {
		notes = "fasten off"
		stitchType = "finish"
	}

	// Mark if it's a range instruction
	if isRange {
		notes = "repeated rounds"
	}

	return models.Round{
		Number:       roundNum,
		Instructions: instruction,
		StitchCount:  stitchCount,
		StitchType:   stitchType,
		Repeats:      repeats,
		Notes:        notes,
	}, nil
}

// isPartHeader detects if a line is a part header
func (p *Parser) isPartHeader(line string) bool {
	upper := strings.ToUpper(line)
	headers := []string{"HEAD", "BODY", "ARM", "LEG", "TAIL", "EAR", "WING", "BEAK", "FEET", "CREST", "FEATHERS"}

	for _, header := range headers {
		if strings.Contains(upper, header) && !strings.Contains(upper, "RND") {
			// Check if it's all caps or title-like
			if upper == line || strings.Title(line) == line {
				return true
			}
		}
	}

	return false
}

// ParseInstruction breaks down an instruction into components
func (p *Parser) ParseInstruction(instruction string) InstructionComponents {
	components := InstructionComponents{
		Raw: instruction,
	}

	lower := strings.ToLower(instruction)

	// Count specific stitch types
	components.SingleCrochet = strings.Count(lower, "sc")
	components.Increase = strings.Count(lower, "inc")
	components.Decrease = strings.Count(lower, "dec")
	components.HalfDouble = strings.Count(lower, "hdc")
	components.DoubleCrochet = strings.Count(lower, "dc")
	components.SlipStitch = strings.Count(lower, "sl st")
	components.Chain = strings.Count(lower, "ch")

	// Extract multiplier
	if match := regexp.MustCompile(`x\s*(\d+)`).FindStringSubmatch(lower); match != nil {
		components.Multiplier, _ = strconv.Atoi(match[1])
	} else if match := regexp.MustCompile(`repeat\s+(\d+)`).FindStringSubmatch(lower); match != nil {
		components.Multiplier, _ = strconv.Atoi(match[1])
	} else {
		components.Multiplier = 1
	}

	return components
}

// InstructionComponents breaks down instruction details
type InstructionComponents struct {
	Raw           string
	SingleCrochet int
	Increase      int
	Decrease      int
	HalfDouble    int
	DoubleCrochet int
	SlipStitch    int
	Chain         int
	Multiplier    int
}

// ParsePartFromText parses a single part section
func (p *Parser) ParsePartFromText(name, text string) (models.Part, error) {
	lines := strings.Split(text, "\n")

	part := models.Part{
		Name:         name,
		Rounds:       []models.Round{},
		Color:        "main color",
		StartingType: "magic ring",
		Notes:        []string{},
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}

		// Try to parse as round
		if round, err := p.ParseRound(line); err == nil {
			part.Rounds = append(part.Rounds, round)
		} else {
			// Might be a note or instruction
			if !p.isPartHeader(line) {
				part.Notes = append(part.Notes, line)
			}
		}
	}

	return part, nil
}

// ValidatePattern checks if a parsed pattern is valid
func (p *Parser) ValidatePattern(pattern *models.Pattern) []string {
	issues := []string{}

	if len(pattern.Parts) == 0 {
		issues = append(issues, "Pattern has no parts")
	}

	for i, part := range pattern.Parts {
		if len(part.Rounds) == 0 {
			issues = append(issues, fmt.Sprintf("Part %d (%s) has no rounds", i, part.Name))
		}

		// Check stitch progression makes sense
		for j := 1; j < len(part.Rounds); j++ {
			prev := part.Rounds[j-1]
			curr := part.Rounds[j]

			if prev.StitchCount > 0 && curr.StitchCount > 0 {
				diff := curr.StitchCount - prev.StitchCount
				// Flag unrealistic jumps (more than 50% change)
				if diff > prev.StitchCount/2 || diff < -prev.StitchCount/2 {
					issues = append(issues, fmt.Sprintf("Part %s, Round %d: Large stitch change from %d to %d",
						part.Name, curr.Number, prev.StitchCount, curr.StitchCount))
				}
			}
		}
	}

	return issues
}
