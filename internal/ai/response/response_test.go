package response

import (
	"testing"

	"github.com/ArmyClaw/open-think-reflex/pkg/ai"
)

func TestParser_NewParser(t *testing.T) {
	p := NewParser()
	if p == nil {
		t.Error("NewParser returned nil")
	}
	if p.format != FormatAuto {
		t.Errorf("Expected default format FormatAuto, got %v", p.format)
	}
}

func TestParser_WithOptions(t *testing.T) {
	p := NewParser(WithFormat(FormatJSON))
	if p.format != FormatJSON {
		t.Errorf("Expected format FormatJSON, got %v", p.format)
	}
}

func TestParser_Parse_NilResponse(t *testing.T) {
	p := NewParser()
	_, err := p.Parse(nil)
	if err == nil {
		t.Error("Expected error for nil response, got nil")
	}
}

func TestParser_ParseJSON_ValidJSON(t *testing.T) {
	p := NewParser(WithFormat(FormatJSON))
	content := `{
		"content": "Hello, World!",
		"thoughts": [
			{
				"thought": "Greeting the user",
				"action": "Generate greeting",
				"score": 0.9
			}
		]
	}`

	result, err := p.ParseJSON(content)
	if err != nil {
		t.Errorf("ParseJSON failed: %v", err)
	}
	if result.Content != "Hello, World!" {
		t.Errorf("Expected content 'Hello, World!', got '%s'", result.Content)
	}
	if len(result.ThoughtSteps) != 1 {
		t.Errorf("Expected 1 thought step, got %d", len(result.ThoughtSteps))
	}
	if result.ThoughtSteps[0].Thought != "Greeting the user" {
		t.Errorf("Expected thought 'Greeting the user', got '%s'", result.ThoughtSteps[0].Thought)
	}
	if result.ThoughtSteps[0].Score != 0.9 {
		t.Errorf("Expected score 0.9, got %f", result.ThoughtSteps[0].Score)
	}
}

func TestParser_ParseJSON_InvalidJSON(t *testing.T) {
	p := NewParser(WithFormat(FormatJSON))
	content := `not valid json {`

	result, err := p.ParseJSON(content)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
	// Should still return a result with original content
	if result.Content != content {
		t.Errorf("Expected content to be original string, got '%s'", result.Content)
	}
}

func TestParser_parseThoughtChain_ThoughtAction(t *testing.T) {
	p := NewParser(WithFormat(FormatThoughtChain))
	content := `Thought: I need to help the user
Action: Generate a response
Observation: Response generated successfully`

	result, err := p.parseThoughtChain(content, &ParseResult{
		Content:     content,
		ThoughtSteps: []ai.ThoughtStep{},
		Metadata:    make(map[string]interface{}),
	})
	if err != nil {
		t.Errorf("parseThoughtChain failed: %v", err)
	}
	if len(result.ThoughtSteps) != 1 {
		t.Errorf("Expected 1 thought step, got %d", len(result.ThoughtSteps))
	}
	if result.ThoughtSteps[0].Thought != "I need to help the user" {
		t.Errorf("Expected thought 'I need to help the user', got '%s'", result.ThoughtSteps[0].Thought)
	}
	if result.ThoughtSteps[0].Action != "Generate a response" {
		t.Errorf("Expected action 'Generate a response', got '%s'", result.ThoughtSteps[0].Action)
	}
	if result.ThoughtSteps[0].Observation != "Response generated successfully" {
		t.Errorf("Expected observation 'Response generated successfully', got '%s'", result.ThoughtSteps[0].Observation)
	}
}

