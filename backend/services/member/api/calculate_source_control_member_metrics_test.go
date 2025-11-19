package api

import (
	"context"
	"errors"
	"testing"
	"time"

	liberrors "ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/member/types"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateSourceControlMemberMetrics(t *testing.T) {
	ctx := context.Background()
	orgID := "org-1"
	memberID := "member-1"
	titleID := "title-1"
	accountID1 := "account-1"
	peerMemberID := "peer-member-1"
	peerAccountID := "peer-account-1"

	tests := []struct {
		name                    string
		organizationID          string
		memberID                string
		params                  sourcecontroltypes.MemberMetricsParams
		mockExternalAccounts    []types.ExternalAccount
		mockExternalError       error
		mockMember              *types.OrganizationMember
		mockMemberError         error
		mockOrgMembers          []types.OrganizationMember
		mockOrgMembersError     error
		mockPeerAccounts        []types.ExternalAccount
		mockPeerAccountsError   error
		mockMetricsResponse     *sourcecontroltypes.MetricsResponse
		mockMetricsError        error
		expectedResponse        *sourcecontroltypes.MetricsResponse
		expectedError           error
	}{
		{
			name:           "success",
			organizationID: orgID,
			memberID:       memberID,
			params: sourcecontroltypes.MemberMetricsParams{
				StartDate: timePtr(time.Now()),
				EndDate:   timePtr(time.Now()),
				Interval:  "daily",
			},
			mockExternalAccounts: []types.ExternalAccount{
				{ID: accountID1},
			},
			mockMember: &types.OrganizationMember{
				ID:      memberID,
				TitleID: &titleID,
			},
			mockOrgMembers: []types.OrganizationMember{
				{ID: memberID},
				{ID: peerMemberID},
			},
			mockPeerAccounts: []types.ExternalAccount{
				{ID: peerAccountID},
			},
			mockMetricsResponse: &sourcecontroltypes.MetricsResponse{},
			expectedResponse:    &sourcecontroltypes.MetricsResponse{},
		},
		{
			name:           "error - no external accounts",
			organizationID: orgID,
			memberID:       memberID,
			params: sourcecontroltypes.MemberMetricsParams{},
			mockExternalAccounts: []types.ExternalAccount{},
			expectedError:       liberrors.NewNotFoundError("no source control accounts found for member"),
		},
		{
			name:           "error - get external accounts fails",
			organizationID: orgID,
			memberID:       memberID,
			params: sourcecontroltypes.MemberMetricsParams{},
			mockExternalError: errors.New("database error"),
			expectedError:   errors.New("database error"),
		},
		{
			name:           "error - get member fails",
			organizationID: orgID,
			memberID:       memberID,
			params: sourcecontroltypes.MemberMetricsParams{},
			mockExternalAccounts: []types.ExternalAccount{
				{ID: accountID1},
			},
			mockMemberError: errors.New("database error"),
			expectedError:   errors.New("database error"),
		},
		{
			name:           "error - get org members fails",
			organizationID: orgID,
			memberID:       memberID,
			params: sourcecontroltypes.MemberMetricsParams{},
			mockExternalAccounts: []types.ExternalAccount{
				{ID: accountID1},
			},
			mockMember: &types.OrganizationMember{
				ID:      memberID,
				TitleID: &titleID,
			},
			mockOrgMembersError: errors.New("database error"),
			expectedError:       errors.New("database error"),
		},
		{
			name:           "error - calculate metrics fails",
			organizationID: orgID,
			memberID:       memberID,
			params: sourcecontroltypes.MemberMetricsParams{},
			mockExternalAccounts: []types.ExternalAccount{
				{ID: accountID1},
			},
			mockMember: &types.OrganizationMember{
				ID:      memberID,
				TitleID: &titleID,
			},
			mockOrgMembers: []types.OrganizationMember{
				{ID: memberID},
			},
			mockPeerAccounts: []types.ExternalAccount{},
			mockMetricsError: errors.New("metrics error"),
			expectedError:    errors.New("metrics error"),
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

			sourceControlType := "sourcecontrol"
			mockDB.On("GetExternalAccounts", ctx, &types.ExternalAccountParams{
				OrganizationID: tt.organizationID,
				MemberIDs:      []string{tt.memberID},
				AccountType:    &sourceControlType,
			}).Return(tt.mockExternalAccounts, tt.mockExternalError)

			if tt.mockExternalError == nil && len(tt.mockExternalAccounts) > 0 {
				mockDB.On("GetOrganizationMemberByID", ctx, tt.memberID).Return(tt.mockMember, tt.mockMemberError)

				if tt.mockMemberError == nil && tt.mockMember != nil {
					mockDB.On("GetOrganizationMembers", tt.organizationID, &types.OrganizationMemberParams{
						TitleIDs: []string{*tt.mockMember.TitleID},
					}).Return(tt.mockOrgMembers, tt.mockOrgMembersError)

					if tt.mockOrgMembersError == nil {
						// Mock GetMemberManager for each member (called by GetOrganizationMembers)
						for _, m := range tt.mockOrgMembers {
							mockDirectsAPI.On("GetMemberManager", ctx, m.ID, tt.organizationID).Return(nil, nil)
						}
						
						peerMemberIDs := []string{}
						for _, m := range tt.mockOrgMembers {
							peerMemberIDs = append(peerMemberIDs, m.ID)
						}
						mockDB.On("GetExternalAccounts", ctx, &types.ExternalAccountParams{
							OrganizationID: tt.organizationID,
							MemberIDs:      peerMemberIDs,
							AccountType:    &sourceControlType,
						}).Return(tt.mockPeerAccounts, tt.mockPeerAccountsError)

						if tt.mockPeerAccountsError == nil {
							mockSourceControlAPI.On("CalculateMetrics", ctx, mock.Anything).Return(tt.mockMetricsResponse, tt.mockMetricsError)
						}
					}
				}
			}

			result, err := api.CalculateSourceControlMemberMetrics(ctx, tt.organizationID, tt.memberID, tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockDB.AssertExpectations(t)
			mockSourceControlAPI.AssertExpectations(t)
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
