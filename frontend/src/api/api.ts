import { Organization } from '../types/organization'
import { OrganizationConflictError } from './errors/organizations'
import { CreateIntegrationConfigRequest, IntegrationConfig, UpdateIntegrationConfigRequest } from '../types/integration'
import { Title, CreateTitleRequest, UpdateTitleRequest } from '../types/title'

export default class Api {
  private static API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'
  private static accessToken: string | null = null

  static setAccessToken(token: string | null) {
    console.log("setting access token", token)
    this.accessToken = token
  }

  private static async get(path: string): Promise<any> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await fetch(`${this.API_BASE_URL}${path}`, {
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    if (!response.ok) {
      throw new Error('Failed to fetch data')
    }

    return response.json()
  }

  private static async post(path: string, data: any): Promise<any> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await fetch(`${this.API_BASE_URL}${path}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      throw new Error('Failed to create data')
    }

    return response.json()
  }

  private static async put(path: string, data: any): Promise<any> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await fetch(`${this.API_BASE_URL}${path}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      throw new Error('Failed to update data')
    }

    return response.json()
  }

  private static async delete(path: string): Promise<void> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await fetch(`${this.API_BASE_URL}${path}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    if (!response.ok) {
      throw new Error('Failed to delete data')
    }
  }

  static async getOrganizations(): Promise<Organization[]> {
    const response = await this.get('/organizations')
    return response.organizations
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

  static async getOrganizationIntegrations(organizationId: string): Promise<IntegrationConfig[]> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await fetch(`${this.API_BASE_URL}/organizations/${organizationId}/integrations`, {
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    if (!response.ok) {
      throw new Error('Failed to fetch integrations')
    }

    const data = await response.json()
    return data.integrations
  }

  static async createIntegrationConfig(organizationId: string, config: CreateIntegrationConfigRequest): Promise<IntegrationConfig> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await fetch(`${this.API_BASE_URL}/organizations/${organizationId}/integrations`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(config),
    })

    if (!response.ok) {
      throw new Error('Failed to create integration')
    }

    const data = await response.json()
    return data.integration
  }

  static async updateIntegrationConfig(organizationId: string, integrationId: string, config: UpdateIntegrationConfigRequest): Promise<IntegrationConfig> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await fetch(`${this.API_BASE_URL}/organizations/${organizationId}/integrations/${integrationId}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(config),
    })

    if (!response.ok) {
      throw new Error('Failed to update integration')
    }

    const data = await response.json()
    return data.integration
  }

  static async deleteIntegrationConfig(organizationId: string, integrationId: string): Promise<void> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await fetch(`${this.API_BASE_URL}/organizations/${organizationId}/integrations/${integrationId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    if (!response.ok) {
      throw new Error('Failed to delete integration')
    }
  }

  // Title API functions
  static async getOrganizationTitles(organizationId: string): Promise<Title[]> {
    const response = await this.get(`/organizations/${organizationId}/titles`)
    return response.titles
  }

  static async createTitle(organizationId: string, request: CreateTitleRequest): Promise<Title> {
    const response = await this.post(`/organizations/${organizationId}/titles`, request)
    return response.title
  }

  static async updateTitle(organizationId: string, titleId: string, request: UpdateTitleRequest): Promise<Title> {
    const response = await this.put(`/organizations/${organizationId}/titles/${titleId}`, request)
    return response.title
  }

  static async deleteTitle(organizationId: string, titleId: string): Promise<void> {
    await this.delete(`/organizations/${organizationId}/titles/${titleId}`)
  }
}

