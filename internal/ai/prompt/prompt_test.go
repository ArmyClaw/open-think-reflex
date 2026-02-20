package prompt

import (
	"testing"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

func TestBuilder_BuildRequest(t *testing.T) {
	builder := NewBuilder()

	patterns := []*models.Pattern{
		{
			Trigger:   "hello",
			Response:  "Hi there! How can I help you?",
			Tags:      []string{"greeting", "hello"},
			Project:   "general",
			Strength:  80.0,
		},
		{
			Trigger:   "help",
			Response:  "I'm here to help!",
			Tags:      []string{"help"},
			Project:   "general",
			Strength:  60.0,
		},
	}

	input := "Hello, I need some help"
	result := builder.BuildRequest(input, patterns)

	// Verify system prompt is present
	if !contains(result, "Open-Think-Reflex") {
		t.Error("Expected system prompt to be present")
	}

	// Verify patterns are included
	if !contains(result, "hello") {
		t.Error("Expected pattern trigger to be present")
	}

	// Verify user input is present
	if !contains(result, input) {
		t.Error("Expected user input to be present")
	}
}

func TestBuilder_BuildReflexPrompt(t *testing.T) {
	builder := NewBuilder()

	patterns := []*models.Pattern{
		{
			Trigger:   "test",
			Response:  "This is a test response",
			Strength:  70.0,
		},
	}

	input := "Run the test"
	result := builder.BuildReflexPrompt(input, patterns)

	// Verify system prompt
	if !contains(result, "Open-Think-Reflex") {
		t.Error("Expected system prompt to be present")
	}

	// Verify user input
	if !contains(result, input) {
		t.Error("Expected user input to be present")
	}

	// Verify patterns
	if !contains(result, "test") {
		t.Error("Expected pattern to be present")
	}
}

func TestBuilder_WithSystemPrompt(t *testing.T) {
	customPrompt := "Custom system prompt"
	builder := NewBuilder()
	builder = builder.WithOptions(WithSystemPrompt(customPrompt))

	result := builder.BuildRequest("test input", nil)

	if !contains(result, customPrompt) {
		t.Error("Expected custom system prompt to be present")
	}
}

func TestBuilder_BuildSystemPrompt(t *testing.T) {
	builder := NewBuilder()

	context := "User is working on project X"
	result := builder.BuildSystemPrompt(context)

	if !contains(result, "Open-Think-Reflex") {
		t.Error("Expected system prompt to be present")
	}

	if !contains(result, context) {
		t.Error("Expected user context to be present")
	}
}

func TestBuilder_EmptyPatterns(t *testing.T) {
	builder := NewBuilder()

	input := "Test input"
	result := builder.BuildRequest(input, nil)

	if !contains(result, input) {
		t.Error("Expected user input to be present even with no patterns")
	}

	// Should not have pattern section
	if contains(result, "Relevant Patterns") {
		t.Error("Should not have pattern section when no patterns")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		(len(s) >= len(substr)) &&
		(s == substr || 
		 findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
