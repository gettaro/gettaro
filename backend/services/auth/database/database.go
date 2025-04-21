package database

import (
	"context"

	"ems.dev/backend/services/auth/types"
	"gorm.io/gorm"
)

// AuthDB defines the interface for authentication-related database operations
type AuthDB interface {
	// GetExternalAuth retrieves an auth provider by its provider ID
	GetExternalAuth(ctx context.Context, providerID string) (*types.AuthProvider, error)
	// CreateExternalAuth creates a new external auth provider entry
	CreateExternalAuth(ctx context.Context, authProvider *types.AuthProvider) error
}

type authDB struct {
	db *gorm.DB
}

// New creates a new instance of the auth database
func New(db *gorm.DB) AuthDB {
	return &authDB{db}
}

// GetExternalAuth retrieves an auth provider by its provider ID
func (d *authDB) GetExternalAuth(ctx context.Context, providerID string) (*types.AuthProvider, error) {
	var authProvider types.AuthProvider
	err := d.db.WithContext(ctx).
		Where("provider_id = ?", providerID).
		First(&authProvider).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &authProvider, nil
}

// CreateExternalAuth creates a new external auth provider entry
func (d *authDB) CreateExternalAuth(ctx context.Context, authProvider *types.AuthProvider) error {
	return d.db.WithContext(ctx).Create(authProvider).Error
}
