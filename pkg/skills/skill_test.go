package skills

import (
	"testing"
	"time"

	otrmodels "github.com/ArmyClaw/open-think-reflex/pkg/models"
)

func TestSkillValidate(t *testing.T) {
	tests := []struct {
		name    string
		skill   *Skill
		wantErr bool
	}{
		{
			name: "valid skill",
			skill: &Skill{
				Name:        "test-trigger",
				Description: "Test description",
				Version:     "1.0",
				Trigger:     "test-trigger",
				Response:    "Test response content",
				Source:      "otr",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			skill: &Skill{
				Name:        "",
				Description: "Test description",
				Version:     "1.0",
				Trigger:     "test-trigger",
				Response:    "Test response",
			},
			wantErr: true,
		},
		{
			name: "missing trigger",
			skill: &Skill{
				Name:        "test-skill",
				Description: "Test description",
				Version:     "1.0",
				Trigger:     "",
				Response:    "Test response",
			},
			wantErr: true,
		},
		{
			name: "missing response",
			skill: &Skill{
				Name:        "test-skill",
				Description: "Test description",
				Version:     "1.0",
				Trigger:     "test-trigger",
				Response:    "",
			},
			wantErr: true,
		},
		{
			name: "with strength values",
			skill: &Skill{
				Name:        "test-skill",
				Description: "Test description",
				Version:     "1.0",
				Trigger:     "test-trigger",
				Response:    "Test response",
				Strength:    85.5,
				Threshold:  50.0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.skill.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Skill.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertPatternToSkill(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		pattern   *otrmodels.Pattern
		spaceName string
		wantName  string
		wantSpace string
	}{
		{
			name: "basic pattern conversion",
			pattern: &otrmodels.Pattern{
				ID:        "pattern-1",
				Trigger:   "hello",
				Response:  "Hello! How can I help you?",
				Strength:  80.0,
				Threshold: 30.0,
				Project:   "default",
				CreatedAt: now,
				UpdatedAt: now,
			},
			spaceName: "default",
			wantName:  "hello",
			wantSpace: "default",
		},
		{
			name: "pattern with nil last used",
			pattern: &otrmodels.Pattern{
				ID:        "pattern-2",
				Trigger:   "test",
				Response:  "Test response",
				Strength:  50.0,
				Threshold: 25.0,
			},
			spaceName: "work",
			wantName:  "test",
			wantSpace: "work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := ConvertPatternToSkill(tt.pattern, tt.spaceName)
			if skill.Name != tt.wantName {
				t.Errorf("Name = %v, want %v", skill.Name, tt.wantName)
			}
			if skill.SpaceName != tt.wantSpace {
				t.Errorf("SpaceName = %v, want %v", skill.SpaceName, tt.wantSpace)
			}
			if skill.Trigger != tt.pattern.Trigger {
				t.Errorf("Trigger = %v, want %v", skill.Trigger, tt.pattern.Trigger)
			}
			if skill.Response != tt.pattern.Response {
				t.Errorf("Response = %v, want %v", skill.Response, tt.pattern.Response)
			}
			if skill.Strength != tt.pattern.Strength {
				t.Errorf("Strength = %v, want %v", skill.Strength, tt.pattern.Strength)
			}
		})
	}
}

func TestConvertSkillToPattern(t *testing.T) {
	tests := []struct {
		name   string
		skill  *Skill
		wantID string
	}{
		{
			name: "basic skill conversion",
			skill: &Skill{
				Name:        "greeting",
				Description: "A friendly greeting",
				Version:     "1.0",
				Trigger:     "greeting",
				Response:    "Hello there!",
				Strength:    75.0,
				Threshold:   25.0,
				SpaceName:  "default",
			},
			wantID: "greeting",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := ConvertSkillToPattern(tt.skill)
			if pattern.Trigger != tt.skill.Trigger {
				t.Errorf("Trigger = %v, want %v", pattern.Trigger, tt.skill.Trigger)
			}
			if pattern.Response != tt.skill.Response {
				t.Errorf("Response = %v, want %v", pattern.Response, tt.skill.Response)
			}
			if pattern.Strength != tt.skill.Strength {
				t.Errorf("Strength = %v, want %v", pattern.Strength, tt.skill.Strength)
			}
			if pattern.Threshold != tt.skill.Threshold {
				t.Errorf("Threshold = %v, want %v", pattern.Threshold, tt.skill.Threshold)
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	original := &otrmodels.Pattern{
		ID:         "round-trip-test",
		Trigger:    "test-trigger",
		Response:   "Test response content with special chars: <>&\"'",
		Strength:   88.5,
		Threshold:  40.0,
		Project:    "default",
	}

	// Pattern -> Skill
	skill := ConvertPatternToSkill(original, "default")

	// Skill -> Pattern
	converted := ConvertSkillToPattern(skill)

	// Verify
	if converted.Trigger != original.Trigger {
		t.Errorf("Trigger mismatch: got %v, want %v", converted.Trigger, original.Trigger)
	}
	if converted.Response != original.Response {
		t.Errorf("Response mismatch: got %v, want %v", converted.Response, original.Response)
	}
	if converted.Strength != original.Strength {
		t.Errorf("Strength mismatch: got %v, want %v", converted.Strength, original.Strength)
	}
	if converted.Threshold != original.Threshold {
		t.Errorf("Threshold mismatch: got %v, want %v", converted.Threshold, original.Threshold)
	}
}
