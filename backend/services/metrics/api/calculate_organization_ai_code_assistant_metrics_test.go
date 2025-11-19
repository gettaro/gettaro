package api

import (
	"context"
	"errors"
	"testing"
	"time"

	aicodeassistanttypes "ems.dev/backend/services/aicodeassistant/types"
	membertypes "ems.dev/backend/services/member/types"
	"ems.dev/backend/services/metrics/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateOrganizationAICodeAssistantMetrics(t *testing.T) {
	ctx := context.Background()
	orgID := "org-123"
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name                  string
		params                types.OrganizationMetricsParams
		mockExternalAccounts  []membertypes.ExternalAccount
		mockExternalAccountsError error
		mockMetricsResponse   *aicodeassistanttypes.MetricsResponse
		mockMetricsError      error
		expectedResponse      *aicodeassistanttypes.MetricsResponse
		expectedError         error
	}{
		{
			name: "success - with external accounts",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "weekly",
			},
			mockExternalAccounts: []membertypes.ExternalAccount{
				{ID: "account-1", AccountType: "ai-code-assistant"},
				{ID: "account-2", AccountType: "ai-code-assistant"},
			},
			mockMetricsResponse: &aicodeassistanttypes.MetricsResponse{
				SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
				GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
			},
			expectedResponse: &aicodeassistanttypes.MetricsResponse{
				SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
				GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
			},
		},
		{
			name: "success - no external accounts",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "weekly",
			},
			mockExternalAccounts: []membertypes.ExternalAccount{},
			expectedResponse: &aicodeassistanttypes.MetricsResponse{
				SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
				GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
			},
		},
		{
			name: "success - default interval",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "",
			},
			mockExternalAccounts: []membertypes.ExternalAccount{
				{ID: "account-1", AccountType: "ai-code-assistant"},
			},
			mockMetricsResponse: &aicodeassistanttypes.MetricsResponse{
				SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
				GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
			},
			expectedResponse: &aicodeassistanttypes.MetricsResponse{
				SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
				GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
			},
		},
		{
			name: "error - get external accounts fails",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "weekly",
			},
			mockExternalAccountsError: errors.New("failed to get external accounts"),
			expectedError:              errors.New("failed to get external accounts"),
		},
		{
			name: "error - calculate metrics fails",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "weekly",
			},
			mockExternalAccounts: []membertypes.ExternalAccount{
				{ID: "account-1", AccountType: "ai-code-assistant"},
			},
			mockMetricsError: errors.New("failed to calculate metrics"),
			expectedError:    errors.New("failed to calculate metrics"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberAPI := new(MockMemberAPI)
			mockTeamAPI := new(MockTeamAPI)
			mockSourceControlAPI := new(MockSourceControlAPI)
			mockAICodeAssistantAPI := new(MockAICodeAssistantAPI)

			api := &Api{
				memberApi:         mockMemberAPI,
				teamApi:           mockTeamAPI,
				sourceControlApi:  mockSourceControlAPI,
				aiCodeAssistantApi: mockAICodeAssistantAPI,
			}

			// Setup external accounts mock
			mockMemberAPI.On("GetExternalAccounts", ctx, mock.MatchedBy(func(params *membertypes.ExternalAccountParams) bool {
				return params.OrganizationID == orgID && params.AccountType != nil && *params.AccountType == "ai-code-assistant"
			})).Return(tt.mockExternalAccounts, tt.mockExternalAccountsError).Once()

			// Setup metrics calculation mock (only if we have external accounts and no error getting them)
			if tt.mockExternalAccountsError == nil && len(tt.mockExternalAccounts) > 0 {
				mockAICodeAssistantAPI.On("CalculateMetrics", ctx, mock.Anything).
					Return(tt.mockMetricsResponse, tt.mockMetricsError).Once()
			}

			result, err := api.CalculateOrganizationAICodeAssistantMetrics(ctx, tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResponse.SnapshotMetrics, result.SnapshotMetrics)
				assert.Equal(t, tt.expectedResponse.GraphMetrics, result.GraphMetrics)
			}

			mockMemberAPI.AssertExpectations(t)
			if tt.mockExternalAccountsError == nil && len(tt.mockExternalAccounts) > 0 {
				mockAICodeAssistantAPI.AssertExpectations(t)
			}
		})
	}
}
