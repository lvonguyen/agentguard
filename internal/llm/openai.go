package llm

import (
	"context"
	"fmt"
)

// OpenAIConfig holds configuration for the OpenAI provider
type OpenAIConfig struct {
	APIKey       string
	Model        string
	MaxTokens    int
	Organization string
	BaseURL      string // For Azure OpenAI or compatible APIs
}

// OpenAIProvider implements the LLM Provider interface for OpenAI
type OpenAIProvider struct {
	config OpenAIConfig
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(cfg OpenAIConfig) (*OpenAIProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("openai API key is required")
	}

	if cfg.Model == "" {
		cfg.Model = "gpt-4-turbo"
	}

	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 4096
	}

	return &OpenAIProvider{
		config: cfg,
	}, nil
}

// Complete sends a chat completion request to OpenAI
func (p *OpenAIProvider) Complete(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// TODO: Implement OpenAI completion
	// Use openai-go SDK or direct HTTP calls
	return nil, fmt.Errorf("OpenAI provider not yet implemented")
}

// StreamComplete sends a streaming chat completion request
func (p *OpenAIProvider) StreamComplete(ctx context.Context, req ChatRequest, callback func(chunk string) error) error {
	// TODO: Implement OpenAI streaming
	return fmt.Errorf("OpenAI streaming not yet implemented")
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Model returns the model being used
func (p *OpenAIProvider) Model() string {
	return p.config.Model
}
