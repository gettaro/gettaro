package api

import (
	"context"

	"ems.dev/backend/libraries/errors"
	directstypes "ems.dev/backend/services/directs/types"
	"ems.dev/backend/services/member/types"
	usertypes "ems.dev/backend/services/user/types"
)

// AddOrganizationMember adds a user as a member to an organization
func (a *Api) AddOrganizationMember(ctx context.Context, req types.AddMemberRequest, member *types.OrganizationMember) (*types.OrganizationMember, error) {
	// Note: Title validation will be handled by the database foreign key constraint

	// Look up user by email
	user, err := a.userApi.FindUser(usertypes.UserSearchParams{Email: &member.Email})
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = a.userApi.CreateUser(&usertypes.User{
			Email: member.Email,
		})

		if err != nil {
			return nil, err
		}
	}

	// Check for duplicate member
	existingMember, err := a.db.GetOrganizationMember(member.OrganizationID, user.ID)
	if err != nil {
		return nil, err
	}

	if existingMember != nil {
		return nil, errors.NewConflictError("user already a member of organization")
	}

	member.UserID = user.ID
	member.TitleID = &req.TitleID // Set the title ID directly

	// Add user as member
	err = a.db.AddOrganizationMember(member)
	if err != nil {
		return nil, err
	}

	// Now get the member ID
	createdMember, err := a.db.GetOrganizationMember(member.OrganizationID, user.ID)
	if err != nil {
		return nil, err
	}

	if createdMember != nil {
		// Update external account if provided
		if req.ExternalAccountID != "" {
			externalAccount, err := a.GetExternalAccount(ctx, req.ExternalAccountID)
			if err != nil {
				return nil, err
			}
			if externalAccount == nil {
				return nil, errors.NewNotFoundError("external account not found")
			}

			externalAccount.MemberID = &createdMember.ID
			err = a.UpdateExternalAccount(ctx, externalAccount)
			if err != nil {
				return nil, err
			}
		}

		// Create manager relationship if specified
		if req.ManagerID != nil && *req.ManagerID != "" {
			// Get the manager's member record
			managerMember, err := a.db.GetOrganizationMemberByID(ctx, *req.ManagerID)
			if err != nil {
				return nil, err
			}
			if managerMember == nil {
				return nil, errors.NewNotFoundError("manager not found")
			}

			// Create direct report relationship using member IDs
			_, err = a.directsApi.CreateDirectReport(ctx, directstypes.CreateDirectReportParams{
				ManagerMemberID: managerMember.ID,
				ReportMemberID:  createdMember.ID,
				OrganizationID:  member.OrganizationID,
				Depth:           1, // Direct report
			})
			if err != nil {
				// Log the error but don't fail the member creation
				// The member is already created, just the manager relationship failed
				return nil, err
			}
		}
	}

	return createdMember, nil
}
