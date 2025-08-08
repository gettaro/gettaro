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

// GetTitle retrieves a title by ID
func (s *Api) GetTitle(ctx context.Context, id string) (*types.Title, error) {
	return s.db.GetTitle(id)
}

// ListTitles retrieves all titles for an organization
func (s *Api) ListTitles(ctx context.Context, orgID string) ([]types.Title, error) {
	return s.db.ListTitles(orgID)
}

// UpdateTitle updates an existing title
func (s *Api) UpdateTitle(ctx context.Context, id string, request types.Title) (*types.Title, error) {
	title, err := s.db.GetTitle(id)
	if err != nil {
		return nil, err
	}
	if title == nil {
		return nil, nil
	}

	title.Name = request.Name
	err = s.db.UpdateTitle(*title)
	if err != nil {
		return nil, err
	}

	return title, nil
}

// DeleteTitle deletes a title
func (s *Api) DeleteTitle(ctx context.Context, id string) error {
	return s.db.DeleteTitle(id)
}

// AssignUserTitle assigns a title to a user within an organization
func (s *Api) AssignUserTitle(ctx context.Context, userTitle types.UserTitle) error {
	return s.db.AssignUserTitle(userTitle)
}

// GetUserTitle retrieves a user's title assignment
func (s *Api) GetUserTitle(ctx context.Context, userID string, orgID string) (*types.UserTitle, error) {
	return s.db.GetUserTitle(userID, orgID)
}

// RemoveUserTitle removes a user's title assignment
func (s *Api) RemoveUserTitle(ctx context.Context, userID string, orgID string) error {
	return s.db.RemoveUserTitle(userID, orgID)
}
