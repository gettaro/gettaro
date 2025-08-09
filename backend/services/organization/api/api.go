package api

import (
	"context"

	orgdb "ems.dev/backend/services/organization/database"
	"ems.dev/backend/services/organization/types"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	titleapi "ems.dev/backend/services/title/api"
	userapi "ems.dev/backend/services/user/api"
)

// OrganizationAPI defines the interface for organization operations
type OrganizationAPI interface {
	CreateOrganization(ctx context.Context, org *types.Organization, ownerID string) error
	GetMemberOrganizations(ctx context.Context, userID string) ([]types.Organization, error)
	GetOrganizations(ctx context.Context) ([]types.Organization, error)
	GetOrganizationByID(ctx context.Context, id string) (*types.Organization, error)
	UpdateOrganization(ctx context.Context, org *types.Organization) error
	DeleteOrganization(ctx context.Context, id string) error
}

type Api struct {
	db               orgdb.DB
	userApi          userapi.UserAPI
	titleApi         titleapi.TitleAPI
	sourceControlApi sourcecontrolapi.SourceControlAPI
}

func NewApi(orgDb orgdb.DB, userApi userapi.UserAPI, titleApi titleapi.TitleAPI, sourceControlApi sourcecontrolapi.SourceControlAPI) *Api {
	return &Api{
		db:               orgDb,
		userApi:          userApi,
		titleApi:         titleApi,
		sourceControlApi: sourceControlApi,
	}
}

// GetOrganizations returns all organizations in the system
func (a *Api) GetOrganizations(ctx context.Context) ([]types.Organization, error) {
	return a.db.GetOrganizations()
}
