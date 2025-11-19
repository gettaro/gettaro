package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/directs/types"
	"github.com/stretchr/testify/assert"
)

func TestGetManagerDirectReports(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		managerID      string
		orgID          string
		mockResult     []types.DirectReport
		mockError      error
		expectedResult []types.DirectReport
		expectedError  error
	}{
		{
			name:      "successful retrieval",
			managerID: "manager-1",
			orgID:     "org-1",
			mockResult: []types.DirectReport{
				{
					ID:              "dr-1",
					ManagerMemberID: "manager-1",
					ReportMemberID:  "report-1",
					OrganizationID:  "org-1",
					Depth:           1,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
				{
					ID:              "dr-2",
					ManagerMemberID: "manager-1",
					ReportMemberID:  "report-2",
					OrganizationID:  "org-1",
					Depth:           1,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			expectedResult: []types.DirectReport{
				{ID: "dr-1", ManagerMemberID: "manager-1", ReportMemberID: "report-1", OrganizationID: "org-1", Depth: 1},
				{ID: "dr-2", ManagerMemberID: "manager-1", ReportMemberID: "report-2", OrganizationID: "org-1", Depth: 1},
			},
		},
		{
			name:           "empty result",
			managerID:      "manager-nonexistent",
			orgID:          "org-1",
			mockResult:     []types.DirectReport{},
			expectedResult: []types.DirectReport{},
		},
		{
			name:          "database error",
			managerID:     "manager-1",
			orgID:         "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:      "multiple depths",
			managerID: "manager-1",
			orgID:     "org-1",
			mockResult: []types.DirectReport{
				{
					ID:              "dr-1",
					ManagerMemberID: "manager-1",
					ReportMemberID:  "report-1",
					OrganizationID:  "org-1",
					Depth:           1,
				},
				{
					ID:              "dr-2",
					ManagerMemberID: "manager-1",
					ReportMemberID:  "report-2",
					OrganizationID:  "org-1",
					Depth:           2,
				},
			},
			expectedResult: []types.DirectReport{
				{ID: "dr-1", ManagerMemberID: "manager-1", ReportMemberID: "report-1", OrganizationID: "org-1", Depth: 1},
				{ID: "dr-2", ManagerMemberID: "manager-1", ReportMemberID: "report-2", OrganizationID: "org-1", Depth: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			mockDB.On("GetManagerDirectReports", ctx, tt.managerID, tt.orgID).Return(tt.mockResult, tt.mockError)

			result, err := api.GetManagerDirectReports(ctx, tt.managerID, tt.orgID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(result))
				for i, expected := range tt.expectedResult {
					assert.Equal(t, expected.ID, result[i].ID)
					assert.Equal(t, expected.ManagerMemberID, result[i].ManagerMemberID)
					assert.Equal(t, expected.ReportMemberID, result[i].ReportMemberID)
					assert.Equal(t, expected.OrganizationID, result[i].OrganizationID)
					assert.Equal(t, expected.Depth, result[i].Depth)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
