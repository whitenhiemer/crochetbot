package mesh

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// LoadSTL loads an STL file (ASCII or Binary format) and returns a Mesh
func LoadSTL(filepath string) (*Mesh, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read first 5 bytes to determine format
	header := make([]byte, 5)
	_, err = file.Read(header)
	if err != nil {
		return nil, fmt.Errorf("failed to read file header: %w", err)
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Check if ASCII (starts with "solid")
	if string(header) == "solid" {
		return loadSTLAscii(file)
	}

	return loadSTLBinary(file)
}

// loadSTLAscii loads ASCII format STL file
func loadSTLAscii(file *os.File) (*Mesh, error) {
	mesh := &Mesh{
		Vertices: []Vertex{},
		Faces:    []Face{},
	}

	scanner := bufio.NewScanner(file)
	var currentTriangle []Vertex
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and solid/endsolid
		if line == "" || strings.HasPrefix(line, "solid") || strings.HasPrefix(line, "endsolid") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "facet":
			// Start of a new triangle, normal follows
			currentTriangle = []Vertex{}

		case "vertex":
			// vertex x y z
			if len(parts) != 4 {
				return nil, fmt.Errorf("line %d: invalid vertex format", lineNum)
			}
			x, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid x coordinate: %w", lineNum, err)
			}
			y, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid y coordinate: %w", lineNum, err)
			}
			z, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid z coordinate: %w", lineNum, err)
			}
			currentTriangle = append(currentTriangle, Vertex{X: x, Y: y, Z: z})

		case "endfacet":
			// End of triangle, add to mesh
			if len(currentTriangle) != 3 {
				return nil, fmt.Errorf("line %d: triangle must have exactly 3 vertices", lineNum)
			}

			// Add vertices
			v1Idx := len(mesh.Vertices)
			mesh.Vertices = append(mesh.Vertices, currentTriangle[0])
			mesh.Vertices = append(mesh.Vertices, currentTriangle[1])
			mesh.Vertices = append(mesh.Vertices, currentTriangle[2])

			// Add face
			mesh.Faces = append(mesh.Faces, Face{
				V1: v1Idx,
				V2: v1Idx + 1,
				V3: v1Idx + 2,
			})

			currentTriangle = nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(mesh.Vertices) == 0 {
		return nil, fmt.Errorf("no vertices found in STL file")
	}

	// Calculate bounding box
	mesh.CalculateBounds()

	return mesh, nil
}

// loadSTLBinary loads binary format STL file
func loadSTLBinary(file *os.File) (*Mesh, error) {
	mesh := &Mesh{
		Vertices: []Vertex{},
		Faces:    []Face{},
	}

	// Binary STL format:
	// 80 bytes header
	// 4 bytes uint32 triangle count
	// For each triangle:
	//   12 bytes (3x float32) - normal
	//   12 bytes (3x float32) - vertex 1
	//   12 bytes (3x float32) - vertex 2
	//   12 bytes (3x float32) - vertex 3
	//   2 bytes uint16 - attribute byte count (unused)

	// Skip 80-byte header
	header := make([]byte, 80)
	_, err := io.ReadFull(file, header)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Read triangle count
	var triangleCount uint32
	err = binary.Read(file, binary.LittleEndian, &triangleCount)
	if err != nil {
		return nil, fmt.Errorf("failed to read triangle count: %w", err)
	}

	// Read each triangle
	for i := uint32(0); i < triangleCount; i++ {
		// Read normal (3 float32s) - we don't use it
		var normal [3]float32
		err = binary.Read(file, binary.LittleEndian, &normal)
		if err != nil {
			return nil, fmt.Errorf("failed to read normal for triangle %d: %w", i, err)
		}

		// Read 3 vertices
		var vertices [3][3]float32
		for j := 0; j < 3; j++ {
			err = binary.Read(file, binary.LittleEndian, &vertices[j])
			if err != nil {
				return nil, fmt.Errorf("failed to read vertex %d for triangle %d: %w", j, i, err)
			}
		}

		// Read attribute byte count (unused)
		var attributeByteCount uint16
		err = binary.Read(file, binary.LittleEndian, &attributeByteCount)
		if err != nil {
			return nil, fmt.Errorf("failed to read attribute byte count for triangle %d: %w", i, err)
		}

		// Add vertices to mesh
		v1Idx := len(mesh.Vertices)
		mesh.Vertices = append(mesh.Vertices, Vertex{
			X: float64(vertices[0][0]),
			Y: float64(vertices[0][1]),
			Z: float64(vertices[0][2]),
		})
		mesh.Vertices = append(mesh.Vertices, Vertex{
			X: float64(vertices[1][0]),
			Y: float64(vertices[1][1]),
			Z: float64(vertices[1][2]),
		})
		mesh.Vertices = append(mesh.Vertices, Vertex{
			X: float64(vertices[2][0]),
			Y: float64(vertices[2][1]),
			Z: float64(vertices[2][2]),
		})

		// Add face
		mesh.Faces = append(mesh.Faces, Face{
			V1: v1Idx,
			V2: v1Idx + 1,
			V3: v1Idx + 2,
		})
	}

	if len(mesh.Vertices) == 0 {
		return nil, fmt.Errorf("no vertices found in STL file")
	}

	// Calculate bounding box
	mesh.CalculateBounds()

	return mesh, nil
}

// MergeVertices merges duplicate vertices in the mesh (useful for STL files)
// Returns a new mesh with merged vertices
func (m *Mesh) MergeVertices(tolerance float64) *Mesh {
	if tolerance == 0 {
		tolerance = 1e-6 // Default tolerance
	}

	vertexMap := make(map[string]int)
	newVertices := []Vertex{}
	vertexMapping := make([]int, len(m.Vertices))

	// Helper to create a key for vertex matching
	makeKey := func(v Vertex) string {
		x := math.Round(v.X/tolerance) * tolerance
		y := math.Round(v.Y/tolerance) * tolerance
		z := math.Round(v.Z/tolerance) * tolerance
		return fmt.Sprintf("%.6f,%.6f,%.6f", x, y, z)
	}

	// Build vertex map and mapping
	for i, v := range m.Vertices {
		key := makeKey(v)
		if existingIdx, exists := vertexMap[key]; exists {
			vertexMapping[i] = existingIdx
		} else {
			newIdx := len(newVertices)
			vertexMap[key] = newIdx
			vertexMapping[i] = newIdx
			newVertices = append(newVertices, v)
		}
	}

	// Remap faces
	newFaces := make([]Face, len(m.Faces))
	for i, face := range m.Faces {
		newFaces[i] = Face{
			V1: vertexMapping[face.V1],
			V2: vertexMapping[face.V2],
			V3: vertexMapping[face.V3],
		}
	}

	newMesh := &Mesh{
		Vertices: newVertices,
		Faces:    newFaces,
	}
	newMesh.CalculateBounds()

	return newMesh
}
