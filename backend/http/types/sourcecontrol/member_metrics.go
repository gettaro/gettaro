package sourcecontrol

// GetMemberMetricsRequest represents the request parameters for getting member metrics
type GetMemberMetricsRequest struct {
	StartDate string `form:"startDate" binding:"omitempty"` // YYYY-MM-DD format
	EndDate   string `form:"endDate" binding:"omitempty"`   // YYYY-MM-DD format
	Interval  string `form:"interval" binding:"omitempty"`  // daily, weekly, monthly
}

// SnapshotMetric represents a single metric in the snapshot
type SnapshotMetric struct {
	Label      string  `json:"label"`
	Value      float64 `json:"value"`
	PeersValue float64 `json:"peersValue"`
	Unit       string  `json:"unit"` // "count", "time" etc.
}

// SnapshotCategory represents a category of metrics in the snapshot
type SnapshotCategory struct {
	Category string           `json:"category"`
	Metrics  []SnapshotMetric `json:"metrics"`
}

// TimeSeriesDataPoint represents a single data point in a time series
type TimeSeriesDataPoint struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

// TimeSeriesEntry represents a single time entry in a time series
type TimeSeriesEntry struct {
	Date string                `json:"date"`
	Data []TimeSeriesDataPoint `json:"data"`
}

// GraphMetric represents a single metric in the graph data
type GraphMetric struct {
	Label      string            `json:"label"`
	Type       string            `json:"type"`
	TimeSeries []TimeSeriesEntry `json:"timeSeries"`
}

// GraphCategory represents a category of metrics in the graph data
type GraphCategory struct {
	Category string        `json:"category"`
	Metrics  []GraphMetric `json:"metrics"`
}

// GetMemberMetricsResponse represents the response for getting member metrics
type GetMemberMetricsResponse struct {
	SnapshotMetrics []SnapshotCategory `json:"snapshotMetrics"`
	GraphMetrics    []GraphCategory    `json:"graphMetrics"`
}
