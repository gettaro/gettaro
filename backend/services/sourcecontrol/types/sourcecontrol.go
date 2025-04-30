package types

import (
	"time"

	"gorm.io/datatypes"
)

type SourceControlAccount struct {
	ID             string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         *string
	OrganizationID *string
	ProviderName   string
	ProviderID     string
	Username       string
	Metadata       datatypes.JSON
	LastSyncedAt   *time.Time
}

// PullRequest represents a pull request in our system
type PullRequest struct {
	ID                     string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SourceControlAccountID string
	ProviderID             string
	RepositoryName         string
	OrganizationID         string
	Title                  string
	Description            string
	URL                    string
	Status                 string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	MergedAt               *time.Time
	LastUpdatedAt          time.Time
	Comments               int
	ReviewComments         int
	Additions              int
	Deletions              int
	ChangedFiles           int
	Metadata               datatypes.JSON
}

// PRComment represents a comment on a pull request
type PRComment struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PRID       string
	AuthorID   string
	ProviderID string
	Body       string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

// PullRequestParams represents the parameters for querying pull requests
type PullRequestParams struct {
	ProviderID     string
	OrganizationID *string
	ProviderName   string
	RepositoryName string
}
