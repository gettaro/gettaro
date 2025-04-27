package types

import (
	"time"
)

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID                 int        `json:"id"`
	Number             int        `json:"number"`
	State              string     `json:"state"`
	Title              string     `json:"title"`
	Body               string     `json:"body"`
	URL                string     `json:"html_url"`
	CommitsURL         string     `json:"commits_url"`
	ReviewCommentsURL  string     `json:"review_comments_url"`
	StatusesURL        string     `json:"statuses_url"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ClosedAt           *time.Time `json:"closed_at"`
	MergedAt           *time.Time `json:"merged_at"`
	MergeCommitSHA     string     `json:"merge_commit_sha"`
	Assignee           *User      `json:"assignee"`
	Assignees          []User     `json:"assignees"`
	RequestedReviewers []User     `json:"requested_reviewers"`
	Draft              bool       `json:"draft"`
	Merged             bool       `json:"merged"`
	MergedBy           *User      `json:"merged_by"`
	Comments           int        `json:"comments"`
	ReviewComments     int        `json:"review_comments"`
	Commits            int        `json:"commits"`
	Additions          int        `json:"additions"`
	Deletions          int        `json:"deletions"`
	ChangedFiles       int        `json:"changed_files"`
	User               User       `json:"user"`
}

// ReviewComment represents a GitHub pull request review comment
type ReviewComment struct {
	ID        int       `json:"id"`
	NodeID    string    `json:"node_id"`
	User      User      `json:"user"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	NodeID    string `json:"node_id"`
	AvatarURL string `json:"avatar_url"`
}
