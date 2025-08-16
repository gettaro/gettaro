package metrics

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/database"
	"ems.dev/backend/services/sourcecontrol/metrics/engine"
	metrictypes "ems.dev/backend/services/sourcecontrol/metrics/types"
	"ems.dev/backend/services/sourcecontrol/types"
)

type Engine struct {
	Metrics         []metrictypes.MetricRule `json:"metrics"`
	sourceControlDB database.DB
}

func NewEngine(sourceControlDB database.DB) *Engine {
	// TODO: Eventually I'll have a get engine for organization and remove this. No need to inject it from main.go for now.
	return &Engine{
		Metrics: []metrictypes.MetricRule{
			engine.NewTimeToMergeRule(metrictypes.BaseMetricRule{
				ID:        "median_time_to_merge",
				Name:      "Median time to merge",
				Unit:      types.UnitSeconds,
				Category:  "Efficiency",
				Dimension: metrictypes.MetricDimensionTimeToMerge,
				Operation: metrictypes.MetricOperationMedian,
			}, sourceControlDB),
			engine.NewPRsMergedRule(metrictypes.BaseMetricRule{
				ID:        "prs_merged_count",
				Name:      "PRs Merged",
				Unit:      types.UnitCount,
				Category:  "Activity",
				Dimension: metrictypes.MetricDimensionMergedPRs,
				Operation: metrictypes.MetricOperationCount,
			}, sourceControlDB),
		},
		sourceControlDB: sourceControlDB,
	}
}

func (e *Engine) CalculateMetrics(ctx context.Context, params types.MetricRuleParams) (*types.MetricsResponse, error) {
	metrics := []*types.SnapshotCategory{}
	graphMetrics := []*types.GraphCategory{}

	for _, rule := range e.Metrics {
		snapshotMetric, graphMetric, err := rule.Calculate(ctx, params)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &types.SnapshotCategory{
			Category: rule.Category(),
			Metrics:  []types.SnapshotMetric{*snapshotMetric},
		})
		graphMetrics = append(graphMetrics, &types.GraphCategory{
			Category: rule.Category(),
			Metrics:  []types.GraphMetric{*graphMetric},
		})
	}
	return &types.MetricsResponse{
		SnapshotMetrics: metrics,
		GraphMetrics:    graphMetrics,
	}, nil
}
