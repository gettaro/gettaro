package api

import (
	"context"

	"ems.dev/backend/services/conversation/types"
)

// GetConversation retrieves a conversation by ID
func (a *ConversationAPI) GetConversation(ctx context.Context, id string) (*types.Conversation, error) {
	return a.db.GetConversation(ctx, id)
}
