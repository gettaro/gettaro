package api

import (
	"context"

	"ems.dev/backend/services/directs/types"
)

// GetMemberManagementChain retrieves the full management chain for a member
func (a *DirectReportsAPIImpl) GetMemberManagementChain(ctx context.Context, reportMemberID, orgID string) ([]types.ManagementChain, error) {
	return a.db.GetMemberManagementChain(ctx, reportMemberID, orgID)
}
