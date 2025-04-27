export type IntegrationProvider = 'github'

export interface IntegrationConfig {
  id: string
  organizationId: string
  providerName: IntegrationProvider
  providerType: string
  encryptedToken: string
  metadata: Record<string, any>
  lastSyncedAt?: string
  createdAt: string
  updatedAt: string
}

export interface CreateIntegrationConfigRequest {
  providerName: IntegrationProvider
  token: string
  metadata?: Record<string, any>
}

export interface UpdateIntegrationConfigRequest {
  token?: string
  metadata?: Record<string, any>
} 