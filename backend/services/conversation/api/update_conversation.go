package api

import (
	"context"

	"ems.dev/backend/services/conversation/types"
)

// UpdateConversation updates a conversation
func (a *ConversationAPI) UpdateConversation(ctx context.Context, id string, req *types.UpdateConversationRequest) error {
	return a.db.UpdateConversation(ctx, id, req)
}
