package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Pattern represents a reflex pattern that maps triggers to AI responses.
// Patterns form the basis of the reflex system - they are matched against
// user input and used to generate quick responses.
type Pattern struct {
	// ==================== Identification ====================
	// Unique identifier for the pattern (UUID v4)
	ID string `json:"id" db:"id"`
	// Timestamp when the pattern was created
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// Timestamp when the pattern was last updated
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// ==================== Core Fields ====================
	// Trigger is the input pattern that activates this reflex
	// (e.g., "用户API" triggers a user API template)
	Trigger string `json:"trigger" db:"trigger"`
	// Response is the AI-generated content returned when triggered
	Response string `json:"response" db:"response"`

	// ==================== Space (v2.0) ====================
	// SpaceID identifies which space this pattern belongs to
	SpaceID string `json:"space_id" db:"space_id"`

	// ==================== Strength Management ====================
	// Strength is the current activation strength (0-100).
	// Increases with use (reinforce), decreases over time (decay).
	Strength float64 `json:"strength" db:"strength"`
	// Threshold is the minimum strength required for activation (0-100).
	// Pattern only matches when Strength >= Threshold.
	Threshold float64 `json:"threshold" db:"threshold"`
	// DecayRate is the rate at which strength decreases per time unit (0-1).
	// Higher values mean faster decay.
	DecayRate float64 `json:"decay_rate" db:"decay_rate"`
	// DecayEnabled indicates whether automatic decay is active.
	DecayEnabled bool `json:"decay_enabled" db:"decay_enabled"`

	// ==================== Usage Statistics ====================
	// ReinforceCnt is the number of times this pattern was reinforced
	ReinforceCnt int `json:"reinforcement_count" db:"reinforcement_count"`
	// DecayCnt is the number of times decay was applied
	DecayCnt int `json:"decay_count" db:"decay_count"`
	// LastUsedAt is the timestamp of the last reinforcement
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`

	// ==================== Metadata ====================
	// Connections links this pattern to other patterns (for thought chains)
	Connections []string `json:"connections,omitempty" db:"connections"`
	// Tags categorize patterns for filtering (e.g., "api", "frontend")
	Tags []string `json:"tags,omitempty" db:"tags"`
	// Project groups patterns under a common project name
	Project string `json:"project,omitempty" db:"project"`
	// UserID identifies the owner of this pattern
	UserID string `json:"user_id,omitempty" db:"user_id"`

	// ==================== Soft Delete ====================
	// DeletedAt marks the pattern as deleted (if not nil)
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// NewPattern creates a new pattern with default values.
// The trigger and response are required. Other fields are initialized
// with sensible defaults:
//   - Strength: 0.0 (starts inactive)
//   - Threshold: 50.0 (activates at 50% strength)
//   - DecayRate: 0.01 (1% decay per period)
//   - DecayEnabled: true (automatic decay on)
func NewPattern(trigger, response string) *Pattern {
	now := time.Now()
	return &Pattern{
		ID:           uuid.New().String(),
		CreatedAt:    now,
		UpdatedAt:    now,
		Trigger:      trigger,
		Response:     response,
		SpaceID:      "global", // Default to global space
		Strength:     0.0,
		Threshold:    50.0,
		DecayRate:    0.01,
		DecayEnabled: true,
		ReinforceCnt: 0,
		DecayCnt:     0,
	}
}

// Validate validates all fields of the pattern.
// Returns an error if any validation fails.
// Validates trigger, response, and strength fields.
func (p *Pattern) Validate() error {
	if err := p.ValidateTrigger(); err != nil {
		return err
	}
	if err := p.ValidateResponse(); err != nil {
		return err
	}
	if err := p.ValidateStrength(); err != nil {
		return err
	}
	return nil
}

// ValidateTrigger validates the trigger field
func (p *Pattern) ValidateTrigger() error {
	trimmed := strings.TrimSpace(p.Trigger)
	if trimmed == "" {
		return ErrValidation{Field: "trigger", Reason: "cannot be empty"}
	}
	if len(trimmed) > 500 {
		return ErrValidation{Field: "trigger", Reason: "exceeds 500 characters"}
	}
	p.Trigger = trimmed
	return nil
}

// ValidateResponse validates the response field
func (p *Pattern) ValidateResponse() error {
	if strings.TrimSpace(p.Response) == "" {
		return ErrValidation{Field: "response", Reason: "cannot be empty"}
	}
	return nil
}

// ValidateStrength validates strength and threshold
func (p *Pattern) ValidateStrength() error {
	if p.Strength < 0 || p.Strength > 100 {
		return ErrValidation{Field: "strength", Reason: "must be between 0 and 100"}
	}
	if p.Threshold < 0 || p.Threshold > 100 {
		return ErrValidation{Field: "threshold", Reason: "must be between 0 and 100"}
	}
	if p.DecayRate < 0 || p.DecayRate > 1 {
		return ErrValidation{Field: "decay_rate", Reason: "must be between 0 and 1"}
	}
	return nil
}

// ErrValidation represents a validation error for a specific field.
// Used to provide detailed feedback about which field failed validation
// and why.
type ErrValidation struct {
	Field  string // The field name that failed validation
	Reason string  // Human-readable reason for the failure
}

// Error implements the error interface for ErrValidation.
func (e ErrValidation) Error() string {
	return fmt.Sprintf("validation error: %s %s", e.Field, e.Reason)
}
