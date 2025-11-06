import { create, StateCreator } from 'zustand'
import { Organization } from '../types/organization'
import Api from '../api/api'

interface OrganizationState {
  organizations: Organization[]
  currentOrganization: Organization | null
  isLoading: boolean
  error: string | null
  fetchOrganizations: () => Promise<void>
  setCurrentOrganization: (org: Organization | null) => void
  createOrganization: (name: string, slug: string) => Promise<Organization>
}

const logger = <T extends OrganizationState>(
  config: StateCreator<T>
): StateCreator<T> => (set, get, store) => {
  const setWithLog = (partial: Partial<T> | ((state: T) => Partial<T>)) => {
    const previousState = get()
    const nextState = typeof partial === 'function' ? partial(previousState) : partial
    console.log('Organization Store - State Changed:', {
      previous: previousState,
      next: nextState,
      action: new Error().stack?.split('\n')[2]?.trim() || 'unknown'
    })
    set(nextState)
  }
  return config(setWithLog, get, store)
}

export const useOrganizationStore = create<OrganizationState>()(
  logger((set, get) => ({
    organizations: [],
    currentOrganization: null,
    isLoading: false,
    error: null,

    fetchOrganizations: async () => {
      // Do not fetch organizations if they are already being fetched
      if (get().isLoading) {
        return
      }

      set({ isLoading: true, error: null })
      try {
        const orgs = await Api.getOrganizations()
        set({ organizations: orgs })
        
        // If we have organizations but no current organization, set the first one
        if (orgs.length > 0 && !get().currentOrganization) {
          set({ currentOrganization: orgs[0] })
        }
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to fetch organizations'
        console.error('Error fetching organizations:', err)
        // Don't set error state for auth errors that might be temporary
        if (errorMessage.includes('Authentication failed')) {
          set({ error: null }) // Clear error, will retry on next auth check
        } else {
          set({ error: errorMessage })
        }
      } finally {
        set({ isLoading: false })
      }
    },

    setCurrentOrganization: (org) => {
      set({ currentOrganization: org })
    },

    createOrganization: async (name: string, slug: string) => {
      set({ isLoading: true, error: null })
      try {
        const newOrg = await Api.createOrganization(name, slug)
        set((state) => ({
          organizations: [...state.organizations, newOrg],
          currentOrganization: newOrg,
        }))
        return newOrg
      } catch (err) {
        console.error('Error creating organization:', err)
        set({ error: err instanceof Error ? err.message : 'Failed to create organization' })
        throw err
      } finally {
        set({ isLoading: false })
      }
    },
  }))
)