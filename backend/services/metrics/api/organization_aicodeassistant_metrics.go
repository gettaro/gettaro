package api

import (
	"context"
	"encoding/json"

	aicodeassistanttypes "ems.dev/backend/services/aicodeassistant/types"
	membertypes "ems.dev/backend/services/member/types"
	"ems.dev/backend/services/metrics/types"
	"gorm.io/datatypes"
)

// CalculateOrganizationAICodeAssistantMetrics calculates aggregated AI code assistant metrics for an organization
// Params:
// - ctx: The context for the request
// - params: Parameters containing organization ID, date range, and interval
// Returns:
// - MetricsResponse: Aggregated AI code assistant metrics
// - error: If any error occurs during calculation
func (a *Api) CalculateOrganizationAICodeAssistantMetrics(ctx context.Context, params types.OrganizationMetricsParams) (*aicodeassistanttypes.MetricsResponse, error) {
	// Get all external accounts for the organization (filter by ai-code-assistant type)
	accountType := "ai-code-assistant"
	externalAccounts, err := a.memberApi.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
		OrganizationID: params.OrganizationID,
		AccountType:    &accountType,
	})
	if err != nil {
		return nil, err
	}

	externalAccountIDs := []string{}
	for _, account := range externalAccounts {
		externalAccountIDs = append(externalAccountIDs, account.ID)
	}

	if len(externalAccountIDs) == 0 {
		// Return empty metrics response instead of error
		return &aicodeassistanttypes.MetricsResponse{
			SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
			GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
		}, nil
	}

	// For organization-level metrics, we don't need peer comparison
	// Create the metric params with the external account IDs
	metricParamsMap := map[string]interface{}{
		"organizationId":          params.OrganizationID,
		"externalAccountIDs":      externalAccountIDs,
		"peersExternalAccountIDs": []string{}, // No peer comparison at org level
	}

	// Marshal to JSON bytes
	metricParamsJSON, err := json.Marshal(metricParamsMap)
	if err != nil {
		return nil, err
	}

	interval := params.Interval
	if interval == "" {
		interval = "weekly" // default
	}

	metricParams := aicodeassistanttypes.MetricRuleParams{
		MetricParams: datatypes.JSON(metricParamsJSON),
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Interval:     interval,
	}

	return a.aiCodeAssistantApi.CalculateMetrics(ctx, metricParams)
}

// CalculateTeamAICodeAssistantMetrics calculates AI code assistant metrics for a specific team
// by getting all team members and their external accounts
func (a *Api) CalculateTeamAICodeAssistantMetrics(ctx context.Context, organizationID string, teamID string, params types.OrganizationMetricsParams) (*aicodeassistanttypes.MetricsResponse, error) {
	// Get team to access team members
	team, err := a.teamApi.GetTeamByOrganization(ctx, teamID, organizationID)
	if err != nil {
		return nil, err
	}

	// Extract member IDs from team
	memberIDs := make([]string, 0, len(team.Members))
	for _, teamMember := range team.Members {
		memberIDs = append(memberIDs, teamMember.MemberID)
	}

	if len(memberIDs) == 0 {
		// Return empty metrics response if team has no members
		return &aicodeassistanttypes.MetricsResponse{
			SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
			GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
		}, nil
	}

	// Get external accounts for team members (filter by ai-code-assistant type)
	accountType := "ai-code-assistant"
	externalAccounts, err := a.memberApi.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
		OrganizationID: organizationID,
		MemberIDs:      memberIDs,
		AccountType:    &accountType,
	})
	if err != nil {
		return nil, err
	}

	externalAccountIDs := []string{}
	for _, account := range externalAccounts {
		externalAccountIDs = append(externalAccountIDs, account.ID)
	}

	if len(externalAccountIDs) == 0 {
		// Return empty metrics response if team members have no AI code assistant accounts
		return &aicodeassistanttypes.MetricsResponse{
			SnapshotMetrics: []*aicodeassistanttypes.SnapshotCategory{},
			GraphMetrics:    []*aicodeassistanttypes.GraphCategory{},
		}, nil
	}

	// Create the metric params with the external account IDs
	metricParamsMap := map[string]interface{}{
		"organizationId":          organizationID,
		"externalAccountIDs":      externalAccountIDs,
		"peersExternalAccountIDs": []string{}, // No peer comparison at team level
	}

	// Marshal to JSON bytes
	metricParamsJSON, err := json.Marshal(metricParamsMap)
	if err != nil {
		return nil, err
	}

	interval := params.Interval
	if interval == "" {
		interval = "weekly" // default
	}

	metricParams := aicodeassistanttypes.MetricRuleParams{
		MetricParams: datatypes.JSON(metricParamsJSON),
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Interval:     interval,
	}

	return a.aiCodeAssistantApi.CalculateMetrics(ctx, metricParams)
}
