package metrics

import (
	"context"
	"sort"

	"ems.dev/backend/services/aicodeassistant/database"
	"ems.dev/backend/services/aicodeassistant/metrics/engine"
	metrictypes "ems.dev/backend/services/aicodeassistant/metrics/types"
	"ems.dev/backend/services/aicodeassistant/types"
)

type MetricsEngine interface {
	CalculateMetrics(ctx context.Context, params types.MetricRuleParams) (*types.MetricsResponse, error)
}

type Engine struct {
	Metrics        []metrictypes.MetricRule `json:"metrics"`
	aicodeassistantDB database.DB
}

func NewEngine(aicodeassistantDB database.DB) *Engine {
	// Define metric categories
	categories := map[string]types.MetricRuleCategory{
		"Usage": {
			Name:     "Usage",
			Priority: 1,
		},
		"Efficiency": {
			Name:     "Efficiency",
			Priority: 2,
		},
	}

	return &Engine{
		Metrics: []metrictypes.MetricRule{
			engine.NewLinesOfCodeAcceptedRule(metrictypes.BaseMetricRule{
				ID:             "lines_of_code_accepted_count",
				Name:           "Lines of Code Accepted",
				Description:    "Total lines of code accepted from AI suggestions. This metric counts all lines from accepted AI code suggestions across all tools. Peer comparison shows the median lines accepted across other organization members.",
				Unit:           types.UnitCount,
				Category:       categories["Usage"],
				Dimension:      metrictypes.MetricDimensionLinesOfCodeAccepted,
				Operation:      metrictypes.MetricOperationCount,
				IconIdentifier: "check-circle",
				IconColor:      "green",
			}, aicodeassistantDB),
			engine.NewLinesOfCodeSuggestedRule(metrictypes.BaseMetricRule{
				ID:             "lines_of_code_suggested_count",
				Name:           "Lines of Code Suggested",
				Description:    "Total lines of code suggested by AI assistants. This metric counts all lines suggested by AI code assistants across all tools. Peer comparison shows the median lines suggested across other organization members.",
				Unit:           types.UnitCount,
				Category:       categories["Usage"],
				Dimension:      metrictypes.MetricDimensionLinesOfCodeSuggested,
				Operation:      metrictypes.MetricOperationCount,
				IconIdentifier: "code",
				IconColor:      "blue",
			}, aicodeassistantDB),
			engine.NewActiveSessionsRule(metrictypes.BaseMetricRule{
				ID:             "active_sessions_count",
				Name:           "Active Sessions",
				Description:    "Total number of active AI code assistant sessions. This metric counts all active sessions across all tools. Peer comparison shows the median active sessions across other organization members.",
				Unit:           types.UnitCount,
				Category:       categories["Usage"],
				Dimension:      metrictypes.MetricDimensionActiveSessions,
				Operation:      metrictypes.MetricOperationCount,
				IconIdentifier: "activity",
				IconColor:      "purple",
			}, aicodeassistantDB),
			engine.NewAcceptRateRule(metrictypes.BaseMetricRule{
				ID:             "accept_rate_percent",
				Name:           "Accept Rate",
				Description:    "Average percentage of AI suggestions that were accepted. Calculated as (lines accepted / lines suggested) * 100. Peer comparison shows the median accept rate across other organization members.",
				Unit:           types.UnitPercent,
				Category:       categories["Efficiency"],
				Dimension:      metrictypes.MetricDimensionAcceptRate,
				Operation:      metrictypes.MetricOperationAverage,
				IconIdentifier: "trending-up",
				IconColor:      "orange",
			}, aicodeassistantDB),
		},
		aicodeassistantDB: aicodeassistantDB,
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
		if snapshotCategoriesMap[category.Name] == nil {
			snapshotCategoriesMap[category.Name] = &types.SnapshotCategory{
				Category: category,
				Metrics:  []types.SnapshotMetric{},
			}
		}
		snapshotCategoriesMap[category.Name].Metrics = append(snapshotCategoriesMap[category.Name].Metrics, *snapshotMetric)

		// Group graph metrics by category
		if graphCategoriesMap[category.Name] == nil {
			graphCategoriesMap[category.Name] = &types.GraphCategory{
				Category: category,
				Metrics:  []types.GraphMetric{},
			}
		}
		graphCategoriesMap[category.Name].Metrics = append(graphCategoriesMap[category.Name].Metrics, *graphMetric)
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

	// Sort categories by priority (ascending)
	sort.Slice(snapshotMetrics, func(i, j int) bool {
		return snapshotMetrics[i].Category.Priority < snapshotMetrics[j].Category.Priority
	})

	sort.Slice(graphMetrics, func(i, j int) bool {
		return graphMetrics[i].Category.Priority < graphMetrics[j].Category.Priority
	})

	return &types.MetricsResponse{
		SnapshotMetrics: snapshotMetrics,
		GraphMetrics:    graphMetrics,
	}, nil
}

