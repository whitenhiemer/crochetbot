package mesh

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOBJ(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		objContent    string
		wantVertices  int
		wantFaces     int
		wantErr       bool
	}{
		{
			name: "simple triangle",
			objContent: `# Simple triangle
v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 0.5 1.0 0.0
f 1 2 3
`,
			wantVertices: 3,
			wantFaces:    1,
			wantErr:      false,
		},
		{
			name: "quad (triangulated)",
			objContent: `v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 1.0 1.0 0.0
v 0.0 1.0 0.0
f 1 2 3 4
`,
			wantVertices: 4,
			wantFaces:    2, // Quad gets triangulated into 2 faces
			wantErr:      false,
		},
		{
			name: "with texture and normals",
			objContent: `v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 0.5 1.0 0.0
vt 0.0 0.0
vt 1.0 0.0
vt 0.5 1.0
vn 0.0 0.0 1.0
f 1/1/1 2/2/1 3/3/1
`,
			wantVertices: 3,
			wantFaces:    1,
			wantErr:      false,
		},
		{
			name: "with comments and groups",
			objContent: `# This is a comment
o MyObject
g MyGroup
v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 0.5 1.0 0.0
# Another comment
f 1 2 3
`,
			wantVertices: 3,
			wantFaces:    1,
			wantErr:      false,
		},
		{
			name: "negative indices",
			objContent: `v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 0.5 1.0 0.0
f -3 -2 -1
`,
			wantVertices: 3,
			wantFaces:    1,
			wantErr:      false,
		},
		{
			name: "empty file",
			objContent: `# Only comments
`,
			wantVertices: 0,
			wantFaces:    0,
			wantErr:      true, // Should error with no vertices
		},
		{
			name: "invalid vertex",
			objContent: `v 0.0 0.0
v 1.0 0.0 0.0
f 1 2
`,
			wantVertices: 0,
			wantFaces:    0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tmpDir, tt.name+".obj")
			if err := os.WriteFile(testFile, []byte(tt.objContent), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Load OBJ
			mesh, err := LoadOBJ(testFile)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadOBJ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return // Expected error, test passed
			}

			// Check vertex count
			if len(mesh.Vertices) != tt.wantVertices {
				t.Errorf("LoadOBJ() got %d vertices, want %d", len(mesh.Vertices), tt.wantVertices)
			}

			// Check face count
			if len(mesh.Faces) != tt.wantFaces {
				t.Errorf("LoadOBJ() got %d faces, want %d", len(mesh.Faces), tt.wantFaces)
			}

			// Check bounding box was calculated
			if tt.wantVertices > 0 {
				if mesh.Bounds.MaxX == 0 && mesh.Bounds.MaxY == 0 && mesh.Bounds.MaxZ == 0 {
					// Check if all vertices are actually at origin
					allZero := true
					for _, v := range mesh.Vertices {
						if v.X != 0 || v.Y != 0 || v.Z != 0 {
							allZero = false
							break
						}
					}
					if !allZero {
						t.Error("Bounding box not calculated")
					}
				}
			}
		})
	}
}

func TestLoadOBJFile(t *testing.T) {
	// Test with non-existent file
	_, err := LoadOBJ("/nonexistent/file.obj")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestCalculateBounds(t *testing.T) {
	mesh := &Mesh{
		Vertices: []Vertex{
			{X: -1.0, Y: -2.0, Z: -3.0},
			{X: 1.0, Y: 2.0, Z: 3.0},
			{X: 0.0, Y: 0.0, Z: 0.0},
		},
	}

	mesh.CalculateBounds()

	if mesh.Bounds.MinX != -1.0 || mesh.Bounds.MaxX != 1.0 {
		t.Errorf("X bounds incorrect: got [%f, %f], want [-1.0, 1.0]", mesh.Bounds.MinX, mesh.Bounds.MaxX)
	}
	if mesh.Bounds.MinY != -2.0 || mesh.Bounds.MaxY != 2.0 {
		t.Errorf("Y bounds incorrect: got [%f, %f], want [-2.0, 2.0]", mesh.Bounds.MinY, mesh.Bounds.MaxY)
	}
	if mesh.Bounds.MinZ != -3.0 || mesh.Bounds.MaxZ != 3.0 {
		t.Errorf("Z bounds incorrect: got [%f, %f], want [-3.0, 3.0]", mesh.Bounds.MinZ, mesh.Bounds.MaxZ)
	}
}
