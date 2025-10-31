package types

import (
	"time"

	"gorm.io/datatypes"
)

type SourceControlAccount struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	MemberID       *string        `json:"member_id,omitempty"`
	OrganizationID *string        `json:"organization_id,omitempty"`
	ProviderName   string         `json:"provider_name"`
	ProviderID     string         `json:"provider_id"`
	Username       string         `json:"username"`
	Metadata       datatypes.JSON `json:"metadata,omitempty"`
	LastSyncedAt   *time.Time     `json:"last_synced_at,omitempty"`
}

// PullRequest represents a pull request in our system
type PullRequest struct {
	ID                     string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SourceControlAccountID string         `json:"source_control_account_id"`
	ProviderID             string         `json:"provider_id"`
	RepositoryName         string         `json:"repository_name"`
	Title                  string         `json:"title"`
	Description            string         `json:"description"`
	URL                    string         `json:"url"`
	Status                 string         `json:"status"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	MergedAt               *time.Time     `json:"merged_at"`
	LastUpdatedAt          time.Time      `json:"last_updated_at"`
	Comments               int            `json:"comments"`
	ReviewComments         int            `json:"review_comments"`
	Additions              int            `json:"additions"`
	Deletions              int            `json:"deletions"`
	ChangedFiles           int            `json:"changed_files"`
	Metrics                datatypes.JSON `json:"metrics"`
	Metadata               datatypes.JSON `json:"metadata"`
}

// PRComment represents a comment on a pull request
type PRComment struct {
	ID                     string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PRID                   string
	SourceControlAccountID string
	ProviderID             string
	Body                   string
	Type                   string
	CreatedAt              time.Time
	UpdatedAt              *time.Time
}

// PullRequestParams represents the parameters for querying pull requests
type PullRequestParams struct {
	ProviderIDs    []string
	OrganizationID *string
	RepositoryName string
	UserIDs        []string
	StartDate      *time.Time
	EndDate        *time.Time
	Status         string // "open", "closed", "merged"
}

type PullRequestMetrics struct {
	MergedPRsCount         int
	OpenPRsCount           int
	MeanPublishToMergeTime float64
	MeanTimeToFirstReview  float64
}

type SourceControlAccountParams struct {
	SourceControlAccountIDs []string `json:"source_control_account_ids"`
	OrganizationID          string   `json:"organization_id"`
	Usernames               []string `json:"usernames"`
	MemberIDs               []string `json:"member_ids"`
}

// MemberActivity represents a single activity item in the timeline
type MemberActivity struct {
	ID               string         `json:"id"`
	Type             string         `json:"type"` // "pull_request", "pr_review", "pr_comment"
	Title            string         `json:"title"`
	Description      string         `json:"description,omitempty"`
	URL              string         `json:"url,omitempty"`
	Repository       string         `json:"repository,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	Metadata         datatypes.JSON `json:"metadata,omitempty"`
	AuthorUsername   string         `json:"author_username,omitempty"`
	PRTitle          string         `json:"pr_title,omitempty"`           // For comments/reviews: the PR title
	PRAuthorUsername string         `json:"pr_author_username,omitempty"` // For comments/reviews: the PR author
	PRMetrics        datatypes.JSON `json:"pr_metrics,omitempty"`         // Added PRMetrics
}

// PullRequestWithComments represents a pull request with its comments
type PullRequestWithComments struct {
	*PullRequest
	Comments []*PRComment `json:"comments,omitempty"`
}

// MemberPullRequestParams represents the parameters for getting member pull requests
type MemberPullRequestParams struct {
	MemberID        string     `json:"member_id"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	IncludeComments *bool      `json:"include_comments,omitempty"` // Optional include PR comments with body
}

// MemberPullRequestReviewsParams represents the parameters for getting member pull request reviews
type MemberPullRequestReviewsParams struct {
	MemberID  string     `json:"member_id"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	HasBody   *bool      `json:"has_body,omitempty"` // Optional filter for reviews with body content
}

// MemberMetricsParams represents the parameters for getting member metrics
type MemberMetricsParams struct {
	MemberID  string     `json:"member_id"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Interval  string     `json:"interval,omitempty"` // daily, weekly, monthly
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
