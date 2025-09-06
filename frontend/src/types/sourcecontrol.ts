export interface SourceControlAccount {
  id: string
  memberId?: string
  organizationId?: string
  providerName: string
  providerId: string
  username: string
  lastSyncedAt?: string
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
  createdAt: string
  metadata?: Record<string, any>
  authorUsername?: string
  prTitle?: string          // For comments/reviews: the PR title
  prAuthorUsername?: string // For comments/reviews: the PR author
  prMetrics?: Record<string, any> // PR performance metrics
}

export interface GetMemberActivityParams {
  startDate?: string // YYYY-MM-DD format
  endDate?: string   // YYYY-MM-DD format
} 