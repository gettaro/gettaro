export type IntegrationProvider = 'github'

export interface IntegrationConfig {
  id: string
  organizationId: string
  providerName: string
  providerType: string
  encryptedToken: string
  metadata: Record<string, unknown>
  createdAt: string
  updatedAt: string
}

export interface CreateIntegrationRequest {
  providerName: string
  providerType: string
  token: string
  metadata?: Record<string, unknown>
}

export interface UpdateIntegrationRequest {
  token?: string
  metadata?: Record<string, unknown>
} 