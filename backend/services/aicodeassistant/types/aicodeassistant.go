package types

import (
	"time"

	"gorm.io/datatypes"
)

// AICodeAssistantDailyMetric represents aggregated daily metrics at user level
// (overlapping metrics from both Cursor Analytics API and Claude Code Usage Analytics)
type AICodeAssistantDailyMetric struct {
	ID                   string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationID       string         `gorm:"type:uuid;not null" json:"organization_id"`
	ExternalAccountID    string         `gorm:"type:uuid;not null" json:"external_account_id"`
	ToolName             string         `gorm:"type:varchar(255);not null" json:"tool_name"`
	MetricDate           time.Time      `gorm:"type:date;not null" json:"metric_date"`
	LinesOfCodeAccepted  int            `gorm:"default:0" json:"lines_of_code_accepted"`
	TotalSuggestions     int            `gorm:"default:0" json:"total_suggestions"`
	SuggestionsAccepted  int            `gorm:"default:0" json:"suggestions_accepted"`
	SuggestionAcceptRate *float64       `gorm:"type:decimal(5,2)" json:"suggestion_accept_rate,omitempty"`
	ActiveSessions       int            `gorm:"default:0" json:"active_sessions"`
	Metadata             datatypes.JSON  `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt            time.Time      `gorm:"type:timestamp with time zone;default:current_timestamp" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"type:timestamp with time zone;default:current_timestamp" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (AICodeAssistantDailyMetric) TableName() string {
	return "ai_code_assistant_daily_metrics"
}

// AICodeAssistantDailyMetricParams for querying daily metrics (user-level only)
type AICodeAssistantDailyMetricParams struct {
	OrganizationID     string     `json:"organization_id"`
	ExternalAccountIDs []string   `json:"external_account_ids,omitempty"` // Filter by specific users (empty = all users in org)
	ToolName           *string    `json:"tool_name,omitempty"`            // Filter by tool (e.g., "cursor", "claude-code")
	StartDate          *time.Time `json:"start_date,omitempty"`
	EndDate            *time.Time `json:"end_date,omitempty"`
}

// AICodeAssistantMemberMetricsParams for querying member-specific metrics
// This is used by the service layer to build proper params from member ID
type AICodeAssistantMemberMetricsParams struct {
	ToolName  *string    `json:"tool_name,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// AICodeAssistantUsageStats represents aggregated statistics
type AICodeAssistantUsageStats struct {
	TotalLinesAccepted int     `json:"total_lines_accepted"`
	TotalSuggestions   int     `json:"total_suggestions"`
	OverallAcceptRate  float64 `json:"overall_accept_rate"`
	ActiveUsers        int     `json:"active_users"`
}
