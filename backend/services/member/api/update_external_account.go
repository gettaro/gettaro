package api

import (
	"context"

	"ems.dev/backend/services/member/types"
)

// UpdateExternalAccount updates an existing external account
func (a *Api) UpdateExternalAccount(ctx context.Context, account *types.ExternalAccount) error {
	return a.db.UpdateExternalAccount(ctx, account)
}
