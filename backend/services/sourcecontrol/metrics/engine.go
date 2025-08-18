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
			engine.NewPRsMergedRule(metrictypes.BaseMetricRule{
				ID:             "prs_merged_count",
				Name:           "PRs Merged",
				Description:    "Total number of pull requests that were successfully merged. This metric counts PRs with status 'closed' and merged_at timestamp. Peer comparison shows the median PRs merged across other organization members.",
				Unit:           types.UnitCount,
				Category:       "Activity",
				Dimension:      metrictypes.MetricDimensionMergedPRs,
				Operation:      metrictypes.MetricOperationCount,
				IconIdentifier: "git-merge",
				IconColor:      "blue",
			}, sourceControlDB),
			engine.NewLOCAddedRule(metrictypes.BaseMetricRule{
				ID:             "loc_added_count",
				Name:           "Lines of Code Added",
				Description:    "Total lines of code added across all merged pull requests. Extracted from PR metadata additions field. Peer comparison shows the median LOC added across other organization members.",
				Unit:           types.UnitCount,
				Category:       "Activity",
				Dimension:      metrictypes.MetricDimensionLOCAdded,
				Operation:      metrictypes.MetricOperationCount,
				IconIdentifier: "plus-circle",
				IconColor:      "green",
			}, sourceControlDB),
			engine.NewLOCRemovedRule(metrictypes.BaseMetricRule{
				ID:             "loc_removed_count",
				Name:           "Lines of Code Removed",
				Description:    "Total lines of code removed across all merged pull requests. Extracted from PR metadata deletions field. Peer comparison shows the median LOC removed across other organization members.",
				Unit:           types.UnitCount,
				Category:       "Activity",
				Dimension:      metrictypes.MetricDimensionLOCRemoved,
				Operation:      metrictypes.MetricOperationCount,
				IconIdentifier: "minus-circle",
				IconColor:      "red",
			}, sourceControlDB),
			engine.NewPRsReviewedRule(metrictypes.BaseMetricRule{
				ID:             "prs_reviewed_count",
				Name:           "PRs Reviewed",
				Description:    "Total number of unique pull requests that the member has reviewed. This metric counts PRs authored by others that the member has provided review comments on, regardless of whether the PR author is mapped to a member in the system. Peer comparison shows the median PRs reviewed across other organization members.",
				Unit:           types.UnitCount,
				Category:       "Collaboration",
				Dimension:      metrictypes.MetricDimensionReviwedPRs,
				Operation:      metrictypes.MetricOperationCount,
				IconIdentifier: "eye",
				IconColor:      "purple",
			}, sourceControlDB),
			engine.NewTimeToMergeRule(metrictypes.BaseMetricRule{
				ID:             "median_time_to_merge",
				Name:           "Median time to merge",
				Description:    "Median time between PR creation and merge completion. Calculated as the difference between merged_at and created_at timestamps. Peer comparison shows the median time to merge across other organization members.",
				Unit:           types.UnitSeconds,
				Category:       "Efficiency",
				Dimension:      metrictypes.MetricDimensionTimeToMerge,
				Operation:      metrictypes.MetricOperationMedian,
				IconIdentifier: "clock",
				IconColor:      "green",
			}, sourceControlDB),
			engine.NewPRReviewComplexityRule(metrictypes.BaseMetricRule{
				ID:             "avg_pr_review_loc",
				Name:           "Average PR Review LoC",
				Description:    "Average Lines of Code (LoC) in PRs reviewed by the member, calculated as the average of (lines added + lines removed) across all reviewed PRs. Excludes PRs authored by any member in the organization. Peer comparison shows the median PR review complexity across other organization members.",
				Unit:           types.UnitCount,
				Category:       "Collaboration",
				Dimension:      metrictypes.MetricDimensionPRReviewComplexity,
				Operation:      metrictypes.MetricOperationAverage,
				IconIdentifier: "bar-chart-3",
				IconColor:      "orange",
			}, sourceControlDB),
		},
		sourceControlDB: sourceControlDB,
	}
}

func (e *Engine) CalculateMetrics(ctx context.Context, params types.MetricRuleParams) (*types.MetricsResponse, error) {
	// Use maps to group metrics by category
	snapshotCategoriesMap := make(map[string]*types.SnapshotCategory)
	graphCategoriesMap := make(map[string]*types.GraphCategory)

	for _, rule := range e.Metrics {
		snapshotMetric, graphMetric, err := rule.Calculate(ctx, params)
		if err != nil {
			return nil, err
		}

		category := rule.Category()

		// Group snapshot metrics by category
		if snapshotCategoriesMap[category] == nil {
			snapshotCategoriesMap[category] = &types.SnapshotCategory{
				Category: category,
				Metrics:  []types.SnapshotMetric{},
			}
		}
		snapshotCategoriesMap[category].Metrics = append(snapshotCategoriesMap[category].Metrics, *snapshotMetric)

		// Group graph metrics by category
		if graphCategoriesMap[category] == nil {
			graphCategoriesMap[category] = &types.GraphCategory{
				Category: category,
				Metrics:  []types.GraphMetric{},
			}
		}
		graphCategoriesMap[category].Metrics = append(graphCategoriesMap[category].Metrics, *graphMetric)
	}

	// Convert maps to slices
	snapshotMetrics := make([]*types.SnapshotCategory, 0, len(snapshotCategoriesMap))
	graphMetrics := make([]*types.GraphCategory, 0, len(graphCategoriesMap))

	for _, category := range snapshotCategoriesMap {
		snapshotMetrics = append(snapshotMetrics, category)
	}

	for _, category := range graphCategoriesMap {
		graphMetrics = append(graphMetrics, category)
	}

	return &types.MetricsResponse{
		SnapshotMetrics: snapshotMetrics,
		GraphMetrics:    graphMetrics,
	}, nil
}
