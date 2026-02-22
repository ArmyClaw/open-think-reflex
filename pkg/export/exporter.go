// Package export provides data export and import functionality for patterns.
package export

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/ArmyClaw/open-think-reflex/pkg/skills"
	"gopkg.in/yaml.v3"
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

// SpaceExportData represents the structure of an exported Space file.
type SpaceExportData struct {
	Version    string            `json:"version"`
	ExportedAt time.Time        `json:"exported_at"`
	Space      *models.Space    `json:"space"`
	Patterns   []models.Pattern `json:"patterns"`
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

// ExportSpaceToJSON exports a Space with its patterns to a JSON file.
func (e *Exporter) ExportSpaceToJSON(ctx context.Context, space *models.Space, patterns []*models.Pattern, filepath string) error {
	// Convert pointers to values for JSON serialization
	patternValues := make([]models.Pattern, len(patterns))
	for i, p := range patterns {
		patternValues[i] = *p
	}

	exportData := SpaceExportData{
		Version:    "2.0",
		ExportedAt: time.Now(),
		Space:      space,
		Patterns:   patternValues,
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

// ExportToYAML exports all patterns to a YAML file.
func (e *Exporter) ExportToYAML(ctx context.Context, patterns []*models.Pattern, filepath string) error {
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

	data, err := yaml.Marshal(exportData)
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ExportSpaceToYAML exports a Space with patterns to a YAML file.
func (e *Exporter) ExportSpaceToYAML(ctx context.Context, space *models.Space, patterns []*models.Pattern, filepath string) error {
	patternValues := make([]models.Pattern, len(patterns))
	for i, p := range patterns {
		patternValues[i] = *p
	}

	exportData := SpaceExportData{
		Version:    "2.0",
		ExportedAt: time.Now(),
		Space:      space,
		Patterns:   patternValues,
	}

	data, err := yaml.Marshal(exportData)
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

// SpaceImportData represents the structure of an imported Space file.
type SpaceImportData struct {
	Version  string             `json:"version"`
	ExportedAt time.Time        `json:"exported_at"`
	Space    *models.Space     `json:"space"`
	Patterns []models.Pattern  `json:"patterns"`
}

// NoteExportData represents the structure of an exported notes file.
type NoteExportData struct {
	Version    string          `json:"version"`
	ExportedAt time.Time       `json:"exported_at"`
	NoteCount  int             `json:"note_count"`
	Notes      []models.Note   `json:"notes"`
}

// NoteImportData represents the structure of an imported notes file.
type NoteImportData struct {
	Version   string         `json:"version"`
	ExportedAt time.Time    `json:"exported_at"`
	NoteCount int           `json:"note_count"`
	Notes     []models.Note `json:"notes"`
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

// ImportSpaceFromJSON imports a Space with patterns from a JSON file.
func (i *Importer) ImportSpaceFromJSON(ctx context.Context, filepath string) (*SpaceImportData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var importData SpaceImportData
	if err := json.Unmarshal(data, &importData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	// Validate import data
	if importData.Patterns == nil {
		importData.Patterns = make([]models.Pattern, 0)
	}

	return &importData, nil
}

// ExportNotesToJSON exports notes to a JSON file.
func (e *Exporter) ExportNotesToJSON(ctx context.Context, notes []*models.Note, filepath string) error {
	noteValues := make([]models.Note, len(notes))
	for i, n := range notes {
		noteValues[i] = *n
	}

	exportData := NoteExportData{
		Version:    "1.0",
		ExportedAt: time.Now(),
		NoteCount:  len(noteValues),
		Notes:      noteValues,
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

// ExportNotesToYAML exports notes to a YAML file.
func (e *Exporter) ExportNotesToYAML(ctx context.Context, notes []*models.Note, filepath string) error {
	noteValues := make([]models.Note, len(notes))
	for i, n := range notes {
		noteValues[i] = *n
	}

	exportData := NoteExportData{
		Version:    "1.0",
		ExportedAt: time.Now(),
		NoteCount:  len(noteValues),
		Notes:      noteValues,
	}

	data, err := yaml.Marshal(exportData)
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ImportNotesFromJSON imports notes from a JSON file.
func (i *Importer) ImportNotesFromJSON(ctx context.Context, filepath string) (*NoteImportData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var importData NoteImportData
	if err := json.Unmarshal(data, &importData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	if importData.Notes == nil {
		importData.Notes = make([]models.Note, 0)
	}

	return &importData, nil
}

// ImportNotesFromYAML imports notes from a YAML file.
func (i *Importer) ImportNotesFromYAML(ctx context.Context, filepath string) (*NoteImportData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var importData NoteImportData
	if err := yaml.Unmarshal(data, &importData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	if importData.Notes == nil {
		importData.Notes = make([]models.Note, 0)
	}

	return &importData, nil
}

// ExportSkillToJSON exports a single skill to a JSON file.
func (e *Exporter) ExportSkillToJSON(skill *skills.Skill, filepath string) error {
	data, err := json.MarshalIndent(skill, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal skill: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ExportSkillToYAML exports a single skill to a YAML file.
func (e *Exporter) ExportSkillToYAML(skill *skills.Skill, filepath string) error {
	data, err := yaml.Marshal(skill)
	if err != nil {
		return fmt.Errorf("failed to marshal skill: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
