package ai

import (
	"context"
	"io"
)

// Provider defines the interface for AI providers
type Provider interface {
	// Name returns the provider name
	Name() string
	
	// Generate generates a response for the given prompt
	Generate(ctx context.Context, req *Request) (*Response, error)
	
	// GenerateStream generates a streaming response
	GenerateStream(ctx context.Context, req *Request) (io.ReadCloser, error)
	
	// ValidateKey validates the API key
	ValidateKey(ctx context.Context) error
}

// Request represents an AI generation request
type Request struct {
	Prompt      string
	Model       string
	MaxTokens   int
	Temperature float64
	Stream      bool
	System      string
}

// Response represents an AI generation response
type Response struct {
	Content   string
	Model     string
	Usage     *Usage
	FinishReason string
}

// Usage represents token usage information
type Usage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens int
}

// ThoughtStep represents a step in the thought chain
type ThoughtStep struct {
	Thought   string
	Action    string
	Observation string
	Score     float64
}

// Config holds provider configuration
type Config struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
	Endpoint    string
}

// Option is a functional option for configuring the provider
type Option func(*Config)

// WithAPIKey sets the API key
func WithAPIKey(key string) Option {
	return func(c *Config) {
		c.APIKey = key
	}
}

// WithModel sets the model
func WithModel(model string) Option {
	return func(c *Config) {
		c.Model = model
	}
}

// WithMaxTokens sets the max tokens
func WithMaxTokens(tokens int) Option {
	return func(c *Config) {
		c.MaxTokens = tokens
	}
}

// WithTemperature sets the temperature
func WithTemperature(temp float64) Option {
	return func(c *Config) {
		c.Temperature = temp
	}
}

// WithEndpoint sets a custom endpoint
func WithEndpoint(endpoint string) Option {
	return func(c *Config) {
		c.Endpoint = endpoint
	}
}
