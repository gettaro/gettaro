package api

import (
	"context"
	"encoding/json"

	"ems.dev/backend/libraries/errors"
	directsapi "ems.dev/backend/services/directs/api"
	directstypes "ems.dev/backend/services/directs/types"
	memberdb "ems.dev/backend/services/member/database"
	"ems.dev/backend/services/member/types"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
	titleapi "ems.dev/backend/services/title/api"
	userapi "ems.dev/backend/services/user/api"
	usertypes "ems.dev/backend/services/user/types"
	"gorm.io/datatypes"
)

// MemberAPI defines the interface for member operations
type MemberAPI interface {
	AddOrganizationMember(ctx context.Context, req types.AddMemberRequest, member *types.OrganizationMember) error
	RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error
	GetOrganizationMembers(ctx context.Context, orgID string, params *types.OrganizationMemberParams) ([]types.OrganizationMember, error)
	GetOrganizationMemberByID(ctx context.Context, memberID string) (*types.OrganizationMember, error)
	IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error)
	UpdateOrganizationMember(ctx context.Context, orgID string, memberID string, req types.UpdateMemberRequest) error
	CalculateSourceControlMemberMetrics(ctx context.Context, organizationID string, memberID string, params sourcecontroltypes.MemberMetricsParams) (*sourcecontroltypes.MetricsResponse, error)
}

type Api struct {
	db               memberdb.DB
	userApi          userapi.UserAPI
	sourceControlApi sourcecontrolapi.SourceControlAPI
	titleApi         titleapi.TitleAPI
	directsApi       directsapi.DirectReportsAPI
}

func NewApi(memberDb memberdb.DB, userApi userapi.UserAPI, sourceControlApi sourcecontrolapi.SourceControlAPI, titleApi titleapi.TitleAPI, directsApi directsapi.DirectReportsAPI) *Api {
	return &Api{
		db:               memberDb,
		userApi:          userApi,
		sourceControlApi: sourceControlApi,
		titleApi:         titleApi,
		directsApi:       directsApi,
	}
}

// AddOrganizationMember adds a user as a member to an organization
func (a *Api) AddOrganizationMember(ctx context.Context, req types.AddMemberRequest, member *types.OrganizationMember) error {
	// Note: Title validation will be handled by the database foreign key constraint
	// Get the source control account
	sourceControlAccount, err := a.sourceControlApi.GetSourceControlAccount(ctx, req.SourceControlAccountID)
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
	member.TitleID = &req.TitleID // Set the title ID directly

	// Add user as member
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

		// Create manager relationship if specified
		if req.ManagerID != nil && *req.ManagerID != "" {
			// Get the manager's member record
			managerMember, err := a.db.GetOrganizationMemberByID(ctx, *req.ManagerID)
			if err != nil {
				return err
			}
			if managerMember == nil {
				return errors.NewNotFoundError("manager not found")
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
				return err
			}
		}
	}

	return nil
}

// RemoveOrganizationMember removes a user from an organization
func (a *Api) RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error {
	return a.db.RemoveOrganizationMember(orgID, userID)
}

// GetOrganizationMembers returns all members of an organization
func (a *Api) GetOrganizationMembers(ctx context.Context, orgID string, params *types.OrganizationMemberParams) ([]types.OrganizationMember, error) {
	return a.db.GetOrganizationMembers(orgID, params)
}

// GetOrganizationMemberByID retrieves a member by their ID
func (a *Api) GetOrganizationMemberByID(ctx context.Context, memberID string) (*types.OrganizationMember, error) {
	return a.db.GetOrganizationMemberByID(ctx, memberID)
}

// IsOrganizationOwner checks if a user is the owner of an organization
func (a *Api) IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error) {
	return a.db.IsOrganizationOwner(orgID, userID)
}

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

	// Verify the new source control account exists and belongs to the organization
	newSourceControlAccount, err := a.sourceControlApi.GetSourceControlAccount(ctx, req.SourceControlAccountID)
	if err != nil {
		return errors.NewNotFoundError("source control account not found")
	}

	// Check if the source control account belongs to the organization
	if newSourceControlAccount.OrganizationID == nil || *newSourceControlAccount.OrganizationID != orgID {
		return errors.NewNotFoundError("source control account does not belong to this organization")
	}

	// Update the username and title_id in the organization_members table
	err = a.db.UpdateOrganizationMember(orgID, existingMember.UserID, req.Username, &req.TitleID)
	if err != nil {
		return err
	}

	// First, remove the member_id from any existing source control accounts for this member
	// This ensures we don't have multiple accounts pointing to the same member
	existingAccounts, err := a.sourceControlApi.GetSourceControlAccounts(ctx, &sourcecontroltypes.SourceControlAccountParams{
		OrganizationID: orgID,
	})
	if err != nil {
		return err
	}

	for _, account := range existingAccounts {
		if account.MemberID != nil && *account.MemberID == memberID {
			// Clear the member_id from this account
			account.MemberID = nil
			err = a.sourceControlApi.UpdateSourceControlAccount(ctx, &account)
			if err != nil {
				return err
			}
		}
	}

	// Now assign the new source control account to this member
	newSourceControlAccount.MemberID = &memberID
	err = a.sourceControlApi.UpdateSourceControlAccount(ctx, newSourceControlAccount)
	if err != nil {
		return err
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

// CalculateSourceControlMemberMetrics retrieves source control metrics for a specific member
func (a *Api) CalculateSourceControlMemberMetrics(ctx context.Context, organizationID string, memberID string, params sourcecontroltypes.MemberMetricsParams) (*sourcecontroltypes.MetricsResponse, error) {
	sourceControlAccounts, err := a.sourceControlApi.GetSourceControlAccounts(ctx, &sourcecontroltypes.SourceControlAccountParams{
		OrganizationID: organizationID,
		MemberIDs:      []string{memberID},
	})
	if err != nil {
		return nil, err
	}

	sourceControlAccountIDs := []string{}
	for _, account := range sourceControlAccounts {
		sourceControlAccountIDs = append(sourceControlAccountIDs, account.ID)
	}

	member, err := a.GetOrganizationMemberByID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	orgMembers, err := a.GetOrganizationMembers(ctx, organizationID, &types.OrganizationMemberParams{
		TitleIDs: []string{*member.TitleID},
	})
	if err != nil {
		return nil, err
	}

	peerMemberIDs := []string{}
	for _, orgMember := range orgMembers {
		peerMemberIDs = append(peerMemberIDs, orgMember.ID)
	}

	peerSourceControlAccounts, err := a.sourceControlApi.GetSourceControlAccounts(ctx, &sourcecontroltypes.SourceControlAccountParams{
		OrganizationID: organizationID,
		MemberIDs:      peerMemberIDs,
	})
	if err != nil {
		return nil, err
	}

	peerSourceControlAccountIDs := []string{}
	for _, account := range peerSourceControlAccounts {
		peerSourceControlAccountIDs = append(peerSourceControlAccountIDs, account.ID)
	}
	// Create the metric params with the source control account IDs
	metricParamsMap := map[string]interface{}{
		"organizationId":               organizationID,
		"sourceControlAccountIDs":      sourceControlAccountIDs,
		"peersSourceControlAccountIDs": peerSourceControlAccountIDs,
	}

	// Marshal to JSON bytes
	metricParamsJSON, err := json.Marshal(metricParamsMap)
	if err != nil {
		return nil, err
	}

	metricParams := sourcecontroltypes.MetricRuleParams{
		MetricParams: datatypes.JSON(metricParamsJSON),
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Interval:     params.Interval,
	}

	return a.sourceControlApi.CalculateMetrics(ctx, metricParams)
}
