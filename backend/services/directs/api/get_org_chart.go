package api

import (
	"context"

	"ems.dev/backend/services/directs/types"
)

// GetOrgChart retrieves the complete organizational chart
func (a *DirectReportsAPIImpl) GetOrgChart(ctx context.Context, orgID string) ([]types.OrgChartNode, error) {
	return a.db.GetOrgChart(ctx, orgID)
}
