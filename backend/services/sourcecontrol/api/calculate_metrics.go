package api

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

// CalculateMetrics calculates source control metrics
func (a *Api) CalculateMetrics(ctx context.Context, params types.MetricRuleParams) (*types.MetricsResponse, error) {
	return a.metricsEngine.CalculateMetrics(ctx, params)
}
