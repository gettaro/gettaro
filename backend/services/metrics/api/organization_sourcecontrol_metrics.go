package api

import (
	"context"
	"encoding/json"

	membertypes "ems.dev/backend/services/member/types"
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
	var memberIDs []string

	// Get teams if we need them
	if len(params.TeamIDs) > 0 {
		// Get all teams for the organization in one query
		var err error
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

		// Filter teams to only the ones we need and collect member IDs
		memberIDMap := make(map[string]bool)
		for _, team := range allTeams {
			if requestedTeamIDs[team.ID] {
				filteredTeams = append(filteredTeams, team)
				for _, teamMember := range team.Members {
					if !memberIDMap[teamMember.MemberID] {
						memberIDMap[teamMember.MemberID] = true
						memberIDs = append(memberIDs, teamMember.MemberID)
					}
				}
			}
		}
	} else {
		// Get all organization members
		orgMembers, err := a.memberApi.GetOrganizationMembers(ctx, params.OrganizationID, &membertypes.OrganizationMemberParams{})
		if err != nil {
			return nil, err
		}
		for _, member := range orgMembers {
			memberIDs = append(memberIDs, member.ID)
		}
	}

	if len(memberIDs) == 0 {
		// Return empty metrics response
		return &types.OrganizationMetricsResponse{
			SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
			GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
			TeamsBreakdown:  []types.TeamMetricsBreakdown{},
		}, nil
	}

	// Calculate cumulative metrics for all members
	cumulativeMetrics, err := a.calculateMetricsForMembers(ctx, params, memberIDs)
	if err != nil {
		return nil, err
	}

	// Calculate per-team breakdown if teams were specified
	var teamsBreakdown []types.TeamMetricsBreakdown
	if len(filteredTeams) > 0 {
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

// calculateMetricsForMembers is a helper function that calculates metrics for a set of member IDs
func (a *Api) calculateMetricsForMembers(ctx context.Context, params types.OrganizationMetricsParams, memberIDs []string) (*sourcecontroltypes.MetricsResponse, error) {
	if len(memberIDs) == 0 {
		return &sourcecontroltypes.MetricsResponse{
			SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
			GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
		}, nil
	}

	// Get source control accounts for these members
	sourceControlAccounts, err := a.sourceControlApi.GetSourceControlAccounts(ctx, &sourcecontroltypes.SourceControlAccountParams{
		OrganizationID: params.OrganizationID,
		MemberIDs:      memberIDs,
	})
	if err != nil {
		return nil, err
	}

	sourceControlAccountIDs := []string{}
	for _, account := range sourceControlAccounts {
		sourceControlAccountIDs = append(sourceControlAccountIDs, account.ID)
	}

	if len(sourceControlAccountIDs) == 0 {
		// Return empty metrics response
		return &sourcecontroltypes.MetricsResponse{
			SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
			GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
		}, nil
	}

	// Create the metric params with the source control account IDs
	metricParamsMap := map[string]interface{}{
		"organizationId":          params.OrganizationID,
		"sourceControlAccountIDs": sourceControlAccountIDs,
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
	// Get all PRs with this prefix in the organization
	prs, err := a.sourceControlApi.GetPullRequests(ctx, &sourcecontroltypes.PullRequestParams{
		OrganizationID: &params.OrganizationID,
		Prefix:         &prefix,
		StartDate:      params.StartDate,
		EndDate:        params.EndDate,
	})
	if err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		// Return empty metrics response
		return &sourcecontroltypes.MetricsResponse{
			SnapshotMetrics: []*sourcecontroltypes.SnapshotCategory{},
			GraphMetrics:    []*sourcecontroltypes.GraphCategory{},
		}, nil
	}

	// Get source control account IDs from the PRs
	sourceControlAccountIDs := make([]string, 0, len(prs))
	accountIDMap := make(map[string]bool)
	for _, pr := range prs {
		if !accountIDMap[pr.SourceControlAccountID] {
			accountIDMap[pr.SourceControlAccountID] = true
			sourceControlAccountIDs = append(sourceControlAccountIDs, pr.SourceControlAccountID)
		}
	}

	// Create the metric params with the source control account IDs
	metricParamsMap := map[string]interface{}{
		"organizationId":          params.OrganizationID,
		"sourceControlAccountIDs": sourceControlAccountIDs,
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
