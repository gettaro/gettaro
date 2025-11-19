package api

import (
	"context"

	"ems.dev/backend/services/conversation/types"
)

// ListConversationsWithDetails retrieves conversations with related data
func (a *ConversationAPI) ListConversationsWithDetails(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.ConversationWithDetails, error) {
	return a.db.ListConversationsWithDetails(ctx, organizationID, query)
}
