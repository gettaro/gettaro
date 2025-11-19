package api

import (
	"context"

	"ems.dev/backend/services/directs/types"
)

// GetDirectReport retrieves a direct report by ID
func (a *DirectReportsAPIImpl) GetDirectReport(ctx context.Context, id string) (*types.DirectReport, error) {
	return a.db.GetDirectReport(ctx, id)
}
