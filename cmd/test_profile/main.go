package main

import (
	"fmt"
	"os"
	"github.com/whitenhiemer/crochetbot/internal/mesh"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: test_profile <stl_file>")
		return
	}

	m, err := mesh.LoadSTL(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	m.CalculateBounds()
	
	// Get radius profile for 40 slices
	profile := m.GetRadiusProfile(40)
	
	fmt.Printf("Radius profile (%d slices):\n", len(profile))
	
	// Find max
	maxRadius := 0.0
	for _, r := range profile {
		if r > maxRadius {
			maxRadius = r
		}
	}
	fmt.Printf("Max radius: %.2f\n\n", maxRadius)
	
	// Show normalized with visual bar
	fmt.Println("Normalized profile (bottom to top):")
	for i, radius := range profile {
		pct := (radius / maxRadius) * 100
		bars := int(pct / 2) // Scale for display
		barStr := ""
		for j := 0; j < bars; j++ {
			barStr += "█"
		}
		fmt.Printf("Slice %2d: %3.0f%% %s\n", i, pct, barStr)
	}
}
