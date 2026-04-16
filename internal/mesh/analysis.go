package mesh

import "math"

// GetDimensions returns the width, height, and depth of the mesh
func (m *Mesh) GetDimensions() (width, height, depth float64) {
	width = m.Bounds.MaxX - m.Bounds.MinX
	height = m.Bounds.MaxY - m.Bounds.MinY
	depth = m.Bounds.MaxZ - m.Bounds.MinZ
	return
}

// GetCenter returns the center point of the mesh
func (m *Mesh) GetCenter() Vertex {
	return Vertex{
		X: (m.Bounds.MinX + m.Bounds.MaxX) / 2,
		Y: (m.Bounds.MinY + m.Bounds.MaxY) / 2,
		Z: (m.Bounds.MinZ + m.Bounds.MinZ) / 2,
	}
}

// IsApproximatelySphere checks if the mesh is roughly spherical
func (m *Mesh) IsApproximatelySphere() bool {
	width, height, depth := m.GetDimensions()

	// Calculate ratios between dimensions
	maxDim := math.Max(width, math.Max(height, depth))
	minDim := math.Min(width, math.Min(height, depth))

	if maxDim == 0 {
		return false
	}

	// If ratio between max and min dimension is close to 1, it's sphere-like
	ratio := minDim / maxDim
	return ratio > 0.8 // Within 20% tolerance
}

// IsApproximatelyCylinder checks if the mesh is roughly cylindrical
func (m *Mesh) IsApproximatelyCylinder() bool {
	width, height, depth := m.GetDimensions()

	// Sort dimensions
	dims := []float64{width, height, depth}
	for i := 0; i < len(dims)-1; i++ {
		for j := i + 1; j < len(dims); j++ {
			if dims[i] < dims[j] {
				dims[i], dims[j] = dims[j], dims[i]
			}
		}
	}

	// For cylinder: one dimension significantly larger, other two similar
	if dims[0] == 0 {
		return false
	}

	lengthRatio := dims[0] / dims[1]
	crossRatio := dims[1] / dims[2]

	// Long in one direction, circular in cross-section
	return lengthRatio > 1.5 && crossRatio > 0.8 && crossRatio < 1.25
}

// GetAverageRadius returns the average distance from center to vertices
func (m *Mesh) GetAverageRadius() float64 {
	if len(m.Vertices) == 0 {
		return 0
	}

	center := m.GetCenter()
	totalDist := 0.0

	for _, v := range m.Vertices {
		dx := v.X - center.X
		dy := v.Y - center.Y
		dz := v.Z - center.Z
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
		totalDist += dist
	}

	return totalDist / float64(len(m.Vertices))
}

// SurfaceArea calculates approximate surface area using triangle faces
func (m *Mesh) SurfaceArea() float64 {
	area := 0.0

	for _, face := range m.Faces {
		if face.V1 >= len(m.Vertices) || face.V2 >= len(m.Vertices) || face.V3 >= len(m.Vertices) {
			continue
		}

		v1 := m.Vertices[face.V1]
		v2 := m.Vertices[face.V2]
		v3 := m.Vertices[face.V3]

		// Calculate triangle area using cross product
		// AB = v2 - v1, AC = v3 - v1
		abX := v2.X - v1.X
		abY := v2.Y - v1.Y
		abZ := v2.Z - v1.Z

		acX := v3.X - v1.X
		acY := v3.Y - v1.Y
		acZ := v3.Z - v1.Z

		// Cross product
		crossX := abY*acZ - abZ*acY
		crossY := abZ*acX - abX*acZ
		crossZ := abX*acY - abY*acX

		// Area = |cross| / 2
		crossMag := math.Sqrt(crossX*crossX + crossY*crossY + crossZ*crossZ)
		area += crossMag / 2.0
	}

	return area
}

// EstimateVolume calculates approximate volume (signed volume method)
func (m *Mesh) EstimateVolume() float64 {
	volume := 0.0

	for _, face := range m.Faces {
		if face.V1 >= len(m.Vertices) || face.V2 >= len(m.Vertices) || face.V3 >= len(m.Vertices) {
			continue
		}

		v1 := m.Vertices[face.V1]
		v2 := m.Vertices[face.V2]
		v3 := m.Vertices[face.V3]

		// Signed volume of tetrahedron formed by origin and triangle
		volume += v1.X*(v2.Y*v3.Z-v2.Z*v3.Y) -
			v1.Y*(v2.X*v3.Z-v2.Z*v3.X) +
			v1.Z*(v2.X*v3.Y-v2.Y*v3.X)
	}

	return math.Abs(volume) / 6.0
}

// GetRadiusProfile calculates average radius at different heights
// Returns slice of radii from bottom to top
func (m *Mesh) GetRadiusProfile(numSlices int) []float64 {
	if len(m.Vertices) == 0 || numSlices <= 0 {
		return []float64{}
	}

	center := m.GetCenter()
	_, height, _ := m.GetDimensions()
	sliceHeight := height / float64(numSlices)

	radii := make([]float64, numSlices)
	counts := make([]int, numSlices)

	// For each vertex, determine which slice it belongs to and accumulate radius
	for _, v := range m.Vertices {
		// Height relative to bottom
		relativeHeight := v.Y - m.Bounds.MinY
		sliceIdx := int(relativeHeight / sliceHeight)
		if sliceIdx >= numSlices {
			sliceIdx = numSlices - 1
		}
		if sliceIdx < 0 {
			sliceIdx = 0
		}

		// Calculate distance from center in XZ plane
		dx := v.X - center.X
		dz := v.Z - center.Z
		radius := math.Sqrt(dx*dx + dz*dz)

		radii[sliceIdx] += radius
		counts[sliceIdx]++
	}

	// Average the radii
	for i := 0; i < numSlices; i++ {
		if counts[i] > 0 {
			radii[i] /= float64(counts[i])
		}
	}

	return radii
}
