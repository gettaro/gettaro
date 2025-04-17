package api

import (
	"ems.dev/backend/services/user/types"
)

// FindUser searches for a user by ID or email
func (s *Api) FindUser(params types.UserSearchParams) (*types.User, error) {
	return s.db.FindUser(params)
}
