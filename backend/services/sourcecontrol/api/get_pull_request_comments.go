package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

// GetPullRequestComments retrieves all comments for a specific pull request
func (a *Api) GetPullRequestComments(ctx context.Context, prID string) ([]*types.PRComment, error) {
	return a.db.GetPullRequestComments(ctx, prID)
}
