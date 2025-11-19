package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/directs/types"
	membertypes "ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestGetManagerTree(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		managerID      string
		orgID          string
		mockResult     []types.OrgChartNode
		mockError      error
		expectedResult []types.OrgChartNode
		expectedError  error
	}{
		{
			name:      "successful retrieval - single level",
			managerID: "manager-1",
			orgID:     "org-1",
			mockResult: []types.OrgChartNode{
				{
					Member: membertypes.OrganizationMember{
						ID:             "report-1",
						UserID:         "user-1",
						Email:          "report1@example.com",
						OrganizationID: "org-1",
					},
					DirectReports: []types.OrgChartNode{},
					Depth:         0,
				},
			},
			expectedResult: []types.OrgChartNode{
				{
					Member:        membertypes.OrganizationMember{ID: "report-1", UserID: "user-1", Email: "report1@example.com", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{},
					Depth:         0,
				},
			},
		},
		{
			name:      "successful retrieval - multi level",
			managerID: "manager-1",
			orgID:     "org-1",
			mockResult: []types.OrgChartNode{
				{
					Member: membertypes.OrganizationMember{
						ID:             "report-1",
						UserID:         "user-1",
						Email:          "report1@example.com",
						OrganizationID: "org-1",
					},
					DirectReports: []types.OrgChartNode{
						{
							Member: membertypes.OrganizationMember{
								ID:             "report-2",
								UserID:         "user-2",
								Email:          "report2@example.com",
								OrganizationID: "org-1",
							},
							DirectReports: []types.OrgChartNode{},
							Depth:         1,
						},
					},
					Depth: 0,
				},
			},
			expectedResult: []types.OrgChartNode{
				{
					Member: membertypes.OrganizationMember{ID: "report-1", UserID: "user-1", Email: "report1@example.com", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{
						{
							Member:        membertypes.OrganizationMember{ID: "report-2", UserID: "user-2", Email: "report2@example.com", OrganizationID: "org-1"},
							DirectReports: []types.OrgChartNode{},
							Depth:         1,
						},
					},
					Depth: 0,
				},
			},
		},
		{
			name:           "empty tree",
			managerID:      "manager-nonexistent",
			orgID:          "org-1",
			mockResult:     []types.OrgChartNode{},
			expectedResult: []types.OrgChartNode{},
		},
		{
			name:          "database error",
			managerID:     "manager-1",
			orgID:         "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:      "deep hierarchy",
			managerID: "manager-1",
			orgID:     "org-1",
			mockResult: []types.OrgChartNode{
				{
					Member: membertypes.OrganizationMember{ID: "report-1", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{
						{
							Member: membertypes.OrganizationMember{ID: "report-2", OrganizationID: "org-1"},
							DirectReports: []types.OrgChartNode{
								{
									Member:        membertypes.OrganizationMember{ID: "report-3", OrganizationID: "org-1"},
									DirectReports: []types.OrgChartNode{},
									Depth:         2,
								},
							},
							Depth: 1,
						},
					},
					Depth: 0,
				},
			},
			expectedResult: []types.OrgChartNode{
				{
					Member: membertypes.OrganizationMember{ID: "report-1", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{
						{
							Member: membertypes.OrganizationMember{ID: "report-2", OrganizationID: "org-1"},
							DirectReports: []types.OrgChartNode{
								{
									Member:        membertypes.OrganizationMember{ID: "report-3", OrganizationID: "org-1"},
									DirectReports: []types.OrgChartNode{},
									Depth:         2,
								},
							},
							Depth: 1,
						},
					},
					Depth: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			mockDB.On("GetManagerTree", ctx, tt.managerID, tt.orgID).Return(tt.mockResult, tt.mockError)

			result, err := api.GetManagerTree(ctx, tt.managerID, tt.orgID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(result))
				// Verify tree structure recursively
				verifyOrgChartNode(t, tt.expectedResult, result)
			}

		mockDB.AssertExpectations(t)
	})
	}
}
