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