package matcher

import (
	"context"
	"testing"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

func TestExactMatcher(t *testing.T) {
	matcher := NewExactMatcher()
	ctx := context.Background()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "test", Response: "response 1"},
		{ID: "2", Trigger: "python setup", Response: "response 2"},
		{ID: "3", Trigger: "docker run", Response: "response 3"},
	}

	tests := []struct {
		name     string
		query    string
		wantLen  int
		wantID   string
		wantConf float64
	}{
		{"exact match", "test", 1, "1", 100.0},
		{"exact match python", "python setup", 1, "2", 100.0},
		{"no match", "unknown", 0, "", 0},
		{"partial match", "python", 0, "", 0}, // Not exact
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := matcher.Match(ctx, tt.query, patterns)

			if tt.wantLen == 0 {
				if len(results) != 0 {
					t.Errorf("Expected 0 results, got %d", len(results))
				}
				return
			}

			if len(results) != tt.wantLen {
				t.Errorf("Expected %d results, got %d", tt.wantLen, len(results))
				return
			}

			if results[0].Pattern.ID != tt.wantID {
				t.Errorf("Expected pattern ID %s, got %s", tt.wantID, results[0].Pattern.ID)
			}

			if results[0].Confidence != tt.wantConf {
				t.Errorf("Expected confidence %f, got %f", tt.wantConf, results[0].Confidence)
			}
		})
	}
}

func TestKeywordMatcher(t *testing.T) {
	matcher := NewKeywordMatcher()
	ctx := context.Background()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "test", Response: "response 1"},
		{ID: "2", Trigger: "python setup", Response: "response 2"},
		{ID: "3", Trigger: "docker run container", Response: "response 3"},
	}

	tests := []struct {
		name    string
		query   string
		wantLen int
	}{
		{"single word match", "test", 1},
		{"partial match", "python", 1},
		{"multi word match", "python setup", 1},
		{"no match", "unknown", 0},
		{"multi word query", "docker container", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := matcher.Match(ctx, tt.query, patterns)

			if len(results) != tt.wantLen {
				t.Errorf("Expected %d results, got %d", tt.wantLen, len(results))
			}

			// Verify confidence is reasonable
			for _, r := range results {
				if r.Confidence > 100 || r.Confidence < 0 {
					t.Errorf("Invalid confidence: %f", r.Confidence)
				}
			}
		})
	}
}

func TestEngineMatch(t *testing.T) {
	engine := NewEngine()
	ctx := context.Background()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "test", Response: "response 1"},
		{ID: "2", Trigger: "python setup", Response: "response 2"},
		{ID: "3", Trigger: "python", Response: "response 3"},
	}

	// Test exact first priority
	results := engine.Match(ctx, "python", patterns, contracts.MatchOptions{
		Threshold:  0,
		Limit:      10,
		ExactFirst: true,
	})

	if len(results) == 0 {
		t.Fatal("Expected matches, got none")
	}

	// First result should be exact match
	if results[0].Branch != "exact" {
		t.Errorf("Expected first result to be exact, got %s", results[0].Branch)
	}

	if results[0].Confidence != 100.0 {
		t.Errorf("Expected confidence 100, got %f", results[0].Confidence)
	}
}

func TestEngineMatchOne(t *testing.T) {
	engine := NewEngine()
	ctx := context.Background()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "test", Response: "response 1"},
	}

	result := engine.MatchOne(ctx, "test", patterns)

	if result == nil {
		t.Fatal("Expected match, got nil")
	}

	if result.Pattern.ID != "1" {
		t.Errorf("Expected pattern ID 1, got %s", result.Pattern.ID)
	}
}

func TestEngineMatchNoResults(t *testing.T) {
	engine := NewEngine()
	ctx := context.Background()

	patterns := []*models.Pattern{
		{ID: "1", Trigger: "test", Response: "response 1"},
	}

	result := engine.MatchOne(ctx, "unknown query", patterns)

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestExtractWords(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"test", []string{"test"}},
		{"python setup", []string{"python", "setup"}},
		{"docker-run", []string{"docker", "run"}},
		{"file.txt", []string{"file", "txt"}},
		{"hello world test", []string{"hello", "world", "test"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractWords(tt.input)
			
			if len(got) != len(tt.want) {
				t.Errorf("extractWords(%q) = %v, want %v", tt.input, got, tt.want)
				return
			}

			for i, w := range got {
				if w != tt.want[i] {
					t.Errorf("extractWords(%q)[%d] = %q, want %q", tt.input, i, w, tt.want[i])
				}
			}
		})
	}
}
