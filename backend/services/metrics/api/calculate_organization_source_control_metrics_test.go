package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/metrics/types"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
	teamtypes "ems.dev/backend/services/team/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateOrganizationSourceControlMetrics(t *testing.T) {
	ctx := context.Background()
	orgID := "org-123"
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name                    string
		params                  types.OrganizationMetricsParams
		mockCumulativeMetrics   *sourcecontroltypes.MetricsResponse
		mockCumulativeError     error
		mockTeams               []teamtypes.Team
		mockTeamsError          error
		mockTeamMetrics         map[string]*sourcecontroltypes.MetricsResponse
		mockTeamMetricsError    error
		expectedResponse        *types.OrganizationMetricsResponse
		expectedError           error
	}{
		{
			name: "success - no team breakdown",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "monthly",
			},
			mockCumulativeMetrics: &sourcecontroltypes.MetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
			},
			expectedResponse: &types.OrganizationMetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
				TeamsBreakdown:  []types.TeamMetricsBreakdown{},
			},
		},
		{
			name: "success - with team breakdown",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				TeamIDs:        []string{"team-1", "team-2"},
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "monthly",
			},
			mockCumulativeMetrics: &sourcecontroltypes.MetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
			},
			mockTeams: []teamtypes.Team{
				{
					ID:       "team-1",
					Name:     "Team 1",
					PRPrefix: stringPtr("TEAM-1"),
				},
				{
					ID:       "team-2",
					Name:     "Team 2",
					PRPrefix: stringPtr("TEAM-2"),
				},
			},
			mockTeamMetrics: map[string]*sourcecontroltypes.MetricsResponse{
				"TEAM-1": {
					SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
					GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
				},
				"TEAM-2": {
					SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
					GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
				},
			},
			expectedResponse: &types.OrganizationMetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
				TeamsBreakdown: []types.TeamMetricsBreakdown{
					{
						TeamID:          "team-1",
						TeamName:        "Team 1",
						SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
						GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
					},
					{
						TeamID:          "team-2",
						TeamName:        "Team 2",
						SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
						GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
					},
				},
			},
		},
		{
			name: "success - team without prefix",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				TeamIDs:        []string{"team-1"},
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "monthly",
			},
			mockCumulativeMetrics: &sourcecontroltypes.MetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
			},
			mockTeams: []teamtypes.Team{
				{
					ID:       "team-1",
					Name:     "Team 1",
					PRPrefix: nil,
				},
			},
			expectedResponse: &types.OrganizationMetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
				TeamsBreakdown: []types.TeamMetricsBreakdown{
					{
						TeamID:          "team-1",
						TeamName:        "Team 1",
						SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
						GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
					},
				},
			},
		},
		{
			name: "error - cumulative metrics calculation fails",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "monthly",
			},
			mockCumulativeError: errors.New("failed to calculate metrics"),
			expectedError:       errors.New("failed to calculate metrics"),
		},
		{
			name: "error - list teams fails",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				TeamIDs:        []string{"team-1"},
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "monthly",
			},
			mockCumulativeMetrics: &sourcecontroltypes.MetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
			},
			mockTeamsError: errors.New("failed to list teams"),
			expectedError:  errors.New("failed to list teams"),
		},
		{
			name: "error - team metrics calculation fails",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				TeamIDs:        []string{"team-1"},
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "monthly",
			},
			mockCumulativeMetrics: &sourcecontroltypes.MetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
			},
			mockTeams: []teamtypes.Team{
				{
					ID:       "team-1",
					Name:     "Team 1",
					PRPrefix: stringPtr("TEAM-1"),
				},
			},
			mockTeamMetricsError: errors.New("failed to calculate team metrics"),
			expectedError:         errors.New("failed to calculate team metrics"),
		},
		{
			name: "success - empty teams list",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				TeamIDs:        []string{"team-1"},
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "monthly",
			},
			mockCumulativeMetrics: &sourcecontroltypes.MetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
			},
			mockTeams: []teamtypes.Team{},
			expectedResponse: &types.OrganizationMetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
				TeamsBreakdown:  []types.TeamMetricsBreakdown{},
			},
		},
		{
			name: "success - team not in filtered list",
			params: types.OrganizationMetricsParams{
				OrganizationID: orgID,
				TeamIDs:        []string{"team-1"},
				StartDate:      &startDate,
				EndDate:        &endDate,
				Interval:       "monthly",
			},
			mockCumulativeMetrics: &sourcecontroltypes.MetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
			},
			mockTeams: []teamtypes.Team{
				{
					ID:       "team-2",
					Name:     "Team 2",
					PRPrefix: stringPtr("TEAM-2"),
				},
			},
			expectedResponse: &types.OrganizationMetricsResponse{
				SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
				GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
				TeamsBreakdown:  []types.TeamMetricsBreakdown{},
			},
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

			// Setup cumulative metrics mock - this is always called first
			mockSourceControlAPI.On("CalculateMetrics", ctx, mock.Anything).
				Return(tt.mockCumulativeMetrics, tt.mockCumulativeError).
				Once()

			// Setup teams mock if needed
			if len(tt.params.TeamIDs) > 0 {
				mockTeamAPI.On("ListTeams", ctx, teamtypes.TeamSearchParams{
					OrganizationID: &tt.params.OrganizationID,
				}).Return(tt.mockTeams, tt.mockTeamsError).Once()

				// Setup team metrics mocks - these are called after cumulative
				if tt.mockTeamsError == nil {
					teamCallCount := 0
					for _, team := range tt.mockTeams {
						if contains(tt.params.TeamIDs, team.ID) && team.PRPrefix != nil && *team.PRPrefix != "" {
							if tt.mockTeamMetricsError != nil && teamCallCount == 0 {
								// Only fail on first team to avoid multiple error calls
								mockSourceControlAPI.On("CalculateMetrics", ctx, mock.Anything).
									Return(nil, tt.mockTeamMetricsError).Once()
								break
							} else if metrics, ok := tt.mockTeamMetrics[*team.PRPrefix]; ok {
								mockSourceControlAPI.On("CalculateMetrics", ctx, mock.Anything).
									Return(metrics, nil).Once()
								teamCallCount++
							}
						}
					}
				}
			}

			result, err := api.CalculateOrganizationSourceControlMetrics(ctx, tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResponse.SnapshotMetrics, result.SnapshotMetrics)
				assert.Equal(t, tt.expectedResponse.GraphMetrics, result.GraphMetrics)
				assert.Equal(t, len(tt.expectedResponse.TeamsBreakdown), len(result.TeamsBreakdown))
				for i, expectedTeam := range tt.expectedResponse.TeamsBreakdown {
					if i < len(result.TeamsBreakdown) {
						assert.Equal(t, expectedTeam.TeamID, result.TeamsBreakdown[i].TeamID)
						assert.Equal(t, expectedTeam.TeamName, result.TeamsBreakdown[i].TeamName)
					}
				}
			}

			mockSourceControlAPI.AssertExpectations(t)
			if len(tt.params.TeamIDs) > 0 {
				mockTeamAPI.AssertExpectations(t)
			}
		})
	}
}

// Helper function to check if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
