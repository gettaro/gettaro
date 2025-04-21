package api

import (
	"context"

	auth0client "ems.dev/backend/libraries/auth0"
	"ems.dev/backend/services/auth/database"
	"ems.dev/backend/services/auth/types"
)

// AuthAPI defines the interface for authentication-related operations
type AuthAPI interface {
	// GetUserInfo retrieves user information from Auth0
	GetUserInfo(ctx context.Context, accessToken string) (*auth0client.UserInfo, error)
	// GetExternalAuth retrieves an auth provider by its provider ID
	GetExternalAuth(ctx context.Context, providerID string) (*types.AuthProvider, error)
	// CreateExternalAuth creates a new external auth provider entry
	CreateExternalAuth(ctx context.Context, authProvider *types.AuthProvider) error
}

type Api struct {
	auth0Client auth0client.Auth0Client
	authDB      database.AuthDB
}

// NewApi creates a new instance of the auth API
func NewApi(auth0Client auth0client.Auth0Client, authDB database.AuthDB) *Api {
	return &Api{
		auth0Client: auth0Client,
		authDB:      authDB,
	}
}

// GetUserInfo retrieves user information from Auth0
func (a *Api) GetUserInfo(ctx context.Context, accessToken string) (*auth0client.UserInfo, error) {
	return a.auth0Client.GetUserInfo(ctx, accessToken)
}

// GetExternalAuth retrieves an auth provider by its provider ID
func (a *Api) GetExternalAuth(ctx context.Context, providerID string) (*types.AuthProvider, error) {
	return a.authDB.GetExternalAuth(ctx, providerID)
}

// CreateExternalAuth creates a new external auth provider entry
func (a *Api) CreateExternalAuth(ctx context.Context, authProvider *types.AuthProvider) error {
	return a.authDB.CreateExternalAuth(ctx, authProvider)
}
