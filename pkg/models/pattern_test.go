package models

import (
	"testing"
)

func TestNewPattern(t *testing.T) {
	p := NewPattern("test trigger", "test response")

	if p.Trigger != "test trigger" {
		t.Errorf("Expected trigger 'test trigger', got '%s'", p.Trigger)
	}

	if p.Response != "test response" {
		t.Errorf("Expected response 'test response', got '%s'", p.Response)
	}

	if p.Strength != 0.0 {
		t.Errorf("Expected strength 0.0, got %f", p.Strength)
	}

	if p.Threshold != 50.0 {
		t.Errorf("Expected threshold 50.0, got %f", p.Threshold)
	}

	if p.DecayRate != 0.01 {
		t.Errorf("Expected decay_rate 0.01, got %f", p.DecayRate)
	}

	if p.DecayEnabled != true {
		t.Error("Expected DecayEnabled to be true")
	}

	if p.ID == "" {
		t.Error("Expected non-empty ID")
	}

	if p.CreatedAt.IsZero() {
		t.Error("Expected non-zero CreatedAt")
	}
}

func TestPatternValidateTrigger(t *testing.T) {
	tests := []struct {
		name    string
		trigger string
		wantErr bool
	}{
		{"valid trigger", "valid trigger", false},
		{"empty trigger", "", true},
		{"trigger with spaces", "  test  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pattern{Trigger: tt.trigger}
			err := p.ValidateTrigger()

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTrigger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatternValidateResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{"valid response", "test response", false},
		{"empty response", "", true},
		{"response with spaces", "  content  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pattern{Response: tt.response}
			err := p.ValidateResponse()

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatternValidateStrength(t *testing.T) {
	tests := []struct {
		name      string
		strength  float64
		threshold float64
		decayRate float64
		wantErr   bool
	}{
		{"valid values", 50.0, 50.0, 0.01, false},
		{"zero strength", 0.0, 50.0, 0.01, false},
		{"max strength", 100.0, 50.0, 0.01, false},
		{"strength too low", -1.0, 50.0, 0.01, true},
		{"strength too high", 101.0, 50.0, 0.01, true},
		{"threshold too low", 50.0, -1.0, 0.01, true},
		{"threshold too high", 50.0, 101.0, 0.01, true},
		{"decay rate too high", 50.0, 50.0, 1.01, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pattern{
				Strength:   tt.strength,
				Threshold:  tt.threshold,
				DecayRate: tt.decayRate,
			}
			err := p.ValidateStrength()

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStrength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatternFullValidation(t *testing.T) {
	p := &Pattern{
		Trigger:   "test",
		Response:  "response",
		Strength:  50.0,
		Threshold: 50.0,
		DecayRate: 0.01,
	}

	if err := p.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}
}

func TestPatternInvalidValidation(t *testing.T) {
	p := &Pattern{
		Trigger:   "",
		Response:  "response",
		Strength:  50.0,
		Threshold: 50.0,
		DecayRate: 0.01,
	}

	if err := p.Validate(); err == nil {
		t.Error("Validate() expected error for empty trigger")
	}
}
