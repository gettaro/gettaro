package database

import (
	"context"

	"ems.dev/backend/services/member/types"
	"gorm.io/gorm"
)

// GetExternalAccounts retrieves external accounts based on the given parameters
func (d *MemberDB) GetExternalAccounts(ctx context.Context, params *types.ExternalAccountParams) ([]types.ExternalAccount, error) {
	var accounts []types.ExternalAccount
	query := d.db.WithContext(ctx).Model(&types.ExternalAccount{})

	// Query by external account IDs if provided
	if len(params.ExternalAccountIDs) > 0 {
		query = query.Where("id IN ?", params.ExternalAccountIDs)
	}

	// Query by usernames if provided
	if len(params.Usernames) > 0 {
		query = query.Where("username IN ?", params.Usernames)
	}

	// Filter by organization ID if provided
	if params.OrganizationID != "" {
		query = query.Where("organization_id = ?", params.OrganizationID)
	}

	// Filter by member IDs if provided
	if len(params.MemberIDs) > 0 {
		query = query.Where("member_id IN ?", params.MemberIDs)
	}

	// Filter by account type if provided
	if params.AccountType != nil && *params.AccountType != "" {
		query = query.Where("account_type = ?", *params.AccountType)
	}

	if err := query.Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

// CreateExternalAccounts creates multiple external accounts
func (d *MemberDB) CreateExternalAccounts(ctx context.Context, accounts []*types.ExternalAccount) error {
	return d.db.WithContext(ctx).Create(accounts).Error
}

// GetExternalAccount retrieves an external account by ID
func (d *MemberDB) GetExternalAccount(ctx context.Context, id string) (*types.ExternalAccount, error) {
	var account types.ExternalAccount
	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// UpdateExternalAccount updates an existing external account
func (d *MemberDB) UpdateExternalAccount(ctx context.Context, account *types.ExternalAccount) error {
	// Use explicit field updates to handle nil values properly
	updates := map[string]interface{}{
		"member_id":       account.MemberID,
		"organization_id": account.OrganizationID,
		"account_type":    account.AccountType,
		"provider_name":   account.ProviderName,
		"provider_id":     account.ProviderID,
		"username":        account.Username,
		"metadata":        account.Metadata,
		"last_synced_at":  account.LastSyncedAt,
		"updated_at":      account.UpdatedAt,
	}

	return d.db.WithContext(ctx).Model(account).Updates(updates).Error
}

