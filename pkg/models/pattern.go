package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Pattern represents a reflex pattern
type Pattern struct {
	// Identification
	ID        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Core fields
	Trigger   string  `json:"trigger" db:"trigger"`
	Response  string  `json:"response" db:"response"`

	// Strength management
	Strength    float64 `json:"strength" db:"strength"`
	Threshold   float64 `json:"threshold" db:"threshold"`
	DecayRate   float64 `json:"decay_rate" db:"decay_rate"`
	DecayEnabled bool   `json:"decay_enabled" db:"decay_enabled"`

	// Statistics
	ReinforceCnt int        `json:"reinforcement_count" db:"reinforcement_count"`
	DecayCnt     int        `json:"decay_count" db:"decay_count"`
	LastUsedAt   *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`

	// Metadata
	Connections []string `json:"connections,omitempty" db:"connections"`
	Tags        []string `json:"tags,omitempty" db:"tags"`
	Project     string   `json:"project,omitempty" db:"project"`
	UserID      string   `json:"user_id,omitempty" db:"user_id"`

	// Soft delete
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// NewPattern creates a new pattern with defaults
func NewPattern(trigger, response string) *Pattern {
	now := time.Now()
	return &Pattern{
		ID:           uuid.New().String(),
		CreatedAt:    now,
		UpdatedAt:    now,
		Trigger:      trigger,
		Response:     response,
		Strength:     0.0,
		Threshold:    50.0,
		DecayRate:    0.01,
		DecayEnabled: true,
		ReinforceCnt: 0,
		DecayCnt:     0,
	}
}

// Validate validates the pattern
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

// ErrValidation represents a validation error
type ErrValidation struct {
	Field   string
	Reason  string
}

func (e ErrValidation) Error() string {
	return fmt.Sprintf("validation error: %s %s", e.Field, e.Reason)
}
