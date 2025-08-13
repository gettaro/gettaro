package api

import (
	"context"
	"encoding/json"
	"time"

	"ems.dev/backend/services/sourcecontrol/database"
	"ems.dev/backend/services/sourcecontrol/types"
)

// SourceControlAPI defines the interface for source control operations
type SourceControlAPI interface {
	// Source Control Accounts
	GetSourceControlAccountsByUsernames(ctx context.Context, usernames []string) (map[string]*types.SourceControlAccount, error)
	CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error
	GetSourceControlAccount(ctx context.Context, id string) (*types.SourceControlAccount, error)
	UpdateSourceControlAccount(ctx context.Context, account *types.SourceControlAccount) error
	GetSourceControlAccountsByOrganization(ctx context.Context, orgID string) ([]*types.SourceControlAccount, error)

	// Pull Requests
	GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error)
	CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error
	UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error
	GetPullRequestMetrics(ctx context.Context, params *types.PullRequestParams) (*types.PullRequestMetrics, error)

	// Comments
	CreatePRComments(ctx context.Context, comments []*types.PRComment) error
	GetPullRequestComments(ctx context.Context, prID string) ([]*types.PRComment, error)

	// Member Activity
	GetMemberActivity(ctx context.Context, params *types.MemberActivityParams) ([]*types.MemberActivity, error)
}

type Api struct {
	db database.DB
}

func NewAPI(db database.DB) SourceControlAPI {
	return &Api{
		db: db,
	}
}

func (a *Api) GetSourceControlAccountsByUsernames(ctx context.Context, usernames []string) (map[string]*types.SourceControlAccount, error) {
	return a.db.GetSourceControlAccountsByUsernames(ctx, usernames)
}

func (a *Api) CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error {
	return a.db.CreateSourceControlAccounts(ctx, accounts)
}

func (a *Api) GetSourceControlAccount(ctx context.Context, id string) (*types.SourceControlAccount, error) {
	return a.db.GetSourceControlAccount(ctx, id)
}

func (a *Api) UpdateSourceControlAccount(ctx context.Context, account *types.SourceControlAccount) error {
	return a.db.UpdateSourceControlAccount(ctx, account)
}

func (a *Api) GetSourceControlAccountsByOrganization(ctx context.Context, orgID string) ([]*types.SourceControlAccount, error) {
	return a.db.GetSourceControlAccountsByOrganization(ctx, orgID)
}

func (a *Api) GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error) {
	return a.db.GetPullRequests(ctx, params)
}

func (a *Api) CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error {
	return a.db.CreatePullRequests(ctx, prs)
}

func (a *Api) CreatePRComments(ctx context.Context, comments []*types.PRComment) error {
	return a.db.CreatePRComments(ctx, comments)
}

// GetPullRequestComments retrieves all comments for a specific pull request
func (a *Api) GetPullRequestComments(ctx context.Context, prID string) ([]*types.PRComment, error) {
	return a.db.GetPullRequestComments(ctx, prID)
}

func (a *Api) UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error {
	return a.db.UpdatePullRequest(ctx, pr)
}

// GetPullRequestMetrics handles the retrieval of pull request metrics for an organization.
// Params:
// - ctx: The context for the request, used for cancellation and timeouts
// - params: The parameters for filtering pull requests
// Returns:
// - *types.PullRequestMetrics: The calculated metrics including:
//   - Number of merged PRs
//   - Number of open PRs
//   - Mean time from publish to merge (in hours)
//   - Mean time to first review (in hours)
//
// - error: If any error occurs during the calculation
// Side Effects:
// - Makes a database query to fetch pull requests
// - Performs calculations on the fetched data
func (s *Api) GetPullRequestMetrics(
	ctx context.Context,
	params *types.PullRequestParams,
) (*types.PullRequestMetrics, error) {
	// Get pull requests from database
	prs, err := s.db.GetPullRequests(ctx, params)
	if err != nil {
		return nil, err
	}

	metrics := &types.PullRequestMetrics{}
	var totalPublishToMergeTime time.Duration
	var totalTimeToFirstReview time.Duration
	var mergedCount, reviewedCount int

	for _, pr := range prs {
		// Count open and merged PRs
		if pr.MergedAt != nil && pr.Status == "closed" {
			metrics.MergedPRsCount++

			var prMetrics map[string]interface{}
			if err := json.Unmarshal(pr.Metrics, &prMetrics); err != nil {
				return nil, err
			}

			// Calculate publish to merge time
			if seconds, ok := prMetrics["time_to_merge_seconds"].(float64); ok && seconds > 0 {
				totalPublishToMergeTime += time.Duration(seconds) * time.Second
				mergedCount++
			}

			// Calculate time to first review
			if seconds, ok := prMetrics["time_to_first_non_bot_review_seconds"].(float64); ok && seconds > 0 {
				totalTimeToFirstReview += time.Duration(seconds) * time.Second
				reviewedCount++
			}
		} else if pr.Status == "open" {
			metrics.OpenPRsCount++
		}
	}

	// Calculate mean times
	if mergedCount > 0 {
		metrics.MeanPublishToMergeTime = totalPublishToMergeTime.Hours() / float64(mergedCount)
	}
	if reviewedCount > 0 {
		metrics.MeanTimeToFirstReview = totalTimeToFirstReview.Hours() / float64(reviewedCount)
	}

	return metrics, nil
}

// GetMemberActivity handles the retrieval of source control activity timeline for a specific member.
// Params:
// - ctx: The context for the request, used for cancellation and timeouts
// - params: The parameters containing member ID and optional date range filters
// Returns:
// - []*types.MemberActivity: A timeline of activities including pull requests and comments
// - error: If any error occurs during the retrieval
// Side Effects:
// - Makes a database query to fetch member activities
func (a *Api) GetMemberActivity(ctx context.Context, params *types.MemberActivityParams) ([]*types.MemberActivity, error) {
	return a.db.GetMemberActivity(ctx, params)
}
