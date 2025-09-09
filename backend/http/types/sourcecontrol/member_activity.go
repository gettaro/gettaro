package sourcecontrol

import (
	"time"

	"gorm.io/datatypes"
)

// GetMemberActivityRequest represents the request parameters for getting member activity
type GetMemberActivityRequest struct {
	StartDate string `form:"startDate" binding:"omitempty"`
	EndDate   string `form:"endDate" binding:"omitempty"`
}

// MemberActivity represents a single activity item in the timeline
type MemberActivity struct {
	ID               string                 `json:"id"`
	Type             string                 `json:"type"` // "pull_request", "pr_review", "pr_comment"
	Title            string                 `json:"title"`
	Description      string                 `json:"description,omitempty"`
	URL              string                 `json:"url,omitempty"`
	Repository       string                 `json:"repository,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	AuthorUsername   string                 `json:"author_username,omitempty"`
	PRTitle          string                 `json:"pr_title,omitempty"`           // For comments/reviews: the PR title
	PRAuthorUsername string                 `json:"pr_author_username,omitempty"` // For comments/reviews: the PR author
	PRMetrics        datatypes.JSON         `json:"pr_metrics,omitempty"`
}

// GetMemberActivityResponse represents the response for getting member activity
type GetMemberActivityResponse struct {
	Activities []MemberActivity `json:"activities"`
}
