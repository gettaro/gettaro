package api

import (
	"context"
	"errors"
	"testing"

	directstypes "ems.dev/backend/services/directs/types"
	"ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestGetOrganizationMembers(t *testing.T) {
	ctx := context.Background()
	orgID := "org-1"
	memberID1 := "member-1"
	memberID2 := "member-2"

	tests := []struct {
		name          string
		orgID         string
		params        *types.OrganizationMemberParams
		mockMembers   []types.OrganizationMember
		mockError     error
		mockManager   *directstypes.DirectReport
		mockManagerError error
		expectedMembers []types.OrganizationMember
		expectedError error
	}{
		{
			name:  "success - empty list",
			orgID:  orgID,
			params: nil,
			mockMembers: []types.OrganizationMember{},
			expectedMembers: []types.OrganizationMember{},
		},
		{
			name:  "success - with members",
			orgID:  orgID,
			params: nil,
			mockMembers: []types.OrganizationMember{
				{ID: memberID1, OrganizationID: orgID},
				{ID: memberID2, OrganizationID: orgID},
			},
			mockManager: nil,
			expectedMembers: []types.OrganizationMember{
				{ID: memberID1, OrganizationID: orgID},
				{ID: memberID2, OrganizationID: orgID},
			},
		},
		{
			name:  "success - with manager",
			orgID:  orgID,
			params: nil,
			mockMembers: []types.OrganizationMember{
				{ID: memberID1, OrganizationID: orgID},
			},
			mockManager: &directstypes.DirectReport{
				ManagerMemberID: "manager-1",
			},
			expectedMembers: []types.OrganizationMember{
				{ID: memberID1, OrganizationID: orgID, ManagerID: stringPtr("manager-1")},
			},
		},
		{
			name:        "database error",
			orgID:       orgID,
			params:      nil,
			mockError:   errors.New("database error"),
			expectedError: errors.New("database error"),
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

			mockDB.On("GetOrganizationMembers", tt.orgID, tt.params).Return(tt.mockMembers, tt.mockError)

			if tt.mockError == nil {
				for _, member := range tt.mockMembers {
					if tt.mockManager != nil && member.ID == memberID1 {
						mockDirectsAPI.On("GetMemberManager", ctx, member.ID, orgID).Return(tt.mockManager, nil)
					} else {
						mockDirectsAPI.On("GetMemberManager", ctx, member.ID, orgID).Return(nil, nil)
					}
				}
			}

			result, err := api.GetOrganizationMembers(ctx, tt.orgID, tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedMembers), len(result))
				if len(tt.expectedMembers) > 0 {
					assert.Equal(t, tt.expectedMembers[0].ID, result[0].ID)
					if tt.expectedMembers[0].ManagerID != nil {
						assert.Equal(t, *tt.expectedMembers[0].ManagerID, *result[0].ManagerID)
					}
				}
			}

			mockDB.AssertExpectations(t)
			mockDirectsAPI.AssertExpectations(t)
		})
	}
}
