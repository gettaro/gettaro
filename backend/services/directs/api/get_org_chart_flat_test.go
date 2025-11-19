package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/directs/types"
	"github.com/stretchr/testify/assert"
)

func TestGetOrgChartFlat(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		orgID          string
		mockResult     []types.DirectReport
		mockError      error
		expectedResult []types.DirectReport
		expectedError  error
	}{
		{
			name:  "successful retrieval",
			orgID: "org-1",
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
				{
					ID:              "dr-3",
					ManagerMemberID: "manager-2",
					ReportMemberID:  "report-3",
					OrganizationID:  "org-1",
					Depth:           1,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			expectedResult: []types.DirectReport{
				{ID: "dr-1", ManagerMemberID: "manager-1", ReportMemberID: "report-1", OrganizationID: "org-1", Depth: 1},
				{ID: "dr-2", ManagerMemberID: "manager-1", ReportMemberID: "report-2", OrganizationID: "org-1", Depth: 1},
				{ID: "dr-3", ManagerMemberID: "manager-2", ReportMemberID: "report-3", OrganizationID: "org-1", Depth: 1},
			},
		},
		{
			name:           "empty org chart",
			orgID:          "org-empty",
			mockResult:     []types.DirectReport{},
			expectedResult: []types.DirectReport{},
		},
		{
			name:          "database error",
			orgID:         "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:  "multiple depths",
			orgID: "org-1",
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
					ManagerMemberID: "report-1",
					ReportMemberID:  "report-2",
					OrganizationID:  "org-1",
					Depth:           2,
				},
				{
					ID:              "dr-3",
					ManagerMemberID: "report-2",
					ReportMemberID:  "report-3",
					OrganizationID:  "org-1",
					Depth:           3,
				},
			},
			expectedResult: []types.DirectReport{
				{ID: "dr-1", ManagerMemberID: "manager-1", ReportMemberID: "report-1", OrganizationID: "org-1", Depth: 1},
				{ID: "dr-2", ManagerMemberID: "report-1", ReportMemberID: "report-2", OrganizationID: "org-1", Depth: 2},
				{ID: "dr-3", ManagerMemberID: "report-2", ReportMemberID: "report-3", OrganizationID: "org-1", Depth: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			mockDB.On("GetOrgChartFlat", ctx, tt.orgID).Return(tt.mockResult, tt.mockError)

			result, err := api.GetOrgChartFlat(ctx, tt.orgID)

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
