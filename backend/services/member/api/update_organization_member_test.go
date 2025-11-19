package api

import (
	"context"
	"errors"
	"testing"

	liberrors "ems.dev/backend/libraries/errors"
	directstypes "ems.dev/backend/services/directs/types"
	"ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateOrganizationMember(t *testing.T) {
	ctx := context.Background()
	orgID := "org-1"
	memberID := "member-1"
	userID := "user-1"
	titleID := "title-1"
	managerID := "manager-member-1"
	externalAccountID := "ext-account-1"

	tests := []struct {
		name                string
		orgID               string
		memberID            string
		req                 types.UpdateMemberRequest
		mockMember          *types.OrganizationMember
		mockMemberError     error
		mockUpdateError     error
		mockExternalAccount *types.ExternalAccount
		mockExternalError   error
		mockExistingAccounts []types.ExternalAccount
		mockAccountsError   error
		mockManagerMember   *types.OrganizationMember
		mockManagerError    error
		mockExistingManager *directstypes.DirectReport
		mockManagerGetError error
		mockCreateDirectError error
		mockDeleteDirectError error
		expectedError       error
	}{
		{
			name:     "success - basic update",
			orgID:    orgID,
			memberID: memberID,
			req: types.UpdateMemberRequest{
				Username: "newusername",
				TitleID:  titleID,
			},
			mockMember: &types.OrganizationMember{
				ID:             memberID,
				UserID:         userID,
				OrganizationID: orgID,
			},
		},
		{
			name:     "error - member not found",
			orgID:    orgID,
			memberID: memberID,
			req: types.UpdateMemberRequest{
				Username: "newusername",
				TitleID:  titleID,
			},
			mockMember:      nil,
			expectedError:   liberrors.NewNotFoundError("member not found"),
		},
		{
			name:     "error - member belongs to different org",
			orgID:    orgID,
			memberID: memberID,
			req: types.UpdateMemberRequest{
				Username: "newusername",
				TitleID:  titleID,
			},
			mockMember: &types.OrganizationMember{
				ID:             memberID,
				UserID:         userID,
				OrganizationID: "different-org",
			},
			expectedError:   liberrors.NewNotFoundError("member not found in this organization"),
		},
		{
			name:     "error - database update fails",
			orgID:    orgID,
			memberID: memberID,
			req: types.UpdateMemberRequest{
				Username: "newusername",
				TitleID:  titleID,
			},
			mockMember: &types.OrganizationMember{
				ID:             memberID,
				UserID:         userID,
				OrganizationID: orgID,
			},
			mockUpdateError: errors.New("database error"),
			expectedError:   errors.New("database error"),
		},
		{
			name:     "success - with external account",
			orgID:    orgID,
			memberID: memberID,
			req: types.UpdateMemberRequest{
				Username:           "newusername",
				TitleID:            titleID,
				ExternalAccountID:  externalAccountID,
			},
			mockMember: &types.OrganizationMember{
				ID:             memberID,
				UserID:         userID,
				OrganizationID: orgID,
			},
			mockExternalAccount: &types.ExternalAccount{
				ID:             externalAccountID,
				OrganizationID: &orgID,
			},
			mockExistingAccounts: []types.ExternalAccount{},
		},
		{
			name:     "success - with manager",
			orgID:    orgID,
			memberID: memberID,
			req: types.UpdateMemberRequest{
				Username: "newusername",
				TitleID:  titleID,
				ManagerID: stringPtr(managerID),
			},
			mockMember: &types.OrganizationMember{
				ID:             memberID,
				UserID:         userID,
				OrganizationID: orgID,
			},
			mockManagerMember: &types.OrganizationMember{
				ID: managerID,
			},
			mockExistingManager: nil,
		},
		{
			name:     "success - change manager",
			orgID:    orgID,
			memberID: memberID,
			req: types.UpdateMemberRequest{
				Username: "newusername",
				TitleID:  titleID,
				ManagerID: stringPtr(managerID),
			},
			mockMember: &types.OrganizationMember{
				ID:             memberID,
				UserID:         userID,
				OrganizationID: orgID,
			},
			mockManagerMember: &types.OrganizationMember{
				ID: managerID,
			},
			mockExistingManager: &directstypes.DirectReport{
				ID:              "direct-1",
				ManagerMemberID: "old-manager",
			},
		},
		{
			name:     "success - remove manager",
			orgID:    orgID,
			memberID: memberID,
			req: types.UpdateMemberRequest{
				Username: "newusername",
				TitleID:  titleID,
				ManagerID: nil,
			},
			mockMember: &types.OrganizationMember{
				ID:             memberID,
				UserID:         userID,
				OrganizationID: orgID,
			},
			mockExistingManager: &directstypes.DirectReport{
				ID: "direct-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockMemberDB)
			mockUserAPI := new(MockUserAPI)
			mockDirectsAPI := new(MockDirectReportsAPI)
			mockSourceControlAPI := new(MockSourceControlAPI)
			mockTitleAPI := new(MockTitleAPI)

			api := NewApi(mockDB, mockUserAPI, mockSourceControlAPI, mockTitleAPI, mockDirectsAPI)

			mockDB.On("GetOrganizationMemberByID", ctx, tt.memberID).Return(tt.mockMember, tt.mockMemberError)

			if tt.mockMember != nil && tt.mockMember.OrganizationID == tt.orgID {
				mockDB.On("UpdateOrganizationMember", tt.orgID, tt.mockMember.UserID, tt.req.Username, &tt.req.TitleID).Return(tt.mockUpdateError)
			}

			// Handle external account
			if tt.req.ExternalAccountID != "" && tt.mockMember != nil && tt.mockMember.OrganizationID == tt.orgID && tt.mockUpdateError == nil {
				mockDB.On("GetExternalAccount", ctx, tt.req.ExternalAccountID).Return(tt.mockExternalAccount, tt.mockExternalError)
				if tt.mockExternalAccount != nil && tt.mockExternalAccount.OrganizationID != nil && *tt.mockExternalAccount.OrganizationID == tt.orgID {
					mockDB.On("GetExternalAccounts", ctx, mock.AnythingOfType("*types.ExternalAccountParams")).Return(tt.mockExistingAccounts, tt.mockAccountsError)
					if tt.mockAccountsError == nil {
						mockDB.On("UpdateExternalAccount", ctx, mock.AnythingOfType("*types.ExternalAccount")).Return(nil).Maybe()
					}
				}
			}

			// Handle manager
			if tt.req.ManagerID != nil && *tt.req.ManagerID != "" && tt.mockMember != nil && tt.mockMember.OrganizationID == tt.orgID && tt.mockUpdateError == nil {
				mockDB.On("GetOrganizationMemberByID", ctx, *tt.req.ManagerID).Return(tt.mockManagerMember, tt.mockManagerError)
				if tt.mockManagerMember != nil {
					mockDirectsAPI.On("GetMemberManager", ctx, tt.memberID, tt.orgID).Return(tt.mockExistingManager, tt.mockManagerGetError)
					if tt.mockExistingManager != nil && tt.mockExistingManager.ManagerMemberID != tt.mockManagerMember.ID {
						mockDirectsAPI.On("DeleteDirectReport", ctx, tt.mockExistingManager.ID).Return(tt.mockDeleteDirectError)
						if tt.mockDeleteDirectError == nil {
							mockDirectsAPI.On("CreateDirectReport", ctx, mock.MatchedBy(func(params directstypes.CreateDirectReportParams) bool {
								return params.ManagerMemberID == tt.mockManagerMember.ID && params.ReportMemberID == tt.memberID
							})).Return(&directstypes.DirectReport{}, tt.mockCreateDirectError)
						}
					} else if tt.mockExistingManager == nil {
						mockDirectsAPI.On("CreateDirectReport", ctx, mock.MatchedBy(func(params directstypes.CreateDirectReportParams) bool {
							return params.ManagerMemberID == tt.mockManagerMember.ID && params.ReportMemberID == tt.memberID
						})).Return(&directstypes.DirectReport{}, tt.mockCreateDirectError)
					}
				}
			} else if tt.req.ManagerID == nil && tt.mockMember != nil && tt.mockMember.OrganizationID == tt.orgID && tt.mockUpdateError == nil {
				mockDirectsAPI.On("GetMemberManager", ctx, tt.memberID, tt.orgID).Return(tt.mockExistingManager, tt.mockManagerGetError)
				if tt.mockExistingManager != nil {
					mockDirectsAPI.On("DeleteDirectReport", ctx, tt.mockExistingManager.ID).Return(tt.mockDeleteDirectError)
				}
			}

			err := api.UpdateOrganizationMember(ctx, tt.orgID, tt.memberID, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
			mockDirectsAPI.AssertExpectations(t)
		})
	}
}
