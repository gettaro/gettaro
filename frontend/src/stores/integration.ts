import { create } from 'zustand'
import Api from '../api/api'
import { IntegrationConfig } from '../types/integration'

interface IntegrationState {
  integrations: IntegrationConfig[]
  isLoading: boolean
  error: string | null
  fetchIntegrations: (organizationId: string) => Promise<void>
}

export const useIntegrationStore = create<IntegrationState>((set) => ({
  integrations: [],
  isLoading: false,
  error: null,
  fetchIntegrations: async (organizationId: string) => {
    try {
      set({ isLoading: true, error: null })
      const integrations = await Api.getOrganizationIntegrations(organizationId)
      set({ integrations, isLoading: false })
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to fetch integrations', isLoading: false })
    }
  },
})) 