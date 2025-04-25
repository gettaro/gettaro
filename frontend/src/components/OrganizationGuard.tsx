import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import { useOrganizationStore } from '../stores/organization'

export default function OrganizationGuard({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate()
  const { isAuthenticated, isLoading: isAuthLoading, getToken } = useAuth()
  const { currentOrganization, organizations, isLoading, fetchOrganizations } = useOrganizationStore()

  useEffect(() => {
    // organizationAPI.setAuth(getAccessTokenSilently)
  }, [getAccessTokenSilently])

  useEffect(() => {
    if (isAuthenticated) {
      fetchOrganizations()
    }
  }, [isAuthenticated, fetchOrganizations])

  useEffect(() => {
    if (!isAuthLoading && !isLoading) {
      if (!isAuthenticated) {
        return
      }
      if (organizations.length === 0) {
        navigate('/create-organization', { replace: true })
      } else if (!currentOrganization) {
        navigate('/select-organization', { replace: true })
      }
    }
  }, [isAuthenticated, isAuthLoading, isLoading, organizations.length, currentOrganization, navigate])

  // Show loading state while checking auth or fetching organizations
  if (isAuthLoading || isLoading) {
    return null
  }

  // If not authenticated, let the auth guard handle it
  if (!isAuthenticated) {
    return <>{children}</>
  }

  // If no organizations exist or no organization is selected, show nothing (navigation is handled in useEffect)
  if (organizations.length === 0 || !currentOrganization) {
    return null
  }

  return <>{children}</>
} 