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
	ReviewComments int        `json:"review_comments"`  // Number of review comments
	Additions      int        `json:"additions"`        // Additions
	Deletions      int        `json:"deletions"`        // Deletions
	ChangedFiles   int        `json:"changed_files"`    // Changed files
	User           User       `json:"user"`             // Author
	Number         int        `json:"number"`           // PR number
	ClosedAt       *time.Time `json:"closed_at"`        // When PR was closed
	MergeCommitSHA string     `json:"merge_commit_sha"` // SHA of the merge commit
	Draft          bool       `json:"draft"`            // Whether PR is a draft
	Commits        int        `json:"commits"`          // Number of commits
	Head           Ref        `json:"head"`             // The head ref
	Base           Ref        `json:"base"`             // The base ref
	Links          Links      `json:"_links"`           // Hypermedia links
}

// Commnent types
type CommentType string

const (
	CommentTypeComment       CommentType = "COMMENT"
	CommentTypeReview        CommentType = "REVIEW"
	CommentTypeReviewComment CommentType = "REVIEW_COMMENT"
)

// ReviewComment represents a GitHub pull request review comment
type ReviewComment struct {
	ID        int       `json:"id"`
	User      User      `json:"user"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Review represents a GitHub pull request review
type Review struct {
	ID          int       `json:"id"`
	User        User      `json:"user"`
	Body        string    `json:"body"`
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
}

// Ref represents a Git reference
type Ref struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	Sha   string `json:"sha"`
	User  User   `json:"user"`
	Repo  Repo   `json:"repo"`
}

// Repo represents a GitHub repository
type Repo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}

// Links represents hypermedia links for a pull request
type Links struct {
	Self           Link `json:"self"`
	HTML           Link `json:"html"`
	Issue          Link `json:"issue"`
	Comments       Link `json:"comments"`
	ReviewComments Link `json:"review_comments"`
	Commits        Link `json:"commits"`
	Statuses       Link `json:"statuses"`
}

// Link represents a hypermedia link
type Link struct {
	Href string `json:"href"`
}

// Commit represents a GitHub commit
type Commit struct {
	Sha    string       `json:"sha"`
	Commit CommitDetail `json:"commit"`
}

// CommitDetail represents the commit details
type CommitDetail struct {
	Message string `json:"message"`
	Author  struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"author"`
	Committer struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"committer"`
}
