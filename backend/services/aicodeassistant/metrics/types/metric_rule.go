package types

import (
	"context"

	"ems.dev/backend/services/aicodeassistant/types"
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
	MetricDimensionLinesOfCodeAccepted  MetricDimension = "LINES_OF_CODE_ACCEPTED"
	MetricDimensionLinesOfCodeSuggested MetricDimension = "LINES_OF_CODE_SUGGESTED"
	MetricDimensionActiveSessions       MetricDimension = "ACTIVE_SESSIONS"
	MetricDimensionAcceptRate           MetricDimension = "ACCEPT_RATE"
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
