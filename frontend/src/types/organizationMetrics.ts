import { GetMemberMetricsResponse, SnapshotCategory, GraphCategory } from './memberMetrics'

// Team metrics breakdown
export interface TeamMetricsBreakdown {
  team_id: string
  team_name: string
  snapshot_metrics: SnapshotCategory[]
  graph_metrics: GraphCategory[]
}

// Organization metrics response with team breakdown
export interface OrganizationMetricsResponse extends GetMemberMetricsResponse {
  teams_breakdown?: TeamMetricsBreakdown[]
}

