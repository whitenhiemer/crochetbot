# Pattern API Documentation

API endpoints for parsing, formatting, validating, and comparing crochet patterns.

## Endpoints

### POST /api/pattern/parse

Parse text pattern into structured JSON.

**Request:**
```json
{
  "text": "HEAD & BODY\n\nRnd 1. 6 sc in magic ring (6)\nRnd 2. 6 inc (12)..."
}
```

**Response:**
```json
{
  "success": true,
  "pattern": {
    "id": "pattern-123456",
    "name": "",
    "parts": [
      {
        "name": "HEAD & BODY",
        "type": "sphere",
        "rounds": [
          {
            "number": 1,
            "instructions": "6 sc in magic ring",
            "stitch_count": 6,
            "stitch_type": "sc",
            "repeats": 1,
            "notes": "magic ring start"
          }
        ],
        "color": "main color",
        "starting_type": "magic ring"
      }
    ]
  }
}
```

**Supported Pattern Notation:**
- `Rnd 1.` - Single round
- `Rnds 5-7.` - Round range
- `[sc, inc] x 6` - Pattern with repeats
- `[2 sc, dec] repeat 6 times` - Verbose repeat notation
- `(18)` - Total stitch count
- Part headers: HEAD, BODY, WINGS, FEET, etc. (uppercase)
- Color changes: `With pink yarn.` or `(switch to white yarn)`

---

### POST /api/pattern/format

Convert structured pattern to human-readable text.

**Request:**
```json
{
  "pattern": { /* Pattern object */ },
  "compact_mode": false
}
```

**Response:**
```json
{
  "success": true,
  "text": "HEAD & BODY\nWith main color.\n\nRnd 1. 6 sc in magic ring (6)...",
  "compact_text": "Test Pattern: Body (3 rnds)"
}
```

**Output Format:**
- Woobles-style pattern text
- Includes materials, abbreviations, notes
- Groups consecutive identical rounds: `Rnds 5-7. 24 sc (24)`
- Adds stuffing reminders for 3D parts
- Includes accuracy metrics if available

---

### POST /api/pattern/validate

Validate pattern quality and realism.

**Request:**
```json
{
  "pattern": { /* Pattern object */ },
  "reference_pattern": { /* Optional reference for comparison */ }
}
```

**Response:**
```json
{
  "success": true,
  "validation_result": {
    "is_valid": true,
    "score": 88.0,
    "issues": [
      {
        "severity": "warning",
        "location": "Part HEAD & BODY, Round 2",
        "message": "Unrealistic increase rate",
        "details": "100.0% increase (6 → 12) exceeds typical crochet constraints"
      }
    ],
    "warnings": [
      "Part HEAD & BODY, Round 9: Stitch count 15 not divisible by 6 - may be irregular"
    ],
    "stitch_progression": {
      "average_increase_rate": 38.2,
      "average_decrease_rate": 30.0,
      "max_jump": 6,
      "max_drop": 8,
      "unrealistic_changes": 6,
      "smooth_transitions": 10,
      "total_transitions": 16
    },
    "terminology_score": 100.0,
    "structural_score": 100.0,
    "realism_score": 70.0
  },
  "comparison_metrics": {
    "structural_similarity": 0.95,
    "length_ratio": 1.1,
    "stitch_count_drift": 0.05,
    "progression_match": 0.92,
    "terminology_match": 1.0
  }
}
```

**Validation Criteria:**

**Terminology Score (0-100):**
- Checks for standard crochet abbreviations (sc, inc, dec, sl st, ch, etc.)
- Flags non-standard instructions
- Deducts for vague instructions like "increase evenly"

**Structural Score (0-100):**
- Verifies proper start (magic ring with 6 stitches)
- Checks for proper closing on 3D shapes
- Evaluates balance between increases/decreases
- Ensures parts have rounds

**Realism Score (0-100):**
- Penalizes unrealistic stitch changes (>30% per round)
- Checks for smooth transitions
- Validates stitch counts are reasonable (6-200)
- Rewards consistent progression

**Stitch Progression Metrics:**
- Average increase/decrease rates
- Max jump/drop in single round
- Count of smooth vs unrealistic transitions
- Total transitions analyzed

---

### POST /api/pattern/compare

Compare two patterns for similarity.

**Request:**
```json
{
  "generated": { /* Pattern object */ },
  "reference": { /* Pattern object */ }
}
```

