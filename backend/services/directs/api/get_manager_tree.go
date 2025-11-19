package api

import (
	"context"

	"ems.dev/backend/services/directs/types"
)

// GetManagerTree retrieves the full management tree for a manager
func (a *DirectReportsAPIImpl) GetManagerTree(ctx context.Context, managerMemberID, orgID string) ([]types.OrgChartNode, error) {
	return a.db.GetManagerTree(ctx, managerMemberID, orgID)
}
