package types

import "fmt"

// MessageRequest represents a request to the Anthropic Messages API
type MessageRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Messages    []Message `json:"messages"`
	System      string    `json:"system,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MessageResponse represents the response from the Anthropic Messages API
type MessageResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence *string        `json:"stop_sequence"`
	Usage        Usage          `json:"usage"`
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

// APIError represents an error response from the Anthropic API
type APIError struct {
	Type    string      `json:"type"`
	Details ErrorDetail `json:"error"`
}

// ErrorDetail represents error details
type ErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("anthropic API error [%s]: %s", e.Details.Type, e.Details.Message)
}

// GetTextContent extracts the text content from a MessageResponse
func (resp *MessageResponse) GetTextContent() string {
	if len(resp.Content) == 0 {
		return ""
	}

	var textContent string
	for _, block := range resp.Content {
		if block.Type == "text" {
			textContent += block.Text
		}
	}

	return textContent
}

