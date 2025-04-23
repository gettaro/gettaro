import { Organization } from '../types/organization'
// import { getAccessTokenSilently } from '@auth0/auth0-react'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

export class OrganizationConflictError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'OrganizationConflictError'
  }
}

export async function getOrganizations(getToken: () => Promise<string>): Promise<Organization[]> {
  const token = await getToken()
  const response = await fetch(`${API_BASE_URL}/organizations`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error('Failed to fetch organizations')
  }

  const data = await response.json()
  return data.organizations
}

export async function getOrganization(id: string, token: string): Promise<Organization> {
  const response = await fetch(`${API_BASE_URL}/organizations/${id}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error('Failed to fetch organization')
  }

  const data = await response.json()
  return data.organization
}

export async function createOrganization(
  name: string,
  slug: string,
  getToken: () => Promise<string>
): Promise<Organization> {
  const token = await getToken()
  const response = await fetch(`${API_BASE_URL}/organizations`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, slug }),
  })

  if (response.status === 409) {
    throw new OrganizationConflictError('An organization with this name already exists')
  }

  if (!response.ok) {
    throw new Error('Failed to create organization')
  }

  const data = await response.json()
  return data.organization
} 