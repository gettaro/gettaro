package api

import (
	userdb "ems.dev/backend/services/user/database"
	"ems.dev/backend/services/user/types"
)

// UserAPI defines the interface for user-related operations
type UserAPI interface {
	// FindUser searches for a user by ID or email
	FindUser(params types.UserSearchParams) (*types.User, error)

	// CreateUser creates a new user in the system and returns the created user
	CreateUser(user *types.User) (*types.User, error)
}

type Api struct {
	db userdb.DB
}

func NewApi(userDb userdb.DB) *Api {
	return &Api{
		db: userDb,
	}
}

// CreateUser creates a new user in the system and returns the created user
func (s *Api) CreateUser(user *types.User) (*types.User, error) {
	createdUser, err := s.db.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
