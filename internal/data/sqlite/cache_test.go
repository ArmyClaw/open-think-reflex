package sqlite

import (
	"testing"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

func TestQueryCache(t *testing.T) {
	cache := NewQueryCache(time.Minute, 10)

	pattern := &models.Pattern{
		ID:       "test-1",
		Trigger:  "test trigger",
		Response: "test response",
	}

	// Test miss
	if _, ok := cache.Get("test-1"); ok {
		t.Error("Expected cache miss")
	}

	// Test set and hit
	cache.Set("test-1", pattern)
	if p, ok := cache.Get("test-1"); !ok {
		t.Error("Expected cache hit")
	} else if p.ID != pattern.ID {
		t.Errorf("Expected %s, got %s", pattern.ID, p.ID)
	}

	// Test stats
	stats := cache.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}

	// Test invalidate
	cache.Invalidate("test-1")
	if _, ok := cache.Get("test-1"); ok {
		t.Error("Expected cache miss after invalidate")
	}

	// Test clear
	cache.Set("test-2", pattern)
	cache.Clear()
	if _, ok := cache.Get("test-2"); ok {
		t.Error("Expected cache miss after clear")
	}
}

func TestQueryCacheEviction(t *testing.T) {
	cache := NewQueryCache(time.Minute, 3)

	for i := 0; i < 5; i++ {
		pattern := &models.Pattern{
			ID:       "test-" + string(rune(i)),
			Trigger:  "trigger",
		}
		cache.Set(pattern.ID, pattern)
	}

	// Should have evicted at least one
	stats := cache.Stats()
	if stats.Evicts == 0 {
		t.Error("Expected evictions when cache is full")
	}
}
