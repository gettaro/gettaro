// Request parameters for getting member metrics
export interface GetMemberMetricsParams {
  startDate?: string // YYYY-MM-DD format
  endDate?: string   // YYYY-MM-DD format
  interval?: string  // daily, weekly, monthly
}

// Snapshot metric for comparison
export interface SnapshotMetric {
  label: string
  value: number
  peersValue: number
  unit: string // "count", "time", "loc", etc.
}

// Category of snapshot metrics
export interface SnapshotCategory {
  category: string
  metrics: SnapshotMetric[]
}

// Time series data point
export interface TimeSeriesDataPoint {
  key: string
  value: number
}

// Time series entry
export interface TimeSeriesEntry {
  date: string
  data: TimeSeriesDataPoint[]
}

// Graph metric for visualization
export interface GraphMetric {
  label: string
  type: string
  timeSeries: TimeSeriesEntry[]
}

// Category of graph metrics
export interface GraphCategory {
  category: string
  metrics: GraphMetric[]
}

// Response for member metrics
export interface GetMemberMetricsResponse {
  snapshotMetrics: SnapshotCategory[]
  graphMetrics: GraphCategory[]
} 