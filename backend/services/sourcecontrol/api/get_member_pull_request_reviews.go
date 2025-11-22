package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

// GetMemberPullRequestReviews handles the retrieval of pull request reviews for a specific member.
// Params:
// - ctx: The context for the request, used for cancellation and timeouts
// - params: The parameters containing member ID and optional date range filters
// Returns:
// - []*types.MemberActivity: A list of pull request reviews by the member, ordered by created_at descending
// - error: If any error occurs during the retrieval
// Side Effects:
// - Makes a database query to fetch member pull request reviews
func (a *Api) GetMemberPullRequestReviews(ctx context.Context, params *types.MemberPullRequestReviewsParams) ([]*types.MemberActivity, error) {
	return a.db.GetMemberPullRequestReviews(ctx, params)
}
