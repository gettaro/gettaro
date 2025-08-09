export interface SourceControlAccount {
  id: string
  userId?: string
  organizationId?: string
  providerName: string
  providerId: string
  username: string
  lastSyncedAt?: string
} 