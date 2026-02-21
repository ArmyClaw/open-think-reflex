package sqlite

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// TestStorage_ConcurrentReadWrite tests concurrent read/write operations (Iter 47)
func TestStorage_ConcurrentReadWrite(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	var wg sync.WaitGroup

	// Create initial patterns
	for i := 0; i < 10; i++ {
		p := models.NewPattern("trigger", "response")
		p.Project = "test"
		if err := storage.SavePattern(ctx, p); err != nil {
			t.Fatalf("Setup failed: %v", err)
		}
	}

	// Reset stats before test
	storage.ResetStats()

	// Concurrent readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 10})
				if err != nil {
					t.Logf("Read error: %v", err)
				}
				if len(patterns) > 0 {
					// Read individual pattern
					storage.GetPattern(ctx, patterns[0].ID)
				}
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// Concurrent writers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				p := models.NewPattern("concurrent-trigger", "response")
				p.Project = "concurrent"
				if err := storage.SavePattern(ctx, p); err != nil {
					t.Logf("Write error: %v", err)
				}
				time.Sleep(time.Millisecond)
			}
		}()
	}

	wg.Wait()

	// Check stats
	stats := storage.Stats()
	t.Logf("Storage stats after concurrent test: %+v", stats)

	// Verify data integrity
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 100})
	if err != nil {
		t.Fatalf("ListPatterns failed: %v", err)
	}

	// Should have initial 10 + 30 concurrent writes = ~40 patterns
	if len(patterns) < 30 {
		t.Errorf("Expected at least 30 patterns, got %d", len(patterns))
	}

	t.Logf("Final pattern count: %d", len(patterns))
}

// TestStorage_ConcurrentStats tests that Stats() works correctly under concurrent access
func TestStorage_ConcurrentStats(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	var wg sync.WaitGroup

	// Perform many operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				p := models.NewPattern("stat-test", "response")
				storage.SavePattern(ctx, p)
				storage.ListPatterns(ctx, contracts.ListOptions{Limit: 10})
				storage.Stats() // Read stats concurrently
			}
		}()
	}

	wg.Wait()

	stats := storage.Stats()
	t.Logf("Stats after concurrent operations: ReadOps=%d, WriteOps=%d",
		stats.ReadOps, stats.WriteOps)
}

// TestStorage_GetPatternByTrigger_Concurrent tests GetPatternByTrigger under concurrency
func TestStorage_GetPatternByTrigger_Concurrent(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test patterns
	triggers := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	for _, trigger := range triggers {
		p := models.NewPattern(trigger, "response for "+trigger)
		if err := storage.SavePattern(ctx, p); err != nil {
			t.Fatal(err)
		}
	}

	var wg sync.WaitGroup
	results := make(chan string, 100)

	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for _, trigger := range triggers {
				p, err := storage.GetPatternByTrigger(ctx, trigger)
				if err == nil && p != nil {
					results <- p.Trigger
				}
				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	close(results)

	// Verify we got results
	count := 0
	for range results {
		count++
	}

	if count == 0 {
		t.Error("Expected some successful reads")
	}

	t.Logf("Successful concurrent reads: %d", count)
}

// TestStorage_SearchPatterns_Concurrent tests SearchPatterns under concurrency
func TestStorage_SearchPatterns_Concurrent(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create patterns with searchable content
	keywords := []string{"golang", "python", "rust", "javascript", "typescript"}
	for _, kw := range keywords {
		for j := 0; j < 5; j++ {
			p := models.NewPattern(kw+"-"+string(rune('a'+j)), "Response about "+kw)
			p.Project = "search-test"
			if err := storage.SavePattern(ctx, p); err != nil {
				t.Fatal(err)
			}
		}
	}

	var wg sync.WaitGroup

	// Concurrent searches
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for _, kw := range keywords {
				results, err := storage.SearchPatterns(ctx, kw, contracts.ListOptions{Limit: 10})
				if err != nil {
					t.Logf("Search error: %v", err)
					return
				}
				if len(results) == 0 {
					t.Logf("No results for keyword: %s", kw)
				}
			}
		}(i)
	}

	wg.Wait()
	t.Log("Concurrent search test completed")
}

// TestStorage_GetTopPatterns tests GetTopPatterns
func TestStorage_GetTopPatterns(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create patterns with different strengths
	for i := 0; i < 20; i++ {
		p := models.NewPattern("trigger-"+string(rune('a'+i)), "response")
		p.Strength = float64(20 - i) // Different strengths
		if err := storage.SavePattern(ctx, p); err != nil {
			t.Fatal(err)
		}
	}

	patterns, err := storage.GetTopPatterns(ctx, 5)
	if err != nil {
		t.Fatalf("GetTopPatterns failed: %v", err)
	}

	if len(patterns) != 5 {
		t.Errorf("Expected 5 patterns, got %d", len(patterns))
	}

	// Verify ordering (highest strength first)
	for i := 1; i < len(patterns); i++ {
		if patterns[i].Strength > patterns[i-1].Strength {
			t.Error("Patterns not ordered by strength descending")
		}
	}

	t.Logf("Top pattern strength: %.2f", patterns[0].Strength)
}

// TestStorage_CountPatterns tests CountPatterns
func TestStorage_CountPatterns(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create patterns in different projects
	for i := 0; i < 10; i++ {
		p := models.NewPattern("trigger", "response")
		if i < 5 {
			p.Project = "project-a"
		} else {
			p.Project = "project-b"
		}
		if err := storage.SavePattern(ctx, p); err != nil {
			t.Fatal(err)
		}
	}

	// Count all
	total, err := storage.CountPatterns(ctx, contracts.ListOptions{})
	if err != nil {
		t.Fatalf("CountPatterns failed: %v", err)
	}

	if total != 10 {
		t.Errorf("Expected 10 patterns, got %d", total)
	}

	// Count by project
	projectACount, err := storage.CountPatterns(ctx, contracts.ListOptions{Project: "project-a"})
	if err != nil {
		t.Fatal(err)
	}

	if projectACount != 5 {
		t.Errorf("Expected 5 patterns in project-a, got %d", projectACount)
	}

	t.Logf("Total: %d, Project-A: %d", total, projectACount)
}
