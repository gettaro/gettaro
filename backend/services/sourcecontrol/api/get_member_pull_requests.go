package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

// GetMemberPullRequests handles the retrieval of pull requests for a specific member.
// Params:
// - ctx: The context for the request, used for cancellation and timeouts
// - params: The parameters containing member ID and optional date range filters
// Returns:
// - []*types.PullRequestWithComments: A list of pull requests created by the member with optional comments, ordered by created_at descending
// - error: If any error occurs during the retrieval
// Side Effects:
// - Makes a database query to fetch member pull requests and optionally their comments
func (a *Api) GetMemberPullRequests(ctx context.Context, params *types.MemberPullRequestParams) ([]*types.PullRequestWithComments, error) {
	return a.db.GetMemberPullRequests(ctx, params)
}
