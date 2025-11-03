package api

import (
	"context"

	memberapi "ems.dev/backend/services/member/api"
	"ems.dev/backend/services/metrics/types"
	aicodeassistantapi "ems.dev/backend/services/aicodeassistant/api"
	aicodeassistanttypes "ems.dev/backend/services/aicodeassistant/types"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	teamapi "ems.dev/backend/services/team/api"
)

// MetricsAPI defines the interface for metrics operations
type MetricsAPI interface {
	// CalculateOrganizationSourceControlMetrics calculates aggregated source control metrics for an organization
	// Returns cumulative metrics and optionally a breakdown by team
	CalculateOrganizationSourceControlMetrics(ctx context.Context, params types.OrganizationMetricsParams) (*types.OrganizationMetricsResponse, error)
	// CalculateOrganizationAICodeAssistantMetrics calculates aggregated AI code assistant metrics for an organization
	CalculateOrganizationAICodeAssistantMetrics(ctx context.Context, params types.OrganizationMetricsParams) (*aicodeassistanttypes.MetricsResponse, error)
	// CalculateTeamAICodeAssistantMetrics calculates AI code assistant metrics for a specific team
	CalculateTeamAICodeAssistantMetrics(ctx context.Context, organizationID string, teamID string, params types.OrganizationMetricsParams) (*aicodeassistanttypes.MetricsResponse, error)
}

type Api struct {
	memberApi        memberapi.MemberAPI
	teamApi          teamapi.TeamAPI
	sourceControlApi sourcecontrolapi.SourceControlAPI
	aiCodeAssistantApi aicodeassistantapi.AICodeAssistantAPI
}

func NewApi(memberApi memberapi.MemberAPI, teamApi teamapi.TeamAPI, sourceControlApi sourcecontrolapi.SourceControlAPI, aiCodeAssistantApi aicodeassistantapi.AICodeAssistantAPI) MetricsAPI {
	return &Api{
		memberApi:        memberApi,
		teamApi:          teamApi,
		sourceControlApi: sourceControlApi,
		aiCodeAssistantApi: aiCodeAssistantApi,
	}
}
