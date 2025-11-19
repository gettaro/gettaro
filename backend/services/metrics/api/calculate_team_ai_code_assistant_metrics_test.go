package api

import (
	"context"
	"errors"
	"testing"
	"time"

	aicodeassistanttypes "ems.dev/backend/services/aicodeassistant/types"
	membertypes "ems.dev/backend/services/member/types"
	"ems.dev/backend/services/metrics/types"
	teamtypes "ems.dev/backend/services/team/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateTeamAICodeAssistantMetrics(t *testing.T) {
	ctx := context.Background()
	orgID := "org-123"
	teamID := "team-123"
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name                    string
		organizationID          string
		teamID                  string
		params                  types.OrganizationMetricsParams
		mockTeam                *teamtypes.Team
		mockTeamError           error
		mockExternalAccounts    []membertypes.ExternalAccount
		mockExternalAccountsError error
		mockMetricsResponse     *aicodeassistanttypes.MetricsResponse
		mockMetricsError        error
		expectedResponse        *aicodeassistanttypes.MetricsResponse
		expectedError           error
	}{
		{
			name:           "success - with team members and external accounts",
			organizationID: orgID,
			teamID:         teamID,
			params: types.OrganizationMetricsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Interval:  "weekly",
			},
			mockTeam: &teamtypes.Team{
				ID:   teamID,
				Name: "Test Team",
				Members: []teamtypes.TeamMember{
					{MemberID: "member-1"},
					{MemberID: "member-2"},
				},
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
			name:           "success - team with no members",
			organizationID: orgID,
			teamID:         teamID,
			params: types.OrganizationMetricsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Interval:  "weekly",
			},
			mockTeam: &teamtypes.Team{
				ID:      teamID,
				Name:    "Test Team",
				Members: []teamtypes.TeamMember{},
			},
			expectedResponse: &aicodeassistanttypes.MetricsResponse{
				SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
				GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
			},
		},
		{
			name:           "success - team members with no external accounts",
			organizationID: orgID,
			teamID:         teamID,
			params: types.OrganizationMetricsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Interval:  "weekly",
			},
			mockTeam: &teamtypes.Team{
				ID:   teamID,
				Name: "Test Team",
				Members: []teamtypes.TeamMember{
					{MemberID: "member-1"},
				},
			},
			mockExternalAccounts: []membertypes.ExternalAccount{},
			expectedResponse: &aicodeassistanttypes.MetricsResponse{
				SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
				GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
			},
		},
		{
			name:           "success - default interval",
			organizationID: orgID,
			teamID:         teamID,
			params: types.OrganizationMetricsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Interval:  "",
			},
			mockTeam: &teamtypes.Team{
				ID:   teamID,
				Name: "Test Team",
				Members: []teamtypes.TeamMember{
					{MemberID: "member-1"},
				},
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
			name:           "error - get team fails",
			organizationID: orgID,
			teamID:         teamID,
			params: types.OrganizationMetricsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Interval:  "weekly",
			},
			mockTeamError: errors.New("team not found"),
			expectedError:  errors.New("team not found"),
		},
		{
			name:           "error - get external accounts fails",
			organizationID: orgID,
			teamID:         teamID,
			params: types.OrganizationMetricsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Interval:  "weekly",
			},
			mockTeam: &teamtypes.Team{
				ID:   teamID,
				Name: "Test Team",
				Members: []teamtypes.TeamMember{
					{MemberID: "member-1"},
				},
			},
			mockExternalAccountsError: errors.New("failed to get external accounts"),
			expectedError:               errors.New("failed to get external accounts"),
		},
		{
			name:           "error - calculate metrics fails",
			organizationID: orgID,
			teamID:         teamID,
			params: types.OrganizationMetricsParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Interval:  "weekly",
			},
			mockTeam: &teamtypes.Team{
				ID:   teamID,
				Name: "Test Team",
				Members: []teamtypes.TeamMember{
					{MemberID: "member-1"},
				},
			},
			mockExternalAccounts: []membertypes.ExternalAccount{
				{ID: "account-1", AccountType: "ai-code-assistant"},
			},
			mockMetricsError: errors.New("failed to calculate metrics"),
			expectedError:     errors.New("failed to calculate metrics"),
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

			// Setup team mock
			mockTeamAPI.On("GetTeamByOrganization", ctx, tt.teamID, tt.organizationID).
				Return(tt.mockTeam, tt.mockTeamError).Once()

			// Setup external accounts mock (only if team retrieval succeeds and team has members)
			if tt.mockTeamError == nil && tt.mockTeam != nil && len(tt.mockTeam.Members) > 0 {
				mockMemberAPI.On("GetExternalAccounts", ctx, mock.MatchedBy(func(params *membertypes.ExternalAccountParams) bool {
					return params.OrganizationID == tt.organizationID &&
						params.AccountType != nil &&
						*params.AccountType == "ai-code-assistant" &&
						len(params.MemberIDs) > 0
				})).Return(tt.mockExternalAccounts, tt.mockExternalAccountsError).Once()

				// Setup metrics calculation mock (only if we have external accounts and no error getting them)
				if tt.mockExternalAccountsError == nil && len(tt.mockExternalAccounts) > 0 {
					mockAICodeAssistantAPI.On("CalculateMetrics", ctx, mock.Anything).
						Return(tt.mockMetricsResponse, tt.mockMetricsError).Once()
				}
			}

			result, err := api.CalculateTeamAICodeAssistantMetrics(ctx, tt.organizationID, tt.teamID, tt.params)

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

			mockTeamAPI.AssertExpectations(t)
			if tt.mockTeamError == nil && tt.mockTeam != nil && len(tt.mockTeam.Members) > 0 {
				mockMemberAPI.AssertExpectations(t)
				if tt.mockExternalAccountsError == nil && len(tt.mockExternalAccounts) > 0 {
					mockAICodeAssistantAPI.AssertExpectations(t)
				}
			}
		})
	}
}
