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
	db database.TitleDB
}

// NewApi creates a new TitleAPI instance
func NewApi(db database.TitleDB) *Api {
	return &Api{
		db: db,
	}
}
