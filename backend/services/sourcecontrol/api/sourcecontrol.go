package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/database"
	"ems.dev/backend/services/sourcecontrol/metrics"
	"ems.dev/backend/services/sourcecontrol/types"
)

// SourceControlAPI defines the interface for source control operations
type SourceControlAPI interface {
	// Source Control Accounts
	GetSourceControlAccounts(ctx context.Context, params *types.SourceControlAccountParams) ([]types.SourceControlAccount, error)
	CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error
	GetSourceControlAccount(ctx context.Context, id string) (*types.SourceControlAccount, error)
	UpdateSourceControlAccount(ctx context.Context, account *types.SourceControlAccount) error

	// Pull Requests
	GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error)
	CreatePullRequest(ctx context.Context, pr *types.PullRequest) (*types.PullRequest, error)
	UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error

	// Comments
	CreatePRComments(ctx context.Context, comments []*types.PRComment) error
	GetPullRequestComments(ctx context.Context, prID string) ([]*types.PRComment, error)

	// Member Activity
	GetMemberActivity(ctx context.Context, params *types.MemberActivityParams) ([]*types.MemberActivity, error)
	GetMemberPullRequests(ctx context.Context, params *types.MemberPullRequestParams) ([]*types.PullRequest, error)

	// CalculateMetrics calculates source control metrics
	CalculateMetrics(ctx context.Context, params types.MetricRuleParams) (*types.MetricsResponse, error)
}

type Api struct {
	db            database.DB
	metricsEngine metrics.MetricsEngine
}

func NewAPI(db database.DB) SourceControlAPI {
	return &Api{
		db:            db,
		metricsEngine: metrics.NewEngine(db),
	}
}

func (a *Api) GetSourceControlAccounts(ctx context.Context, params *types.SourceControlAccountParams) ([]types.SourceControlAccount, error) {
	return a.db.GetSourceControlAccounts(ctx, params)
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

func (a *Api) GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error) {
	return a.db.GetPullRequests(ctx, params)
}

func (a *Api) CreatePullRequest(ctx context.Context, pr *types.PullRequest) (*types.PullRequest, error) {
	return a.db.CreatePullRequest(ctx, pr)
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

// GetMemberPullRequests handles the retrieval of pull requests for a specific member.
// Params:
// - ctx: The context for the request, used for cancellation and timeouts
// - params: The parameters containing member ID and optional date range filters
// Returns:
// - []*types.PullRequest: A list of pull requests created by the member, ordered by created_at descending
// - error: If any error occurs during the retrieval
// Side Effects:
// - Makes a database query to fetch member pull requests
func (a *Api) GetMemberPullRequests(ctx context.Context, params *types.MemberPullRequestParams) ([]*types.PullRequest, error) {
	return a.db.GetMemberPullRequests(ctx, params)
}

// CalculateMetrics calculates source control metrics
func (a *Api) CalculateMetrics(ctx context.Context, params types.MetricRuleParams) (*types.MetricsResponse, error) {
	return a.metricsEngine.CalculateMetrics(ctx, params)
}
