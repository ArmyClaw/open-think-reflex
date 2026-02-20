package response

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ArmyClaw/open-think-reflex/pkg/ai"
)

// Parser parses AI responses into structured data
type Parser struct {
	format ResponseFormat
}

// ResponseFormat defines how to parse the response
type ResponseFormat int

const (
	FormatAuto ResponseFormat = iota
	FormatJSON
	FormatText
	FormatThoughtChain
)

// Option is a functional option for Parser
type Option func(*Parser)

// WithFormat sets the response format
func WithFormat(format ResponseFormat) Option {
	return func(p *Parser) {
		p.format = format
	}
}

// NewParser creates a new response parser
func NewParser(opts ...Option) *Parser {
	p := &Parser{
		format: FormatAuto,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// ParseResult represents the parsed result
type ParseResult struct {
	Content     string
	ThoughtSteps []ai.ThoughtStep
	Metadata    map[string]interface{}
	Format      ResponseFormat
}

// Parse parses an AI response into structured data
func (p *Parser) Parse(resp *ai.Response) (*ParseResult, error) {
	if resp == nil {
		return nil, fmt.Errorf("nil response")
	}

	result := &ParseResult{
		Content:     resp.Content,
		ThoughtSteps: []ai.ThoughtStep{},
		Metadata:    make(map[string]interface{}),
	}

	// Add usage metadata
	if resp.Usage != nil {
		result.Metadata["input_tokens"] = resp.Usage.InputTokens
		result.Metadata["output_tokens"] = resp.Usage.OutputTokens
		result.Metadata["total_tokens"] = resp.Usage.TotalTokens
	}
	result.Metadata["model"] = resp.Model
	result.Metadata["finish_reason"] = resp.FinishReason

	// Auto-detect format or use specified format
	format := p.format
	if format == FormatAuto {
		format = detectFormat(resp.Content)
	}
	result.Format = format

	// Parse based on detected format
	switch format {
	case FormatJSON:
		return p.parseJSON(resp.Content, result)
	case FormatThoughtChain:
		return p.parseThoughtChain(resp.Content, result)
	case FormatText:
		// Text format - no special parsing needed
		return result, nil
	default:
		// Try JSON first, then thought chain, then plain text
		if jsonResult, err := p.parseJSON(resp.Content, result); err == nil {
			return jsonResult, nil
		}
		if tcResult, err := p.parseThoughtChain(resp.Content, result); err == nil {
			return tcResult, nil
		}
		return result, nil
	}
}

// detectFormat auto-detects the response format
func detectFormat(content string) ResponseFormat {
	content = strings.TrimSpace(content)

	// Check if it's JSON
	if strings.HasPrefix(content, "{") || strings.HasPrefix(content, "[") {
		if _, err := json.MarshalIndent(nil, "", ""); err == nil {
			// Try to parse as JSON
			var js json.RawMessage
			if json.Unmarshal([]byte(content), &js) == nil {
				return FormatJSON
			}
		}
	}

	// Check for thought chain patterns
	thoughtChainIndicators := []string{
		"Thought:",
		"Action:",
		"Observation:",
		"Step 1:",
		"Step 2:",
		"Step 3:",
		"**Thought**",
		"**Action**",
	}
	contentLower := strings.ToLower(content)
	for _, indicator := range thoughtChainIndicators {
		if strings.Contains(contentLower, strings.ToLower(indicator)) {
			return FormatThoughtChain
		}
	}

	return FormatText
}

// parseJSON parses JSON formatted responses
func (p *Parser) ParseJSON(content string) (*ParseResult, error) {
	result := &ParseResult{
		Content:     content,
		ThoughtSteps: []ai.ThoughtStep{},
		Metadata:    make(map[string]interface{}),
	}
	return p.parseJSON(content, result)
}

func (p *Parser) parseJSON(content string, result *ParseResult) (*ParseResult, error) {
	// Try to unmarshal as JSON
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return result, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Try to extract thought steps if present
	if steps, ok := data["thoughts"].([]interface{}); ok {
		for _, step := range steps {
			if stepMap, ok := step.(map[string]interface{}); ok {
				thoughtStep := p.extractThoughtStep(stepMap)
				result.ThoughtSteps = append(result.ThoughtSteps, thoughtStep)
			}
		}
	}

	// Try to extract content
	if contentField, ok := data["content"].(string); ok {
		result.Content = contentField
	} else if responseField, ok := data["response"].(string); ok {
		result.Content = responseField
	} else if textField, ok := data["text"].(string); ok {
		result.Content = textField
	}

	// Store remaining fields as metadata
	for k, v := range data {
		if k != "thoughts" && k != "content" && k != "response" && k != "text" {
			result.Metadata[k] = v
		}
	}

	result.Format = FormatJSON
	return result, nil
}

// parseThoughtChain parses thought chain formatted responses
func (p *Parser) parseThoughtChain(content string, result *ParseResult) (*ParseResult, error) {
	result.Format = FormatThoughtChain

	// Pattern to match thought-action-observation blocks
	// Matches formats like:
	// Thought: ... 
	// Action: ...
	// Observation: ...
	// or
	// **Thought** ...
	// **Action** ...
	thoughtPattern := regexp.MustCompile(`(?i)(?:Thought|思考):\s*(.+?)(?:\n|$)`)
	actionPattern := regexp.MustCompile(`(?i)(?:Action|行动):\s*(.+?)(?:\n|$)`)
	observationPattern := regexp.MustCompile(`(?i)(?:Observation|观察):\s*(.+?)(?:\n|$)`)

	// Also match numbered steps
	stepPattern := regexp.MustCompile(`(?i)(?:Step\s*(\d+)|[-•*]\s*)\s*(.+?)(?:\n|$)`)

	lines := strings.Split(content, "\n")
	currentStep := ai.ThoughtStep{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try to match thought-action-observation
		if match := thoughtPattern.FindStringSubmatch(line); len(match) > 1 {
			if currentStep.Thought != "" {
				result.ThoughtSteps = append(result.ThoughtSteps, currentStep)
				currentStep = ai.ThoughtStep{}
			}
			currentStep.Thought = strings.TrimSpace(match[1])
		} else if match := actionPattern.FindStringSubmatch(line); len(match) > 1 {
			currentStep.Action = strings.TrimSpace(match[1])
		} else if match := observationPattern.FindStringSubmatch(line); len(match) > 1 {
			currentStep.Observation = strings.TrimSpace(match[1])
			result.ThoughtSteps = append(result.ThoughtSteps, currentStep)
			currentStep = ai.ThoughtStep{}
		} else if match := stepPattern.FindStringSubmatch(line); len(match) > 2 {
			// Handle numbered/bullet steps
			if currentStep.Thought != "" {
				result.ThoughtSteps = append(result.ThoughtSteps, currentStep)
				currentStep = ai.ThoughtStep{}
			}
			currentStep.Thought = strings.TrimSpace(match[2])
		}
	}

	// Add last step if exists
	if currentStep.Thought != "" {
		result.ThoughtSteps = append(result.ThoughtSteps, currentStep)
	}

	// If no structured steps found, treat entire content as single step
	if len(result.ThoughtSteps) == 0 {
		result.ThoughtSteps = []ai.ThoughtStep{
			{Thought: content},
		}
	}

	return result, nil
}

// extractThoughtStep extracts a ThoughtStep from a map
func (p *Parser) extractThoughtStep(m map[string]interface{}) ai.ThoughtStep {
	step := ai.ThoughtStep{}

	if thought, ok := m["thought"].(string); ok {
		step.Thought = thought
	} else if t, ok := m["Thought"].(string); ok {
		step.Thought = t
	}

	if action, ok := m["action"].(string); ok {
		step.Action = action
	} else if a, ok := m["Action"].(string); ok {
		step.Action = a
	}

	if observation, ok := m["observation"].(string); ok {
		step.Observation = observation
	} else if o, ok := m["Observation"].(string); ok {
		step.Observation = o
	}

	// Try to extract score
	if score, ok := m["score"].(float64); ok {
		step.Score = score
	} else if score, ok := m["Score"].(float64); ok {
		step.Score = score
	}

	return step
}

// ParseText parses plain text responses
func (p *Parser) ParseText(content string) *ParseResult {
	return &ParseResult{
		Content:      content,
		ThoughtSteps: []ai.ThoughtStep{},
		Metadata:     make(map[string]interface{}),
		Format:       FormatText,
	}
}

// ParseToThoughtSteps extracts just the thought steps from a response
func (p *Parser) ParseToThoughtSteps(resp *ai.Response) ([]ai.ThoughtStep, error) {
	result, err := p.Parse(resp)
	if err != nil {
		return nil, err
	}
	return result.ThoughtSteps, nil
}

// FormatAsJSON formats the response as JSON string
func (p *Parser) FormatAsJSON(result *ParseResult) (string, error) {
	data := map[string]interface{}{
		"content":      result.Content,
		"format":       result.Format.String(),
		"metadata":     result.Metadata,
		"thoughtSteps": result.ThoughtSteps,
	}
	bytes, err := json.MarshalIndent(data, "", "  ")
	return string(bytes), err
}

// String returns the string representation of ResponseFormat
func (f ResponseFormat) String() string {
	switch f {
	case FormatJSON:
		return "JSON"
	case FormatText:
		return "Text"
	case FormatThoughtChain:
		return "ThoughtChain"
	case FormatAuto:
		return "Auto"
	default:
		return "Unknown"
	}
}
