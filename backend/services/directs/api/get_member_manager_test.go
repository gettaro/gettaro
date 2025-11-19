package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/directs/types"
	membertypes "ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestGetMemberManager(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		reportID       string
		orgID          string
		mockResult     *types.DirectReport
		mockError      error
		expectedResult *types.DirectReport
		expectedError  error
	}{
		{
			name:     "successful retrieval",
			reportID: "report-1",
			orgID:    "org-1",
			mockResult: &types.DirectReport{
				ID:              "dr-1",
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Manager: membertypes.OrganizationMember{
					ID:             "manager-1",
					UserID:         "user-1",
					Email:          "manager@example.com",
					OrganizationID: "org-1",
				},
				Report: membertypes.OrganizationMember{
					ID:             "report-1",
					UserID:         "user-2",
					Email:          "report@example.com",
					OrganizationID: "org-1",
				},
			},
			expectedResult: &types.DirectReport{
				ID:              "dr-1",
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
			},
		},
		{
			name:           "member has no manager",
			reportID:       "report-top-level",
			orgID:          "org-1",
			mockResult:     nil,
			expectedResult: nil,
		},
		{
			name:          "database error",
			reportID:      "report-1",
			orgID:         "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			mockDB.On("GetMemberManager", ctx, tt.reportID, tt.orgID).Return(tt.mockResult, tt.mockError)

			result, err := api.GetMemberManager(ctx, tt.reportID, tt.orgID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectedResult == nil {
					assert.Nil(t, result)
				} else {
					assert.Equal(t, tt.expectedResult.ID, result.ID)
					assert.Equal(t, tt.expectedResult.ManagerMemberID, result.ManagerMemberID)
					assert.Equal(t, tt.expectedResult.ReportMemberID, result.ReportMemberID)
					assert.Equal(t, tt.expectedResult.OrganizationID, result.OrganizationID)
					assert.Equal(t, tt.expectedResult.Depth, result.Depth)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
