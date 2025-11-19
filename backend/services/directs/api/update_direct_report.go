package api

import (
	"context"

	"ems.dev/backend/services/directs/types"
)

// UpdateDirectReport updates a direct report relationship
func (a *DirectReportsAPIImpl) UpdateDirectReport(ctx context.Context, id string, params types.UpdateDirectReportParams) error {
	return a.db.UpdateDirectReport(ctx, id, params)
}
