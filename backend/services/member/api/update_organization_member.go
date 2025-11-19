package api

import (
	"context"

	"ems.dev/backend/libraries/errors"
	directstypes "ems.dev/backend/services/directs/types"
	"ems.dev/backend/services/member/types"
)

// UpdateOrganizationMember updates a member's details in an organization
func (a *Api) UpdateOrganizationMember(ctx context.Context, orgID string, memberID string, req types.UpdateMemberRequest) error {
	// Check if the member exists by member ID
	existingMember, err := a.db.GetOrganizationMemberByID(ctx, memberID)
	if err != nil {
		return err
	}
	if existingMember == nil {
		return errors.NewNotFoundError("member not found")
	}

	// Verify the member belongs to the specified organization
	if existingMember.OrganizationID != orgID {
		return errors.NewNotFoundError("member not found in this organization")
	}

	// Update the username and title_id in the organization_members table
	err = a.db.UpdateOrganizationMember(orgID, existingMember.UserID, req.Username, &req.TitleID)
	if err != nil {
		return err
	}

	// Update external account if provided
	if req.ExternalAccountID != "" {
		// Verify the new external account exists and belongs to the organization
		newExternalAccount, err := a.GetExternalAccount(ctx, req.ExternalAccountID)
		if err != nil {
			return err
		}
		if newExternalAccount == nil {
			return errors.NewNotFoundError("external account not found")
		}

		// Check if the external account belongs to the organization
		if newExternalAccount.OrganizationID == nil || *newExternalAccount.OrganizationID != orgID {
			return errors.NewNotFoundError("external account does not belong to this organization")
		}

		// First, remove the member_id from any existing external accounts for this member
		// This ensures we don't have multiple accounts pointing to the same member
		// Filter by account type sourcecontrol since that's what we're managing here
		sourceControlType := "sourcecontrol"
		existingAccounts, err := a.GetExternalAccounts(ctx, &types.ExternalAccountParams{
			OrganizationID: orgID,
			AccountType:    &sourceControlType,
		})
		if err != nil {
			return err
		}

		for _, account := range existingAccounts {
			if account.MemberID != nil && *account.MemberID == memberID {
				// Clear the member_id from this account
				account.MemberID = nil
				err = a.UpdateExternalAccount(ctx, &account)
				if err != nil {
					return err
				}
			}
		}

		// Now assign the new external account to this member
		newExternalAccount.MemberID = &memberID
		err = a.UpdateExternalAccount(ctx, newExternalAccount)
		if err != nil {
			return err
		}
	}

	// Update manager relationship if specified
	if req.ManagerID != nil && *req.ManagerID != "" {
		// Get the manager's member record
		managerMember, err := a.db.GetOrganizationMemberByID(ctx, *req.ManagerID)
		if err != nil {
			return err
		}
		if managerMember == nil {
			return errors.NewNotFoundError("manager not found")
		}

		// Check if there's an existing manager relationship using member ID
		existingManager, err := a.directsApi.GetMemberManager(ctx, existingMember.ID, orgID)
		if err != nil {
			return err
		}

		if existingManager != nil {
			// Check if the manager is actually changing
			if existingManager.ManagerMemberID != managerMember.ID {
				// Manager is changing, delete old relationship and create new one
				err = a.directsApi.DeleteDirectReport(ctx, existingManager.ID)
				if err != nil {
					return err
				}
				// Create new relationship with new manager
				_, err = a.directsApi.CreateDirectReport(ctx, directstypes.CreateDirectReportParams{
					ManagerMemberID: managerMember.ID,
					ReportMemberID:  existingMember.ID,
					OrganizationID:  orgID,
					Depth:           1, // Direct report
				})
				if err != nil {
					return err
				}
			}
			// If manager is the same, no action needed
		} else {
			// Create new manager relationship
			_, err = a.directsApi.CreateDirectReport(ctx, directstypes.CreateDirectReportParams{
				ManagerMemberID: managerMember.ID,
				ReportMemberID:  existingMember.ID,
				OrganizationID:  orgID,
				Depth:           1, // Direct report
			})
			if err != nil {
				return err
			}
		}
	} else {
		// Remove manager relationship if ManagerID is empty
		existingManager, err := a.directsApi.GetMemberManager(ctx, existingMember.ID, orgID)
		if err != nil {
			return err
		}
		if existingManager != nil {
			err = a.directsApi.DeleteDirectReport(ctx, existingManager.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
