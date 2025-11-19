package api

import (
	"context"

	"ems.dev/backend/services/conversation/types"
)

// GetConversationWithDetails retrieves a conversation with related data
func (a *ConversationAPI) GetConversationWithDetails(ctx context.Context, id string) (*types.ConversationWithDetails, error) {
	return a.db.GetConversationWithDetails(ctx, id)
}
