# Getting Started with CrochetBot

Complete guide to running CrochetBot locally.

## Prerequisites

- **Go 1.21+** (backend)
- **Node.js 18+** (frontend)
- **npm** (comes with Node.js)
- **Git** (for cloning)

## Quick Start (5 minutes)

### 1. Clone Repository

```bash
git clone git@github.com:whitenhiemer/crochetbot.git
cd crochetbot
```

### 2. Start Backend

```bash
# Build the Go server
make build

# Run the server
make run
```

Backend will start on **http://localhost:8080**

Verify it's running:
```bash
curl http://localhost:8080/health
# Should return: {"status":"ok","service":"crochetbot"}
```

### 3. Start Frontend

Open a new terminal:

```bash
cd web/frontend

# Install dependencies (first time only)
npm install

# Start dev server
npm start
```

Frontend will start on **http://localhost:3000**

Your browser should automatically open to http://localhost:3000

---

## Using the Application

### Step 1: Upload a 3D Model

1. Drag and drop an .obj file onto the upload area
2. Or click "Browse Files" to select a file
3. Test file available at: `test/sphere.obj`

**Supported formats:**
- .obj (Wavefront OBJ) ✅
- .stl (STL) - coming soon

**File size limit:** 10 MB

### Step 2: Preview Model

1. Your 3D model will render in the viewer
2. Use mouse to interact:
   - **Left click + drag** - Rotate
   - **Scroll** - Zoom in/out
   - **Right click + drag** - Pan
3. Review model information (vertices, faces, size)
4. Click **"Generate Crochet Pattern"**

### Step 3: View & Download Pattern

1. Pattern displays with:
   - Materials needed
   - Round-by-round instructions
   - Stitch counts
   - Notes and tips
2. **Download Pattern** - Save as .txt file
3. **Download JSON** - Save structured data
4. **New Pattern** - Start over

---

## Testing the API Directly

### Health Check

```bash
curl http://localhost:8080/health
```

### Upload File

```bash
curl -X POST \
  -F "file=@test/sphere.obj" \
  http://localhost:8080/api/upload
```

### Generate Pattern

```bash
# Use filename from upload response
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"filename": "YOUR_FILENAME_HERE"}' \
  http://localhost:8080/api/generate
```

### Get Pattern

```bash
# Use pattern ID from generate response
curl http://localhost:8080/api/pattern/PATTERN_ID_HERE
```

---

## Running Tests

### Backend Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test -v ./internal/mesh/
```

### Frontend Tests

```bash
cd web/frontend
npm test
```

### API Integration Tests

```bash
# Ensure backend is running first
./test/test_api.sh
```

---

## Building for Production

### Backend

```bash
# Build optimized binary
make build-prod

# Output: bin/crochetbot-linux

# Run production build
./bin/crochetbot-linux
```

### Frontend

```bash
cd web/frontend

# Build optimized bundle
npm run build

# Output: build/ directory

# Serve with any static file server
npx serve -s build
```

---

## Docker

### Build Images

```bash
# Backend
docker build -t crochetbot-backend .

# Frontend
cd web/frontend
docker build -t crochetbot-frontend .
```

### Run Containers

```bash
# Backend
docker run -p 8080:8080 crochetbot-backend

