package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/whitenhiemer/crochetbot/internal/mesh"
	"github.com/whitenhiemer/crochetbot/internal/pattern"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/test-obj/main.go <path-to-obj-file>")
		fmt.Println("Example: go run cmd/test-obj/main.go test/sphere.obj")
		os.Exit(1)
	}

	objFile := os.Args[1]
	fmt.Printf("Loading OBJ file: %s\n", objFile)

	// Load mesh
	m, err := mesh.LoadOBJ(objFile)
	if err != nil {
		log.Fatalf("Failed to load OBJ: %v", err)
	}

	fmt.Printf("\nMesh loaded successfully!\n")
	fmt.Printf("  Vertices: %d\n", len(m.Vertices))
	fmt.Printf("  Faces: %d\n", len(m.Faces))
	fmt.Printf("  Bounds: X[%.2f, %.2f] Y[%.2f, %.2f] Z[%.2f, %.2f]\n",
		m.Bounds.MinX, m.Bounds.MaxX,
		m.Bounds.MinY, m.Bounds.MaxY,
		m.Bounds.MinZ, m.Bounds.MaxZ)

	// Generate pattern
	fmt.Printf("\nGenerating crochet pattern...\n")
	gen := pattern.NewGenerator()
	pat, err := gen.Generate(m)
	if err != nil {
		log.Fatalf("Failed to generate pattern: %v", err)
	}

	// Print pattern as JSON
	patternJSON, err := json.MarshalIndent(pat, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal pattern: %v", err)
	}

	fmt.Printf("\n=== Generated Pattern ===\n")
	fmt.Println(string(patternJSON))

	// Print human-readable pattern
	fmt.Printf("\n=== Human-Readable Pattern ===\n")
	fmt.Printf("Pattern: %s\n", pat.Name)
	fmt.Printf("Difficulty: %s\n\n", pat.Difficulty)

	fmt.Printf("Materials:\n")
	fmt.Printf("  - Yarn: %s weight, ~%d yards\n", pat.Materials.YarnWeight, pat.Materials.YarnYardage)
	fmt.Printf("  - Hook: %s\n", pat.Materials.HookSize)
	for _, supply := range pat.Materials.OtherSupplies {
		fmt.Printf("  - %s\n", supply)
	}

	for _, part := range pat.Parts {
		fmt.Printf("\n%s (%s):\n", part.Name, part.Type)
		fmt.Printf("Starting: %s\n\n", part.StartingType)

		for _, round := range part.Rounds {
			fmt.Printf("Round %d: %s (%d sts)\n", round.Number, round.Instructions, round.StitchCount)
			if round.Notes != "" {
				fmt.Printf("  Note: %s\n", round.Notes)
			}
		}

		if len(part.Notes) > 0 {
			fmt.Printf("\nNotes:\n")
			for _, note := range part.Notes {
				fmt.Printf("  - %s\n", note)
			}
		}
	}
}