func TestParser_parseThoughtChain_MultipleSteps(t *testing.T) {
	p := NewParser(WithFormat(FormatThoughtChain))
	content := `Thought: First step
Action: Do action 1

Thought: Second step
Action: Do action 2
Observation: Done`

	result, err := p.parseThoughtChain(content, &ParseResult{
		Content:     content,
		ThoughtSteps: []ai.ThoughtStep{},
		Metadata:    make(map[string]interface{}),
	})
	if err != nil {
		t.Errorf("parseThoughtChain failed: %v", err)
	}
	if len(result.ThoughtSteps) != 2 {
		t.Errorf("Expected 2 thought steps, got %d", len(result.ThoughtSteps))
	}
	if result.ThoughtSteps[0].Thought != "First step" {
		t.Errorf("Expected first thought 'First step', got '%s'", result.ThoughtSteps[0].Thought)
	}
	if result.ThoughtSteps[1].Thought != "Second step" {
		t.Errorf("Expected second thought 'Second step', got '%s'", result.ThoughtSteps[1].Thought)
	}
}

func TestParser_parseThoughtChain_NumberedSteps(t *testing.T) {
	p := NewParser(WithFormat(FormatThoughtChain))
	content := `Step 1: Analyze the input
Step 2: Generate response
Step 3: Return result`

	result, err := p.parseThoughtChain(content, &ParseResult{
		Content:     content,
		ThoughtSteps: []ai.ThoughtStep{},
		Metadata:    make(map[string]interface{}),
	})
	if err != nil {
		t.Errorf("parseThoughtChain failed: %v", err)
	}
	if len(result.ThoughtSteps) != 3 {
		t.Errorf("Expected 3 thought steps, got %d", len(result.ThoughtSteps))
	}
}

func TestParser_parseThoughtChain_NoStructuredContent(t *testing.T) {
	p := NewParser(WithFormat(FormatThoughtChain))
	content := `This is just plain text without any structure.`

	result, err := p.parseThoughtChain(content, &ParseResult{
		Content:     content,
		ThoughtSteps: []ai.ThoughtStep{},
		Metadata:    make(map[string]interface{}),
	})
	if err != nil {
		t.Errorf("parseThoughtChain failed: %v", err)
	}
	// Should treat entire content as single step
	if len(result.ThoughtSteps) != 1 {
		t.Errorf("Expected 1 thought step for plain text, got %d", len(result.ThoughtSteps))
	}
}

func TestParser_ParseText(t *testing.T) {
	p := NewParser()
	content := "This is a plain text response"

	result := p.ParseText(content)
	if result.Content != content {
		t.Errorf("Expected content '%s', got '%s'", content, result.Content)
	}
	if result.Format != FormatText {
		t.Errorf("Expected format FormatText, got %v", result.Format)
	}
}

func TestParser_ParseToThoughtSteps(t *testing.T) {
	p := NewParser()
	resp := &ai.Response{
		Content: `Thought: Test thought
Action: Test action`,
		Model: "claude-3",
	}

	steps, err := p.ParseToThoughtSteps(resp)
	if err != nil {
		t.Errorf("ParseToThoughtSteps failed: %v", err)
	}
	if len(steps) == 0 {
		t.Error("Expected at least one thought step")
	}
}

func TestParser_FormatAsJSON(t *testing.T) {
	p := NewParser()
	result := &ParseResult{
		Content: "Test content",
		ThoughtSteps: []ai.ThoughtStep{
			{Thought: "Test thought", Action: "Test action", Score: 0.8},
		},
		Metadata: map[string]interface{}{
			"model": "claude-3",
		},
		Format: FormatText,
	}

	jsonStr, err := p.FormatAsJSON(result)
	if err != nil {
		t.Errorf("FormatAsJSON failed: %v", err)
	}
	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}
}

func TestParser_Parse_WithUsage(t *testing.T) {
	p := NewParser()
	resp := &ai.Response{
		Content:   "Test response",
		Model:     "claude-3-sonnet",
		Usage:     &ai.Usage{InputTokens: 100, OutputTokens: 50, TotalTokens: 150},
		FinishReason: "stop",
	}

	result, err := p.Parse(resp)
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	}
	if result.Metadata["input_tokens"] != 100 {
		t.Errorf("Expected input_tokens 100, got %v", result.Metadata["input_tokens"])
	}
	if result.Metadata["output_tokens"] != 50 {
		t.Errorf("Expected output_tokens 50, got %v", result.Metadata["output_tokens"])
	}
	if result.Metadata["total_tokens"] != 150 {
		t.Errorf("Expected total_tokens 150, got %v", result.Metadata["total_tokens"])
	}
}

