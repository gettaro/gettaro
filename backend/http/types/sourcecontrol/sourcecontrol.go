package sourcecontrol

import (
	"time"

	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
)

type ListOrganizationPullRequestsQuery struct {
	UserIDs        []string `form:"userIds" binding:"omitempty"`
	RepositoryName string   `form:"repositoryName" binding:"omitempty"`
	StartDate      string   `form:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate        string   `form:"endDate" binding:"omitempty,datetime=2006-01-02"`
	Status         string   `form:"status" binding:"omitempty,oneof=open closed merged"`
}

type ListOrganizationPullRequestsResponse struct {
	PullRequests []PullRequest `json:"pull_requests"`
}

type PullRequest struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	URL            string     `json:"url"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	MergedAt       *time.Time `json:"merged_at,omitempty"`
	Comments       int        `json:"comments"`
	ReviewComments int        `json:"review_comments"`
	Additions      int        `json:"additions"`
	Deletions      int        `json:"deletions"`
	ChangedFiles   int        `json:"changed_files"`
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
