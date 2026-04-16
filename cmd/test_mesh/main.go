package main

import (
	"fmt"
	"os"
	"github.com/whitenhiemer/crochetbot/internal/mesh"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: test_mesh <stl_file>")
		return
	}

	m, err := mesh.LoadSTL(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	m.CalculateBounds()
	width, height, depth := m.GetDimensions()
	
	fmt.Printf("Width: %.2f\n", width)
	fmt.Printf("Height: %.2f\n", height)
	fmt.Printf("Depth: %.2f\n", depth)
	fmt.Printf("\nSphere check: %v\n", m.IsApproximatelySphere())
	fmt.Printf("Cylinder check: %v\n", m.IsApproximatelyCylinder())
	
	// Check ratios
	fmt.Printf("\nWidth/Height: %.2f\n", width/height)
	fmt.Printf("Width/Depth: %.2f\n", width/depth)
	fmt.Printf("Height/Width: %.2f\n", height/width)
	fmt.Printf("Height/Depth: %.2f\n", height/depth)
}
