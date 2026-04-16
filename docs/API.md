# CrochetBot API Documentation

## Base URL

```
http://localhost:8080
```

## Endpoints

### Health Check

Check if the server is running.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "ok",
  "service": "crochetbot"
}
```

---

### Upload 3D Model

Upload a 3D model file (.obj or .stl) for pattern generation.

**Endpoint:** `POST /api/upload`

**Content-Type:** `multipart/form-data`

**Parameters:**
- `file` (required): The 3D model file to upload

**Supported Formats:**
- `.obj` (Wavefront OBJ)
- `.stl` (STL - not yet implemented)

**Constraints:**
- Maximum file size: 50 MB
- Only one file per request

**Success Response (200):**
```json
{
  "success": true,
  "message": "File uploaded successfully: sphere.obj (513 bytes)",
  "file": {
    "id": "mesh-1776305652",
    "filename": "1776305652_sphere.obj",
    "uploaded_at": "2026-04-15T19:14:12.56241-07:00",
    "vertices": 0,
    "faces": 0,
    "format": "obj"
  }
}
```

**Error Response (400):**
```json
{
  "success": false,
  "error": "Invalid file type. Allowed: [.obj .stl]"
}
```

**Example (curl):**
```bash
curl -X POST \
  -F "file=@sphere.obj" \
  http://localhost:8080/api/upload
```

**Example (JavaScript):**
```javascript
const formData = new FormData();
formData.append('file', fileInput.files[0]);

const response = await fetch('http://localhost:8080/api/upload', {
  method: 'POST',
  body: formData
});

const result = await response.json();
console.log(result.file.filename); // Use this for generate request
```

---

### Generate Crochet Pattern

Generate a crochet pattern from an uploaded 3D model.

**Endpoint:** `POST /api/generate`

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "filename": "1776305652_sphere.obj",
  "file_id": "mesh-1776305652"
}
```

**Parameters:**
- `filename` (required): The filename returned from upload endpoint
- `file_id` (optional): The file ID for reference

**Success Response (200):**
```json
{
  "success": true,
  "message": "Pattern generated successfully with 1 part(s)",
  "pattern": {
    "id": "pattern-1776305652",
    "name": "Generated sphere Pattern",
    "created_at": "2026-04-15T19:14:12.613524-07:00",
    "description": "Auto-generated pattern from 3D model",
    "difficulty": "beginner",
    "parts": [
      {
        "name": "Body",
        "type": "sphere",
        "rounds": [
          {
            "number": 1,
            "instructions": "6 sc in magic ring",
            "stitch_count": 6,
            "stitch_type": "sc",
            "repeats": 1,
            "notes": "Pull tight to close"
          }
        ],
        "color": "main color",
        "starting_type": "magic ring",
        "notes": ["Use stitch marker to track rounds"]
      }
    ],
    "materials": {
      "yarn_weight": "worsted",
      "yarn_yardage": 50,
      "hook_size": "3.5mm",
      "colors": [{"name": "main color", "amount": 50}],
      "other_supplies": ["stuffing", "yarn needle", "stitch marker"]
    },
    "assembly_instructions": []
  }
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "File not found"
}
```

**Error Response (501):**
```json
{
  "success": false,
  "error": "STL files not yet supported"
}
```

**Example (curl):**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"filename": "1776305652_sphere.obj"}' \
  http://localhost:8080/api/generate
```

**Example (JavaScript):**
```javascript
const response = await fetch('http://localhost:8080/api/generate', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    filename: uploadResult.file.filename,
    file_id: uploadResult.file.id
  })
});

