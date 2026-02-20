package contracts

import (
	"context"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// Matcher defines the interface for pattern matching
type Matcher interface {
	Match(ctx context.Context, query string, patterns []*models.Pattern) []MatchResult
}

// MatchOptions contains options for matching
type MatchOptions struct {
	Threshold  float64
	Limit      int
	Tags       []string
	SpaceID    string
	ExactFirst bool
}

// MatchResult represents a match result
type MatchResult struct {
	Pattern    *models.Pattern
	Confidence float64
	Branch     string
}
