package api

import (
	"context"

	"ems.dev/backend/services/member/types"
)

// GetExternalAccount retrieves an external account by ID
func (a *Api) GetExternalAccount(ctx context.Context, id string) (*types.ExternalAccount, error) {
	return a.db.GetExternalAccount(ctx, id)
}
