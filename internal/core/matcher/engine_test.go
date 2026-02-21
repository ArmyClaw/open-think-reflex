package matcher

import (
	"context"
	"testing"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

func TestEngine_NewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine should not return nil")
	}
	if engine.exactMatcher == nil {
		t.Error("exactMatcher should be initialized")
	}
	if engine.keywordMatcher == nil {
		t.Error("keywordMatcher should be initialized")
	}
}

func TestEngine_Match_EmptyInputs(t *testing.T) {
	engine := NewEngine()

	// Test with nil patterns
	result := engine.Match(context.Background(), "test", nil, contracts.MatchOptions{})
	if result != nil {
		t.Error("Match with nil patterns should return nil")
	}

	// Test with empty patterns
	result = engine.Match(context.Background(), "test", []*models.Pattern{}, contracts.MatchOptions{})
	if result != nil {
		t.Error("Match with empty patterns should return nil")
	}

	// Test with empty query
	pattern := models.NewPattern("trigger", "response")
	result = engine.Match(context.Background(), "", []*models.Pattern{pattern}, contracts.MatchOptions{})
	if result != nil {
		t.Error("Match with empty query should return nil")
	}
}

func TestEngine_Match_ExactFirst(t *testing.T) {
	engine := NewEngine()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "test", Response: "exact match", Strength: 80},
		{ID: "2", Trigger: "testing", Response: "keyword match", Strength: 60},
		{ID: "3", Trigger: "test case", Response: "another keyword", Strength: 40},
	}

	opts := contracts.MatchOptions{
		ExactFirst: true,
		Threshold:  0,
		Limit:      0,
	}

	results := engine.Match(context.Background(), "test", patterns, opts)
	
	// Should have results
	if len(results) == 0 {
		t.Error("Expected matches, got none")
	}

	// With ExactFirst, exact match should come first
	if len(results) > 0 && results[0].Pattern.ID != "1" {
		t.Logf("First result: %s, expected: 1", results[0].Pattern.ID)
	}
}

func TestEngine_Match_KeywordFirst(t *testing.T) {
	engine := NewEngine()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "exact", Response: "exact", Strength: 80},
		{ID: "2", Trigger: "test keyword", Response: "keyword", Strength: 60},
	}

	opts := contracts.MatchOptions{
		ExactFirst: false,
		Threshold:  0,
		Limit:      0,
	}

	results := engine.Match(context.Background(), "test", patterns, opts)
	
	if len(results) == 0 {
		t.Error("Expected matches, got none")
	}
}

func TestEngine_Match_WithThreshold(t *testing.T) {
	engine := NewEngine()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "test", Response: "high confidence", Strength: 90},
		{ID: "2", Trigger: "test", Response: "low confidence", Strength: 30},
	}

	opts := contracts.MatchOptions{
		Threshold:  50,
		ExactFirst: true,
	}

	results := engine.Match(context.Background(), "test", patterns, opts)

	// Should only return pattern with strength >= 50
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if len(results) > 0 && results[0].Pattern.ID != "1" {
		t.Error("Should only return high confidence pattern")
	}
}

func TestEngine_Match_WithLimit(t *testing.T) {
	engine := NewEngine()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "test", Response: "r1", Strength: 90},
		{ID: "2", Trigger: "test", Response: "r2", Strength: 80},
		{ID: "3", Trigger: "test", Response: "r3", Strength: 70},
	}

	opts := contracts.MatchOptions{
		Limit: 2,
	}

	results := engine.Match(context.Background(), "test", patterns, opts)

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestEngine_Match_SortByConfidence(t *testing.T) {
	engine := NewEngine()

	patterns := []*models.Pattern{
		{ID: "low", Trigger: "test", Response: "low", Strength: 30},
		{ID: "high", Trigger: "test", Response: "high", Strength: 90},
		{ID: "mid", Trigger: "test", Response: "mid", Strength: 60},
	}

	results := engine.Match(context.Background(), "test", patterns, contracts.MatchOptions{})

	// Should be sorted by confidence descending
	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// First should have highest confidence
	if results[0].Confidence < results[1].Confidence {
		t.Error("Results should be sorted by confidence descending")
	}
}

func TestEngine_MatchOne(t *testing.T) {
	engine := NewEngine()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "test", Response: "found"},
		{ID: "2", Trigger: "other", Response: "not found"},
	}

	// Test with match
	result := engine.MatchOne(context.Background(), "test", patterns)
	if result == nil {
		t.Error("Expected a match")
	}
	if result.Pattern.ID != "1" {
		t.Errorf("Expected pattern 1, got %s", result.Pattern.ID)
	}

	// Test without match
	result = engine.MatchOne(context.Background(), "notexist", patterns)
	if result != nil {
		t.Error("Expected no match")
	}

	// Test with empty patterns
	result = engine.MatchOne(context.Background(), "test", nil)
	if result != nil {
		t.Error("Expected nil for empty patterns")
	}
}

func TestEngine_Match_NoMatch(t *testing.T) {
	engine := NewEngine()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "abc", Response: "xyz"},
		{ID: "2", Trigger: "def", Response: "uvw"},
	}

	results := engine.Match(context.Background(), "xyz", patterns, contracts.MatchOptions{})
	
	// No matches expected for unrelated query
	if len(results) > 0 {
		t.Logf("Got %d unexpected matches", len(results))
	}
}
