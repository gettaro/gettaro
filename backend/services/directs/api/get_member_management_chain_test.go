package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/directs/types"
	membertypes "ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestGetMemberManagementChain(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		reportID       string
		orgID          string
		mockResult     []types.ManagementChain
		mockError      error
		expectedResult []types.ManagementChain
		expectedError  error
	}{
		{
			name:     "successful retrieval - single level",
			reportID: "report-1",
			orgID:    "org-1",
			mockResult: []types.ManagementChain{
				{
					Member: membertypes.OrganizationMember{
						ID:             "report-1",
						UserID:         "user-1",
						Email:          "report1@example.com",
						OrganizationID: "org-1",
					},
					Manager: &membertypes.OrganizationMember{
						ID:             "manager-1",
						UserID:         "user-2",
						Email:          "manager1@example.com",
						OrganizationID: "org-1",
					},
					Depth: 0,
				},
			},
			expectedResult: []types.ManagementChain{
				{
					Member:  membertypes.OrganizationMember{ID: "report-1", UserID: "user-1", Email: "report1@example.com", OrganizationID: "org-1"},
					Manager: &membertypes.OrganizationMember{ID: "manager-1", UserID: "user-2", Email: "manager1@example.com", OrganizationID: "org-1"},
					Depth:   0,
				},
			},
		},
		{
			name:     "successful retrieval - multi level chain",
			reportID: "report-1",
			orgID:    "org-1",
			mockResult: []types.ManagementChain{
				{
					Member: membertypes.OrganizationMember{
						ID:             "report-1",
						UserID:         "user-1",
						Email:          "report1@example.com",
						OrganizationID: "org-1",
					},
					Manager: &membertypes.OrganizationMember{
						ID:             "manager-1",
						UserID:         "user-2",
						Email:          "manager1@example.com",
						OrganizationID: "org-1",
					},
					Depth: 0,
				},
				{
					Member: membertypes.OrganizationMember{
						ID:             "manager-1",
						UserID:         "user-2",
						Email:          "manager1@example.com",
						OrganizationID: "org-1",
					},
					Manager: &membertypes.OrganizationMember{
						ID:             "manager-2",
						UserID:         "user-3",
						Email:          "manager2@example.com",
						OrganizationID: "org-1",
					},
					Depth: 1,
				},
			},
			expectedResult: []types.ManagementChain{
				{
					Member:  membertypes.OrganizationMember{ID: "report-1", UserID: "user-1", Email: "report1@example.com", OrganizationID: "org-1"},
					Manager: &membertypes.OrganizationMember{ID: "manager-1", UserID: "user-2", Email: "manager1@example.com", OrganizationID: "org-1"},
					Depth:   0,
				},
				{
					Member:  membertypes.OrganizationMember{ID: "manager-1", UserID: "user-2", Email: "manager1@example.com", OrganizationID: "org-1"},
					Manager: &membertypes.OrganizationMember{ID: "manager-2", UserID: "user-3", Email: "manager2@example.com", OrganizationID: "org-1"},
					Depth:   1,
				},
			},
		},
		{
			name:           "member has no manager",
			reportID:       "report-top-level",
			orgID:          "org-1",
			mockResult:     []types.ManagementChain{},
			expectedResult: []types.ManagementChain{},
		},
		{
			name:          "database error",
			reportID:      "report-1",
			orgID:         "org-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:     "deep chain",
			reportID: "report-1",
			orgID:    "org-1",
			mockResult: []types.ManagementChain{
				{
					Member:  membertypes.OrganizationMember{ID: "report-1", OrganizationID: "org-1"},
					Manager: &membertypes.OrganizationMember{ID: "manager-1", OrganizationID: "org-1"},
					Depth:   0,
				},
				{
					Member:  membertypes.OrganizationMember{ID: "manager-1", OrganizationID: "org-1"},
					Manager: &membertypes.OrganizationMember{ID: "manager-2", OrganizationID: "org-1"},
					Depth:   1,
				},
				{
					Member:  membertypes.OrganizationMember{ID: "manager-2", OrganizationID: "org-1"},
					Manager: &membertypes.OrganizationMember{ID: "manager-3", OrganizationID: "org-1"},
					Depth:   2,
				},
			},
			expectedResult: []types.ManagementChain{
				{
					Member:  membertypes.OrganizationMember{ID: "report-1", OrganizationID: "org-1"},
					Manager: &membertypes.OrganizationMember{ID: "manager-1", OrganizationID: "org-1"},
					Depth:   0,
				},
				{
					Member:  membertypes.OrganizationMember{ID: "manager-1", OrganizationID: "org-1"},
					Manager: &membertypes.OrganizationMember{ID: "manager-2", OrganizationID: "org-1"},
					Depth:   1,
				},
				{
					Member:  membertypes.OrganizationMember{ID: "manager-2", OrganizationID: "org-1"},
					Manager: &membertypes.OrganizationMember{ID: "manager-3", OrganizationID: "org-1"},
					Depth:   2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			mockDB.On("GetMemberManagementChain", ctx, tt.reportID, tt.orgID).Return(tt.mockResult, tt.mockError)

			result, err := api.GetMemberManagementChain(ctx, tt.reportID, tt.orgID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(result))
				for i, expected := range tt.expectedResult {
					assert.Equal(t, expected.Member.ID, result[i].Member.ID)
					assert.Equal(t, expected.Depth, result[i].Depth)
					if expected.Manager != nil {
						assert.NotNil(t, result[i].Manager)
						assert.Equal(t, expected.Manager.ID, result[i].Manager.ID)
					} else {
						assert.Nil(t, result[i].Manager)
					}
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
