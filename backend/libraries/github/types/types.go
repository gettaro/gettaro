package types

import (
	"time"
)

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID             int        `json:"id"`               // ProviderID
	State          string     `json:"state"`            // Status
	Title          string     `json:"title"`            // Title
	Body           string     `json:"body"`             // Description
	URL            string     `json:"html_url"`         // URL
	CreatedAt      time.Time  `json:"created_at"`       // Created at
	UpdatedAt      time.Time  `json:"updated_at"`       // Updated at
	MergedAt       *time.Time `json:"merged_at"`        // Merged at
	Comments       int        `json:"comments"`         // Number of comments
	ReviewComments int        `json:"review_comments"`  // Number of comments (need to check)
	Additions      int        `json:"additions"`        // Additions
	Deletions      int        `json:"deletions"`        // Deletions
	ChangedFiles   int        `json:"changed_files"`    // Changed files
	User           User       `json:"user"`             // Author
	Number         int        `json:"number"`           // Metadata
	ClosedAt       *time.Time `json:"closed_at"`        // Metadata
	MergeCommitSHA string     `json:"merge_commit_sha"` //Metadata
	Draft          bool       `json:"draft"`            // Metadata
	Commits        int        `json:"commits"`          // Metadata
}

// ReviewComment represents a GitHub pull request review comment
type ReviewComment struct {
	ID        int       `json:"id"`
	User      User      `json:"user"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
}
