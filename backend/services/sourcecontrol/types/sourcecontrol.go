package types

import (
	"time"

	"gorm.io/datatypes"
)

type SourceControlAccount struct {
	ID             string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string
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
	URL                    string
	Title                  string
	Status                 string
	CreatedAt              time.Time
	MergedAt               *time.Time
	LastUpdatedAt          time.Time
}

// PRComment represents a comment on a pull request
type PRComment struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PRID      string
	AuthorID  *string
	Body      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type PRReviewer struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PRID        string
	ReviewerID  string
	ReviewedAt  time.Time
	ReviewState string
}
