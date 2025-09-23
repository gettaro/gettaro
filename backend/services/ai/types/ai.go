package types

import (
	"time"
)

// AIQueryRequest represents a request to query the AI assistant
type AIQueryRequest struct {
	EntityType     string                 `json:"entity_type"`     // "member", "team", "organization", "project"
	EntityID       string                 `json:"entity_id"`       // ID of the specific entity
	OrganizationID string                 `json:"organization_id"` // Always required for context
	Query          string                 `json:"query"`
	Context        string                 `json:"context,omitempty"`         // "performance", "conversations", "overview"
	AdditionalData map[string]interface{} `json:"additional_data,omitempty"` // For extra context
}

// AIQueryResponse represents the response from the AI assistant
type AIQueryResponse struct {
	Answer      string                 `json:"answer"`
	Sources     []string               `json:"sources"`
	Confidence  float64                `json:"confidence"`
	RelatedData map[string]interface{} `json:"related_data,omitempty"`
	Suggestions []string               `json:"suggestions,omitempty"` // Follow-up questions
}

// AIQueryHistory represents a stored query for history and analytics
type AIQueryHistory struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationID string    `json:"organization_id"`
	UserID         string    `json:"user_id"`
	EntityType     string    `json:"entity_type"`
	EntityID       string    `json:"entity_id"`
	Query          string    `json:"query"`
	Answer         string    `json:"answer"`
	Context        string    `json:"context"`
	Confidence     float64   `json:"confidence"`
	Sources        []string  `json:"sources"`
	CreatedAt      time.Time `json:"created_at"`
}

// AIQueryStats represents statistics about AI usage
type AIQueryStats struct {
	TotalQueries      int              `json:"total_queries"`
	QueriesByEntity   map[string]int   `json:"queries_by_entity"`
	QueriesByContext  map[string]int   `json:"queries_by_context"`
	AverageConfidence float64          `json:"average_confidence"`
	RecentQueries     []AIQueryHistory `json:"recent_queries"`
}

// EntityData represents aggregated data for an entity
type EntityData struct {
	EntityType     string                 `json:"entity_type"`
	EntityID       string                 `json:"entity_id"`
	OrganizationID string                 `json:"organization_id"`
	Data           map[string]interface{} `json:"data"`
	LastUpdated    time.Time              `json:"last_updated"`
}

// AIServiceConfig represents configuration for the AI service
type AIServiceConfig struct {
	Provider          string  `json:"provider"` // "openai", "anthropic", "local"
	Model             string  `json:"model"`    // "gpt-4", "claude-3", etc.
	MaxTokens         int     `json:"max_tokens"`
	Temperature       float64 `json:"temperature"`
	MaxContextSize    int     `json:"max_context_size"`
	EnableHistory     bool    `json:"enable_history"`
	EnableSuggestions bool    `json:"enable_suggestions"`
}
