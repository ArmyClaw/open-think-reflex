package sqlite

import (
	"context"
	"testing"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// TestSavePatternsBatch tests batch saving of patterns (Iter 44)
func TestSavePatternsBatch(t *testing.T) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	storage := NewStorage(db)

	// Create multiple patterns
	patterns := []*models.Pattern{
		{
			ID:       "batch-1",
			Trigger:  "trigger1",
			Response: "response1",
			Strength: 50,
		},
		{
			ID:       "batch-2",
			Trigger:  "trigger2",
			Response: "response2",
			Strength: 60,
		},
		{
			ID:       "batch-3",
			Trigger:  "trigger3",
			Response: "response3",
			Strength: 70,
		},
	}

	// Test batch save
	err = storage.SavePatternsBatch(ctx, patterns)
	if err != nil {
		t.Fatalf("SavePatternsBatch failed: %v", err)
	}

	// Verify all patterns were saved
	for _, p := range patterns {
		got, err := storage.GetPattern(ctx, p.ID)
		if err != nil {
			t.Errorf("failed to get pattern %s: %v", p.ID, err)
			continue
		}
		if got.Trigger != p.Trigger {
			t.Errorf("expected trigger %s, got %s", p.Trigger, got.Trigger)
		}
	}
	t.Logf("✅ SavePatternsBatch: saved and verified %d patterns", len(patterns))
}

// TestDeletePatternsBatch tests batch deletion of patterns (Iter 44)
func TestDeletePatternsBatch(t *testing.T) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	storage := NewStorage(db)

	// First save some patterns
	patterns := []*models.Pattern{
		{ID: "del-1", Trigger: "trigger1", Response: "response1"},
		{ID: "del-2", Trigger: "trigger2", Response: "response2"},
		{ID: "del-3", Trigger: "trigger3", Response: "response3"},
	}

	err = storage.SavePatternsBatch(ctx, patterns)
	if err != nil {
		t.Fatalf("failed to save patterns: %v", err)
	}

	// Delete first two
	ids := []string{"del-1", "del-2"}
	err = storage.DeletePatternsBatch(ctx, ids)
	if err != nil {
		t.Fatalf("DeletePatternsBatch failed: %v", err)
	}

	// Verify deletion
	_, err = storage.GetPattern(ctx, "del-1")
	if err == nil {
		t.Error("expected error for deleted pattern del-1")
	}

	_, err = storage.GetPattern(ctx, "del-2")
	if err == nil {
		t.Error("expected error for deleted pattern del-2")
	}

	// Verify third pattern still exists
	_, err = storage.GetPattern(ctx, "del-3")
	if err != nil {
		t.Errorf("pattern del-3 should still exist: %v", err)
	}

	t.Logf("✅ DeletePatternsBatch: deleted %d patterns", len(ids))
}

// TestUpdatePatternsBatch tests batch update of patterns (Iter 44)
func TestUpdatePatternsBatch(t *testing.T) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	storage := NewStorage(db)

	// Save initial patterns
	patterns := []*models.Pattern{
		{ID: "update-1", Trigger: "old1", Response: "old-response-1"},
		{ID: "update-2", Trigger: "old2", Response: "old-response-2"},
	}

	err = storage.SavePatternsBatch(ctx, patterns)
	if err != nil {
		t.Fatalf("failed to save patterns: %v", err)
	}

	// Update patterns
	updated := []*models.Pattern{
		{ID: "update-1", Trigger: "new1", Response: "new-response-1", Strength: 80},
		{ID: "update-2", Trigger: "new2", Response: "new-response-2", Strength: 90},
	}

	err = storage.UpdatePatternsBatch(ctx, updated)
	if err != nil {
		t.Fatalf("UpdatePatternsBatch failed: %v", err)
	}

	// Verify updates
	for _, p := range updated {
		got, err := storage.GetPattern(ctx, p.ID)
		if err != nil {
			t.Errorf("failed to get pattern %s: %v", p.ID, err)
			continue
		}
		if got.Trigger != p.Trigger {
			t.Errorf("expected trigger %s, got %s", p.Trigger, got.Trigger)
		}
		if got.Strength != p.Strength {
			t.Errorf("expected strength %f, got %f", p.Strength, got.Strength)
		}
	}

	t.Logf("✅ UpdatePatternsBatch: updated %d patterns", len(updated))
}

// TestBatchEmptyInput tests that batch operations handle empty input gracefully
func TestBatchEmptyInput(t *testing.T) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	storage := NewStorage(db)

	// Test empty batch operations should not error
	err = storage.SavePatternsBatch(ctx, nil)
	if err != nil {
		t.Errorf("SavePatternsBatch(nil) should not error: %v", err)
	}

	err = storage.SavePatternsBatch(ctx, []*models.Pattern{})
	if err != nil {
		t.Errorf("SavePatternsBatch([]) should not error: %v", err)
	}

	err = storage.DeletePatternsBatch(ctx, nil)
	if err != nil {
		t.Errorf("DeletePatternsBatch(nil) should not error: %v", err)
	}

	err = storage.DeletePatternsBatch(ctx, []string{})
	if err != nil {
		t.Errorf("DeletePatternsBatch([]) should not error: %v", err)
	}

	err = storage.UpdatePatternsBatch(ctx, nil)
	if err != nil {
		t.Errorf("UpdatePatternsBatch(nil) should not error: %v", err)
	}

	err = storage.UpdatePatternsBatch(ctx, []*models.Pattern{})
	if err != nil {
		t.Errorf("UpdatePatternsBatch([]) should not error: %v", err)
	}

	t.Log("✅ Batch operations handle empty input gracefully")
}
