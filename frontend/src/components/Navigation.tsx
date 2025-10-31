import { useState, useEffect } from 'react'
import OrganizationDropdown from './OrganizationDropdown'
import { useAuth } from '../hooks/useAuth'
import { useTheme } from '../hooks/useTheme'
// import CreateOrganizationModal from './CreateOrganizationModal'
// import { Organization } from '../types/organization'
// import { createOrganization } from '../api/organizations'

export default function Navigation() {
  const { isAuthenticated, user, login, logout } = useAuth()
  const { theme, toggleTheme } = useTheme()
  const [isProfileDropdownOpen, setIsProfileDropdownOpen] = useState(false)
  
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
            <OrganizationDropdown />
            <button
              onClick={toggleTheme}
              className="flex items-center justify-center w-8 h-8 rounded-md hover:bg-muted transition-colors"
              title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
            >
              {theme === 'dark' ? (
                <svg
                  className="w-5 h-5 text-muted-foreground"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
                  />
                </svg>
              ) : (
                <svg
                  className="w-5 h-5 text-muted-foreground"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
                  />
                </svg>
              )}
            </button>
            <button
              onClick={() => window.location.href = '/settings'}
              className="flex items-center justify-center w-8 h-8 rounded-md hover:bg-muted transition-colors"
              title="Settings"
            >
              <svg
                className="w-5 h-5 text-muted-foreground"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
              </svg>
            </button>
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