package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/database"
	"ems.dev/backend/services/sourcecontrol/metrics"
	"ems.dev/backend/services/sourcecontrol/types"
)

// SourceControlAPI defines the interface for source control operations
type SourceControlAPI interface {
	// Pull Requests
	GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error)
	CreatePullRequest(ctx context.Context, pr *types.PullRequest) (*types.PullRequest, error)
	UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error

	// Comments
	CreatePRComments(ctx context.Context, comments []*types.PRComment) error
	GetPullRequestComments(ctx context.Context, prID string) ([]*types.PRComment, error)

	// Member Activity
	GetMemberPullRequests(ctx context.Context, params *types.MemberPullRequestParams) ([]*types.PullRequestWithComments, error)
	GetMemberPullRequestReviews(ctx context.Context, params *types.MemberPullRequestReviewsParams) ([]*types.MemberActivity, error)

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
