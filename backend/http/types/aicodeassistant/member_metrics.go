package aicodeassistant

import "ems.dev/backend/services/aicodeassistant/types"

// GetMemberMetricsRequest represents the request parameters for getting member metrics
type GetMemberMetricsRequest struct {
	StartDate string `form:"startDate" binding:"omitempty"` // YYYY-MM-DD format
	EndDate   string `form:"endDate" binding:"omitempty"`   // YYYY-MM-DD format
	Interval  string `form:"interval" binding:"omitempty"`  // daily, weekly, monthly
}

// SnapshotMetric represents a single metric in the snapshot
type SnapshotMetric struct {
	Label      string     `json:"label"`
	Value      float64    `json:"value"`
	PeersValue float64    `json:"peers_value"`
	Unit       types.Unit `json:"unit"` // "count", "time" etc.
}

// SnapshotCategory represents a category of metrics in the snapshot
type SnapshotCategory struct {
	Category string           `json:"category"`
	Metrics  []SnapshotMetric `json:"metrics"`
}

// Reuse TimeSeriesDataPoint and TimeSeriesEntry from service types
type TimeSeriesDataPoint = types.TimeSeriesDataPoint
type TimeSeriesEntry = types.TimeSeriesEntry

// GraphMetric represents a single metric in the graph data
type GraphMetric struct {
	Label      string                      `json:"label"`
	TimeSeries []types.TimeSeriesEntry     `json:"time_series"`
}

// GraphCategory represents a category of metrics in the graph data
type GraphCategory struct {
	Category string        `json:"category"`
	Metrics  []GraphMetric `json:"metrics"`
}

// GetMemberMetricsResponse represents the response for getting member metrics
type GetMemberMetricsResponse struct {
	SnapshotMetrics []SnapshotCategory `json:"snapshot_metrics"`
	GraphMetrics    []GraphCategory    `json:"graph_metrics"`
}

// MarshalMetricsResponse converts a service MetricsResponse to HTTP response type
func MarshalMetricsResponse(metrics *types.MetricsResponse) *GetMemberMetricsResponse {
	response := &GetMemberMetricsResponse{
		SnapshotMetrics: make([]SnapshotCategory, len(metrics.SnapshotMetrics)),
		GraphMetrics:    make([]GraphCategory, len(metrics.GraphMetrics)),
	}

	for i, snapshotCategory := range metrics.SnapshotMetrics {
		response.SnapshotMetrics[i] = SnapshotCategory{
			Category: snapshotCategory.Category.Name,
			Metrics:  make([]SnapshotMetric, len(snapshotCategory.Metrics)),
		}
		for j, metric := range snapshotCategory.Metrics {
			response.SnapshotMetrics[i].Metrics[j] = SnapshotMetric{
				Label:      metric.Label,
				Value:      metric.Value,
				PeersValue: metric.PeersValue,
				Unit:       metric.Unit,
			}
		}
	}

	for i, graphCategory := range metrics.GraphMetrics {
		response.GraphMetrics[i] = GraphCategory{
			Category: graphCategory.Category.Name,
			Metrics:  make([]GraphMetric, len(graphCategory.Metrics)),
		}
		for j, metric := range graphCategory.Metrics {
			// Reuse TimeSeriesEntry directly from service types (they're aliased)
			response.GraphMetrics[i].Metrics[j] = GraphMetric{
				Label:      metric.Label,
				TimeSeries: metric.TimeSeries,
			}
		}
	}

	return response
}

