export interface SourceControlAccount {
  id: string
  memberId?: string
  organizationId?: string
  providerName: string
  providerId: string
  username: string
  lastSyncedAt?: string
} 