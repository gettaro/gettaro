import { Organization } from '../types/organization'
import { OrganizationConflictError } from './errors/organizatinos'

export default class Api {
  private static API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'
  private static accessToken: string | null = null

  static setAccessToken(token: string | null) {
    console.log("setting access token", token)
    this.accessToken = token
  }

  static async getOrganizations(): Promise<Organization[]> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await fetch(`${this.API_BASE_URL}/organizations`, {
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    if (!response.ok) {
      throw new Error('Failed to fetch organizations')
    }

    const data = await response.json()
    return data.organizations
  }

  static async getOrganization(id: string): Promise<Organization> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }
    const response = await fetch(`${this.API_BASE_URL}/organizations/${id}`, {
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    if (!response.ok) {
      throw new Error('Failed to fetch organization')
    }

    const data = await response.json()
    return data.organization
  }

  static async createOrganization(name: string, slug: string): Promise<Organization> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }
    const response = await fetch(`${this.API_BASE_URL}/organizations`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
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
}

