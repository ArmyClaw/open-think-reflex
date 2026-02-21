// Package ai provides interfaces and implementations for AI provider integration.
// Currently supports Anthropic's Claude API. Designed to be extensible for
// additional providers (OpenAI, Google, etc.).
package ai

import (
	"context"
	"io"
)

// Provider defines the interface for AI model providers.
// Implementations must be thread-safe and handle concurrent requests.
type Provider interface {
	// Name returns the provider identifier (e.g., "claude", "openai")
	Name() string

	// Generate creates a complete response for the given prompt.
	// Blocks until the full response is received.
	Generate(ctx context.Context, req *Request) (*Response, error)

	// GenerateStream creates a streaming response.
	// Returns an io.ReadCloser that yields response chunks.
	// The caller is responsible for closing the reader.
	GenerateStream(ctx context.Context, req *Request) (io.ReadCloser, error)

	// ValidateKey checks if the configured API key is valid.
	// Returns nil if valid, error otherwise.
	ValidateKey(ctx context.Context) error
}

// Request represents a request to generate AI content.
// All fields are optional unless otherwise specified.
type Request struct {
	// Prompt is the user input / prompt (required)
	Prompt string

	// Model is the model identifier (e.g., "claude-3-sonnet")
	// Default: provider-specific default
	Model string

	// MaxTokens limits the maximum tokens in the response.
	// Default: 1024
	MaxTokens int

	// Temperature controls randomness (0.0 = deterministic, 2.0 = very random).
	// Default: 0.7
	Temperature float64

	// Stream enables streaming response mode.
	// Default: false (blocking)
	Stream bool

	// System is the system prompt / instructions.
	// Default: ""
	System string
}

// Response represents a complete AI generation response.
type Response struct {
	// Content is the generated text
	Content string

	// Model is the model that generated the response
	Model string

	// Usage contains token usage statistics
	Usage *Usage

	// FinishReason explains why generation stopped
	// (e.g., "stop", "length", "content_filtered")
	FinishReason string
}

// Usage represents token consumption information.
type Usage struct {
	InputTokens  int // Tokens in the request
	OutputTokens int // Tokens in the response
	TotalTokens  int // Input + Output
}

// ThoughtStep represents a single step in a thought chain (ReAct pattern).
// Used for structured reasoning and action planning.
type ThoughtStep struct {
	Thought      string // The reasoning/thinking
	Action       string // The action to take
	Observation  string // The result of the action
	Score        float64 // Confidence score (0-1)
}

// Config holds provider-specific configuration.
// Use functional options to configure.
type Config struct {
	APIKey      string // Provider API key
	Model       string // Model identifier
	MaxTokens   int    // Max response tokens
	Temperature float64 // Randomness factor
	Endpoint    string // Custom API endpoint (optional)
}

// Option is a functional option for configuring a Provider.
// Apply options using the provider constructor:
//   provider := NewClaudeProvider(WithAPIKey("key"), WithModel("model"))
type Option func(*Config)

// WithAPIKey sets the API key for authentication.
func WithAPIKey(key string) Option {
	return func(c *Config) {
		c.APIKey = key
	}
}

// WithModel sets the model identifier.
func WithModel(model string) Option {
	return func(c *Config) {
		c.Model = model
	}
}

// WithMaxTokens sets the maximum tokens in the response.
func WithMaxTokens(tokens int) Option {
	return func(c *Config) {
		c.MaxTokens = tokens
	}
}

// WithTemperature sets the temperature for response randomness.
func WithTemperature(temp float64) Option {
	return func(c *Config) {
		c.Temperature = temp
	}
}

// WithEndpoint sets a custom API endpoint (for proxy/regional endpoints).
func WithEndpoint(endpoint string) Option {
	return func(c *Config) {
		c.Endpoint = endpoint
	}
}
