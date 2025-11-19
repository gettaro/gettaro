package api

import (
	"context"

	"ems.dev/backend/services/directs/types"
)

// GetOrgChartFlat retrieves a flat list of all manager-direct relationships
func (a *DirectReportsAPIImpl) GetOrgChartFlat(ctx context.Context, orgID string) ([]types.DirectReport, error) {
	return a.db.GetOrgChartFlat(ctx, orgID)
}