# Frontend
docker run -p 3000:3000 crochetbot-frontend
```

---

## Configuration

### Backend Environment Variables

Create `.env` file in project root:

```bash
PORT=8080
UPLOAD_DIR=./uploads
PATTERN_STORAGE_DIR=./data/patterns
MAX_FILE_SIZE=10485760
```

### Frontend Environment Variables

Create `web/frontend/.env`:

```bash
REACT_APP_API_URL=http://localhost:8080
```

---

## Directory Structure

```
crochetbot/
├── cmd/
│   ├── server/           # Backend server entry point
│   └── test-obj/         # CLI testing utility
├── internal/
│   ├── api/              # HTTP handlers
│   ├── mesh/             # 3D mesh processing
│   ├── pattern/          # Pattern generation
│   └── models/           # Data structures
├── web/
│   └── frontend/         # React frontend app
├── test/
│   ├── sphere.obj        # Test 3D model
│   └── test_api.sh       # API test script
├── docs/                 # Documentation
├── uploads/              # Uploaded files (created at runtime)
├── data/                 # Pattern storage (created at runtime)
└── bin/                  # Compiled binaries (created at build)
```

---

## Troubleshooting

### Backend Won't Start

**Problem:** Port 8080 already in use

**Solution:**
```bash
# Find process using port 8080
lsof -ti:8080

# Kill the process
kill -9 $(lsof -ti:8080)

# Or change port
PORT=8081 ./bin/crochetbot
```

### Frontend Won't Start

**Problem:** Port 3000 already in use

**Solution:**
```bash
# Kill process
lsof -ti:3000 | xargs kill -9

# Or use different port
PORT=3001 npm start
```

**Problem:** TypeScript errors

**Solution:**
```bash
cd web/frontend
rm -rf node_modules package-lock.json
npm install
```

### CORS Errors

**Problem:** Frontend can't reach backend

**Solution:**
- Ensure backend is running on http://localhost:8080
- Check `REACT_APP_API_URL` in `web/frontend/.env`
- Backend has CORS enabled for all origins

### Upload Fails

**Problem:** File too large

**Solution:**
- Max file size is 10MB
- Simplify your 3D model
- Or increase `MAX_FILE_SIZE` env var

**Problem:** File format not supported

**Solution:**
- Only .obj files fully supported
- Convert other formats to .obj using Blender

### 3D Model Won't Render

**Problem:** Black screen or error

**Solution:**
- Ensure file is valid .obj format
- Check browser console for errors
- Try the test file: `test/sphere.obj`
- Verify WebGL is supported: visit https://get.webgl.org/

---

## Common Tasks

### Clear Uploaded Files

```bash
rm -rf uploads/*
```

### Clear Generated Patterns

```bash
rm -rf data/patterns/*
```

### Rebuild Everything

```bash
# Backend
make clean
make build

# Frontend
cd web/frontend
rm -rf node_modules build
npm install
npm run build
```

### View Backend Logs

```bash
# If running in background
tail -f /tmp/crochetbot-backend.log

# Or run in foreground
make run
```

### View Frontend Logs

```bash
# Development server shows logs in terminal
npm start

# Or check browser console (F12)
```

---

## Development Tips

### Hot Reload

**Backend:** Use `air` for hot reload
```bash
go install github.com/cosmtrek/air@latest
air
```

**Frontend:** Already has hot reload
```bash
npm start
# Edit files, they auto-reload
```

### Code Formatting

**Backend:**
```bash
make fmt
```

**Frontend:**
```bash
cd web/frontend
npx prettier --write src/
```

### Linting

**Backend:**
```bash
make lint
```

**Frontend:**
```bash
cd web/frontend
npm run lint
```

---

## Next Steps

Once you have the app running:

1. **Try the example:** Upload `test/sphere.obj`
2. **Read the docs:** See `docs/` directory
3. **Explore the code:** Start with `cmd/server/main.go` and `web/frontend/src/App.tsx`
4. **Make it yours:** Customize patterns, add features, improve UI

---

## Getting Help

- **Issues:** https://github.com/whitenhiemer/crochetbot/issues
- **Documentation:** See `docs/` directory
- **API Reference:** `docs/API.md`
- **Frontend Guide:** `docs/FRONTEND_IMPLEMENTATION.md`

---

## What's Working

✅ File upload (.obj)  
✅ 3D model preview  
✅ Pattern generation (spheres)  
✅ Pattern display  
✅ Export to .txt/.json  
✅ Full API  
✅ Tests passing  

**Ready to use!** 🎉
