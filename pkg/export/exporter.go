// Package export provides data export and import functionality for patterns.
package export

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// Exporter handles exporting patterns to various formats.
type Exporter struct{}

// NewExporter creates a new Exporter instance.
func NewExporter() *Exporter {
	return &Exporter{}
}

// ExportData represents the structure of an exported data file.
type ExportData struct {
	Version     string            `json:"version"`
	ExportedAt  time.Time         `json:"exported_at"`
	PatternCount int              `json:"pattern_count"`
	Patterns    []models.Pattern `json:"patterns"`
}

// ExportToJSON exports all patterns to a JSON file.
func (e *Exporter) ExportToJSON(ctx context.Context, patterns []*models.Pattern, filepath string) error {
	// Convert pointers to values for JSON serialization
	patternValues := make([]models.Pattern, len(patterns))
	for i, p := range patterns {
		patternValues[i] = *p
	}

	exportData := ExportData{
		Version:     "1.0",
		ExportedAt:  time.Now(),
		PatternCount: len(patternValues),
		Patterns:    patternValues,
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ImportData represents the structure of an imported data file.
type ImportData struct {
	Version     string            `json:"version"`
	ExportedAt  time.Time         `json:"exported_at"`
	PatternCount int              `json:"pattern_count"`
	Patterns    []models.Pattern `json:"patterns"`
}

// Importer handles importing patterns from various formats.
type Importer struct{}

// NewImporter creates a new Importer instance.
func NewImporter() *Importer {
	return &Importer{}
}

// ImportFromJSON imports patterns from a JSON file.
func (i *Importer) ImportFromJSON(ctx context.Context, filepath string) (*ImportData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var importData ImportData
	if err := json.Unmarshal(data, &importData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	// Validate import data
	if importData.Patterns == nil {
		importData.Patterns = make([]models.Pattern, 0)
	}

	return &importData, nil
}
