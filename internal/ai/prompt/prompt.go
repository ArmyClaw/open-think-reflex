package prompt

import (
	"fmt"
	"strings"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// Builder builds prompts for AI generation
type Builder struct {
	systemPrompt string
	 maxTokens   int
}

// NewBuilder creates a new prompt builder
func NewBuilder() *Builder {
	return &Builder{
		systemPrompt: defaultSystemPrompt,
		maxTokens:    1024,
	}
}

// WithOptions applies options to the builder
func (b *Builder) WithOptions(opts ...Option) *Builder {
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Option is a functional option for Builder
type Option func(*Builder)

// WithSystemPrompt sets a custom system prompt
func WithSystemPrompt(prompt string) Option {
	return func(b *Builder) {
		b.systemPrompt = prompt
	}
}

// WithMaxTokens sets max tokens
func WithMaxTokens(tokens int) Option {
	return func(b *Builder) {
		b.maxTokens = tokens
	}
}

// BuildRequest builds a prompt request from input and matched patterns
func (b *Builder) BuildRequest(input string, matchedPatterns []*models.Pattern) string {
	var sb strings.Builder

	// Add system prompt
	sb.WriteString(b.systemPrompt)
	sb.WriteString("\n\n")

	// Add context from matched patterns
	if len(matchedPatterns) > 0 {
		sb.WriteString("## Relevant Patterns\n")
		sb.WriteString("The following patterns are relevant to the user's input:\n\n")

		for i, p := range matchedPatterns {
			sb.WriteString(fmt.Sprintf("### Pattern %d\n", i+1))
			sb.WriteString(fmt.Sprintf("**Trigger**: %s\n", p.Trigger))
			sb.WriteString(fmt.Sprintf("**Response**: %s\n", p.Response))
			
			if len(p.Tags) > 0 {
				sb.WriteString(fmt.Sprintf("**Tags**: %s\n", strings.Join(p.Tags, ", ")))
			}
			if p.Project != "" {
				sb.WriteString(fmt.Sprintf("**Project**: %s\n", p.Project))
			}
			sb.WriteString(fmt.Sprintf("**Strength**: %.1f/100\n", p.Strength))
			sb.WriteString("\n")
		}

		sb.WriteString("---\n\n")
	}

	// Add user input
	sb.WriteString("## User Input\n")
	sb.WriteString(input)

	return sb.String()
}

// BuildSystemPrompt builds just the system prompt with optional context
func (b *Builder) BuildSystemPrompt(userContext string) string {
	var sb strings.Builder
	sb.WriteString(b.systemPrompt)

	if userContext != "" {
		sb.WriteString("\n\n## User Context\n")
		sb.WriteString(userContext)
	}

	return sb.String()
}

// BuildReflexPrompt builds a prompt specifically for reflex generation
func (b *Builder) BuildReflexPrompt(input string, patterns []*models.Pattern) string {
	var sb strings.Builder

	// System prompt for reflex generation
	sb.WriteString(`You are Open-Think-Reflex, an AI input accelerator. 
Your task is to help the user by generating relevant responses based on matched patterns.

Guidelines:
1. Use the provided patterns to generate contextually appropriate responses
2. Be concise and helpful
3. If no patterns match well, generate a reasonable response based on the input
4. Consider the strength of each pattern (higher strength = more relevant)
5. Combine multiple patterns if they are all relevant

`)

	// Add patterns if available
	if len(patterns) > 0 {
		sb.WriteString("## Available Patterns\n\n")
		for i, p := range patterns {
			sb.WriteString(fmt.Sprintf("%d. Trigger: \"%s\" | Response: \"%s\" | Strength: %.1f\n", 
				i+1, p.Trigger, p.Response, p.Strength))
		}
		sb.WriteString("\n")
	}

	// Add user input
	sb.WriteString("## Task\n")
	sb.WriteString("Generate a response for the following input:\n\n")
	sb.WriteString(input)

	return sb.String()
}

// defaultSystemPrompt is the default system prompt
const defaultSystemPrompt = `You are Open-Think-Reflex, an AI-powered input accelerator.

Your role is to help users by:
1. Understanding their input context
2. Generating relevant responses based on matched patterns
3. Providing helpful suggestions

Always be concise, accurate, and context-aware.`
