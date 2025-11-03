package api

import (
	"context"
	"encoding/json"
	"fmt"

	"ems.dev/backend/services/aicodeassistant/database"
	"ems.dev/backend/services/aicodeassistant/metrics"
	"ems.dev/backend/services/aicodeassistant/types"
	memberapi "ems.dev/backend/services/member/api"
	membertypes "ems.dev/backend/services/member/types"
	"gorm.io/datatypes"
)

// AICodeAssistantAPI defines the interface for AI code assistant operations
type AICodeAssistantAPI interface {
	CreateOrUpdateDailyMetric(ctx context.Context, metric *types.AICodeAssistantDailyMetric) error
	GetDailyMetrics(ctx context.Context, params *types.AICodeAssistantDailyMetricParams) ([]*types.AICodeAssistantDailyMetric, error)
	GetUsageStats(ctx context.Context, params *types.AICodeAssistantDailyMetricParams) (*types.AICodeAssistantUsageStats, error)
	GetMemberDailyMetrics(ctx context.Context, organizationID, memberID string, params *types.AICodeAssistantMemberMetricsParams) ([]*types.AICodeAssistantDailyMetric, error)
	GetMemberUsageStats(ctx context.Context, organizationID, memberID string, params *types.AICodeAssistantMemberMetricsParams) (*types.AICodeAssistantUsageStats, error)
	CalculateMetrics(ctx context.Context, params types.MetricRuleParams) (*types.MetricsResponse, error)
	CalculateMemberMetrics(ctx context.Context, organizationID string, memberID string, params types.MemberMetricsParams) (*types.MetricsResponse, error)
}

// Api implements the AICodeAssistantAPI interface
type Api struct {
	db            database.DB
	memberAPI     memberapi.MemberAPI
	metricsEngine metrics.MetricsEngine
}

// NewApi creates a new instance of the AI Code Assistant API
func NewApi(db database.DB, memberAPI memberapi.MemberAPI) AICodeAssistantAPI {
	return &Api{
		db:            db,
		memberAPI:     memberAPI,
		metricsEngine: metrics.NewEngine(db),
	}
}

// CreateOrUpdateDailyMetric creates or updates a daily metric
func (a *Api) CreateOrUpdateDailyMetric(ctx context.Context, metric *types.AICodeAssistantDailyMetric) error {
	return a.db.CreateOrUpdateDailyMetric(ctx, metric)
}

// GetDailyMetrics retrieves daily metrics based on the given parameters
func (a *Api) GetDailyMetrics(ctx context.Context, params *types.AICodeAssistantDailyMetricParams) ([]*types.AICodeAssistantDailyMetric, error) {
	return a.db.GetDailyMetrics(ctx, params)
}

// GetUsageStats calculates aggregated statistics from daily metrics
func (a *Api) GetUsageStats(ctx context.Context, params *types.AICodeAssistantDailyMetricParams) (*types.AICodeAssistantUsageStats, error) {
	return a.db.GetUsageStats(ctx, params)
}

// GetMemberDailyMetrics retrieves daily metrics for a specific member
// This method handles resolving the member's external accounts internally
func (a *Api) GetMemberDailyMetrics(ctx context.Context, organizationID, memberID string, params *types.AICodeAssistantMemberMetricsParams) ([]*types.AICodeAssistantDailyMetric, error) {
	// Get member's external accounts for AI code assistant
	accountType := "ai-code-assistant"
	externalAccounts, err := a.memberAPI.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
		MemberIDs:      []string{memberID},
		OrganizationID: organizationID,
		AccountType:    &accountType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get member external accounts: %w", err)
	}

	// If member has no AI code assistant accounts, return empty array
	if len(externalAccounts) == 0 {
		return []*types.AICodeAssistantDailyMetric{}, nil
	}

	// Extract external account IDs
	externalAccountIDs := make([]string, 0, len(externalAccounts))
	for _, account := range externalAccounts {
		externalAccountIDs = append(externalAccountIDs, account.ID)
	}

	// Build params for database query
	dbParams := &types.AICodeAssistantDailyMetricParams{
		OrganizationID:     organizationID,
		ExternalAccountIDs: externalAccountIDs,
		ToolName:           params.ToolName,
		StartDate:          params.StartDate,
		EndDate:            params.EndDate,
	}

	return a.db.GetDailyMetrics(ctx, dbParams)
}

