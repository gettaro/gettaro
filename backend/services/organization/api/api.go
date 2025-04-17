package api

import (
	"context"

	"ems.dev/backend/services/organization/database"
	"ems.dev/backend/services/organization/types"
)

// OrganizationAPI defines the interface for organization operations
type OrganizationAPI interface {
	CreateOrganization(ctx context.Context, org *types.Organization, ownerID string) error
	GetUserOrganizations(ctx context.Context, userID string) ([]types.Organization, error)
	GetOrganizationByID(ctx context.Context, id string) (*types.Organization, error)
	UpdateOrganization(ctx context.Context, org *types.Organization) error
	DeleteOrganization(ctx context.Context, id string) error
	AddOrganizationMember(ctx context.Context, orgID string, userID string) error
	RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error
	GetOrganizationMembers(ctx context.Context, orgID string) ([]types.OrganizationMember, error)
	IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error)
}

type Api struct {
	db *database.OrganizationDB
}

func NewApi(orgDb *database.OrganizationDB) *Api {
	return &Api{
		db: orgDb,
	}
}

// AddOrganizationMember adds a user as a member to an organization
func (a *Api) AddOrganizationMember(ctx context.Context, orgID string, userID string) error {
	return a.db.AddOrganizationMember(orgID, userID)
}

// RemoveOrganizationMember removes a user from an organization
func (a *Api) RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error {
	return a.db.RemoveOrganizationMember(orgID, userID)
}

// GetOrganizationMembers returns all members of an organization
func (a *Api) GetOrganizationMembers(ctx context.Context, orgID string) ([]types.OrganizationMember, error) {
	return a.db.GetOrganizationMembers(orgID)
}

// IsOrganizationOwner checks if a user is the owner of an organization
func (a *Api) IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error) {
	return a.db.IsOrganizationOwner(orgID, userID)
}
