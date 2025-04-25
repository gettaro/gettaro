import { useState, useEffect } from 'react'
import { Organization } from '../types/organization'
import { useOrganizationStore } from '../stores/organization'
import CreateOrganizationModal from './CreateOrganizationModal'

export default function OrganizationDropdown() {
  const [isOpen, setIsOpen] = useState(false)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const { currentOrganization, organizations, setCurrentOrganization, fetchOrganizations, isLoading: isLoadingOrgs, error } = useOrganizationStore()

  useEffect(() => {
      fetchOrganizations()
  }, [fetchOrganizations])

  const handleSelectOrganization = (orgId: string) => {
    const org = organizations.find((o: Organization) => o.id === orgId)
    if (org) {
      setCurrentOrganization(org)
      setIsOpen(false)
    }
  }

  if (isLoadingOrgs) {
    return <div className="text-foreground">Loading...</div>
  }

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center space-x-2 text-foreground hover:text-primary"
      >
        {currentOrganization ? (
          <>
            {currentOrganization.logo ? (
              <img
                src={currentOrganization.logo}
                alt={currentOrganization.name}
                className="w-6 h-6 rounded"
              />
            ) : (
              <div className="w-6 h-6 rounded bg-primary/10 flex items-center justify-center">
                <span className="text-xs font-medium">{currentOrganization.name[0]}</span>
              </div>
            )}
            <span className="font-medium">{currentOrganization.name}</span>
          </>
        ) : (
          <span className="font-medium">Select Organization</span>
        )}
      </button>
      {isOpen && (
        <div className="absolute right-0 mt-2 w-48 bg-card/50 backdrop-blur-sm rounded-lg border py-1 z-50">
          {organizations.map((org: Organization) => (
            <button
              key={org.id}
              onClick={() => handleSelectOrganization(org.id)}
              className="flex items-center w-full px-4 py-2 text-sm text-foreground hover:bg-primary hover:text-primary-foreground"
            >
              {org.logo ? (
                <img
                  src={org.logo}
                  alt={org.name}
                  className="w-5 h-5 rounded mr-2"
                />
              ) : (
                <div className="w-5 h-5 rounded bg-primary/10 flex items-center justify-center mr-2">
                  <span className="text-xs font-medium">{org.name[0]}</span>
                </div>
              )}
              <span>{org.name}</span>
            </button>
          ))}
          <div className="border-t border-border my-1" />
          <button
            onClick={() => {
              setIsOpen(false)
              setIsCreateModalOpen(true)
            }}
            className="flex items-center w-full px-4 py-2 text-sm text-foreground hover:bg-primary hover:text-primary-foreground"
          >
            <span className="mr-2">+</span>
            Create Organization
          </button>
        </div>
      )}
      <CreateOrganizationModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onCreate={async (name: string, slug: string) => {
          const { createOrganization } = useOrganizationStore.getState()
          await createOrganization(name, slug)
          setIsCreateModalOpen(false)
        }}
      />
    </div>
  )
} 