package ai

// AIQueryRequest represents the HTTP request for AI queries
type AIQueryRequest struct {
	EntityType     string                 `json:"entity_type" binding:"required"` // "member", "team", "organization", "project"
	EntityID       string                 `json:"entity_id" binding:"required"`   // ID of the specific entity
	Query          string                 `json:"query" binding:"required"`
	Context        string                 `json:"context,omitempty"`         // "performance", "conversations", "overview"
	AdditionalData map[string]interface{} `json:"additional_data,omitempty"` // For extra context
}

// AIQueryResponse represents the HTTP response for AI queries
type AIQueryResponse struct {
	Answer      string                 `json:"answer"`
	Sources     []string               `json:"sources"`
	Confidence  float64                `json:"confidence"`
	RelatedData map[string]interface{} `json:"related_data,omitempty"`
	Suggestions []string               `json:"suggestions,omitempty"` // Follow-up questions
}

// AIQueryHistoryResponse represents the HTTP response for query history
type AIQueryHistoryResponse struct {
	Queries []AIQueryHistoryItem `json:"queries"`
}

// AIQueryHistoryItem represents a single query history item
type AIQueryHistoryItem struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	UserID         string   `json:"user_id"`
	EntityType     string   `json:"entity_type"`
	EntityID       string   `json:"entity_id"`
	Query          string   `json:"query"`
	Answer         string   `json:"answer"`
	Context        string   `json:"context"`
	Confidence     float64  `json:"confidence"`
	Sources        []string `json:"sources"`
	CreatedAt      string   `json:"created_at"`
}

// AIQueryStatsResponse represents the HTTP response for query statistics
type AIQueryStatsResponse struct {
	TotalQueries      int                  `json:"total_queries"`
	QueriesByEntity   map[string]int       `json:"queries_by_entity"`
	QueriesByContext  map[string]int       `json:"queries_by_context"`
	AverageConfidence float64              `json:"average_confidence"`
	RecentQueries     []AIQueryHistoryItem `json:"recent_queries"`
}
