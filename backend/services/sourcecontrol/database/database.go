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

	// Pull Requests
	CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error

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

// CreatePullRequests creates multiple pull requests
func (d *SourceControlDB) CreatePullRequests(ctx context.Context, prs []*types.PullRequest) error {
	return d.db.WithContext(ctx).Create(prs).Error
}

// CreatePRComments creates multiple PR comments
func (d *SourceControlDB) CreatePRComments(ctx context.Context, comments []*types.PRComment) error {
	return d.db.WithContext(ctx).Create(comments).Error
}