func TestParser_Parse_WithResponseField(t *testing.T) {
	p := NewParser(WithFormat(FormatJSON))
	content := `{"response": "Hello from response field"}`

	result, err := p.ParseJSON(content)
	if err != nil {
		t.Errorf("ParseJSON failed: %v", err)
	}
	if result.Content != "Hello from response field" {
		t.Errorf("Expected content 'Hello from response field', got '%s'", result.Content)
	}
}

func TestParser_Parse_WithTextField(t *testing.T) {
	p := NewParser(WithFormat(FormatJSON))
	content := `{"text": "Hello from text field"}`

	result, err := p.ParseJSON(content)
	if err != nil {
		t.Errorf("ParseJSON failed: %v", err)
	}
	if result.Content != "Hello from text field" {
		t.Errorf("Expected content 'Hello from text field', got '%s'", result.Content)
	}
}

func TestDetectFormat_JSON(t *testing.T) {
	content := `{"key": "value"}`
	format := detectFormat(content)
	if format != FormatJSON {
		t.Errorf("Expected FormatJSON, got %v", format)
	}
}

func TestDetectFormat_ThoughtChain(t *testing.T) {
	content := `Thought: something
Action: do something`
	format := detectFormat(content)
	if format != FormatThoughtChain {
		t.Errorf("Expected FormatThoughtChain, got %v", format)
	}
}

func TestDetectFormat_Text(t *testing.T) {
	content := `This is just plain text without any special formatting.`
	format := detectFormat(content)
	if format != FormatText {
		t.Errorf("Expected FormatText, got %v", format)
	}
}

func TestResponseFormat_String(t *testing.T) {
	tests := []struct {
		format ResponseFormat
		want   string
	}{
		{FormatJSON, "JSON"},
		{FormatText, "Text"},
		{FormatThoughtChain, "ThoughtChain"},
		{FormatAuto, "Auto"},
		{ResponseFormat(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.format.String(); got != tt.want {
				t.Errorf("ResponseFormat.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_extractThoughtStep(t *testing.T) {
	p := NewParser()
	m := map[string]interface{}{
		"thought":     "Test thought",
		"action":      "Test action",
		"observation": "Test observation",
		"score":       0.75,
	}

	step := p.extractThoughtStep(m)
	if step.Thought != "Test thought" {
		t.Errorf("Expected thought 'Test thought', got '%s'", step.Thought)
	}
	if step.Action != "Test action" {
		t.Errorf("Expected action 'Test action', got '%s'", step.Action)
	}
	if step.Observation != "Test observation" {
		t.Errorf("Expected observation 'Test observation', got '%s'", step.Observation)
	}
	if step.Score != 0.75 {
		t.Errorf("Expected score 0.75, got %f", step.Score)
	}
}

func TestParser_extractThoughtStep_Capitalized(t *testing.T) {
	p := NewParser()
	m := map[string]interface{}{
		"Thought":     "Capitalized thought",
		"Action":      "Capitalized action",
		"Observation": "Capitalized observation",
		"Score":       0.5,
	}

	step := p.extractThoughtStep(m)
	if step.Thought != "Capitalized thought" {
		t.Errorf("Expected thought 'Capitalized thought', got '%s'", step.Thought)
	}
	if step.Action != "Capitalized action" {
		t.Errorf("Expected action 'Capitalized action', got '%s'", step.Action)
	}
	if step.Observation != "Capitalized observation" {
		t.Errorf("Expected observation 'Capitalized observation', got '%s'", step.Observation)
	}
	if step.Score != 0.5 {
		t.Errorf("Expected score 0.5, got %f", step.Score)
	}
}
