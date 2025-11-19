package api

import (
	"context"

	"ems.dev/backend/services/title/types"
)

// CreateTitle creates a new title for an organization
func (s *Api) CreateTitle(ctx context.Context, title types.Title) (*types.Title, error) {
	err := s.db.CreateTitle(&title)
	if err != nil {
		return nil, err
	}

	return &title, nil
}
