package api

import (
	"context"

	"ems.dev/backend/services/directs/types"
)

// GetMemberManager retrieves the manager of a specific member
func (a *DirectReportsAPIImpl) GetMemberManager(ctx context.Context, reportMemberID, orgID string) (*types.DirectReport, error) {
	return a.db.GetMemberManager(ctx, reportMemberID, orgID)
}
