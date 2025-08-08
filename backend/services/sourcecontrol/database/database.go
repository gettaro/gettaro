package database

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
	"gorm.io/gorm"
)

// SourceControlDB defines the interface for source control database operations
type DB interface {
	// Source Control Accounts
	GetSourceControlAccountsByUsernames(ctx context.Context, usernames []string) (map[string]*types.SourceControlAccount, error)
	CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error
	GetSourceControlAccount(ctx context.Context, id string) (*types.SourceControlAccount, error)
	UpdateSourceControlAccount(ctx context.Context, account *types.SourceControlAccount) error
	GetSourceControlAccountsByOrganization(ctx context.Context, orgID string) ([]*types.SourceControlAccount, error)

	// Pull Requests
	GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error)
	CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error
	UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error

	// Comments
	CreatePRComments(ctx context.Context, comments []*types.PRComment) error
}

type SourceControlDB struct {
	db *gorm.DB
}

// NewSourceControlDB creates a new instance of the source control database
func NewSourceControlDB(db *gorm.DB) *SourceControlDB {
	return &SourceControlDB{
		db: db,
	}
}

// GetSourceControlAccountsByUsernames retrieves source control accounts by their usernames
func (d *SourceControlDB) GetSourceControlAccountsByUsernames(ctx context.Context, usernames []string) (map[string]*types.SourceControlAccount, error) {
	var accounts []types.SourceControlAccount
	err := d.db.WithContext(ctx).
		Where("username IN ?", usernames).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]*types.SourceControlAccount)
	for i := range accounts {
		result[accounts[i].Username] = &accounts[i]
	}

	return result, nil
}

// CreateSourceControlAccounts creates multiple source control accounts
func (d *SourceControlDB) CreateSourceControlAccounts(ctx context.Context, accounts []*types.SourceControlAccount) error {
	return d.db.WithContext(ctx).Create(accounts).Error
}

// GetPullRequests retrieves pull requests based on the given parameters
func (d *SourceControlDB) GetPullRequests(ctx context.Context, params *types.PullRequestParams) ([]*types.PullRequest, error) {
	var prs []types.PullRequest
	query := d.db.WithContext(ctx).Model(&types.PullRequest{})

	if params.ProviderID != "" {
		query = query.Where("provider_id = ?", params.ProviderID)
	}
	if params.OrganizationID != nil {
		query = query.Where("organization_id = ?", *params.OrganizationID)
	}
	if params.RepositoryName != "" {
		query = query.Where("repository_name = ?", params.RepositoryName)
	}
	// Add user IDs filter if provided
	if len(params.UserIDs) > 0 {
		query = query.Where("user_id IN ?", params.UserIDs)
	}
	// Add date range filters if provided
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", params.EndDate)
	}

	if err := query.Find(&prs).Error; err != nil {
		return nil, err
	}

	result := make([]*types.PullRequest, len(prs))
	for i := range prs {
		result[i] = &prs[i]
	}
	return result, nil
}

// CreatePullRequests creates multiple pull requests
func (d *SourceControlDB) CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error {
	return d.db.WithContext(ctx).Create(prs).Error
}

// CreatePRComments creates multiple PR comments
func (d *SourceControlDB) CreatePRComments(ctx context.Context, comments []*types.PRComment) error {
	return d.db.WithContext(ctx).Create(comments).Error
}

// UpdatePullRequest updates an existing pull request
func (d *SourceControlDB) UpdatePullRequest(ctx context.Context, pr *types.PullRequest) error {
	return d.db.WithContext(ctx).Model(pr).Updates(pr).Error
}

// UpdateSourceControlAccount updates an existing source control account
func (d *SourceControlDB) UpdateSourceControlAccount(ctx context.Context, account *types.SourceControlAccount) error {
	return d.db.WithContext(ctx).Model(account).Updates(account).Error
}

// GetSourceControlAccount retrieves a source control account by ID
func (d *SourceControlDB) GetSourceControlAccount(ctx context.Context, id string) (*types.SourceControlAccount, error) {
	var account types.SourceControlAccount
	err := d.db.WithContext(ctx).First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetSourceControlAccountsByOrganization retrieves source control accounts for an organization
func (d *SourceControlDB) GetSourceControlAccountsByOrganization(ctx context.Context, orgID string) ([]*types.SourceControlAccount, error) {
	var accounts []types.SourceControlAccount
	err := d.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	result := make([]*types.SourceControlAccount, len(accounts))
	for i := range accounts {
		result[i] = &accounts[i]
	}
	return result, nil
}
