package models

import (
	"testing"
	"time"
)

func TestSpace_DefaultSpaces(t *testing.T) {
	spaces := DefaultSpaces()

	if len(spaces) != 3 {
		t.Errorf("Expected 3 default spaces, got %d", len(spaces))
	}

	// Check global space
	if spaces[0].ID != "global" {
		t.Errorf("Expected first space ID 'global', got '%s'", spaces[0].ID)
	}

	// Check default flag
	if !spaces[0].DefaultSpace {
		t.Error("Global space should be default")
	}
}

func TestSpace_New(t *testing.T) {
	space := &Space{
		ID:          "test-id",
		Name:        "Test Space",
		Description: "A test space",
	}

	if space.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", space.ID)
	}
	if space.Name != "Test Space" {
		t.Errorf("Expected Name 'Test Space', got '%s'", space.Name)
	}
}

func TestSpace_Fields(t *testing.T) {
	now := time.Now()
	space := &Space{
		ID:           "id-1",
		Name:         "Name",
		Description:  "Description",
		CreatedAt:    now,
		UpdatedAt:    now,
		DefaultSpace: true,
		PatternLimit: 100,
		PatternCount:  50,
	}

	// Verify all fields
	if space.ID != "id-1" {
		t.Errorf("ID mismatch")
	}
	if space.Name != "Name" {
		t.Errorf("Name mismatch")
	}
	if !space.DefaultSpace {
		t.Error("DefaultSpace should be true")
	}
	if space.PatternLimit != 100 {
		t.Errorf("PatternLimit expected 100, got %d", space.PatternLimit)
	}
	if space.PatternCount != 50 {
		t.Errorf("PatternCount expected 50, got %d", space.PatternCount)
	}
}
