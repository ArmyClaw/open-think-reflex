package models

import "time"

// Space represents an isolated namespace for patterns
type Space struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Owner (for multi-user support in v2.0)
	Owner string `json:"owner,omitempty" db:"owner"`

	// Settings
	DefaultSpace bool `json:"default_space" db:"is_default"`
	PatternLimit int  `json:"pattern_limit" db:"pattern_limit"`

	// Statistics
	PatternCount int `json:"pattern_count" db:"pattern_count"`
}

// DefaultSpaces returns the default spaces
func DefaultSpaces() []*Space {
	now := time.Now()
	return []*Space{
		{
			ID:           "global",
			Name:         "Global",
			Description:  "Patterns available everywhere",
			Owner:        "",
			DefaultSpace: true,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           "project",
			Name:         "Project",
			Description:  "Project-specific patterns",
			Owner:        "",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           "personal",
			Name:         "Personal",
			Description:  "Personal patterns",
			Owner:        "",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
}
