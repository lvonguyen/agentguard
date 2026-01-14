package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	AnthropicAPIURL     = "https://api.anthropic.com/v1/messages"
	AnthropicAPIVersion = "2023-06-01"
	DefaultModel        = "claude-opus-4-5-20250514"
	DefaultMaxTokens    = 4096
)

// AnthropicConfig holds configuration for the Anthropic provider
type AnthropicConfig struct {
	APIKey    string
	Model     string
	MaxTokens int
}

// AnthropicProvider implements the LLM Provider interface for Claude
type AnthropicProvider struct {
	apiKey    string
	model     string
	maxTokens int
	client    *http.Client
}

// NewAnthropicProvider creates a new Anthropic/Claude provider
func NewAnthropicProvider(cfg AnthropicConfig) (*AnthropicProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("anthropic API key is required")
	}

	model := cfg.Model
	if model == "" {
		model = DefaultModel
	}

	maxTokens := cfg.MaxTokens
	if maxTokens == 0 {
		maxTokens = DefaultMaxTokens
	}

	return &AnthropicProvider{
		apiKey:    cfg.APIKey,
		model:     model,
		maxTokens: maxTokens,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

// Message represents a conversation message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest represents a request to the Anthropic API
type CompletionRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
	Messages  []Message `json:"messages"`
	Stream    bool      `json:"stream,omitempty"`
}

// CompletionResponse represents a response from the Anthropic API
type CompletionResponse struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Role         string `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence,omitempty"`
	Usage        Usage  `json:"usage"`
}

// ContentBlock represents a content block in the response
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Usage represents token usage information
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Messages      []Message
	SystemPrompt  string
	Context       string // Additional context from vector search
	MaxTokens     int
	Stream        bool
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Content      string
	InputTokens  int
	OutputTokens int
	Model        string
}

// Complete sends a completion request to the Anthropic API
func (p *AnthropicProvider) Complete(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = p.maxTokens
	}

	// Build system prompt with context if provided
	systemPrompt := req.SystemPrompt
	if req.Context != "" {
		systemPrompt = fmt.Sprintf("%s\n\nRelevant context:\n%s", systemPrompt, req.Context)
	}

	apiReq := CompletionRequest{
		Model:     p.model,
		MaxTokens: maxTokens,
		System:    systemPrompt,
		Messages:  req.Messages,
		Stream:    req.Stream,
	}

	body, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", AnthropicAPIURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", AnthropicAPIVersion)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var apiResp CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract text from content blocks
	var content string
	for _, block := range apiResp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &ChatResponse{
		Content:      content,
		InputTokens:  apiResp.Usage.InputTokens,
		OutputTokens: apiResp.Usage.OutputTokens,
		Model:        apiResp.Model,
	}, nil
}

// StreamComplete sends a streaming completion request
func (p *AnthropicProvider) StreamComplete(ctx context.Context, req ChatRequest, callback func(chunk string) error) error {
	req.Stream = true
	// TODO: Implement SSE streaming
	// For now, fall back to non-streaming
	resp, err := p.Complete(ctx, req)
	if err != nil {
		return err
	}
	return callback(resp.Content)
}

// Name returns the provider name
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// Model returns the model being used
func (p *AnthropicProvider) Model() string {
	return p.model
}
