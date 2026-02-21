package sqlite

import (
	"context"
	"testing"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// Benchmark tests for storage operations

func BenchmarkSavePattern(b *testing.B) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()
	db.Migrate(ctx)

	storage := NewStorage(db)

	pattern := &models.Pattern{
		ID:          "bench-pattern",
		Trigger:     "test trigger",
		Response:    "test response",
		Strength:    0.8,
		Threshold:   0.5,
		Project:     "bench",
		UserID:      "user1",
		Connections: []string{"node1"},
		Tags:        []string{"benchmark"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pattern.ID = "bench-pattern"
		if err := storage.SavePattern(ctx, pattern); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetPattern(b *testing.B) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()
	db.Migrate(ctx)

	storage := NewStorage(db)

	// Setup
	pattern := &models.Pattern{
		ID:       "bench-get",
		Trigger:  "get trigger",
		Response: "get response",
		Strength: 0.8,
		Project:  "bench",
		UserID:   "user1",
	}
	storage.SavePattern(ctx, pattern)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.GetPattern(ctx, "bench-get")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkListPatterns(b *testing.B) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()
	db.Migrate(ctx)

	storage := NewStorage(db)

	// Setup - insert 100 patterns
	for i := 0; i < 100; i++ {
		pattern := &models.Pattern{
			ID:       "bench-list-" + string(rune(i)),
			Trigger:  "trigger-" + string(rune(i)),
			Response: "response",
			Strength: 0.5 + float64(i)/200,
			Project:  "bench",
			UserID:   "user1",
		}
		storage.SavePattern(ctx, pattern)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 50})
		if err != nil {
			b.Fatal(err)
		}
		if len(patterns) != 50 {
			b.Fatalf("expected 50, got %d", len(patterns))
		}
	}
}

func BenchmarkSearchPatterns(b *testing.B) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()
	db.Migrate(ctx)

	storage := NewStorage(db)

	// Setup - insert 100 patterns
	for i := 0; i < 100; i++ {
		pattern := &models.Pattern{
			ID:       "bench-search-" + string(rune(i)),
			Trigger:  "test trigger " + string(rune(i)),
			Response: "response with keyword " + string(rune(i)),
			Strength: 0.8,
			Project:  "bench",
			UserID:   "user1",
		}
		storage.SavePattern(ctx, pattern)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		patterns, err := storage.SearchPatterns(ctx, "keyword", contracts.ListOptions{Limit: 10})
		if err != nil {
			b.Fatal(err)
		}
		_ = patterns
	}
}

func BenchmarkBatchSave(b *testing.B) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()
	db.Migrate(ctx)

	storage := NewStorage(db)

	patterns := make([]*models.Pattern, 50)
	for i := 0; i < 50; i++ {
		patterns[i] = &models.Pattern{
			ID:       "bench-batch-" + string(rune(i)),
			Trigger:  "trigger",
			Response: "response",
			Strength: 0.8,
			Project:  "bench",
			UserID:   "user1",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := storage.SavePatternsBatch(ctx, patterns); err != nil {
			b.Fatal(err)
		}
	}
}
