export interface Organization {
  id: string
  name: string
  slug: string
  is_owner: boolean
  created_at: string
  updated_at: string
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