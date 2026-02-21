// Package contracts defines the core interfaces for Open-Think-Reflex.
// These interfaces establish the contracts between different layers of the application.
package contracts

import (
	"context"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// Matcher defines the interface for pattern matching.
// Implementations should provide efficient pattern matching algorithms.
type Matcher interface {
	// Match finds patterns matching the given query.
	// Returns a list of match results sorted by confidence (highest first).
	Match(ctx context.Context, query string, patterns []*models.Pattern) []MatchResult
}

// MatchOptions contains optional parameters for matching operations.
// All fields are optional and have sensible defaults.
type MatchOptions struct {
	// Threshold is the minimum confidence score (0.0-1.0) for results.
	// Default: 0.0 (no filter)
	Threshold float64

	// Limit restricts the maximum number of results returned.
	// Default: 0 (no limit)
	Limit int

	// Tags filters patterns by tags. Only patterns containing at least
	// one of these tags will be matched.
	// Default: nil (match all)
	Tags []string

	// SpaceID restricts matching to patterns within a specific space.
	// Default: "" (match all spaces)
	SpaceID string

	// ExactFirst determines whether exact matches are prioritized.
	// When true, exact matches appear before keyword matches.
	// Default: true
	ExactFirst bool
}

// MatchResult represents a single pattern match with its confidence score.
// Confidence scores range from 0.0 (no match) to 1.0 (perfect match).
type MatchResult struct {
	// Pattern is the matched pattern
	Pattern *models.Pattern

	// Confidence is the match confidence score (0.0-1.0)
	Confidence float64

	// Branch is the thought chain branch that matched
	Branch string
}
