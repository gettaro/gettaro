package api

import (
	"context"

	"ems.dev/backend/services/directs/types"
)

// ListDirectReports retrieves direct reports based on search parameters
func (a *DirectReportsAPIImpl) ListDirectReports(ctx context.Context, params types.DirectReportSearchParams) ([]types.DirectReport, error) {
	return a.db.ListDirectReports(ctx, params)
}
