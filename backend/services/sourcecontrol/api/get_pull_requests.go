package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

// GetPullRequests retrieves pull requests based on the given parameters
func (a *Api) GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error) {
	return a.db.GetPullRequests(ctx, params)
}
