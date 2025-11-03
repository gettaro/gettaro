export interface AICodeAssistantDailyMetric {
  id: string
  organization_id: string
  external_account_id: string
  tool_name: string
  metric_date: string
  lines_of_code_accepted: number
  lines_of_code_suggested: number
  suggestion_accept_rate?: number
  active_sessions: number
  metadata?: Record<string, any>
  created_at: string
  updated_at: string
}

export interface GetMemberAICodeAssistantUsageParams {
  toolName?: string
  startDate?: string
  endDate?: string
}

export interface GetMemberAICodeAssistantUsageResponse {
  metrics: AICodeAssistantDailyMetric[]
}

// Metrics types (reusing structure from memberMetrics)
export interface GetMemberAICodeAssistantMetricsParams {
  startDate?: string
  endDate?: string
  interval?: 'daily' | 'weekly' | 'monthly'
}

export interface SnapshotMetric {
  label: string
  value: number
  peers_value: number
  unit: 'count' | 'seconds' | 'percent'
}

export interface SnapshotCategory {
  category: string
  metrics: SnapshotMetric[]
}

export interface TimeSeriesDataPoint {
  key: string
  value: number
}

export interface TimeSeriesEntry {
  date: string
  data: TimeSeriesDataPoint[]
}

export interface GraphMetric {
  label: string
  type?: string
  unit?: string
  time_series: TimeSeriesEntry[]
}

export interface GraphCategory {
  category: string
  metrics: GraphMetric[]
}

export interface GetMemberAICodeAssistantMetricsResponse {
  snapshot_metrics: SnapshotCategory[]
  graph_metrics: GraphCategory[]
}