// GetMemberUsageStats calculates aggregated statistics for a specific member
func (a *Api) GetMemberUsageStats(ctx context.Context, organizationID, memberID string, params *types.AICodeAssistantMemberMetricsParams) (*types.AICodeAssistantUsageStats, error) {
	// Get member's external accounts for AI code assistant
	accountType := "ai-code-assistant"
	externalAccounts, err := a.memberAPI.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
		MemberIDs:      []string{memberID},
		OrganizationID: organizationID,
		AccountType:    &accountType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get member external accounts: %w", err)
	}

	// If member has no AI code assistant accounts, return zero stats
	if len(externalAccounts) == 0 {
		return &types.AICodeAssistantUsageStats{
			TotalLinesAccepted:  0,
			TotalLinesSuggested: 0,
			OverallAcceptRate:   0.0,
			ActiveSessions:      0,
		}, nil
	}

	// Extract external account IDs
	externalAccountIDs := make([]string, 0, len(externalAccounts))
	for _, account := range externalAccounts {
		externalAccountIDs = append(externalAccountIDs, account.ID)
	}

	// Build params for database query
	dbParams := &types.AICodeAssistantDailyMetricParams{
		OrganizationID:     organizationID,
		ExternalAccountIDs: externalAccountIDs,
		ToolName:           params.ToolName,
		StartDate:          params.StartDate,
		EndDate:            params.EndDate,
	}

	return a.db.GetUsageStats(ctx, dbParams)
}

// CalculateMetrics calculates AI code assistant metrics
func (a *Api) CalculateMetrics(ctx context.Context, params types.MetricRuleParams) (*types.MetricsResponse, error) {
	return a.metricsEngine.CalculateMetrics(ctx, params)
}

// CalculateMemberMetrics retrieves AI code assistant metrics for a specific member
func (a *Api) CalculateMemberMetrics(ctx context.Context, organizationID string, memberID string, params types.MemberMetricsParams) (*types.MetricsResponse, error) {
	// Get external accounts for this member (filter by ai-code-assistant type)
	accountType := "ai-code-assistant"
	externalAccounts, err := a.memberAPI.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
		OrganizationID: organizationID,
		MemberIDs:      []string{memberID},
		AccountType:    &accountType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get member external accounts: %w", err)
	}

	externalAccountIDs := []string{}
	for _, account := range externalAccounts {
		externalAccountIDs = append(externalAccountIDs, account.ID)
	}

	if len(externalAccountIDs) == 0 {
		// Return empty metrics response instead of error
		return &types.MetricsResponse{
			SnapshotMetrics: []*types.SnapshotCategory{},
			GraphMetrics:    []*types.GraphCategory{},
		}, nil
	}

	// Get member to find peers (members with same title)
	member, err := a.memberAPI.GetOrganizationMemberByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	// Get peer members (same title)
	peerMemberIDs := []string{}
	if member.TitleID != nil {
		orgMembers, err := a.memberAPI.GetOrganizationMembers(ctx, organizationID, &membertypes.OrganizationMemberParams{
			TitleIDs: []string{*member.TitleID},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get organization members: %w", err)
		}

		for _, orgMember := range orgMembers {
			// Exclude the current member from peers
			if orgMember.ID != memberID {
				peerMemberIDs = append(peerMemberIDs, orgMember.ID)
			}
		}
	}

	// Get peer external accounts
	peerExternalAccountIDs := []string{}
	if len(peerMemberIDs) > 0 {
		peerExternalAccounts, err := a.memberAPI.GetExternalAccounts(ctx, &membertypes.ExternalAccountParams{
			OrganizationID: organizationID,
			MemberIDs:      peerMemberIDs,
			AccountType:    &accountType,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get peer external accounts: %w", err)
		}

		for _, account := range peerExternalAccounts {
			peerExternalAccountIDs = append(peerExternalAccountIDs, account.ID)
		}
	}

	// Create the metric params with the external account IDs
	metricParamsMap := map[string]interface{}{
		"organizationId":          organizationID,
		"externalAccountIDs":      externalAccountIDs,
		"peersExternalAccountIDs": peerExternalAccountIDs,
	}

	// Marshal to JSON bytes
	metricParamsJSON, err := json.Marshal(metricParamsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metric params: %w", err)
	}

	metricParams := types.MetricRuleParams{
		MetricParams: datatypes.JSON(metricParamsJSON),
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Interval:     params.Interval,
	}

	return a.metricsEngine.CalculateMetrics(ctx, metricParams)
}
