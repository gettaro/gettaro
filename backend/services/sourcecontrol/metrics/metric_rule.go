package metrics

import (
	"context"

	"ems.dev/backend/services/sourcecontrol/types"
)

type MetricRule interface {
	ValidateParams(ctx context.Context, params types.MetricRuleParams) error
	Calculate(ctx context.Context, params types.MetricRuleParams) (types.SnapshotMetric, *types.GraphMetric, error)
}
