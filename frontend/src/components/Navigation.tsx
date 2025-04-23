import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import OrganizationDropdown from './OrganizationDropdown'
import CreateOrganizationModal from './CreateOrganizationModal'
import { Organization } from '../types/organization'
import { getOrganizations, createOrganization } from '../api/organizations'

export default function Navigation() {
  const { isAuthenticated, user, loginWithRedirect, logout, getAccessTokenSilently } = useAuth0()
  const [isProfileDropdownOpen, setIsProfileDropdownOpen] = useState(false)
  const [isOrgDropdownOpen, setIsOrgDropdownOpen] = useState(false)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [currentOrganization, setCurrentOrganization] = useState<Organization | null>(null)
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (isAuthenticated) {
      fetchOrganizations()
    }
  }, [isAuthenticated])

  const fetchOrganizations = async () => {
    setIsLoading(true)
    setError(null)
    try {
      const orgs = await getOrganizations(getAccessTokenSilently)
      setOrganizations(orgs)
      if (orgs.length > 0 && !currentOrganization) {
        setCurrentOrganization(orgs[0])
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch organizations')
    } finally {
      setIsLoading(false)
    }
  }

  const handleSelectOrganization = (org: Organization) => {
    setCurrentOrganization(org)
    setIsOrgDropdownOpen(false)
  }

  const handleCreateOrganization = async (name: string, slug: string) => {
    try {
      const newOrg = await createOrganization(name, slug, getAccessTokenSilently)
      setOrganizations([...organizations, newOrg])
      setCurrentOrganization(newOrg)
    } catch (err) {
      throw err
    }
  }

  return (
    <nav className="flex items-center justify-between w-full">
      <div className="flex items-center space-x-4">
        <Link to="/dashboard" className="text-foreground hover:text-primary">Dashboard</Link>
      </div>
      <div className="flex items-center space-x-4">
        {isAuthenticated && (
          <>
            <OrganizationDropdown
              organizations={organizations}
              currentOrganization={currentOrganization}
              onSelectOrganization={handleSelectOrganization}
              onCreateOrganization={() => setIsCreateModalOpen(true)}
              isOpen={isOrgDropdownOpen}
              onToggle={() => setIsOrgDropdownOpen(!isOrgDropdownOpen)}
            />
            <CreateOrganizationModal
              isOpen={isCreateModalOpen}
              onClose={() => setIsCreateModalOpen(false)}
              onCreate={handleCreateOrganization}
            />
            <div className="relative">
              <button
                onClick={() => setIsProfileDropdownOpen(!isProfileDropdownOpen)}
                className="flex items-center focus:outline-none"
              >
                <img
                  src={user?.picture}
                  alt={user?.name}
                  className="w-8 h-8 rounded-full"
                />
              </button>
              {isProfileDropdownOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-card rounded-md shadow-lg py-1 z-50">
                  <div className="px-4 py-2 text-sm text-muted-foreground border-b border-border">
                    {user?.name}
                  </div>
                  <button
                    onClick={() => logout({ logoutParams: { returnTo: window.location.origin } })}
                    className="block w-full text-left px-4 py-2 text-sm text-foreground hover:bg-primary hover:text-primary-foreground"
                  >
                    Log Out
                  </button>
                </div>
              )}
            </div>
          </>
        )}
        {!isAuthenticated && (
          <button
            onClick={() => loginWithRedirect()}
            className="bg-primary text-primary-foreground hover:bg-primary/90 px-4 py-2 rounded-md"
          >
            Log In
          </button>
        )}
      </div>
    </nav>
  )
} 