package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestGetOrganizationMemberByID(t *testing.T) {
	ctx := context.Background()
	memberID := "member-1"

	tests := []struct {
		name          string
		memberID      string
		mockMember    *types.OrganizationMember
		mockError     error
		expectedMember *types.OrganizationMember
		expectedError error
	}{
		{
			name:     "success",
			memberID: memberID,
			mockMember: &types.OrganizationMember{
				ID:             memberID,
				OrganizationID: "org-1",
			},
			expectedMember: &types.OrganizationMember{
				ID:             memberID,
				OrganizationID: "org-1",
			},
		},
		{
			name:          "not found",
			memberID:      memberID,
			mockMember:    nil,
			mockError:     nil,
			expectedMember: nil,
		},
		{
			name:          "database error",
			memberID:      memberID,
			mockError:     errors.New("database error"),
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

			mockDB.On("GetOrganizationMemberByID", ctx, tt.memberID).Return(tt.mockMember, tt.mockError)

			result, err := api.GetOrganizationMemberByID(ctx, tt.memberID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectedMember == nil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Equal(t, tt.expectedMember.ID, result.ID)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
