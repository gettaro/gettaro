export interface MemberActivity {
  id: string
  type: 'pull_request' | 'pr_review' | 'pr_comment'
  title: string
  description?: string
  url?: string
  repository?: string
  createdAt: string
  metadata?: Record<string, any>
  authorUsername?: string
  prTitle?: string          // For comments/reviews: the PR title
  prAuthorUsername?: string // For comments/reviews: the PR author
  prMetrics?: Record<string, any> // PR performance metrics
}

export interface GetMemberActivityResponse {
  activities: MemberActivity[]
}

export interface GetMemberActivityParams {
  startDate?: string // YYYY-MM-DD format
  endDate?: string   // YYYY-MM-DD format
} 