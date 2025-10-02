package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// StringArray is a custom type for handling PostgreSQL string arrays
type StringArray []string

// Scan implements the sql.Scanner interface for database deserialization
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// Handle PostgreSQL array format: {"item1","item2","item3"}
		str := string(v)
		if str == "{}" {
			*s = StringArray{}
			return nil
		}
		// Remove curly braces and split by comma
		str = strings.Trim(str, "{}")
		if str == "" {
			*s = StringArray{}
			return nil
		}
		items := strings.Split(str, ",")
		*s = make(StringArray, len(items))
		for i, item := range items {
			// Remove quotes from each item
			(*s)[i] = strings.Trim(item, `"`)
		}
		return nil
	case string:
		// Handle string input
		if v == "{}" {
			*s = StringArray{}
			return nil
		}
		v = strings.Trim(v, "{}")
		if v == "" {
			*s = StringArray{}
			return nil
		}
		items := strings.Split(v, ",")
		*s = make(StringArray, len(items))
		for i, item := range items {
			(*s)[i] = strings.Trim(item, `"`)
		}
		return nil
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}
}

// Value implements the driver.Valuer interface for database serialization
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}

	// Format as PostgreSQL array: {"item1","item2","item3"}
	items := make([]string, len(s))
	for i, item := range s {
		items[i] = fmt.Sprintf(`"%s"`, item)
	}
	return fmt.Sprintf("{%s}", strings.Join(items, ",")), nil
}

// MarshalJSON implements json.Marshaler
func (s StringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(s))
}

// UnmarshalJSON implements json.Unmarshaler
func (s *StringArray) UnmarshalJSON(data []byte) error {
	var items []string
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	*s = StringArray(items)
	return nil
}

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
	ID             string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationID string      `json:"organization_id"`
	UserID         string      `json:"user_id"`
	EntityType     string      `json:"entity_type"`
	EntityID       string      `json:"entity_id"`
	Query          string      `json:"query"`
	Answer         string      `json:"answer"`
	Context        string      `json:"context"`
	Confidence     float64     `json:"confidence"`
	Sources        StringArray `gorm:"type:text[]" json:"sources"`
	CreatedAt      time.Time   `json:"created_at"`
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

// RetrievalPlan represents a structured plan for data retrieval
type RetrievalPlan struct {
	DataSources []string   `json:"data_sources"` // ["source_control", "conversations", "member_data", "team_data"]
	Time        *TimeRange `json:"time,omitempty"`
	Filters     *Filters   `json:"filters,omitempty"`
	Priority    string     `json:"priority,omitempty"`  // "high", "medium", "low"
	Reasoning   string     `json:"reasoning,omitempty"` // Why this plan was chosen
}

// TimeRange represents a time range for data retrieval
type TimeRange struct {
	From     string `json:"from"`     // ISO date format: "2025-09-01"
	To       string `json:"to"`       // ISO date format: "2025-09-27"
	Interval string `json:"interval"` // "day", "week", "month", "quarter", "year"
}

// Filters represents additional filters for data retrieval
type Filters struct {
	MemberIDs     []string `json:"member_ids,omitempty"`
	TeamIDs       []string `json:"team_ids,omitempty"`
	Statuses      []string `json:"statuses,omitempty"`
	MinConfidence float64  `json:"min_confidence,omitempty"`
	Limit         int      `json:"limit,omitempty"`
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
