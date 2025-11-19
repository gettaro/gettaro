package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/directs/types"
	"github.com/stretchr/testify/assert"
)

func TestListDirectReports(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		params         types.DirectReportSearchParams
		mockResult     []types.DirectReport
		mockError      error
		expectedResult []types.DirectReport
		expectedError  error
	}{
		{
			name: "successful list all",
			params: types.DirectReportSearchParams{},
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
			name: "filter by manager",
			params: types.DirectReportSearchParams{
				ManagerMemberID: stringPtr("manager-1"),
			},
			mockResult: []types.DirectReport{
				{
					ID:              "dr-1",
					ManagerMemberID: "manager-1",
					ReportMemberID:  "report-1",
					OrganizationID:  "org-1",
					Depth:           1,
				},
			},
			expectedResult: []types.DirectReport{
				{ID: "dr-1", ManagerMemberID: "manager-1", ReportMemberID: "report-1", OrganizationID: "org-1", Depth: 1},
			},
		},
		{
			name: "filter by organization",
			params: types.DirectReportSearchParams{
				OrganizationID: stringPtr("org-1"),
			},
			mockResult: []types.DirectReport{
				{
					ID:              "dr-1",
					ManagerMemberID: "manager-1",
					ReportMemberID:  "report-1",
					OrganizationID:  "org-1",
					Depth:           1,
				},
			},
			expectedResult: []types.DirectReport{
				{ID: "dr-1", ManagerMemberID: "manager-1", ReportMemberID: "report-1", OrganizationID: "org-1", Depth: 1},
			},
		},
		{
			name: "empty result",
			params: types.DirectReportSearchParams{
				ManagerMemberID: stringPtr("manager-nonexistent"),
			},
			mockResult:     []types.DirectReport{},
			expectedResult: []types.DirectReport{},
		},
		{
			name: "database error",
			params: types.DirectReportSearchParams{},
			mockError:      errors.New("database error"),
			expectedError:  errors.New("database error"),
		},
		{
			name: "filter by depth",
			params: types.DirectReportSearchParams{
				Depth: intPtr(1),
			},
			mockResult: []types.DirectReport{
				{
					ID:              "dr-1",
					ManagerMemberID: "manager-1",
					ReportMemberID:  "report-1",
					OrganizationID:  "org-1",
					Depth:           1,
				},
			},
			expectedResult: []types.DirectReport{
				{ID: "dr-1", ManagerMemberID: "manager-1", ReportMemberID: "report-1", OrganizationID: "org-1", Depth: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			mockDB.On("ListDirectReports", ctx, tt.params).Return(tt.mockResult, tt.mockError)

			result, err := api.ListDirectReports(ctx, tt.params)

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

