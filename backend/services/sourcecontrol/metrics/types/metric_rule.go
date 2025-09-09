package types

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

type MetricOperation string

const (
	MetricOperationAverage MetricOperation = "AVG"
	MetricOperationCount   MetricOperation = "COUNT"
	MetricOperationMedian  MetricOperation = "MEDIAN"
	MetricOperationMax     MetricOperation = "MAX"
	MetricOperationMin     MetricOperation = "MIN"
)

type MetricDimension string

const (
	MetricDimensionTimeToMerge        MetricDimension = "TIME_TO_MERGE"
	MetricDimensionTimeToFirstReview  MetricDimension = "TIME_TO_FIRST_REVIEW"
	MetricDimensionMergedPRs          MetricDimension = "MERGED_PRS"
	MetricDimensionReviwedPRs         MetricDimension = "REVIEWED_PRS"
	MetricDimensionLOCAdded           MetricDimension = "LOC_ADDED"
	MetricDimensionLOCRemoved         MetricDimension = "LOC_REMOVED"
	MetricDimensionPRReviewComplexity MetricDimension = "PR_REVIEW_COMPLEXITY"
)

type MetricRule interface {
	Category() types.MetricRuleCategory
	Calculate(ctx context.Context, params types.MetricRuleParams) (*types.SnapshotMetric, *types.GraphMetric, error)
}

type BaseMetricRule struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Category       types.MetricRuleCategory `json:"category"`
	Unit           types.Unit               `json:"unit"`
	Dimension      MetricDimension          `json:"dimension"`
	Operation      MetricOperation          `json:"operation"`
	IconIdentifier string                   `json:"icon_identifier"`
	IconColor      string                   `json:"icon_color"`
}
