import { useState, useEffect } from 'react'
import OrganizationDropdown from './OrganizationDropdown'
import { useAuth } from '../hooks/useAuth'
// import CreateOrganizationModal from './CreateOrganizationModal'
// import { Organization } from '../types/organization'
// import { createOrganization } from '../api/organizations'

export default function Navigation() {
  const { isAuthenticated, user, login, logout } = useAuth()
  const [isProfileDropdownOpen, setIsProfileDropdownOpen] = useState(false)
  const [isOrgDropdownOpen, setIsOrgDropdownOpen] = useState(false)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    console.log(user)
  }, [user])

  return (
    <nav className="flex items-center justify-between w-full">
      <div className="flex items-center space-x-4">
      </div>
      <div className="flex items-center space-x-4">
        {isAuthenticated && (
          <>
            <OrganizationDropdown
              isOpen={isOrgDropdownOpen}
              onToggle={() => setIsOrgDropdownOpen(!isOrgDropdownOpen)}
            />
            {/* <CreateOrganizationModal
              isOpen={isCreateModalOpen}
              onClose={() => setIsCreateModalOpen(false)}
            /> */}
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
                    onClick={() => logout()}
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
            onClick={() => login()}
            className="bg-primary text-primary-foreground hover:bg-primary/90 px-4 py-2 rounded-md"
          >
            Log In
          </button>
        )}
      </div>
    </nav>
  )
} 