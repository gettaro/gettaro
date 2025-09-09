export interface SourceControlAccount {
  id: string
  member_id?: string
  organization_id?: string
  provider_name: string
  provider_id: string
  username: string
  last_synced_at?: string
}

export interface PullRequest {
  id: string
  title: string
  description: string
  url: string
  status: string
  created_at: string
  merged_at?: string
  comments: number
  review_comments: number
  additions: number
  deletions: number
  changed_files: number
}

export interface GetMemberPullRequestsParams {
  startDate?: string
  endDate?: string
}

export interface GetMemberPullRequestReviewsParams {
  startDate?: string
  endDate?: string
}

export interface MemberActivity {
  id: string
  type: 'pull_request' | 'pr_review' | 'pr_comment'
  title: string
  description?: string
  url?: string
  repository?: string
  created_at: string
  metadata?: Record<string, any>
  author_username?: string
  pr_title?: string          // For comments/reviews: the PR title
  pr_author_username?: string // For comments/reviews: the PR author
  pr_metrics?: Record<string, any> // PR performance metrics
}

export interface GetMemberActivityParams {
  startDate?: string // YYYY-MM-DD format
  endDate?: string   // YYYY-MM-DD format
} 