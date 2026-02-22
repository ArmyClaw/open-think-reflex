package export

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

func TestExporter(t *testing.T) {
	e := NewExporter()
	if e == nil {
		t.Error("NewExporter returned nil")
	}
}

func TestImporter(t *testing.T) {
	i := NewImporter()
	if i == nil {
		t.Error("NewImporter returned nil")
	}
}

func TestExportToJSON(t *testing.T) {
	e := NewExporter()
	ctx := context.Background()

	patterns := []*models.Pattern{
		{
			ID:       "test-1",
			Trigger:  "test trigger",
			Response: "test response",
		},
	}

	tmpFile, err := os.CreateTemp("", "export-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = e.ExportToJSON(ctx, patterns, tmpFile.Name())
	if err != nil {
		t.Errorf("ExportToJSON failed: %v", err)
	}

	// Verify file exists and has content
	info, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Error("file not created")
	}
	if info.Size() == 0 {
		t.Error("file is empty")
	}
}

func TestExportDataJSON(t *testing.T) {
	data := ExportData{
		Version:     "1.0",
		PatternCount: 1,
		Patterns: []models.Pattern{
			{ID: "test", Trigger: "trigger", Response: "response"},
		},
	}

	// Should be marshalable
	_, err := json.Marshal(data)
	if err != nil {
		t.Errorf("ExportData marshal failed: %v", err)
	}
}

func TestSpaceExportDataJSON(t *testing.T) {
	data := SpaceExportData{
		Version: "2.0",
		Space: &models.Space{
			ID:   "space-1",
			Name: "Test Space",
		},
		Patterns: []models.Pattern{
			{ID: "test", Trigger: "trigger", Response: "response"},
		},
	}

	// Should be marshalable
	_, err := json.Marshal(data)
	if err != nil {
		t.Errorf("SpaceExportData marshal failed: %v", err)
	}
}
