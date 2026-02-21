package sqlite

import (
	"context"
	"os"
	"testing"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

func setupTestDB(t *testing.T) (*Database, func()) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	db, err := NewDatabase(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to create database: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(tmpFile.Name())
	}

	return db, cleanup
}

func TestNewDatabase(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := NewDatabase(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewDatabase failed: %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Fatal("Database should not be nil")
	}
}

func TestDatabase_InvalidPath(t *testing.T) {
	_, err := NewDatabase("/nonexistent/path/test.db")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestStorage_SavePattern(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewStorage(db)
	ctx := context.Background()

	pattern := models.NewPattern("test trigger", "test response")

	err := storage.SavePattern(ctx, pattern)
	if err != nil {
		t.Fatalf("SavePattern failed: %v", err)
	}

	// Verify saved
	if pattern.ID == "" {
		t.Error("Pattern ID should be set after save")
	}
}

func TestStorage_GetPattern(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewStorage(db)
	ctx := context.Background()

	// Create and save pattern
	pattern := models.NewPattern("test trigger", "test response")
	err := storage.SavePattern(ctx, pattern)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve
	retrieved, err := storage.GetPattern(ctx, pattern.ID)
	if err != nil {
		t.Fatalf("GetPattern failed: %v", err)
	}

	if retrieved.ID != pattern.ID {
		t.Errorf("Expected ID %s, got %s", pattern.ID, retrieved.ID)
	}
	if retrieved.Trigger != "test trigger" {
		t.Errorf("Expected trigger 'test trigger', got '%s'", retrieved.Trigger)
	}
}

func TestStorage_GetPattern_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewStorage(db)
	ctx := context.Background()

	_, err := storage.GetPattern(ctx, "nonexistent-id")
	if err == nil {
		t.Error("Expected error for nonexistent pattern")
	}
}

func TestStorage_ListPatterns(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewStorage(db)
	ctx := context.Background()

	// Save multiple patterns
	for i := 0; i < 5; i++ {
		p := models.NewPattern("trigger", "response")
		p.Project = "test-project"
		if err := storage.SavePattern(ctx, p); err != nil {
			t.Fatal(err)
		}
	}

	// List all
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}

	if len(patterns) != 5 {
		t.Errorf("Expected 5 patterns, got %d", len(patterns))
	}
}

func TestStorage_UpdatePattern(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewStorage(db)
	ctx := context.Background()

	// Create pattern
	pattern := models.NewPattern("original", "original response")
	err := storage.SavePattern(ctx, pattern)
	if err != nil {
		t.Fatal(err)
	}

	// Update
	pattern.Trigger = "updated"
	pattern.Response = "updated response"
	err = storage.UpdatePattern(ctx, pattern)
	if err != nil {
		t.Fatal(err)
	}

	// Verify
	updated, err := storage.GetPattern(ctx, pattern.ID)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Trigger != "updated" {
		t.Errorf("Expected trigger 'updated', got '%s'", updated.Trigger)
	}
}

func TestStorage_DeletePattern(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewStorage(db)
	ctx := context.Background()

	// Create pattern
	pattern := models.NewPattern("to delete", "response")
	err := storage.SavePattern(ctx, pattern)
	if err != nil {
		t.Fatal(err)
	}

	// Delete
	err = storage.DeletePattern(ctx, pattern.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Verify deleted
	_, err = storage.GetPattern(ctx, pattern.ID)
	if err == nil {
		t.Error("Expected error after deletion")
	}
}

func TestStorage_CreateSpace(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewStorage(db)
	ctx := context.Background()

	space := &models.Space{
		ID:          "space-1",
		Name:        "Test Space",
		Description: "A test space",
	}

	err := storage.CreateSpace(ctx, space)
	if err != nil {
		t.Fatalf("CreateSpace failed: %v", err)
	}
}

func TestStorage_ListSpaces(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewStorage(db)
	ctx := context.Background()

	// Create spaces
	for i := 0; i < 3; i++ {
		space := &models.Space{
			ID:   "space-" + string(rune('1'+i)),
			Name: "Space " + string(rune('A'+i)),
		}
		if err := storage.CreateSpace(ctx, space); err != nil {
			t.Fatal(err)
		}
	}

	// List
	spaces, err := storage.ListSpaces(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(spaces) != 3 {
		t.Errorf("Expected 3 spaces, got %d", len(spaces))
	}
}

func TestStorage_Close(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := NewDatabase(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	storage := NewStorage(db)

	// Close should not panic
	storage.Close()

	// Double close should also not panic
	storage.Close()
}
