package api

import (
	"context"
	"fmt"

	"ems.dev/backend/services/organization/types"
	userapi "ems.dev/backend/services/user/api"
	usertypes "ems.dev/backend/services/user/types"
)

// OrganizationAPI defines the interface for organization operations
type OrganizationAPI interface {
	CreateOrganization(ctx context.Context, org *types.Organization, ownerID string) error
	GetUserOrganizations(ctx context.Context, userID string) ([]types.Organization, error)
	GetOrganizationByID(ctx context.Context, id string) (*types.Organization, error)
	UpdateOrganization(ctx context.Context, org *types.Organization) error
	DeleteOrganization(ctx context.Context, id string) error
	AddOrganizationMemberByEmail(ctx context.Context, orgID string, email string) error
	RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error
	GetOrganizationMembers(ctx context.Context, orgID string) ([]types.OrganizationMember, error)
	IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error)
}

// OrganizationDB defines the interface for organization database operations
type OrganizationDB interface {
	CreateOrganization(org *types.Organization, ownerID string) error
	GetUserOrganizations(userID string) ([]types.Organization, error)
	GetOrganizationByID(id string) (*types.Organization, error)
	UpdateOrganization(org *types.Organization) error
	DeleteOrganization(id string) error
	AddOrganizationMember(orgID string, userID string) error
	RemoveOrganizationMember(orgID string, userID string) error
	GetOrganizationMembers(orgID string) ([]types.OrganizationMember, error)
	IsOrganizationOwner(orgID string, userID string) (bool, error)
}

type Api struct {
	db      OrganizationDB
	userApi userapi.UserAPI
}

func NewApi(orgDb OrganizationDB, userApi userapi.UserAPI) *Api {
	return &Api{
		db:      orgDb,
		userApi: userApi,
	}
}

// AddOrganizationMemberByEmail adds a user as a member to an organization by their email
func (a *Api) AddOrganizationMemberByEmail(ctx context.Context, orgID string, email string) error {
	// Look up user by email
	user, err := a.userApi.FindUser(usertypes.UserSearchParams{Email: &email})
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Add user as member
	return a.db.AddOrganizationMember(orgID, user.ID)
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
