# CrochetBot Roadmap

## Project Vision
Web application that converts images or 3D models into amigurumi crochet patterns with step-by-step instructions.

## Core Components

### 1. Input Processing
- Image upload and preprocessing
- 3D model file support (.obj, .stl)
- Image-to-3D reconstruction (if starting from 2D image)

### 2. Shape Analysis
- 3D mesh processing and simplification
- Identification of basic shapes (sphere, cylinder, cone)
- Segmentation into crochatable parts
- Symmetry detection

### 3. Pattern Generation
- Convert 3D geometry to crochet rounds
- Calculate stitch counts for increases/decreases
- Generate construction order (bottom-up, top-down)
- Handle complex shapes (limbs, appendages)
- Color planning for multi-color patterns

### 4. Output Formatting
- Human-readable pattern instructions
- Stitch diagrams
- Material requirements (yarn weight, hook size, yardage)
- Assembly instructions
- PDF export

### 5. Web Interface
- File upload (drag-and-drop)
- 3D preview viewer
- Pattern customization options
- Export/download functionality

## Development Phases

### Phase 1: MVP - Basic Shape Conversion (Weeks 1-4)
**Goal:** Convert simple 3D models (sphere, ellipsoid) to basic amigurumi patterns

- [ ] Project setup (repo structure, dependencies)
- [ ] Basic web interface (upload + display)
- [ ] 3D model parser (.obj format)
- [ ] Simple sphere-to-pattern algorithm
- [ ] Basic pattern text output
- [ ] Calculate rounds and stitch counts for spherical shapes

**Tech Stack Decision:**
- Backend: Go (performance, deployment simplicity)
- Frontend: React/Next.js
- 3D Viewer: Three.js
- Deployment: Docker containers

### Phase 2: Enhanced Shape Support (Weeks 5-8)
**Goal:** Handle multiple basic shapes and combinations

- [ ] Support cylinders (arms, legs, bodies)
- [ ] Support cones (noses, tails)
- [ ] Shape combination and attachment logic
- [ ] Part segmentation algorithm
- [ ] Assembly instructions generation
- [ ] Material estimation calculations

### Phase 3: Image-to-Model Pipeline (Weeks 9-12)
**Goal:** Accept 2D images as input

- [ ] Image preprocessing and background removal
- [ ] Depth estimation from single image
- [ ] 2D-to-3D reconstruction
- [ ] Multiple view support (front, side, back)
- [ ] Edge case handling (quality checks)

### Phase 4: Advanced Features (Weeks 13-16)
**Goal:** Complex patterns and user customization

- [ ] Color segmentation for multi-color patterns
- [ ] Surface texture interpretation
- [ ] Size customization (scale patterns)
- [ ] Yarn substitution calculator
- [ ] Pattern difficulty rating
- [ ] Stitch diagram generation

### Phase 5: Polish & Production (Weeks 17-20)
**Goal:** Production-ready application

- [ ] UI/UX improvements
- [ ] Pattern quality validation
- [ ] Export formats (PDF, print-friendly)
- [ ] User authentication (save patterns)
- [ ] Pattern gallery (community sharing)
- [ ] Performance optimization
- [ ] Documentation and tutorials

## Technical Challenges

### Critical Path Items
1. **Stitch Count Algorithm:** Converting 3D curvature to increase/decrease rates
2. **Shape Decomposition:** Breaking complex models into crochatable primitives
3. **Assembly Logic:** Determining optimal construction order
4. **Image-to-3D:** Accurate depth estimation from 2D images

### Research Needed
- Existing crochet pattern formats and standards
- Gaussian curvature to stitch rate mapping
- 3D reconstruction techniques (NeRF, photogrammetry)
- Mesh simplification algorithms

## Success Metrics

### Phase 1 (MVP)
- Generate valid pattern for sphere (tested by actual crochet)
- Processing time < 30 seconds
- Pattern produces recognizable shape

### Phase 2
- Support 5+ basic shapes
- 80% accuracy in shape segmentation
- Assembly instructions understandable

### Phase 3
- Image-to-model conversion success rate > 70%
- Works with phone camera images

### Phase 4+
- Multi-color pattern support
- Pattern completion rate (users finish patterns)
- User satisfaction scores

## Technology Stack

### Backend: Go
**Core Libraries:**
- 3D mesh processing: github.com/fogleman/fauxgl, github.com/qmuntal/gltf
- OBJ parsing: github.com/shabbyrobe/go-obj
- Image processing: github.com/disintegration/imaging, golang.org/x/image
- Math/geometry: gonum.org/v1/gonum
- HTTP server: net/http (stdlib) or github.com/gin-gonic/gin
- JSON handling: encoding/json (stdlib)

**ML Integration (Phase 3):**
- Call Python microservice via gRPC/HTTP for image-to-3D
- Or use ONNX runtime for Go: github.com/yalue/onnxruntime_go
- TensorFlow Go bindings (community maintained)

**Benefits:**
- Fast compilation and execution
- Easy deployment (single binary)
- Excellent concurrency for processing multiple patterns
- Strong stdlib for HTTP/JSON
- Good fit for geometric computations

### Frontend
- Framework: React + Next.js
- 3D Viewer: Three.js or React Three Fiber
- UI Components: Tailwind CSS
- File handling: Drag-and-drop with validation

### Infrastructure
- Hosting: Vercel (frontend), Railway/Render (backend)
- Storage: S3 or equivalent for uploaded files
- Database: PostgreSQL (if user accounts needed)
- Processing: Background jobs for long-running tasks

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Pattern accuracy issues | High | Start with simple shapes, extensive testing with real crochet |
| Image-to-3D poor quality | Medium | Require multiple views, set quality thresholds |
| Processing time too long | Medium | Background jobs, progress indicators, optimize algorithms |
| Complex shapes unsupported | Low | Clear documentation of limitations, graceful degradation |
| User confusion with patterns | High | Include video tutorials, visual guides, test with crocheters |

## Open Questions

1. Target user skill level? (beginner, intermediate, advanced)
2. Pricing model? (free, freemium, subscription)
3. Pattern licensing? (user owns, platform retains rights)
4. Community features? (sharing, comments, modifications)
5. Mobile app or web-only?
6. Real-time preview or batch processing?

## Next Steps

1. **Immediate:** Research existing crochet pattern generators and formats
2. **Week 1:** Set up repository structure and basic web scaffold
3. **Week 1:** Prototype sphere-to-pattern algorithm (Python notebook)
4. **Week 2:** Build minimal web interface for file upload
5. **Week 2:** Integrate pattern generation with simple .obj files

## Resources

### Learning Materials
- Amigurumi pattern basics and terminology
- Computer graphics: mesh processing, curvature calculation
- 3D reconstruction techniques
- Crochet mathematics (increase/decrease rates)

### Similar Projects (Research)
- Existing crochet pattern generators
- 3D-to-craft conversion tools
- Image-to-model services

### Community
- Reddit: r/crochet, r/Amigurumi
- Ravelry: pattern format standards
- Crochet pattern designers (for validation)

---

**Version:** 0.1  
**Last Updated:** 2026-04-15  
**Status:** Planning Phase
