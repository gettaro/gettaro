package api

import (
	"context"

	"ems.dev/backend/services/member/types"
)

// GetExternalAccounts retrieves external accounts based on the given parameters
func (a *Api) GetExternalAccounts(ctx context.Context, params *types.ExternalAccountParams) ([]types.ExternalAccount, error) {
	return a.db.GetExternalAccounts(ctx, params)
}
