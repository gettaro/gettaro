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

func TestGetDirectReport(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		id             string
		mockResult     *types.DirectReport
		mockError      error
		expectedResult *types.DirectReport
		expectedError  error
	}{
		{
			name: "successful retrieval",
			id:   "dr-1",
			mockResult: &types.DirectReport{
				ID:              "dr-1",
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
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
			name:           "not found",
			id:             "dr-nonexistent",
			mockResult:     nil,
			mockError:      nil,
			expectedResult: nil,
		},
		{
			name:          "database error",
			id:            "dr-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name: "with relations",
			id:   "dr-1",
			mockResult: &types.DirectReport{
				ID:              "dr-1",
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			mockDB.On("GetDirectReport", ctx, tt.id).Return(tt.mockResult, tt.mockError)

			result, err := api.GetDirectReport(ctx, tt.id)

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
