# CrochetBot Frontend

React TypeScript application for converting 3D models to amigurumi crochet patterns.

## Features

- **File Upload**: Drag-and-drop or browse to upload .obj files
- **3D Preview**: Interactive 3D model viewer with Three.js
- **Pattern Display**: Beautiful, printable crochet pattern output
- **Export**: Download patterns as text or JSON

## Tech Stack

- **React 18** with TypeScript
- **Three.js** + React Three Fiber - 3D rendering
- **React Dropzone** - File upload
- **Axios** - API communication
- **CSS** - Custom styling

## Quick Start

```bash
# Install dependencies
npm install

# Start dev server
npm start

# Build for production
npm run build
```

## Environment

Create `.env` file:

```
REACT_APP_API_URL=http://localhost:8080
```

## Project Structure

```
src/
├── components/
│   ├── FileUpload.tsx       # Upload UI
│   ├── ModelPreview.tsx     # 3D viewer
│   └── PatternDisplay.tsx   # Pattern output
├── api.ts                   # API client
├── types.ts                 # TypeScript types
├── App.tsx                  # Main app
└── index.tsx                # Entry
```

## Usage

1. Upload .obj file
2. Preview 3D model
3. Generate pattern
4. Download result
