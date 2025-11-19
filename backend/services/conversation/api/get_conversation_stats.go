package api

import (
	"context"
)

// GetConversationStats returns statistics for conversations
func (a *ConversationAPI) GetConversationStats(ctx context.Context, organizationID string, managerMemberID *string) (map[string]int, error) {
	return a.db.GetConversationStats(ctx, organizationID, managerMemberID)
}
