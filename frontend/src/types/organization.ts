export interface Organization {
  id: string
  name: string
  slug: string
  isOwner: boolean
  createdAt: string
  updatedAt: string
  logo?: string
}

export interface OrganizationDropdownProps {
  organizations: Organization[]
  currentOrganization: Organization | null
  onSelectOrganization: (org: Organization) => void
  onCreateOrganization: () => void
  isOpen: boolean
  onToggle: () => void
} 