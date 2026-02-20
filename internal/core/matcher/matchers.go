package matcher

import (
	"context"
	"strings"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// ExactMatcher matches patterns by exact trigger
type ExactMatcher struct{}

func NewExactMatcher() *ExactMatcher {
	return &ExactMatcher{}
}

// Match finds patterns that exactly match the query
func (m *ExactMatcher) Match(ctx context.Context, query string, patterns []*models.Pattern) []contracts.MatchResult {
	var results []contracts.MatchResult
	
	query = strings.TrimSpace(strings.ToLower(query))
	
	for _, p := range patterns {
		trigger := strings.TrimSpace(strings.ToLower(p.Trigger))
		
		if trigger == query {
			results = append(results, contracts.MatchResult{
				Pattern:    p,
				Confidence: 100.0,
				Branch:     "exact",
			})
		}
	}
	
	return results
}

// KeywordMatcher matches patterns by keywords
type KeywordMatcher struct{}

func NewKeywordMatcher() *KeywordMatcher {
	return &KeywordMatcher{}
}

// Match finds patterns that match by keywords
func (m *KeywordMatcher) Match(ctx context.Context, query string, patterns []*models.Pattern) []contracts.MatchResult {
	var results []contracts.MatchResult
	
	queryWords := extractWords(strings.ToLower(query))
	if len(queryWords) == 0 {
		return results
	}
	
	for _, p := range patterns {
		triggerWords := extractWords(strings.ToLower(p.Trigger))
		
		// Count matching words
		matchCount := 0
		for _, qw := range queryWords {
			for _, tw := range triggerWords {
				if strings.Contains(tw, qw) || strings.Contains(qw, tw) {
					matchCount++
					break
				}
			}
		}
		
		// Calculate confidence based on match ratio
		if matchCount > 0 {
			confidence := float64(matchCount) / float64(len(queryWords)) * 100
			results = append(results, contracts.MatchResult{
				Pattern:    p,
				Confidence: confidence,
				Branch:     "keyword",
			})
		}
	}
	
	return results
}

// extractWords extracts words from a string
func extractWords(s string) []string {
	var words []string
	var current strings.Builder
	
	for _, r := range s {
		if r == ' ' || r == '-' || r == '_' || r == '/' || r == '.' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}
	
	if current.Len() > 0 {
		words = append(words, current.String())
	}
	
	return words
}