**Response:**
```json
{
  "success": true,
  "comparison": {
    "structural_similarity": 0.92,
    "length_ratio": 0.95,
    "stitch_count_drift": 0.08,
    "progression_match": 0.88,
    "terminology_match": 0.95
  }
}
```

**Comparison Metrics:**

- **Structural Similarity (0-1):** How similar is the overall shape progression? Compares normalized stitch profiles.
  
- **Length Ratio:** Generated rounds / reference rounds. 1.0 = same length, <1 = shorter, >1 = longer.

- **Stitch Count Drift (0-1):** Average normalized difference in stitch counts at matching positions. Lower is better.

- **Progression Match (0-1):** How well do the increase/decrease patterns align? Compares change sequences.

- **Terminology Match (0-1):** Overlap in crochet terms used. 1.0 = identical terminology.

---

## Error Responses

All endpoints return errors in this format:

```json
{
  "success": false,
  "error": "Error message here"
}
```

**Common HTTP Status Codes:**
- `200 OK` - Success
- `400 Bad Request` - Invalid input (missing fields, malformed JSON, unparseable pattern)
- `405 Method Not Allowed` - Wrong HTTP method (all endpoints require POST)
- `500 Internal Server Error` - Server-side processing error

---

## Usage Examples

### Parse and Format Round-trip

```bash
# Parse text pattern
curl -X POST http://localhost:8080/api/pattern/parse \
  -H "Content-Type: application/json" \
  -d '{"text": "BODY\n\nRnd 1. 6 sc in magic ring (6)\nRnd 2. 6 inc (12)"}' \
  > pattern.json

# Format back to text
curl -X POST http://localhost:8080/api/pattern/format \
  -H "Content-Type: application/json" \
  -d @pattern.json
```

### Validate Generated Pattern

```bash
curl -X POST http://localhost:8080/api/pattern/validate \
  -H "Content-Type: application/json" \
  -d '{"pattern": {...}}'
```

### Compare Generated to Reference

```bash
curl -X POST http://localhost:8080/api/pattern/compare \
  -H "Content-Type: application/json" \
  -d '{
    "generated": {...},
    "reference": {...}
  }'
```

---

## Integration with Pattern Generation

### Full Workflow

1. **Upload 3D Model:** `POST /api/upload` → Get file ID
2. **Generate Pattern:** `POST /api/generate` → Get structured pattern
3. **Validate Quality:** `POST /api/pattern/validate` → Check realism scores
4. **Format for Display:** `POST /api/pattern/format` → Get human-readable text
5. **Compare to Reference (optional):** `POST /api/pattern/compare` → Benchmark accuracy

### Example: Generate and Validate

```javascript
// Generate pattern from STL
const generateResponse = await fetch('/api/generate', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ filename: 'model.stl' })
});
const { pattern } = await generateResponse.json();

// Validate the generated pattern
const validateResponse = await fetch('/api/pattern/validate', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ pattern })
});
const { validation_result } = await validateResponse.json();

console.log(`Pattern score: ${validation_result.score}/100`);
console.log(`Issues: ${validation_result.issues.length}`);

// Format for display
const formatResponse = await fetch('/api/pattern/format', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ pattern })
});
const { text } = await formatResponse.json();

console.log(text);
```

---

## Pattern Object Schema

See `internal/models/pattern.go` for full type definitions:

```go
type Pattern struct {
    ID               string
    Name             string
    Description      string
    Difficulty       string          // "beginner", "intermediate", "advanced"
    Parts            []Part
    Materials        Materials
    Assembly         []string
    FinishedSize     FinishedSize
    AccuracyMetrics  AccuracyMetrics
}

type Part struct {
    Name         string
    Type         string   // "sphere", "cylinder", "cone"
    Rounds       []Round
    Color        string
    StartingType string   // "magic ring", "chain"
    Notes        []string
}

type Round struct {
    Number       int
    Instructions string
    StitchCount  int
    StitchType   string   // "sc", "hdc", "dc", "inc", "dec"
    Repeats      int
    Notes        string
}
```

---

## Testing

Run tests:
```bash
go test ./internal/api -v        # API handler tests
go test ./internal/pattern -v    # Parser/validator tests
```

Test coverage includes:
- ✓ Parse valid patterns (single/multi-part, color changes)
- ✓ Parse edge cases (empty, malformed, ranges)
- ✓ Format patterns to text
- ✓ Round-trip parse→format
- ✓ Validate pattern quality
- ✓ Compare patterns for similarity
- ✓ HTTP error handling
- ✓ Method validation
