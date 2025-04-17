package api

import (
	"ems.dev/backend/services/user/database"
	"ems.dev/backend/services/user/types"
)

// UserAPI defines the interface for user-related operations
type UserAPI interface {
	// FindUser searches for a user by ID or email
	FindUser(params types.UserSearchParams) (*types.User, error)

	// GetOrCreateUserFromAuthProvider checks if a user exists for the given auth provider and ID,
	// and if not, creates a new user and auth provider entry
	GetOrCreateUserFromAuthProvider(provider string, providerID string, email string, name string) (*types.User, error)
}

type Api struct {
	db *database.UserDB
}

func NewApi(userDb *database.UserDB) *Api {
	return &Api{
		db: userDb,
	}
}
