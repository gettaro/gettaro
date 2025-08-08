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
	AssignUserTitle(ctx context.Context, userTitle types.UserTitle) error
	GetUserTitle(ctx context.Context, userID string, orgID string) (*types.UserTitle, error)
	RemoveUserTitle(ctx context.Context, userID string, orgID string) error
}

// Api implements the TitleAPI interface
type Api struct {
	db *database.TitleDBImpl
}

// NewApi creates a new TitleAPI instance
func NewApi(db *database.TitleDBImpl) *Api {
	return &Api{
		db: db,
	}
}
