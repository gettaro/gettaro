package anthropic

import (
	"context"
	"fmt"
	"strings"

	"ems.dev/backend/libraries/anthropic"
	anthropictypes "ems.dev/backend/libraries/anthropic/types"
	"ems.dev/backend/services/ai/types"
)

// Provider implements the AIProviderInterface for Anthropic
type Provider struct {
	client *anthropic.Client
}

// NewProvider creates a new Anthropic AI provider
func NewProvider(apiKey string) *Provider {
	return &Provider{
		client: anthropic.NewClient(apiKey),
	}
}

// Query sends a query to Anthropic and returns a response
func (p *Provider) Query(ctx context.Context, prompt string, config *types.AIServiceConfig) (*types.AIQueryResponse, error) {
	// Create the message request
	req := &anthropictypes.MessageRequest{
		Model:       config.Model,
		MaxTokens:   config.MaxTokens,
		Temperature: config.Temperature,
		Messages: []anthropictypes.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Send the request to Anthropic
	resp, err := p.client.SendMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Anthropic API: %w", err)
	}

	// Extract the text content from the response
	answer := resp.GetTextContent()
	if answer == "" {
		return nil, fmt.Errorf("received empty response from Anthropic API")
	}

	// Calculate confidence based on response characteristics
	confidence := p.calculateConfidence(resp)

	// Determine sources based on the prompt content
	sources := p.determineSources(prompt)

	response := &types.AIQueryResponse{
		Answer:     answer,
		Sources:    sources,
		Confidence: confidence,
	}

	return response, nil
}

// GetProviderName returns the name of the provider
func (p *Provider) GetProviderName() string {
	return "anthropic"
}

// IsAvailable checks if the provider is available and configured
func (p *Provider) IsAvailable() bool {
	return p.client != nil
}

// calculateConfidence calculates a confidence score based on response characteristics
func (p *Provider) calculateConfidence(resp *anthropictypes.MessageResponse) float64 {
	// Base confidence
	confidence := 0.8

	// Adjust based on response length (longer responses might be more confident)
	textLength := len(resp.GetTextContent())
	if textLength > 500 {
		confidence += 0.1
	} else if textLength < 100 {
		confidence -= 0.1
	}

	// Adjust based on token usage efficiency
	if resp.Usage.OutputTokens > 0 && resp.Usage.InputTokens > 0 {
		efficiency := float64(resp.Usage.OutputTokens) / float64(resp.Usage.InputTokens)
		if efficiency > 0.5 {
			confidence += 0.05
		}
	}

	// Ensure confidence is within bounds
	if confidence > 1.0 {
		confidence = 1.0
	} else if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// determineSources determines the data sources used based on prompt content
func (p *Provider) determineSources(prompt string) []string {
	sources := []string{}
	promptLower := strings.ToLower(prompt)

	// Check for different data types in the prompt
	if strings.Contains(promptLower, "conversation") {
		sources = append(sources, "conversations")
	}
	if strings.Contains(promptLower, "performance") || strings.Contains(promptLower, "metric") {
		sources = append(sources, "source_control")
	}
	if strings.Contains(promptLower, "member") || strings.Contains(promptLower, "team") {
		sources = append(sources, "member_data")
	}

	// Default sources if none detected
	if len(sources) == 0 {
		sources = []string{"member_data", "conversations", "source_control"}
	}

	return sources
}
