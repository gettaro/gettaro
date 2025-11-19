package api

import (
	"context"

	"ems.dev/backend/services/conversation/types"
)

// ListConversations retrieves conversations with optional filters
func (a *ConversationAPI) ListConversations(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.Conversation, error) {
	return a.db.ListConversations(ctx, organizationID, query)
}
