package types

import (
	"time"

	"gorm.io/datatypes"
)

type SourceControlAccount struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	MemberID       *string        `json:"memberId,omitempty"`
	OrganizationID *string        `json:"organizationId,omitempty"`
	ProviderName   string         `json:"providerName"`
	ProviderID     string         `json:"providerId"`
	Username       string         `json:"username"`
	Metadata       datatypes.JSON `json:"metadata,omitempty"`
	LastSyncedAt   *time.Time     `json:"lastSyncedAt,omitempty"`
}

// PullRequest represents a pull request in our system
type PullRequest struct {
	ID                     string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SourceControlAccountID string         `json:"sourceControlAccountId"`
	ProviderID             string         `json:"providerId"`
	RepositoryName         string         `json:"repositoryName"`
	Title                  string         `json:"title"`
	Description            string         `json:"description"`
	URL                    string         `json:"url"`
	Status                 string         `json:"status"`
	CreatedAt              time.Time      `json:"createdAt"`
	UpdatedAt              time.Time      `json:"updatedAt"`
	MergedAt               *time.Time     `json:"mergedAt"`
	LastUpdatedAt          time.Time      `json:"lastUpdatedAt"`
	Comments               int            `json:"comments"`
	ReviewComments         int            `json:"reviewComments"`
	Additions              int            `json:"additions"`
	Deletions              int            `json:"deletions"`
	ChangedFiles           int            `json:"changedFiles"`
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
}

type PullRequestMetrics struct {
	MergedPRsCount         int
	OpenPRsCount           int
	MeanPublishToMergeTime float64
	MeanTimeToFirstReview  float64
}

// MemberActivity represents a single activity item in the timeline
type MemberActivity struct {
	ID               string         `json:"id"`
	Type             string         `json:"type"` // "pull_request", "pr_review", "pr_comment"
	Title            string         `json:"title"`
	Description      string         `json:"description,omitempty"`
	URL              string         `json:"url,omitempty"`
	Repository       string         `json:"repository,omitempty"`
	CreatedAt        time.Time      `json:"createdAt"`
	Metadata         datatypes.JSON `json:"metadata,omitempty"`
	AuthorUsername   string         `json:"authorUsername,omitempty"`
	PRTitle          string         `json:"prTitle,omitempty"`          // For comments/reviews: the PR title
	PRAuthorUsername string         `json:"prAuthorUsername,omitempty"` // For comments/reviews: the PR author
	PRMetrics        datatypes.JSON `json:"prMetrics,omitempty"`        // Added PRMetrics
}

// MemberActivityParams represents the parameters for getting member activity
type MemberActivityParams struct {
	MemberID  string     `json:"memberId"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
}

// MemberMetricsParams represents the parameters for getting member metrics
type MemberMetricsParams struct {
	MemberID  string     `json:"memberId"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	Interval  string     `json:"interval,omitempty"` // daily, weekly, monthly
}

// SnapshotMetric represents a single metric in the snapshot
type SnapshotMetric struct {
	Label      string  `json:"label"`
	Value      float64 `json:"value"`
	PeersValue float64 `json:"peersValue"`
	Unit       string  `json:"unit"` // "count", "time", "loc", etc.
}

// SnapshotCategory represents a category of metrics in the snapshot
type SnapshotCategory struct {
	Category string           `json:"category"`
	Metrics  []SnapshotMetric `json:"metrics"`
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

// GraphMetric represents a single metric in the graph data
type GraphMetric struct {
	Label      string            `json:"label"`
	Type       string            `json:"type"`
	TimeSeries []TimeSeriesEntry `json:"timeSeries"`
}

// GraphCategory represents a category of metrics in the graph data
type GraphCategory struct {
	Category string        `json:"category"`
	Metrics  []GraphMetric `json:"metrics"`
}

// MemberMetricsResponse represents the response for getting member metrics
type MemberMetricsResponse struct {
	SnapshotMetrics []SnapshotCategory `json:"snapshotMetrics"`
	GraphMetrics    []GraphCategory    `json:"graphMetrics"`
}
