// Package skills provides AgentSkill format handling for OTR patterns.
package skills

import (
	"context"
	"fmt"
	"time"

	otrmodels "github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/ArmyClaw/open-think-reflex/pkg/ai"
)

// Skill represents an AgentSkill exported from OTR.
// This follows the OpenClaw AgentSkill format.
type Skill struct {
	// Metadata
	Name        string            `json:"name"`         // Skill name (from trigger)
	Description string            `json:"description"` // Skill description (from response summary)
	Version     string            `json:"version"`     // Format version
	ExportedAt  time.Time        `json:"exported_at"` // Export timestamp
	Source      string            `json:"source"`      // Source: "otr"

	// Trigger and Response
	Trigger   string            `json:"trigger"`    // Original trigger
	Response  string            `json:"response"`   // Full response content
	Tags      []string          `json:"tags"`       // Pattern tags

	// Strength metadata (for reflex priority)
	Strength float64 `json:"strength"`    // Current strength (0-100)
	Threshold float64 `json:"threshold"` // Activation threshold

	// Usage stats
	UsageCount int       `json:"usage_count"`    // Times reinforced
	LastUsed   time.Time `json:"last_used"`       // Last used timestamp

	// Space info (v2.0)
	SpaceID   string `json:"space_id,omitempty"` // Source space
	SpaceName string `json:"space_name,omitempty"` // Source space name
}

// ConvertPatternToSkill converts an OTR Pattern to AgentSkill format.
func ConvertPatternToSkill(p *otrmodels.Pattern, spaceName string) *Skill {
	lastUsed := time.Time{}
	if p.LastUsedAt != nil {
		lastUsed = *p.LastUsedAt
	}

	return &Skill{
		Name:        p.Trigger,
		Description: extractDescription(p.Response),
		Version:     "1.0",
		ExportedAt:  time.Now(),
		Source:      "otr",
		Trigger:     p.Trigger,
		Response:    p.Response,
		Tags:        p.Tags,
		Strength:    p.Strength,
		Threshold:   p.Threshold,
		UsageCount:  p.ReinforceCnt,
		LastUsed:    lastUsed,
		SpaceID:     p.SpaceID,
		SpaceName:   spaceName,
	}
}

// extractDescription extracts a short description from response content.
// Takes the first line or first 100 characters.
func extractDescription(response string) string {
	if len(response) == 0 {
		return ""
	}
	// Take first line or first 100 chars
	lines := splitLines(response)
	if len(lines) > 0 && len(lines[0]) > 0 {
		desc := lines[0]
		if len(desc) > 100 {
			desc = desc[:100] + "..."
		}
		return desc
	}
	if len(response) > 100 {
		return response[:100] + "..."
	}
	return response
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// Validate validates the skill data.
func (s *Skill) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("skill name is required")
	}
	if s.Trigger == "" {
		return fmt.Errorf("trigger is required")
	}
	if s.Response == "" {
		return fmt.Errorf("response is required")
	}
	return nil
}

// PolishSkill uses AI to improve the skill response content.
// This enhances the exported skill with better formatting and clarity.
func PolishSkill(ctx context.Context, provider ai.Provider, skill *Skill) error {
	if provider == nil {
		return fmt.Errorf("AI provider is required")
	}

	polishPrompt := fmt.Sprintf(`Improve the following AgentSkill response content. 
Make it more clear, well-structured, and professional. Keep the same meaning but enhance readability.

Skill Name: %s
Current Response:
%s

Return only the improved response content, without explanations.`, skill.Name, skill.Response)

	req := &ai.Request{
		Prompt:    polishPrompt,
		MaxTokens: 2048,
		Temperature: 0.7,
		System:    "You are an expert at improving technical documentation and responses.",
	}

	resp, err := provider.Generate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to polish skill: %w", err)
	}

	skill.Response = resp.Content
	return nil
}

// ConvertSkillToPattern converts an AgentSkill back to OTR Pattern.
func ConvertSkillToPattern(s *Skill) *otrmodels.Pattern {
	p := otrmodels.NewPattern(s.Trigger, s.Response)
	p.Tags = s.Tags
	if s.Strength > 0 {
		p.Strength = s.Strength
	}
	if s.Threshold > 0 {
		p.Threshold = s.Threshold
	}
	if !s.LastUsed.IsZero() {
		p.LastUsedAt = &s.LastUsed
	}
	p.ReinforceCnt = s.UsageCount
	return p
}
