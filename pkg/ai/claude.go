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

// ClaudeProvider implements Provider for Anthropic's Claude
type ClaudeProvider struct {
	config     *Config
	httpClient *http.Client
	client     anthropic.Client
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(opts ...Option) *ClaudeProvider {
	cfg := &Config{
		Model:       "claude-3-sonnet-20240229",
		MaxTokens:   1024,
		Temperature: 0.7,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// Create Anthropic client with API key
	anthropicClient := anthropic.NewClient(
		option.WithAPIKey(cfg.APIKey),
	)

	return &ClaudeProvider{
		config:     cfg,
		client:     anthropicClient,
		httpClient: &http.Client{},
	}
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// Generate generates a response for the given prompt
func (p *ClaudeProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
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

	// Build message request
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

	// Send request
	resp, err := p.client.Messages.New(ctx, messageReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate: %w", err)
	}

	// Extract content from ContentBlockUnion
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

// GenerateStream generates a streaming response
func (p *ClaudeProvider) GenerateStream(ctx context.Context, req *Request) (io.ReadCloser, error) {
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

	// Create streaming response
	stream := p.client.Messages.NewStreaming(ctx, messageReq)

	// Return a reader that wraps the stream
	return newStreamReader(stream), nil
}

// streamReader wraps the SSE stream to implement io.ReadCloser
type streamReader struct {
	stream   *ssestream.Stream[anthropic.MessageStreamEventUnion]
	mu       sync.Mutex
	buffer   []byte
	closed   bool
}

// newStreamReader creates a new stream reader
func newStreamReader(stream *ssestream.Stream[anthropic.MessageStreamEventUnion]) *streamReader {
	return &streamReader{
		stream: stream,
		buffer: make([]byte, 0),
	}
}

// Read implements io.Reader
func (r *streamReader) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return 0, io.EOF
	}

	// If we have buffered data, return it
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

		// Handle different event types using type assertions
		var text string
		switch event.Type {
		case "content_block_delta":
			// Use the AsContentBlockDelta method to get the proper type
			deltaEvent := event.AsContentBlockDelta()
			text = string(deltaEvent.Delta.Text)
		case "message_delta":
			// Final message delta - can be used for usage stats
		case "message_stop":
			r.closed = true
			return 0, io.EOF
		}

		if text != "" {
			// Convert to SSE format
			data := "data: " + text + "\n\n"
			r.buffer = []byte(data)

			copy(p, r.buffer)
			r.buffer = r.buffer[len(p):]
			return len(p), nil
		}
	}
}

// Close implements io.Closer
func (r *streamReader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.closed = true
	if r.stream != nil {
		return r.stream.Close()
	}
	return nil
}

// StreamEvent represents a streaming event
type StreamEvent struct {
	Type      string          `json:"type"`
	Delta     json.RawMessage `json:"delta,omitempty"`
	Index     int             `json:"index,omitempty"`
	Content   string          `json:"content,omitempty"`
	Usage     *Usage          `json:"usage,omitempty"`
	FinishReason string       `json:"finish_reason,omitempty"`
}

// parseStreamEvent parses a line from the stream
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

// ValidateKey validates the API key
func (p *ClaudeProvider) ValidateKey(ctx context.Context) error {
	req := &Request{
		Prompt:   "Hello",
		MaxTokens: 10,
	}

	_, err := p.Generate(ctx, req)
	return err
}
