package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"ems.dev/backend/services/ai/types"
)

// Provider implements the AIProviderInterface for Groq
type Provider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewProvider creates a new Groq AI provider
func NewProvider(apiKey string) *Provider {
	return &Provider{
		apiKey:  apiKey,
		baseURL: "https://api.groq.com/openai/v1",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Query sends a query to Groq and returns a response
func (p *Provider) Query(ctx context.Context, prompt string, config *types.AIServiceConfig) (*types.AIQueryResponse, error) {
	// Create the chat completion request
	req := &GroqChatRequest{
		Model:       p.getModel(config.Model),
		Messages:    []GroqMessage{{Role: "user", Content: prompt}},
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		Stream:      false,
	}

	// Send the request to Groq
	resp, err := p.sendRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Groq API: %w", err)
	}

	// Extract the text content from the response
	answer := resp.GetTextContent()
	if answer == "" {
		return nil, fmt.Errorf("received empty response from Groq API")
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
	return "groq"
}

// IsAvailable checks if the provider is available and configured
func (p *Provider) IsAvailable() bool {
	return p.apiKey != "" && p.client != nil
}

// sendRequest sends a request to the Groq API
func (p *Provider) sendRequest(ctx context.Context, req *GroqChatRequest) (*GroqChatResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("groq API error %d: %s", resp.StatusCode, string(body))
	}

	var groqResp GroqChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &groqResp, nil
}

// getModel maps the generic model name to Groq-specific model names
func (p *Provider) getModel(model string) string {
	// Map common model names to Groq models
	switch strings.ToLower(model) {
	case "gpt-4", "gpt-4-turbo":
		return "llama-3.1-70b-versatile"
	case "gpt-3.5-turbo":
		return "llama-3.1-8b-instant"
	case "claude-3", "claude-3-sonnet":
		return "llama-3.1-70b-versatile"
	case "claude-3-haiku":
		return "llama-3.1-8b-instant"
	default:
		// Default to a fast model
		return "llama-3.1-8b-instant"
	}
}

// calculateConfidence calculates a confidence score based on response characteristics
func (p *Provider) calculateConfidence(resp *GroqChatResponse) float64 {
	// Base confidence
	confidence := 0.8

	// Adjust based on response length
	textLength := len(resp.GetTextContent())
	if textLength > 500 {
		confidence += 0.1
	} else if textLength < 100 {
		confidence -= 0.1
	}

	// Adjust based on token usage efficiency
	if len(resp.Choices) > 0 && resp.Usage.CompletionTokens > 0 && resp.Usage.PromptTokens > 0 {
		efficiency := float64(resp.Usage.CompletionTokens) / float64(resp.Usage.PromptTokens)
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

// Groq API Types

type GroqChatRequest struct {
	Model       string        `json:"model"`
	Messages    []GroqMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// GetTextContent extracts the text content from the response
func (r *GroqChatResponse) GetTextContent() string {
	if len(r.Choices) > 0 && len(r.Choices[0].Message.Content) > 0 {
		return r.Choices[0].Message.Content
	}
	return ""
}

// GetOutputTokens returns the number of output tokens
func (r *GroqChatResponse) GetOutputTokens() int {
	return r.Usage.CompletionTokens
}

// GetInputTokens returns the number of input tokens
func (r *GroqChatResponse) GetInputTokens() int {
	return r.Usage.PromptTokens
}
