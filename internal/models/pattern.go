package models

import "time"

// Pattern represents a complete crochet pattern
type Pattern struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	CreatedAt        time.Time      `json:"created_at"`
	Description      string         `json:"description"`
	Difficulty       string         `json:"difficulty"` // beginner, intermediate, advanced
	Parts            []Part         `json:"parts"`
	Materials        Materials      `json:"materials"`
	Assembly         []string       `json:"assembly_instructions"`
	FinishedSize     FinishedSize   `json:"finished_size"`
	AccuracyMetrics  AccuracyMetrics `json:"accuracy_metrics"`
}

// AccuracyMetrics shows how well the pattern matches the original model
type AccuracyMetrics struct {
	ShapeMatchPercent float64 `json:"shape_match_percent"` // 0-100
	AverageError      float64 `json:"average_error"`       // average radius difference
	MaxError          float64 `json:"max_error"`           // maximum radius difference
	Notes             string  `json:"notes"`
}

// Part represents a single piece of the amigurumi
type Part struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"` // sphere, cylinder, cone, etc.
	Rounds       []Round  `json:"rounds"`
	Color        string   `json:"color"`
	StartingType string   `json:"starting_type"` // magic ring, chain
	Notes        []string `json:"notes"`
}

// Round represents one round/row of crochet
type Round struct {
	Number       int      `json:"number"`
	Instructions string   `json:"instructions"`
	StitchCount  int      `json:"stitch_count"`
	StitchType   string   `json:"stitch_type"` // sc, hdc, dc, inc, dec
	Repeats      int      `json:"repeats"`
	Notes        string   `json:"notes"`
}

// Materials lists required materials for the pattern
type Materials struct {
	YarnWeight    string  `json:"yarn_weight"`    // DK, worsted, etc.
	YarnYardage   int     `json:"yarn_yardage"`   // total yards needed
	HookSize      string  `json:"hook_size"`      // 3.5mm, E/4, etc.
	Colors        []Color `json:"colors"`
	OtherSupplies []string `json:"other_supplies"` // stuffing, safety eyes, etc.
}

// Color represents yarn color info
type Color struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"` // yards
}

// FinishedSize represents the expected dimensions of the finished piece
type FinishedSize struct {
	HeightInches float64 `json:"height_inches"`
	HeightCm     float64 `json:"height_cm"`
	WidthInches  float64 `json:"width_inches"`
	WidthCm      float64 `json:"width_cm"`
}

// Mesh3D represents uploaded 3D model data
type Mesh3D struct {
	ID         string    `json:"id"`
	Filename   string    `json:"filename"`
	UploadedAt time.Time `json:"uploaded_at"`
	Vertices   int       `json:"vertices"`
	Faces      int       `json:"faces"`
	Format     string    `json:"format"` // obj, stl
}
