package api

import (
	"context"

	"ems.dev/backend/libraries/errors"
	memberdb "ems.dev/backend/services/member/database"
	"ems.dev/backend/services/member/types"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	titleapi "ems.dev/backend/services/title/api"
	titletypes "ems.dev/backend/services/title/types"
	userapi "ems.dev/backend/services/user/api"
	usertypes "ems.dev/backend/services/user/types"
)

// MemberAPI defines the interface for member operations
type MemberAPI interface {
	AddOrganizationMember(ctx context.Context, titleID string, sourceControlAccountID string, member *types.OrganizationMember) error
	RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error
	GetOrganizationMembers(ctx context.Context, orgID string) ([]types.OrganizationMember, error)
	IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error)
	UpdateOrganizationMember(ctx context.Context, orgID string, userID string, titleID string, sourceControlAccountID string, username string) error
}

type Api struct {
	db               memberdb.DB
	userApi          userapi.UserAPI
	titleApi         titleapi.TitleAPI
	sourceControlApi sourcecontrolapi.SourceControlAPI
}

func NewApi(memberDb memberdb.DB, userApi userapi.UserAPI, titleApi titleapi.TitleAPI, sourceControlApi sourcecontrolapi.SourceControlAPI) *Api {
	return &Api{
		db:               memberDb,
		userApi:          userApi,
		titleApi:         titleApi,
		sourceControlApi: sourceControlApi,
	}
}

// AddOrganizationMember adds a user as a member to an organization
func (a *Api) AddOrganizationMember(ctx context.Context, titleID string, sourceControlAccountID string, member *types.OrganizationMember) error {
	// Check if the title exists
	title, err := a.titleApi.GetTitle(ctx, titleID)
	if err != nil {
		return errors.NewNotFoundError("title not found")
	}

	// Get the source control account
	sourceControlAccount, err := a.sourceControlApi.GetSourceControlAccount(ctx, sourceControlAccountID)
	if err != nil {
		return errors.NewNotFoundError("source control account not found")
	}

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

	err = a.titleApi.AssignUserTitle(ctx, titletypes.UserTitle{
		TitleID:        title.ID,
		UserID:         user.ID,
		OrganizationID: member.OrganizationID,
	})
	if err != nil {
		return err
	}

	// Add user as member first
	err = a.db.AddOrganizationMember(member)
	if err != nil {
		return err
	}

	// Now get the member ID and update the source control account
	createdMember, err := a.db.GetOrganizationMember(member.OrganizationID, user.ID)
	if err != nil {
		return err
	}

	if createdMember != nil {
		sourceControlAccount.MemberID = &createdMember.ID
		err = a.sourceControlApi.UpdateSourceControlAccount(ctx, sourceControlAccount)
		if err != nil {
			return err
		}
	}

	return nil
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

// UpdateOrganizationMember updates a member's details in an organization
func (a *Api) UpdateOrganizationMember(ctx context.Context, orgID string, userID string, titleID string, sourceControlAccountID string, username string) error {
	// Check if the title exists
	title, err := a.titleApi.GetTitle(ctx, titleID)
	if err != nil {
		return errors.NewNotFoundError("title not found")
	}

	// Get the source control account
	sourceControlAccount, err := a.sourceControlApi.GetSourceControlAccount(ctx, sourceControlAccountID)
	if err != nil {
		return errors.NewNotFoundError("source control account not found")
	}

	// Check if the member exists
	existingMember, err := a.db.GetOrganizationMember(orgID, userID)
	if err != nil {
		return err
	}
	if existingMember == nil {
		return errors.NewNotFoundError("member not found")
	}

	// Update the username in the organization_members table
	err = a.db.UpdateOrganizationMember(orgID, userID, username)
	if err != nil {
		return err
	}

	// Update the title assignment
	err = a.titleApi.RemoveUserTitle(ctx, userID, orgID)
	if err != nil {
		return err
	}

	err = a.titleApi.AssignUserTitle(ctx, titletypes.UserTitle{
		TitleID:        title.ID,
		UserID:         userID,
		OrganizationID: orgID,
	})
	if err != nil {
		return err
	}

	// Update the source control account with member_id
	sourceControlAccount.MemberID = &existingMember.ID
	err = a.sourceControlApi.UpdateSourceControlAccount(ctx, sourceControlAccount)
	if err != nil {
		return err
	}

	return nil
}
