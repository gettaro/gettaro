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
	LinesOfCodeSuggested int            `gorm:"default:0" json:"lines_of_code_suggested"`
	SuggestionAcceptRate *float64       `gorm:"type:decimal(5,2)" json:"suggestion_accept_rate,omitempty"`
	ActiveSessions       int            `gorm:"default:0" json:"active_sessions"`
	Metadata             datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
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

// MemberMetricsParams represents the parameters for getting member metrics
type MemberMetricsParams struct {
	MemberID  string     `json:"member_id"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Interval  string     `json:"interval,omitempty"` // daily, weekly, monthly
}

// AICodeAssistantUsageStats represents aggregated statistics
type AICodeAssistantUsageStats struct {
	TotalLinesAccepted  int     `json:"total_lines_accepted"`
	TotalLinesSuggested int     `json:"total_lines_suggested"`
	OverallAcceptRate   float64 `json:"overall_accept_rate"`
	ActiveSessions      int     `json:"active_sessions"`
}

// MetricRuleParams represents the parameters for a metric rule
type MetricRuleParams struct {
	MetricParams datatypes.JSON `json:"metric_params"`
	StartDate    *time.Time     `json:"start_date,omitempty"`
	EndDate      *time.Time     `json:"end_date,omitempty"`
	Interval     string         `json:"interval,omitempty"` // daily, weekly, monthly
}

// SnapshotMetric represents a single metric in the snapshot
type SnapshotMetric struct {
	Label          string  `json:"label"`
	Description    string  `json:"description"`
	Category       string  `json:"category"`
	Value          float64 `json:"value"`
	PeersValue     float64 `json:"peers_value"`
	Unit           Unit    `json:"unit"` // "count", "time", "loc", etc.
	IconIdentifier string  `json:"icon_identifier"`
	IconColor      string  `json:"icon_color"`
}

// SnapshotCategory represents a category of metrics in the snapshot
type SnapshotCategory struct {
	Category MetricRuleCategory `json:"category"`
	Metrics  []SnapshotMetric   `json:"metrics"`
}

// TimeSeriesDataPoint represents a single data point in a time series
type TimeSeriesDataPoint struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

// TimeSeriesEntry represents a single time entry in a time series
type TimeSeriesEntry struct {
	Date string                `json:"date"`
	Data []TimeSeriesDataPoint `json:"data"`
}

// Unit represents a unit of measurement
type Unit string

const (
	UnitCount   Unit = "count" // nit: This is not a unit of measurement and probably should be renamed
	UnitSeconds Unit = "seconds"
	UnitPercent Unit = "percent"
)

// GraphMetric represents a single metric in the graph data
type GraphMetric struct {
	Label      string            `json:"label"`
	Type       string            `json:"type"`
	Category   string            `json:"category"`
	Unit       Unit              `json:"unit"`
	TimeSeries []TimeSeriesEntry `json:"time_series"`
}

// GraphCategory represents a category of metrics in the graph data
type GraphCategory struct {
	Category MetricRuleCategory `json:"category"`
	Metrics  []GraphMetric      `json:"metrics"`
}

// MetricsResponse represents the response for getting member metrics
type MetricsResponse struct {
	SnapshotMetrics []*SnapshotCategory `json:"snapshot_metrics"`
	GraphMetrics    []*GraphCategory    `json:"graph_metrics"`
}

// MetricRuleCategory represents a category of metric rules
type MetricRuleCategory struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}
