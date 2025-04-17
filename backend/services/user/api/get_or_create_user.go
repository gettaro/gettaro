package api

import "ems.dev/backend/services/user/types"

// GetOrCreateUserFromAuthProvider checks if a user exists for the given auth provider and ID,
// and if not, creates a new user and auth provider entry
func (s *Api) GetOrCreateUserFromAuthProvider(provider string, providerID string, email string, name string) (*types.User, error) {
	return s.db.GetOrCreateUserFromAuthProvider(provider, providerID, email, name)
}
