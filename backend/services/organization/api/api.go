package api

import (
	"context"

	"ems.dev/backend/libraries/errors"
	orgdb "ems.dev/backend/services/organization/database"
	"ems.dev/backend/services/organization/types"
	userapi "ems.dev/backend/services/user/api"
	usertypes "ems.dev/backend/services/user/types"
)

// OrganizationAPI defines the interface for organization operations
type OrganizationAPI interface {
	CreateOrganization(ctx context.Context, org *types.Organization, ownerID string) error
	GetUserOrganizations(ctx context.Context, userID string) ([]types.Organization, error)
	GetOrganizations(ctx context.Context) ([]types.Organization, error)
	GetOrganizationByID(ctx context.Context, id string) (*types.Organization, error)
	UpdateOrganization(ctx context.Context, org *types.Organization) error
	DeleteOrganization(ctx context.Context, id string) error
	AddOrganizationMember(ctx context.Context, member *types.UserOrganization) error
	RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error
	GetOrganizationMembers(ctx context.Context, orgID string) ([]types.UserOrganization, error)
	IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error)
}

type Api struct {
	db      orgdb.DB
	userApi userapi.UserAPI
}

func NewApi(orgDb orgdb.DB, userApi userapi.UserAPI) *Api {
	return &Api{
		db:      orgDb,
		userApi: userApi,
	}
}

// AddOrganizationMember adds a user as a member to an organization
func (a *Api) AddOrganizationMember(ctx context.Context, member *types.UserOrganization) error {
	// Look up user by email
	user, err := a.userApi.FindUser(usertypes.UserSearchParams{Email: &member.Email})
	if err != nil {
		return err
	}

	if user == nil {
		user, err = a.userApi.CreateUser(&usertypes.User{
			Email: member.Email,
		})

		if err != nil {
			return err
		}
	}

	// Check for duplicate member
	existingMember, err := a.db.GetOrganizationMember(member.OrganizationID, user.ID)
	if err != nil {
		return err
	}

	if existingMember != nil {
		return errors.NewConflictError("user already a member of organization")
	}

	member.UserID = user.ID

	// Add user as member
	return a.db.AddOrganizationMember(member)
}

// RemoveOrganizationMember removes a user from an organization
func (a *Api) RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error {
	return a.db.RemoveOrganizationMember(orgID, userID)
}

// GetOrganizationMembers returns all members of an organization
func (a *Api) GetOrganizationMembers(ctx context.Context, orgID string) ([]types.UserOrganization, error) {
	return a.db.GetOrganizationMembers(orgID)
}

// IsOrganizationOwner checks if a user is the owner of an organization
func (a *Api) IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error) {
	return a.db.IsOrganizationOwner(orgID, userID)
}

// GetOrganizations returns all organizations in the system
func (a *Api) GetOrganizations(ctx context.Context) ([]types.Organization, error) {
	return a.db.GetOrganizations()
}
