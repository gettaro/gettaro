package api

import (
	"context"

	"ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/member/types"
)

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
