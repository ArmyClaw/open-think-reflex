package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ClaudeProvider implements Provider for Anthropic's Claude
type ClaudeProvider struct {
	config *Config
	client *http.Client
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(opts ...Option) *ClaudeProvider {
	cfg := &Config{
		Model:     "claude-3-sonnet-20240229",
		MaxTokens: 1024,
		Temperature: 0.7,
	}
	
	for _, opt := range opts {
		opt(cfg)
	}
	
	return &ClaudeProvider{
		config: cfg,
		client: &http.Client{},
	}
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// Generate generates a response for the given prompt
func (p *ClaudeProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	if req.Model == "" {
		req.Model = p.config.Model
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = p.config.MaxTokens
	}
	if req.Temperature == 0 {
		req.Temperature = p.config.Temperature
	}
	
	body := p.buildRequestBody(req)
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.config.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(body))
	}
	
	return p.parseResponse(resp.Body)
}

// GenerateStream generates a streaming response
func (p *ClaudeProvider) GenerateStream(ctx context.Context, req *Request) (io.ReadCloser, error) {
	req.Stream = true
	
	if req.Model == "" {
		req.Model = p.config.Model
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = p.config.MaxTokens
	}
	
	body := p.buildRequestBody(req)
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.config.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Accept", "text/event-stream")
	
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(body))
	}
	
	return resp.Body, nil
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

func (p *ClaudeProvider) buildRequestBody(req *Request) []byte {
	messages := []map[string]string{
		{"role": "user", "content": req.Prompt},
	}
	
	if req.System != "" {
		return map[string]interface{}{
			"model":       req.Model,
			"max_tokens":  req.MaxTokens,
			"temperature": req.Temperature,
			"system":      req.System,
			"messages":    messages,
			"stream":      req.Stream,
		}.toJSON()
	}
	
	return map[string]interface{}{
		"model":       req.Model,
		"max_tokens":  req.MaxTokens,
		"temperature": req.Temperature,
		"messages":    messages,
		"stream":      req.Stream,
	}.toJSON()
}

func (p *ClaudeProvider) parseResponse(body io.Reader) (*Response, error) {
	var resp struct {
		Type     string `json:"type"`
		Content  []struct {
			Type     string `json:"type"`
			Text     string `json:"text"`
		} `json:"content"`
		Model      string `json:"model"`
		Usage      struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
		StopReason string `json:"stop_reason"`
	}
	
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	if resp.Type == "error" {
		return nil, fmt.Errorf("API error: %v", resp.Content)
	}
	
	var content strings.Builder
	for _, c := range resp.Content {
		content.WriteString(c.Text)
	}
	
	return &Response{
		Content:    content.String(),
		Model:      resp.Model,
		Usage:      &Usage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
			TotalTokens:  resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
		FinishReason: resp.StopReason,
	}, nil
}

type jsonMap map[string]interface{}

func (m jsonMap) toJSON() []byte {
	buf := &strings.Builder{}
	encoder := json.NewEncoder(buf)
	encoder.Encode(m)
	return []byte(buf.String())
}
