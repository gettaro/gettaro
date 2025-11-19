package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

// UpdatePullRequest updates an existing pull request
func (a *Api) UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error {
	return a.db.UpdatePullRequest(ctx, pr)
}
