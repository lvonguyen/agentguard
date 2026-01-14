package llm

import "context"

// Provider defines the interface for LLM providers
type Provider interface {
	// Complete sends a chat completion request and returns the response
	Complete(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// StreamComplete sends a streaming chat completion request
	StreamComplete(ctx context.Context, req ChatRequest, callback func(chunk string) error) error

	// Name returns the provider name
	Name() string

	// Model returns the model being used
	Model() string
}
