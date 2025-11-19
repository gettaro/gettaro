package api

import (
	"context"
)

// DeleteDirectReport removes a direct report relationship
func (a *DirectReportsAPIImpl) DeleteDirectReport(ctx context.Context, id string) error {
	return a.db.DeleteDirectReport(ctx, id)
}
