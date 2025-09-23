package providers

import (
	"context"

	"ems.dev/backend/services/ai/types"
)

// AIProviderInterface defines the interface for AI providers
type AIProviderInterface interface {
	// Query sends a query to the AI provider and returns a response
	Query(ctx context.Context, prompt string, config *types.AIServiceConfig) (*types.AIQueryResponse, error)

	// GetProviderName returns the name of the provider
	GetProviderName() string

	// IsAvailable checks if the provider is available and configured
	IsAvailable() bool
}

