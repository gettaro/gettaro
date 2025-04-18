package api

import (
	orgtypes "ems.dev/backend/services/organization/types"
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

// UserDB defines the interface for user database operations
type UserDB interface {
	FindUser(params types.UserSearchParams) (*types.User, error)
	GetOrCreateUserFromAuthProvider(provider string, providerID string, email string, name string) (*types.User, error)
	CreateOrganizationWithOwner(org *orgtypes.Organization, userID string) error
	GetUserOrganizations(userID string) ([]orgtypes.Organization, error)
	CreateUser(user *types.User) error
	UpdateUser(user *types.User) error
	DeleteUser(userID string) error
	GetUserByID(userID string) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	ListUsers() ([]types.User, error)
}

type Api struct {
	db UserDB
}

func NewApi(userDb UserDB) *Api {
	return &Api{
		db: userDb,
	}
}
