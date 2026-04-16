# CrochetBot

Convert images and 3D models into amigurumi crochet patterns.

## Project Structure

```
crochetbot/
├── cmd/
│   └── server/          # HTTP server entry point
├── internal/
│   ├── api/            # HTTP handlers and routing
│   ├── mesh/           # 3D mesh processing
│   ├── pattern/        # Pattern generation logic
│   └── models/         # Data structures
├── pkg/                # Public libraries (if any)
├── web/                # Frontend application
│   └── frontend/       # React app
├── test/               # Integration tests
└── docs/               # Documentation
```

## Quick Start

### Backend (Go)

```bash
# Build and run
make build
make run

# Or build and run manually
go build -o bin/crochetbot cmd/server/main.go
./bin/crochetbot
```

Server runs on http://localhost:8080

### Frontend (React)

```bash
cd web/frontend
npm install
npm start
```

Frontend runs on http://localhost:3000

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- Make (optional)

### Environment Variables
```bash
PORT=8080
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=52428800  # 50MB
```

## API Endpoints

Full API documentation: [docs/API.md](docs/API.md)

### POST /api/upload
Upload 3D model file (.obj, .stl)

**Request:** `multipart/form-data` with `file` field  
**Response:** File metadata with unique filename

### POST /api/generate
Generate crochet pattern from uploaded model

**Request:** JSON with `filename` from upload  
**Response:** Complete crochet pattern with rounds, materials, instructions

### GET /api/pattern/:id
Retrieve previously generated pattern

**Request:** Pattern ID in URL  
**Response:** Full pattern JSON

## Testing the API

```bash
# Start server
make run

# In another terminal, run automated tests
./test/test_api.sh

# Or test manually
curl -X POST -F "file=@test/sphere.obj" http://localhost:8080/api/upload
```

## Current Status

**✅ Phase 1 MVP Complete:**
- OBJ file parsing with full format support
- Sphere pattern generation
- REST API (upload, generate, retrieve)
- Pattern storage (JSON + in-memory cache)
- Comprehensive testing

**✅ Phase 2 Frontend Complete:**
- React TypeScript UI
- Drag-and-drop file upload
- 3D model preview (Three.js)
- Pattern display and export
- End-to-end workflow

See [ROADMAP.md](ROADMAP.md) for full development plan.

## Documentation

- [API Reference](docs/API.md) - Complete API documentation
- [OBJ Implementation](docs/OBJ_IMPLEMENTATION.md) - Technical details on OBJ parsing
- [Frontend Implementation](docs/FRONTEND_IMPLEMENTATION.md) - Frontend architecture and features
- [Implementation Summary](docs/IMPLEMENTATION_SUMMARY.md) - Backend implementation details

## License

TBD
