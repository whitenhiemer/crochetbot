package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/whitenhiemer/crochetbot/internal/models"
)

// PatternStore handles in-memory pattern storage
// TODO: Replace with database in production
type PatternStore struct {
	patterns map[string]*models.Pattern
	mu       sync.RWMutex
}

var store = &PatternStore{
	patterns: make(map[string]*models.Pattern),
}

// Save stores a pattern in memory and optionally to disk
func (s *PatternStore) Save(pattern *models.Pattern) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.patterns[pattern.ID] = pattern

	// Also save to disk for persistence
	return s.saveToDisk(pattern)
}

// Get retrieves a pattern by ID
func (s *PatternStore) Get(id string) (*models.Pattern, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pattern, exists := s.patterns[id]
	if !exists {
		// Try loading from disk
		return s.loadFromDisk(id)
	}

	return pattern, nil
}

// List returns all stored patterns
func (s *PatternStore) List() []*models.Pattern {
	s.mu.RLock()
	defer s.mu.RUnlock()

	patterns := make([]*models.Pattern, 0, len(s.patterns))
	for _, pattern := range s.patterns {
		patterns = append(patterns, pattern)
	}

	return patterns
}

// Delete removes a pattern
func (s *PatternStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.patterns, id)

	// Also delete from disk
	return s.deleteFromDisk(id)
}

// saveToDisk writes pattern to JSON file
func (s *PatternStore) saveToDisk(pattern *models.Pattern) error {
	storageDir := getStorageDir()
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	filename := filepath.Join(storageDir, pattern.ID+".json")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create pattern file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(pattern); err != nil {
		return fmt.Errorf("failed to encode pattern: %w", err)
	}

	return nil
}

// loadFromDisk reads pattern from JSON file
func (s *PatternStore) loadFromDisk(id string) (*models.Pattern, error) {
	storageDir := getStorageDir()
	filename := filepath.Join(storageDir, id+".json")

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("pattern not found: %w", err)
	}
	defer file.Close()

	var pattern models.Pattern
	if err := json.NewDecoder(file).Decode(&pattern); err != nil {
		return nil, fmt.Errorf("failed to decode pattern: %w", err)
	}

	// Cache in memory
	s.patterns[id] = &pattern

	return &pattern, nil
}

// deleteFromDisk removes pattern file
func (s *PatternStore) deleteFromDisk(id string) error {
	storageDir := getStorageDir()
	filename := filepath.Join(storageDir, id+".json")

	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete pattern file: %w", err)
	}

	return nil
}

// getStorageDir returns the pattern storage directory
func getStorageDir() string {
	dir := os.Getenv("PATTERN_STORAGE_DIR")
	if dir == "" {
		dir = "./data/patterns"
	}
	return dir
}
