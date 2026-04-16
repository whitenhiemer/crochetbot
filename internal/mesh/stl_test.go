package mesh

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSTLAscii(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple ASCII STL file (single triangle)
	stlContent := `solid test
  facet normal 0 0 1
    outer loop
      vertex 0.0 0.0 0.0
      vertex 1.0 0.0 0.0
      vertex 0.5 1.0 0.0
    endloop
  endfacet
endsolid test
`

	testFile := filepath.Join(tmpDir, "test_ascii.stl")
	if err := os.WriteFile(testFile, []byte(stlContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load STL
	mesh, err := LoadSTL(testFile)
	if err != nil {
		t.Fatalf("LoadSTL() error = %v", err)
	}

	// Check results
	if len(mesh.Vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(mesh.Vertices))
	}

	if len(mesh.Faces) != 1 {
		t.Errorf("Expected 1 face, got %d", len(mesh.Faces))
	}

	// Check vertex values
	expectedVertices := []Vertex{
		{X: 0.0, Y: 0.0, Z: 0.0},
		{X: 1.0, Y: 0.0, Z: 0.0},
		{X: 0.5, Y: 1.0, Z: 0.0},
	}

	for i, expected := range expectedVertices {
		if mesh.Vertices[i] != expected {
			t.Errorf("Vertex %d: expected %v, got %v", i, expected, mesh.Vertices[i])
		}
	}
}

func TestLoadSTLBinary(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple binary STL file (single triangle)
	testFile := filepath.Join(tmpDir, "test_binary.stl")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	// Write 80-byte header
	header := make([]byte, 80)
	copy(header, []byte("Binary STL test"))
	file.Write(header)

	// Write triangle count (1 triangle)
	binary.Write(file, binary.LittleEndian, uint32(1))

	// Write triangle data
	// Normal
	binary.Write(file, binary.LittleEndian, float32(0.0)) // nx
	binary.Write(file, binary.LittleEndian, float32(0.0)) // ny
	binary.Write(file, binary.LittleEndian, float32(1.0)) // nz

	// Vertex 1
	binary.Write(file, binary.LittleEndian, float32(0.0))
	binary.Write(file, binary.LittleEndian, float32(0.0))
	binary.Write(file, binary.LittleEndian, float32(0.0))

	// Vertex 2
	binary.Write(file, binary.LittleEndian, float32(1.0))
	binary.Write(file, binary.LittleEndian, float32(0.0))
	binary.Write(file, binary.LittleEndian, float32(0.0))

	// Vertex 3
	binary.Write(file, binary.LittleEndian, float32(0.5))
	binary.Write(file, binary.LittleEndian, float32(1.0))
	binary.Write(file, binary.LittleEndian, float32(0.0))

	// Attribute byte count
	binary.Write(file, binary.LittleEndian, uint16(0))

	file.Close()

	// Load STL
	mesh, err := LoadSTL(testFile)
	if err != nil {
		t.Fatalf("LoadSTL() error = %v", err)
	}

	// Check results
	if len(mesh.Vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(mesh.Vertices))
	}

	if len(mesh.Faces) != 1 {
		t.Errorf("Expected 1 face, got %d", len(mesh.Faces))
	}
}

func TestMergeVertices(t *testing.T) {
	// Create mesh with duplicate vertices
	mesh := &Mesh{
		Vertices: []Vertex{
			{X: 0.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 0.0, Z: 0.0},
			{X: 0.0, Y: 1.0, Z: 0.0},
			{X: 0.0, Y: 0.0, Z: 0.0}, // Duplicate of vertex 0
			{X: 1.0, Y: 0.0, Z: 0.0}, // Duplicate of vertex 1
			{X: 0.0, Y: 1.0, Z: 0.0}, // Duplicate of vertex 2
		},
		Faces: []Face{
			{V1: 0, V2: 1, V3: 2},
			{V1: 3, V2: 4, V3: 5},
		},
	}

	merged := mesh.MergeVertices(1e-6)

	// Should have only 3 unique vertices
	if len(merged.Vertices) != 3 {
		t.Errorf("Expected 3 merged vertices, got %d", len(merged.Vertices))
	}

	// Should still have 2 faces
	if len(merged.Faces) != 2 {
		t.Errorf("Expected 2 faces, got %d", len(merged.Faces))
	}

	// Both faces should reference the same 3 vertices
	if merged.Faces[0].V1 != merged.Faces[1].V1 ||
		merged.Faces[0].V2 != merged.Faces[1].V2 ||
		merged.Faces[0].V3 != merged.Faces[1].V3 {
		t.Error("Faces should reference the same merged vertices")
	}
}
