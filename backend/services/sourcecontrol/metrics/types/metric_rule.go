package types

import "ems.dev/backend/services/sourcecontrol/types"

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
	MetricDimensionTimeToMerge       MetricDimension = "TIME_TO_MERGE"
	MetricDimensionTimeToFirstReview MetricDimension = "TIME_TO_FIRST_REVIEW"
	MetricDimensionMergedPRs         MetricDimension = "MERGED_PRS"
	MetricDimensionReviwedPRs        MetricDimension = "REVIEWED_PRS"
	MetricDimensionLOCAdded          MetricDimension = "LOC_ADDED"
	MetricDimensionLOCRemoved        MetricDimension = "LOC_REMOVED"
)

type BaseMetricRule struct {
	Name      string          `json:"name"`
	Category  string          `json:"category"`
	Unit      types.Unit      `json:"unit"`
	Dimension MetricDimension `json:"dimension"`
	Operation MetricOperation `json:"operation"`
}
