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

export interface AICodeAssistantUsageStats {
  total_lines_accepted: number
  total_lines_suggested: number
  overall_accept_rate: number
  active_sessions: number
}

export interface GetMemberAICodeAssistantUsageParams {
  toolName?: string
  startDate?: string
  endDate?: string
}

export interface GetMemberAICodeAssistantUsageResponse {
  metrics: AICodeAssistantDailyMetric[]
}

export interface GetMemberAICodeAssistantUsageStatsResponse {
  stats: AICodeAssistantUsageStats
}


