package api

import (
	"context"
	"encoding/json"

	"ems.dev/backend/services/metrics/types"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
	teamtypes "ems.dev/backend/services/team/types"
	"gorm.io/datatypes"
)

// CalculateOrganizationSourceControlMetrics calculates aggregated source control metrics for an organization
// Params:
// - ctx: The context for the request
// - params: Parameters containing organization ID, optional team IDs, date range, and interval
// Returns:
// - OrganizationMetricsResponse: Aggregated metrics with optional team breakdown
// - error: If any error occurs during calculation
func (a *Api) CalculateOrganizationSourceControlMetrics(ctx context.Context, params types.OrganizationMetricsParams) (*types.OrganizationMetricsResponse, error) {
	var allTeams []teamtypes.Team
	var filteredTeams []teamtypes.Team

	// Calculate metrics for the entire organization
	cumulativeMetrics, err := a.calculateMetricsForOrganization(ctx, params)
	if err != nil {
		return nil, err
	}

	teamsBreakdown := []types.TeamMetricsBreakdown{}
	// Get teams if we need them
	if len(params.TeamIDs) > 0 {
		// Get all teams for the organization in one query
		allTeams, err = a.teamApi.ListTeams(ctx, teamtypes.TeamSearchParams{
			OrganizationID: &params.OrganizationID,
		})
		if err != nil {
			return nil, err
		}

		// Create a map of requested team IDs for fast lookup
		requestedTeamIDs := make(map[string]bool)
		for _, teamID := range params.TeamIDs {
			requestedTeamIDs[teamID] = true
		}

		// Filter teams to only the ones we need
		for _, team := range allTeams {
			if requestedTeamIDs[team.ID] {
				filteredTeams = append(filteredTeams, team)
			}
		}

		// Calculate per-team breakdown for the specified teams
		teamsBreakdown = make([]types.TeamMetricsBreakdown, 0, len(filteredTeams))
		for _, team := range filteredTeams {
			// Always use team prefix for filtering
			if team.PRPrefix != nil && *team.PRPrefix != "" {
				// Calculate metrics for this team using prefix
				teamMetrics, err := a.calculateMetricsForPrefix(ctx, params, *team.PRPrefix)
				if err != nil {
					return nil, err
				}

				teamsBreakdown = append(teamsBreakdown, types.TeamMetricsBreakdown{
					TeamID:          team.ID,
					TeamName:        team.Name,
					SnapshotMetrics: teamMetrics.SnapshotMetrics,
					GraphMetrics:    teamMetrics.GraphMetrics,
				})
			} else {
				// Team has no prefix, add empty metrics
				teamsBreakdown = append(teamsBreakdown, types.TeamMetricsBreakdown{
					TeamID:          team.ID,
					TeamName:        team.Name,
					SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
					GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
				})
			}
		}
	}

	return &types.OrganizationMetricsResponse{
		SnapshotMetrics: cumulativeMetrics.SnapshotMetrics,
		GraphMetrics:    cumulativeMetrics.GraphMetrics,
		TeamsBreakdown:  teamsBreakdown,
	}, nil

}

// calculateMetricsForOrganization calculates metrics for the entire organization without filtering by members
func (a *Api) calculateMetricsForOrganization(ctx context.Context, params types.OrganizationMetricsParams) (*sourcecontroltypes.MetricsResponse, error) {
	// Create the metric params with only organization ID (no sourceControlAccountIDs or pr_prefixes)
	// This will make the metrics engine calculate metrics for all PRs in the organization
	metricParamsMap := map[string]interface{}{
		"organizationId": params.OrganizationID,
	}

	// Marshal to JSON bytes
	metricParamsJSON, err := json.Marshal(metricParamsMap)
	if err != nil {
		return nil, err
	}

	interval := params.Interval
	if interval == "" {
		interval = "monthly" // default
	}

	metricParams := sourcecontroltypes.MetricRuleParams{
		MetricParams: datatypes.JSON(metricParamsJSON),
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Interval:     interval,
	}

	return a.sourceControlApi.CalculateMetrics(ctx, metricParams)
}

// calculateMetricsForPrefix is a helper function that calculates metrics for a team prefix
func (a *Api) calculateMetricsForPrefix(ctx context.Context, params types.OrganizationMetricsParams, prefix string) (*sourcecontroltypes.MetricsResponse, error) {
	// Create the metric params with the pr_prefixes
	metricParamsMap := map[string]interface{}{
		"organizationId": params.OrganizationID,
		"pr_prefixes":    []string{prefix},
	}

	// Marshal to JSON bytes
	metricParamsJSON, err := json.Marshal(metricParamsMap)
	if err != nil {
		return nil, err
	}

	interval := params.Interval
	if interval == "" {
		interval = "monthly" // default
	}

	metricParams := sourcecontroltypes.MetricRuleParams{
		MetricParams: datatypes.JSON(metricParamsJSON),
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Interval:     interval,
	}

	return a.sourceControlApi.CalculateMetrics(ctx, metricParams)
}
