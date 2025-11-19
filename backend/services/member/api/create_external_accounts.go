package api

import (
	"context"

	"ems.dev/backend/services/member/types"
)

// CreateExternalAccounts creates multiple external accounts
func (a *Api) CreateExternalAccounts(ctx context.Context, accounts []*types.ExternalAccount) error {
	return a.db.CreateExternalAccounts(ctx, accounts)
}
