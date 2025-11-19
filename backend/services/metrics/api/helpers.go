package api

import (
	"context"
	"encoding/json"

	"ems.dev/backend/services/metrics/types"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
	"gorm.io/datatypes"
)

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
