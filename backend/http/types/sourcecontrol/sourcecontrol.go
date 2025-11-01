package sourcecontrol

import (
	"time"

	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
)

type ListOrganizationPullRequestsQuery struct {
	UserIDs        []string `form:"userIds" binding:"omitempty"`
	RepositoryName string   `form:"repositoryName" binding:"omitempty"`
	Prefix         string   `form:"prefix" binding:"omitempty"`
	StartDate      string   `form:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate        string   `form:"endDate" binding:"omitempty,datetime=2006-01-02"`
	Status         string   `form:"status" binding:"omitempty,oneof=open closed merged"`
}

type ListOrganizationPullRequestsResponse struct {
	PullRequests []PullRequest `json:"pull_requests"`
}

type PullRequestAuthor struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type PullRequest struct {
	ID             string              `json:"id"`
	Title          string              `json:"title"`
	Description    string              `json:"description"`
	URL            string              `json:"url"`
	Status         string              `json:"status"`
	RepositoryName string              `json:"repository_name"`
	CreatedAt      time.Time           `json:"created_at"`
	MergedAt       *time.Time          `json:"merged_at,omitempty"`
	Comments       int                 `json:"comments"`
	ReviewComments int                 `json:"review_comments"`
	Additions      int                 `json:"additions"`
	Deletions      int                 `json:"deletions"`
	ChangedFiles   int                 `json:"changed_files"`
	Author         *PullRequestAuthor  `json:"author,omitempty"`
}

type ListOrganizationPullRequestsMetricsQuery struct {
	UserIDs        []string `form:"userIds" binding:"omitempty"`
	RepositoryName string   `form:"repositoryName" binding:"omitempty"`
	StartDate      string   `form:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate        string   `form:"endDate" binding:"omitempty,datetime=2006-01-02"`
}

type ListOrganizationSourceControlAccountsResponse struct {
	SourceControlAccounts []sourcecontroltypes.SourceControlAccount `json:"source_control_accounts"`
}

type GetMemberPullRequestsQuery struct {
	StartDate string `form:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"endDate" binding:"omitempty,datetime=2006-01-02"`
}

type GetMemberPullRequestsResponse struct {
	PullRequests []PullRequest `json:"pull_requests"`
}

type GetMemberPullRequestReviewsQuery struct {
	StartDate string `form:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"endDate" binding:"omitempty,datetime=2006-01-02"`
}

type GetMemberPullRequestReviewsResponse struct {
	Reviews []MemberActivity `json:"reviews"`
}

type GetOrganizationMetricsQuery struct {
	StartDate string   `form:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string   `form:"endDate" binding:"omitempty,datetime=2006-01-02"`
	Interval  string   `form:"interval" binding:"omitempty,oneof=daily weekly monthly"`
	TeamIDs   []string `form:"teamIds" binding:"omitempty"`
}
