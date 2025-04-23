import React from 'react'
import { Organization, OrganizationDropdownProps } from '../types/organization'

export default function OrganizationDropdown({
  organizations,
  currentOrganization,
  onSelectOrganization,
  onCreateOrganization,
  isOpen,
  onToggle,
}: OrganizationDropdownProps) {
  return (
    <div className="relative">
      <button
        onClick={onToggle}
        className="flex items-center space-x-2 text-foreground hover:text-primary"
      >
        {currentOrganization?.logo ? (
          <img
            src={currentOrganization.logo}
            alt={currentOrganization.name}
            className="w-6 h-6 rounded"
          />
        ) : (
          <div className="w-6 h-6 rounded bg-primary/10 flex items-center justify-center">
            <span className="text-xs font-medium">{currentOrganization?.name?.[0] || 'O'}</span>
          </div>
        )}
        <span className="font-medium">{currentOrganization?.name || 'Select Organization'}</span>
      </button>
      {isOpen && (
        <div className="absolute right-0 mt-2 w-56 bg-card rounded-md shadow-lg py-1 z-50">
          <div className="px-4 py-2 text-sm text-muted-foreground border-b border-border">
            Organizations
          </div>
          {organizations.map((org) => (
            <button
              key={org.id}
              onClick={() => onSelectOrganization(org)}
              className="flex items-center w-full px-4 py-2 text-sm text-foreground hover:bg-primary hover:text-primary-foreground"
            >
              {org.logo ? (
                <img src={org.logo} alt={org.name} className="w-5 h-5 rounded mr-2" />
              ) : (
                <div className="w-5 h-5 rounded bg-primary/10 flex items-center justify-center mr-2">
                  <span className="text-xs font-medium">{org.name[0]}</span>
                </div>
              )}
              {org.name}
            </button>
          ))}
          <button
            onClick={onCreateOrganization}
            className="w-full px-4 py-2 text-sm text-foreground hover:bg-primary hover:text-primary-foreground border-t border-border"
          >
            Create Organization
          </button>
        </div>
      )}
    </div>
  )
} 