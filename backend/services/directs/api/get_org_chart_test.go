package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/directs/types"
	membertypes "ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestGetOrgChart(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		orgID          string
		mockResult     []types.OrgChartNode
		mockError      error
		expectedResult []types.OrgChartNode
		expectedError  error
	}{
		{
			name:  "successful retrieval - single top level manager",
			orgID: "org-1",
			mockResult: []types.OrgChartNode{
				{
					Member: membertypes.OrganizationMember{
						ID:             "manager-1",
						UserID:         "user-1",
						Email:          "manager1@example.com",
						OrganizationID: "org-1",
					},
					DirectReports: []types.OrgChartNode{
						{
							Member: membertypes.OrganizationMember{
								ID:             "report-1",
								UserID:         "user-2",
								Email:          "report1@example.com",
								OrganizationID: "org-1",
							},
							DirectReports: []types.OrgChartNode{},
							Depth:         0,
						},
					},
					Depth: 0,
				},
			},
			expectedResult: []types.OrgChartNode{
				{
					Member: membertypes.OrganizationMember{ID: "manager-1", UserID: "user-1", Email: "manager1@example.com", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{
						{
							Member:        membertypes.OrganizationMember{ID: "report-1", UserID: "user-2", Email: "report1@example.com", OrganizationID: "org-1"},
							DirectReports: []types.OrgChartNode{},
							Depth:         0,
						},
					},
					Depth: 0,
				},
			},
		},
		{
			name:  "successful retrieval - multiple top level managers",
			orgID: "org-1",
			mockResult: []types.OrgChartNode{
				{
					Member:        membertypes.OrganizationMember{ID: "manager-1", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{},
					Depth:         0,
				},
				{
					Member:        membertypes.OrganizationMember{ID: "manager-2", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{},
					Depth:         0,
				},
			},
			expectedResult: []types.OrgChartNode{
				{
					Member:        membertypes.OrganizationMember{ID: "manager-1", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{},
					Depth:         0,
				},
				{
					Member:        membertypes.OrganizationMember{ID: "manager-2", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{},
					Depth:         0,
				},
			},
		},
		{
			name:           "empty org chart",
			orgID:          "org-empty",
			mockResult:     []types.OrgChartNode{},
			expectedResult: []types.OrgChartNode{},
		},
		{
			name:          "database error",
			orgID:         "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:  "complex org chart",
			orgID: "org-1",
			mockResult: []types.OrgChartNode{
				{
					Member: membertypes.OrganizationMember{ID: "ceo", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{
						{
							Member: membertypes.OrganizationMember{ID: "vp-1", OrganizationID: "org-1"},
							DirectReports: []types.OrgChartNode{
								{
									Member:        membertypes.OrganizationMember{ID: "manager-1", OrganizationID: "org-1"},
									DirectReports: []types.OrgChartNode{},
									Depth:         1,
								},
							},
							Depth: 0,
						},
						{
							Member:        membertypes.OrganizationMember{ID: "vp-2", OrganizationID: "org-1"},
							DirectReports: []types.OrgChartNode{},
							Depth:         0,
						},
					},
					Depth: 0,
				},
			},
			expectedResult: []types.OrgChartNode{
				{
					Member: membertypes.OrganizationMember{ID: "ceo", OrganizationID: "org-1"},
					DirectReports: []types.OrgChartNode{
						{
							Member: membertypes.OrganizationMember{ID: "vp-1", OrganizationID: "org-1"},
							DirectReports: []types.OrgChartNode{
								{
									Member:        membertypes.OrganizationMember{ID: "manager-1", OrganizationID: "org-1"},
									DirectReports: []types.OrgChartNode{},
									Depth:         1,
								},
							},
							Depth: 0,
						},
						{
							Member:        membertypes.OrganizationMember{ID: "vp-2", OrganizationID: "org-1"},
							DirectReports: []types.OrgChartNode{},
							Depth:         0,
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

			mockDB.On("GetOrgChart", ctx, tt.orgID).Return(tt.mockResult, tt.mockError)

			result, err := api.GetOrgChart(ctx, tt.orgID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(result))
				verifyOrgChartNode(t, tt.expectedResult, result)
			}

		mockDB.AssertExpectations(t)
	})
	}
}
