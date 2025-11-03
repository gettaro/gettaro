package api

import (
	"context"

	"ems.dev/backend/libraries/errors"
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

// UpdateExternalAccountMemberID updates the member_id association for an external account
// Validates that the account belongs to the specified organization
func (a *Api) UpdateExternalAccountMemberID(ctx context.Context, organizationID string, accountID string, memberID *string) (*types.ExternalAccount, error) {
	// Get the existing account
	existingAccount, err := a.GetExternalAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if existingAccount == nil {
		return nil, errors.NewNotFoundError("external account not found")
	}

	// Verify the account belongs to the organization
	if existingAccount.OrganizationID == nil || *existingAccount.OrganizationID != organizationID {
		return nil, errors.NewNotFoundError("external account not found in this organization")
	}

	// Update the member_id
	existingAccount.MemberID = memberID

	// Update the account
	if err := a.UpdateExternalAccount(ctx, existingAccount); err != nil {
		return nil, err
	}

	return existingAccount, nil
}

