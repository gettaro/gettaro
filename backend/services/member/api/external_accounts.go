package api

import (
	"context"

	"ems.dev/backend/services/member/types"
)

// GetExternalAccounts retrieves external accounts based on the given parameters
func (a *Api) GetExternalAccounts(ctx context.Context, params *types.ExternalAccountParams) ([]types.ExternalAccount, error) {
	return a.db.GetExternalAccounts(ctx, params)
}

// CreateExternalAccounts creates multiple external accounts
func (a *Api) CreateExternalAccounts(ctx context.Context, accounts []*types.ExternalAccount) error {
	return a.db.CreateExternalAccounts(ctx, accounts)
}

// GetExternalAccount retrieves an external account by ID
func (a *Api) GetExternalAccount(ctx context.Context, id string) (*types.ExternalAccount, error) {
	return a.db.GetExternalAccount(ctx, id)
}

// UpdateExternalAccount updates an existing external account
func (a *Api) UpdateExternalAccount(ctx context.Context, account *types.ExternalAccount) error {
	return a.db.UpdateExternalAccount(ctx, account)
}

