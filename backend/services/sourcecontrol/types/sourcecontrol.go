package types

import (
	"time"

	"gorm.io/datatypes"
)

type SourceControlAccount struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID         *string        `json:"userId,omitempty"`
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
	ProviderID     string
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
