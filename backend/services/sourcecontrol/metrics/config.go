package metrics

import (
	metrictypes "ems.dev/backend/services/sourcecontrol/metrics/types"
	"ems.dev/backend/services/sourcecontrol/types"
)

type Config struct {
	Metrics []metrictypes.BaseMetricRule `json:"metrics"`
}

func NewConfig() *Config {
	return &Config{
		Metrics: []metrictypes.BaseMetricRule{
			{
				Name:      "Median time to merge",
				Unit:      types.UnitSeconds,
				Category:  "Efficiency",
				Dimension: metrictypes.MetricDimensionTimeToMerge,
				Operation: metrictypes.MetricOperationMedian,
			},
		},
	}
}
