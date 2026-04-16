package mesh

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Mesh represents a 3D mesh structure
type Mesh struct {
	Vertices []Vertex
	Faces    []Face
	Bounds   BoundingBox
}

// Vertex represents a 3D point
type Vertex struct {
	X, Y, Z float64
}

// Face represents a triangular face (indices into Vertices)
type Face struct {
	V1, V2, V3 int
}

// BoundingBox represents the mesh bounds
type BoundingBox struct {
	MinX, MinY, MinZ float64
	MaxX, MaxY, MaxZ float64
}

// LoadOBJ loads a .obj file and returns a Mesh
func LoadOBJ(filepath string) (*Mesh, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	mesh := &Mesh{
		Vertices: []Vertex{},
		Faces:    []Face{},
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "v":
			// Vertex: v x y z [w]
			if len(parts) < 4 {
				return nil, fmt.Errorf("line %d: invalid vertex format", lineNum)
			}
			x, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid vertex x coordinate: %w", lineNum, err)
			}
			y, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid vertex y coordinate: %w", lineNum, err)
			}
			z, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid vertex z coordinate: %w", lineNum, err)
			}
			mesh.Vertices = append(mesh.Vertices, Vertex{X: x, Y: y, Z: z})

		case "f":
			// Face: f v1[/vt1][/vn1] v2[/vt2][/vn2] v3[/vt3][/vn3] ...
			if len(parts) < 4 {
				return nil, fmt.Errorf("line %d: face must have at least 3 vertices", lineNum)
			}

			// Parse vertex indices (handle v, v/vt, v/vt/vn, v//vn formats)
			indices := make([]int, 0, len(parts)-1)
			for i := 1; i < len(parts); i++ {
				vertexStr := strings.Split(parts[i], "/")[0]
				idx, err := strconv.Atoi(vertexStr)
				if err != nil {
					return nil, fmt.Errorf("line %d: invalid face vertex index: %w", lineNum, err)
				}

				// OBJ indices are 1-based, convert to 0-based
				if idx > 0 {
					idx--
				} else if idx < 0 {
					// Negative indices count from the end
					idx = len(mesh.Vertices) + idx
				} else {
					return nil, fmt.Errorf("line %d: vertex index cannot be 0", lineNum)
				}

				if idx < 0 || idx >= len(mesh.Vertices) {
					return nil, fmt.Errorf("line %d: vertex index %d out of range", lineNum, idx)
				}

				indices = append(indices, idx)
			}

			// Triangulate face if it has more than 3 vertices
			// Simple fan triangulation from first vertex
			for i := 1; i < len(indices)-1; i++ {
				mesh.Faces = append(mesh.Faces, Face{
					V1: indices[0],
					V2: indices[i],
					V3: indices[i+1],
				})
			}

		case "vt", "vn", "o", "g", "s", "mtllib", "usemtl":
			// Skip texture coords, normals, objects, groups, smoothing, materials
			continue

		default:
			// Skip unknown directives
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(mesh.Vertices) == 0 {
		return nil, fmt.Errorf("no vertices found in OBJ file")
	}

	// Calculate bounding box
	mesh.CalculateBounds()

	return mesh, nil
}

// CalculateBounds computes the bounding box for the mesh
func (m *Mesh) CalculateBounds() {
	if len(m.Vertices) == 0 {
		return
	}

	m.Bounds.MinX, m.Bounds.MaxX = m.Vertices[0].X, m.Vertices[0].X
	m.Bounds.MinY, m.Bounds.MaxY = m.Vertices[0].Y, m.Vertices[0].Y
	m.Bounds.MinZ, m.Bounds.MaxZ = m.Vertices[0].Z, m.Vertices[0].Z

	for _, v := range m.Vertices[1:] {
		if v.X < m.Bounds.MinX {
			m.Bounds.MinX = v.X
		}
		if v.X > m.Bounds.MaxX {
			m.Bounds.MaxX = v.X
		}
		if v.Y < m.Bounds.MinY {
			m.Bounds.MinY = v.Y
		}
		if v.Y > m.Bounds.MaxY {
			m.Bounds.MaxY = v.Y
		}
		if v.Z < m.Bounds.MinZ {
			m.Bounds.MinZ = v.Z
		}
		if v.Z > m.Bounds.MaxZ {
			m.Bounds.MaxZ = v.Z
		}
	}
}

// Volume returns approximate mesh volume
func (m *Mesh) Volume() float64 {
	// TODO: Implement actual volume calculation
	return 0.0
}
