package ai

import (
	"context"
	"fmt"
	"io"
	"net/http"

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

	// Create streaming response - returns *ssestream.Stream[MessageStreamEventUnion]
	stream := p.client.Messages.NewStreaming(ctx, messageReq)

	// Return a reader that wraps the stream
	return &streamReader{stream: stream}, nil
}

// streamReader wraps the SSE stream to implement io.ReadCloser
type streamReader struct {
	stream *ssestream.Stream[anthropic.MessageStreamEventUnion]
}

func (r *streamReader) Read(p []byte) (n int, err error) {
	// This is a simplified implementation
	// In production, you'd properly handle the SSE stream using Next() and Current()
	return 0, io.EOF
}

func (r *streamReader) Close() error {
	if r.stream != nil {
		return r.stream.Close()
	}
	return nil
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