const result = await response.json();
console.log(result.pattern);
```

---

### Get Pattern

Retrieve a previously generated pattern by ID.

**Endpoint:** `GET /api/pattern/{id}`

**Parameters:**
- `id` (path parameter): Pattern ID returned from generate endpoint

**Success Response (200):**
```json
{
  "id": "pattern-1776305652",
  "name": "Generated sphere Pattern",
  "created_at": "2026-04-15T19:14:12.613524-07:00",
  "description": "Auto-generated pattern from 3D model",
  "difficulty": "beginner",
  "parts": [...],
  "materials": {...},
  "assembly_instructions": []
}
```

**Error Response (404):**
```json
{
  "error": "Pattern not found"
}
```

**Example (curl):**
```bash
curl http://localhost:8080/api/pattern/pattern-1776305652
```

**Example (JavaScript):**
```javascript
const response = await fetch(`http://localhost:8080/api/pattern/${patternId}`);
const pattern = await response.json();
```

---

## Data Models

### Pattern

```typescript
interface Pattern {
  id: string;
  name: string;
  created_at: string;
  description: string;
  difficulty: "beginner" | "intermediate" | "advanced";
  parts: Part[];
  materials: Materials;
  assembly_instructions: string[];
}
```

### Part

```typescript
interface Part {
  name: string;
  type: string; // "sphere", "cylinder", "cone", etc.
  rounds: Round[];
  color: string;
  starting_type: string; // "magic ring", "chain"
  notes: string[];
}
```

### Round

```typescript
interface Round {
  number: number;
  instructions: string;
  stitch_count: number;
  stitch_type: string; // "sc", "hdc", "dc", "inc", "dec"
  repeats: number;
  notes: string;
}
```

### Materials

```typescript
interface Materials {
  yarn_weight: string; // "DK", "worsted", etc.
  yarn_yardage: number;
  hook_size: string; // "3.5mm", "E/4", etc.
  colors: Color[];
  other_supplies: string[];
}
```

### Color

```typescript
interface Color {
  name: string;
  amount: number; // yards
}
```

---

## Complete Workflow Example

### 1. Upload a file

```bash
UPLOAD_RESPONSE=$(curl -s -X POST \
  -F "file=@sphere.obj" \
  http://localhost:8080/api/upload)

FILENAME=$(echo "$UPLOAD_RESPONSE" | jq -r '.file.filename')
```

### 2. Generate pattern

```bash
GENERATE_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{\"filename\": \"$FILENAME\"}" \
  http://localhost:8080/api/generate)

PATTERN_ID=$(echo "$GENERATE_RESPONSE" | jq -r '.pattern.id')
```

### 3. Retrieve pattern later

```bash
curl http://localhost:8080/api/pattern/$PATTERN_ID
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request (invalid input) |
| 404 | Not Found (file or pattern doesn't exist) |
| 405 | Method Not Allowed (wrong HTTP method) |
| 500 | Internal Server Error |
| 501 | Not Implemented (feature not yet available) |

---

## CORS

CORS is enabled for all origins (`Access-Control-Allow-Origin: *`).

For production, configure specific origins in the server code.

---

## Rate Limiting

Currently no rate limiting is implemented.

For production deployment, add rate limiting middleware.

---

## Storage

### File Storage
- Uploaded files: `./uploads/` directory
- Filename format: `{timestamp}_{sanitized_original_name}`

### Pattern Storage
- Patterns: `./data/patterns/` directory
- Format: JSON files named `{pattern_id}.json`
- Also cached in memory for fast retrieval

---

## Testing

Run the automated test suite:

```bash
# Start the server
./bin/crochetbot

# In another terminal
./test/test_api.sh
```

Or test manually:

```bash
# Health check
curl http://localhost:8080/health

# Upload
curl -X POST -F "file=@test/sphere.obj" http://localhost:8080/api/upload

# Generate (use filename from upload response)
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"filename": "YOUR_FILENAME_HERE"}' \
  http://localhost:8080/api/generate

# Retrieve (use pattern ID from generate response)
curl http://localhost:8080/api/pattern/YOUR_PATTERN_ID_HERE
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `UPLOAD_DIR` | `./uploads` | Directory for uploaded files |
| `PATTERN_STORAGE_DIR` | `./data/patterns` | Directory for pattern JSON files |
| `MAX_FILE_SIZE` | `52428800` | Max upload size in bytes (50MB) |

---

## Future API Additions

### Planned endpoints:
- `GET /api/patterns` - List all patterns
- `DELETE /api/pattern/{id}` - Delete a pattern
- `POST /api/pattern/{id}/export` - Export as PDF
- `GET /api/mesh/{id}` - Get mesh metadata
- `POST /api/generate/preview` - Preview without saving
