// Package ai provides AI provider implementations.
// This file implements the Anthropic Claude provider using the official SDK.
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

// ClaudeProvider implements the Provider interface for Anthropic's Claude models.
// Uses the official Anthropic SDK for API communication.
// Thread-safe: can handle concurrent requests.
type ClaudeProvider struct {
	config     *Config            // Provider configuration
	httpClient *http.Client       // HTTP client for custom requests
	client     anthropic.Client   // Official Anthropic SDK client
}

// NewClaudeProvider creates a new Claude provider with the given options.
// Default configuration:
//   - Model: claude-3-sonnet-20240229
//   - MaxTokens: 1024
//   - Temperature: 0.7
//
// Example:
//   provider := NewClaudeProvider(
//       WithAPIKey("sk-ant-..."),
//       WithModel("claude-3-5-sonnet-20241022"),
//   )
func NewClaudeProvider(opts ...Option) *ClaudeProvider {
	cfg := &Config{
		Model:       "claude-3-sonnet-20240229",
		MaxTokens:   1024,
		Temperature: 0.7,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(cfg)
	}

	// Initialize Anthropic SDK client
	anthropicClient := anthropic.NewClient(
		option.WithAPIKey(cfg.APIKey),
	)

	return &ClaudeProvider{
		config:     cfg,
		client:     anthropicClient,
		httpClient: &http.Client{},
	}
}

// Name returns the provider identifier.
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// Generate creates a complete response from Claude.
// Blocks until the full response is received.
func (p *ClaudeProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	// Apply defaults from config if not specified in request
	model := req.Model
	if model == "" {
		model = p.config.Model
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = p.config.MaxTokens
	}

	temp := req.Temperature
	if temp == 0 {
		temp = p.config.Temperature
	}

	// Build message request parameters
	messageReq := anthropic.MessageNewParams{
		Model:       anthropic.Model(model),
		MaxTokens:   int64(maxTokens),
		Temperature: anthropic.Float(temp),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(req.Prompt),
			),
		},
	}

	// Add system prompt if provided
	if req.System != "" {
		messageReq.System = []anthropic.TextBlockParam{
			{
				Type: "text",
				Text: req.System,
			},
		}
	}

	// Send request to Claude API
	resp, err := p.client.Messages.New(ctx, messageReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate: %w", err)
	}

	// Extract text content from response blocks
	// Claude can return multiple content blocks (text, images, etc.)
	var content string
	for _, block := range resp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &Response{
		Content: content,
		Model:   string(resp.Model),
		Usage: &Usage{
			InputTokens:  int(resp.Usage.InputTokens),
			OutputTokens: int(resp.Usage.OutputTokens),
			TotalTokens:  int(resp.Usage.InputTokens + resp.Usage.OutputTokens),
		},
		FinishReason: string(resp.StopReason),
	}, nil
}

// GenerateStream creates a streaming response from Claude.
// Returns an io.ReadCloser that yields response chunks as they arrive.
// The caller MUST close the reader to release resources.
//
// The stream format is Server-Sent Events (SSE):
//   data: Hello
//   data: !
//
// Example usage:
//   reader, err := provider.GenerateStream(ctx, req)
//   if err != nil { ... }
//   defer reader.Close()
//   io.Copy(os.Stdout, reader)
func (p *ClaudeProvider) GenerateStream(ctx context.Context, req *Request) (io.ReadCloser, error) {
	// Apply defaults
	model := req.Model
	if model == "" {
		model = p.config.Model
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = p.config.MaxTokens
	}

	// Build streaming message request
	messageReq := anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: int64(maxTokens),
		Temperature: anthropic.Float(req.Temperature),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(req.Prompt),
			),
		},
	}

	if req.System != "" {
		messageReq.System = []anthropic.TextBlockParam{
			{
				Type: "text",
				Text: req.System,
			},
		}
	}

	// Create streaming response from Anthropic SDK
	stream := p.client.Messages.NewStreaming(ctx, messageReq)

	// Wrap in our streamReader for io.ReadCloser compatibility
	return newStreamReader(stream), nil
}

// streamReader wraps the Anthropic SSE stream to implement io.ReadCloser.
// Handles the complex event types from Anthropic's streaming API.
type streamReader struct {
	stream *ssestream.Stream[anthropic.MessageStreamEventUnion]
	mu     sync.Mutex
	buffer []byte   // Buffered data not yet returned
	closed bool     // Whether Close() was called
}

// newStreamReader creates a new stream reader wrapping the Anthropic stream.
func newStreamReader(stream *ssestream.Stream[anthropic.MessageStreamEventUnion]) *streamReader {
	return &streamReader{
		stream: stream,
		buffer: make([]byte, 0),
	}
}

// Read implements io.Reader.
// Reads the next chunk of the streaming response.
// Returns io.EOF when the stream is complete.
func (r *streamReader) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return 0, io.EOF
	}

	// Return buffered data first
	if len(r.buffer) > 0 {
		copy(p, r.buffer)
		r.buffer = r.buffer[len(p):]
		if len(r.buffer) == 0 {
			return len(p), nil
		}
		return len(p), nil
	}

	// Try to read from stream
	for {
		if !r.stream.Next() {
			r.closed = true
			if err := r.stream.Err(); err != nil {
				return 0, err
			}
			return 0, io.EOF
		}

		event := r.stream.Current()

		// Handle different event types
		var text string
		switch event.Type {
		case "content_block_delta":
			// Text delta - most common event type
			deltaEvent := event.AsContentBlockDelta()
			text = string(deltaEvent.Delta.Text)
		case "message_delta":
			// Final message delta - contains usage stats
		case "message_stop":
			// Stream complete
			r.closed = true
			return 0, io.EOF
		}

		if text != "" {
			// Format as SSE (Server-Sent Events)
			data := "data: " + text + "\n\n"
			r.buffer = []byte(data)

			copy(p, r.buffer)
			r.buffer = r.buffer[len(p):]
			return len(p), nil
		}
	}
}

// Close implements io.Closer.
// Must be called when done reading to release resources.
func (r *streamReader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.closed = true
	if r.stream != nil {
		return r.stream.Close()
	}
	return nil
}

// StreamEvent represents a parsed streaming event.
// Used internally for SSE parsing.
type StreamEvent struct {
	Type         string          `json:"type"`
	Delta        json.RawMessage `json:"delta,omitempty"`
	Index        int             `json:"index,omitempty"`
	Content      string          `json:"content,omitempty"`
	Usage        *Usage          `json:"usage,omitempty"`
	FinishReason string          `json:"finish_reason,omitempty"`
}

// parseStreamEvent parses a line from the SSE stream.
// Returns nil if the line is not a data event.
func parseStreamEvent(line string) (*StreamEvent, error) {
	if !strings.HasPrefix(line, "data: ") {
		return nil, nil
	}

	data := strings.TrimPrefix(line, "data: ")
	if data == "[DONE]" {
		return &StreamEvent{Type: "stop"}, nil
	}

	var event StreamEvent
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return nil, fmt.Errorf("failed to parse stream event: %w", err)
	}

	return &event, nil
}

// ValidateKey checks if the API key is valid by making a minimal request.
// Returns nil if the key is valid, error otherwise.
func (p *ClaudeProvider) ValidateKey(ctx context.Context) error {
	// Use minimal request to validate key
	req := &Request{
		Prompt:   "Hello",
		MaxTokens: 10,
	}

	_, err := p.Generate(ctx, req)
	return err
}
