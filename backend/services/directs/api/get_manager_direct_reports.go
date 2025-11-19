package api

import (
	"context"

	"ems.dev/backend/services/directs/types"
)

// GetManagerDirectReports retrieves all direct reports for a specific manager
func (a *DirectReportsAPIImpl) GetManagerDirectReports(ctx context.Context, managerMemberID, orgID string) ([]types.DirectReport, error) {
	return a.db.GetManagerDirectReports(ctx, managerMemberID, orgID)
}
