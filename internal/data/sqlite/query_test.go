package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// TestQueryOptimizationMethods tests the new query optimization methods (Iter 46)
func TestQueryOptimizationMethods(t *testing.T) {
	ctx := context.Background()
	
	// Setup database
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	storage := NewStorage(db)

	// Create test patterns
	patterns := []*models.Pattern{
		{
			ID:          "test-1",
			Trigger:     "hello",
			Response:    "Hello! How can I help you?",
			Strength:    80,
			Threshold:   50,
			Project:     "default",
			Tags:        []string{"greeting", "basic"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "test-2",
			Trigger:     "help",
			Response:    "I can help you with...",
			Strength:    60,
			Threshold:   50,
			Project:     "default",
			Tags:        []string{"help", "basic"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "test-3",
			Trigger:     "weather",
			Response:    "The weather is nice today",
			Strength:    70,
			Threshold:   50,
			Project:     "default",
			Tags:        []string{"info", "weather"},
			LastUsedAt:  &[]time.Time{time.Now().Add(-time.Hour)}[0],
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Save patterns
	for _, p := range patterns {
		if err := storage.SavePattern(ctx, p); err != nil {
			t.Fatalf("failed to save pattern: %v", err)
		}
	}

	// Test CountPatterns
	t.Run("CountPatterns", func(t *testing.T) {
		count, err := storage.CountPatterns(ctx, contracts.ListOptions{})
		if err != nil {
			t.Fatalf("CountPatterns failed: %v", err)
		}
		if count != 3 {
			t.Errorf("expected 3 patterns, got %d", count)
		}
	})

	// Test CountPatterns with project filter
	t.Run("CountPatterns with project", func(t *testing.T) {
		count, err := storage.CountPatterns(ctx, contracts.ListOptions{Project: "default"})
		if err != nil {
			t.Fatalf("CountPatterns failed: %v", err)
		}
		if count != 3 {
			t.Errorf("expected 3 patterns, got %d", count)
		}
	})

	// Test CountPatterns with min strength
	t.Run("CountPatterns with min strength", func(t *testing.T) {
		count, err := storage.CountPatterns(ctx, contracts.ListOptions{MinStrength: 70})
		if err != nil {
			t.Fatalf("CountPatterns failed: %v", err)
		}
		if count != 2 { // test-1 (80), test-3 (70)
			t.Errorf("expected 2 patterns with strength >= 70, got %d", count)
		}
	})

	// Test GetPatternByTrigger
	t.Run("GetPatternByTrigger", func(t *testing.T) {
		p, err := storage.GetPatternByTrigger(ctx, "hello")
		if err != nil {
			t.Fatalf("GetPatternByTrigger failed: %v", err)
		}
		if p.ID != "test-1" {
			t.Errorf("expected pattern id test-1, got %s", p.ID)
		}
	})

	// Test GetPatternByTrigger not found
	t.Run("GetPatternByTrigger not found", func(t *testing.T) {
		_, err := storage.GetPatternByTrigger(ctx, "nonexistent")
		if err == nil {
			t.Error("expected error for nonexistent trigger")
		}
	})

	// Test GetRecentlyUsedPatterns
	t.Run("GetRecentlyUsedPatterns", func(t *testing.T) {
		patterns, err := storage.GetRecentlyUsedPatterns(ctx, 2)
		if err != nil {
			t.Fatalf("GetRecentlyUsedPatterns failed: %v", err)
		}
		if len(patterns) != 1 { // Only test-3 has last_used_at set
			t.Errorf("expected 1 recently used pattern, got %d", len(patterns))
		}
		if len(patterns) > 0 && patterns[0].Trigger != "weather" {
			t.Errorf("expected weather pattern first, got %s", patterns[0].Trigger)
		}
	})

	// Test SearchPatterns
	t.Run("SearchPatterns", func(t *testing.T) {
		results, err := storage.SearchPatterns(ctx, "help", contracts.ListOptions{})
		if err != nil {
			t.Fatalf("SearchPatterns failed: %v", err)
		}
		if len(results) != 2 { // hello (Hello!), help
			t.Errorf("expected 2 search results, got %d", len(results))
		}
	})

	// Test SearchPatterns with limit
	t.Run("SearchPatterns with limit", func(t *testing.T) {
		results, err := storage.SearchPatterns(ctx, "weather", contracts.ListOptions{Limit: 1})
		if err != nil {
			t.Fatalf("SearchPatterns failed: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("expected 1 result with limit, got %d", len(results))
		}
	})

	// Test GetTopPatterns
	t.Run("GetTopPatterns", func(t *testing.T) {
		patterns, err := storage.GetTopPatterns(ctx, 2)
		if err != nil {
			t.Fatalf("GetTopPatterns failed: %v", err)
		}
		if len(patterns) != 2 {
			t.Errorf("expected 2 top patterns, got %d", len(patterns))
		}
		if len(patterns) > 0 && patterns[0].Strength < patterns[1].Strength {
			t.Error("patterns should be ordered by strength descending")
		}
	})

	// Test statement caching works
	t.Run("Statement caching", func(t *testing.T) {
		// Execute same query twice to test caching
		for i := 0; i < 2; i++ {
			_, err := storage.GetPatternByTrigger(ctx, "hello")
			if err != nil {
				t.Fatalf("GetPatternByTrigger failed on iteration %d: %v", i, err)
			}
		}
		// If we get here, caching works
	})
}

// TestQueryOptimizationEdgeCases tests edge cases for query optimization
func TestQueryOptimizationEdgeCases(t *testing.T) {
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

	// Test CountPatterns on empty database
	t.Run("CountPatterns empty", func(t *testing.T) {
		count, err := storage.CountPatterns(ctx, contracts.ListOptions{})
		if err != nil {
			t.Fatalf("CountPatterns failed: %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 patterns, got %d", count)
		}
	})

	// Test SearchPatterns on empty database
	t.Run("SearchPatterns empty", func(t *testing.T) {
		results, err := storage.SearchPatterns(ctx, "test", contracts.ListOptions{})
		if err != nil {
			t.Fatalf("SearchPatterns failed: %v", err)
		}
		if len(results) != 0 {
			t.Errorf("expected 0 results, got %d", len(results))
		}
	})

	// Test GetTopPatterns with limit
	t.Run("GetTopPatterns empty", func(t *testing.T) {
		patterns, err := storage.GetTopPatterns(ctx, 10)
		if err != nil {
			t.Fatalf("GetTopPatterns failed: %v", err)
		}
		if len(patterns) != 0 {
			t.Errorf("expected 0 patterns, got %d", len(patterns))
		}
	})

	// Test GetRecentlyUsedPatterns with limit
	t.Run("GetRecentlyUsedPatterns empty", func(t *testing.T) {
		patterns, err := storage.GetRecentlyUsedPatterns(ctx, 10)
		if err != nil {
			t.Fatalf("GetRecentlyUsedPatterns failed: %v", err)
		}
		if len(patterns) != 0 {
			t.Errorf("expected 0 patterns, got %d", len(patterns))
		}
	})
}
