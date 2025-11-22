package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

// CreatePullRequest creates a new pull request
func (a *Api) CreatePullRequest(ctx context.Context, pr *types.PullRequest) (*types.PullRequest, error) {
	return a.db.CreatePullRequest(ctx, pr)
}
