package api

import (
	"context"

	"ems.dev/backend/services/title/database"
	"ems.dev/backend/services/title/types"
)

// TitleAPI defines the interface for title operations
type TitleAPI interface {
	CreateTitle(ctx context.Context, title types.Title) (*types.Title, error)
	GetTitle(ctx context.Context, id string) (*types.Title, error)
	ListTitles(ctx context.Context, orgID string) ([]types.Title, error)
	UpdateTitle(ctx context.Context, id string, title types.Title) (*types.Title, error)
	DeleteTitle(ctx context.Context, id string) error
	AssignMemberTitle(ctx context.Context, memberTitle types.MemberTitle) error
	GetMemberTitle(ctx context.Context, memberID string, orgID string) (*types.MemberTitle, error)
	RemoveMemberTitle(ctx context.Context, memberID string, orgID string) error
}

// Api implements the TitleAPI interface
type Api struct {
	db *database.DB
}

// NewApi creates a new TitleAPI instance
func NewApi(db *database.DB) *Api {
	return &Api{
		db: db,
	}
}

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

// AssignMemberTitle assigns a title to a member within an organization
func (s *Api) AssignMemberTitle(ctx context.Context, memberTitle types.MemberTitle) error {
	return s.db.AssignMemberTitle(memberTitle)
}

// GetMemberTitle retrieves a member's title assignment
func (s *Api) GetMemberTitle(ctx context.Context, memberID string, orgID string) (*types.MemberTitle, error) {
	return s.db.GetMemberTitle(memberID, orgID)
}

// RemoveMemberTitle removes a member's title assignment
func (s *Api) RemoveMemberTitle(ctx context.Context, memberID string, orgID string) error {
	return s.db.RemoveMemberTitle(memberID, orgID)
}
