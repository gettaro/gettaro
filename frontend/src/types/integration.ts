export type IntegrationProvider = 'github'

export interface IntegrationConfig {
  id: string
  organization_id: string
  provider_name: IntegrationProvider
  provider_type: string
  encrypted_token: string
  metadata: Record<string, any>
  last_synced_at?: string
  created_at: string
  updated_at: string
}

export interface CreateIntegrationConfigRequest {
  provider_name: IntegrationProvider
  token: string
  metadata?: Record<string, any>
}

export interface UpdateIntegrationConfigRequest {
  token?: string
  metadata?: Record<string, any>
} 