package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

// CreatePRComments creates multiple PR comments
func (a *Api) CreatePRComments(ctx context.Context, comments []*types.PRComment) error {
	return a.db.CreatePRComments(ctx, comments)
}
