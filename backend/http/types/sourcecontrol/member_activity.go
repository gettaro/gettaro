package sourcecontrol

import "time"

// GetMemberActivityRequest represents the request parameters for getting member activity
type GetMemberActivityRequest struct {
	StartDate *time.Time `form:"startDate" binding:"omitempty"`
	EndDate   *time.Time `form:"endDate" binding:"omitempty"`
}

// MemberActivity represents a single activity item in the timeline
type MemberActivity struct {
	ID             string                 `json:"id"`
	Type           string                 `json:"type"` // "pull_request", "pr_review", "pr_comment"
	Title          string                 `json:"title"`
	Description    string                 `json:"description,omitempty"`
	URL            string                 `json:"url,omitempty"`
	Repository     string                 `json:"repository,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	AuthorUsername string                 `json:"authorUsername,omitempty"`
}

// GetMemberActivityResponse represents the response for getting member activity
type GetMemberActivityResponse struct {
	Activities []MemberActivity `json:"activities"`
}
