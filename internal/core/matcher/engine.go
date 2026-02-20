package matcher

import (
	"context"
	"sort"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// Engine is the pattern matching engine
type Engine struct {
	exactMatcher  *ExactMatcher
	keywordMatcher *KeywordMatcher
}

// NewEngine creates a new matching engine
func NewEngine() *Engine {
	return &Engine{
		exactMatcher:  NewExactMatcher(),
		keywordMatcher: NewKeywordMatcher(),
	}
}

// Match finds patterns matching the query
func (e *Engine) Match(ctx context.Context, query string, patterns []*models.Pattern, opts contracts.MatchOptions) []contracts.MatchResult {
	if len(patterns) == 0 || query == "" {
		return nil
	}

	var allResults []contracts.MatchResult

	// Exact match first (highest priority)
	if opts.ExactFirst {
		exactResults := e.exactMatcher.Match(ctx, query, patterns)
		allResults = append(allResults, exactResults...)
		
		// Then keyword match
		keywordResults := e.keywordMatcher.Match(ctx, query, patterns)
		allResults = append(allResults, keywordResults...)
	} else {
		// Keyword match first
		keywordResults := e.keywordMatcher.Match(ctx, query, patterns)
		allResults = append(allResults, keywordResults...)
		
		// Then exact match
		exactResults := e.exactMatcher.Match(ctx, query, patterns)
		allResults = append(allResults, exactResults...)
	}

	// Filter by threshold
	var filtered []contracts.MatchResult
	for _, r := range allResults {
		if r.Confidence >= opts.Threshold {
			filtered = append(filtered, r)
		}
	}

	// Sort by confidence descending
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Confidence > filtered[j].Confidence
	})

	// Limit results
	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}

	return filtered
}

// MatchOne finds the best match
func (e *Engine) MatchOne(ctx context.Context, query string, patterns []*models.Pattern) *contracts.MatchResult {
	opts := contracts.MatchOptions{
		Threshold:  0,
		Limit:      1,
		ExactFirst: true,
	}
	
	results := e.Match(ctx, query, patterns, opts)
	if len(results) > 0 {
		return &results[0]
	}
	return nil
}
