import { Organization } from '../types/organization'
import { OrganizationConflictError } from './errors/organizations'
import { CreateIntegrationConfigRequest, IntegrationConfig, UpdateIntegrationConfigRequest } from '../types/integration'
import { Title, CreateTitleRequest, UpdateTitleRequest } from '../types/title'
import { Member, AddMemberRequest, UpdateMemberRequest } from '../types/member'
import { SourceControlAccount, PullRequest, GetMemberPullRequestsParams } from '../types/sourcecontrol'
import { GetMemberActivityParams, MemberActivity } from '../types/memberActivity'
import { GetMemberMetricsParams, GetMemberMetricsResponse } from '../types/memberMetrics'

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

  // Member API functions
  static async getOrganizationMembers(organizationId: string): Promise<Member[]> {
    const response = await this.get(`/organizations/${organizationId}/members`)
    return response.members
  }

  static async addOrganizationMember(organizationId: string, request: AddMemberRequest): Promise<void> {
    await this.post(`/organizations/${organizationId}/members`, request)
  }

  static async updateOrganizationMember(organizationId: string, memberId: string, request: UpdateMemberRequest): Promise<void> {
    await this.put(`/organizations/${organizationId}/members/${memberId}`, request)
  }

  static async deleteOrganizationMember(organizationId: string, memberId: string): Promise<void> {
    await this.delete(`/organizations/${organizationId}/members/${memberId}`)
  }

  // Member Activity API functions
  static async getMemberActivity(organizationId: string, memberId: string, params: GetMemberActivityParams): Promise<MemberActivity[]> {
    const token = this.accessToken
    const queryParams = new URLSearchParams()
    
    if (params.startDate) {
      queryParams.append('startDate', params.startDate)
    }
    if (params.endDate) {
      queryParams.append('endDate', params.endDate)
    }

    const response = await fetch(`${this.API_BASE_URL}/organizations/${organizationId}/members/${memberId}/sourcecontrol/activity?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) {
      throw new Error(`Failed to get member activity: ${response.statusText}`)
    }

    const data = await response.json()
    return data.activities
  }

  // Get member source control metrics
  static async getMemberMetrics(organizationId: string, memberId: string, params: GetMemberMetricsParams): Promise<GetMemberMetricsResponse> {
    const token = this.accessToken
    const queryParams = new URLSearchParams()
    
    if (params.startDate) {
      queryParams.append('startDate', params.startDate)
    }
    if (params.endDate) {
      queryParams.append('endDate', params.endDate)
    }
    if (params.interval) {
      queryParams.append('interval', params.interval)
    }

    const response = await fetch(`${this.API_BASE_URL}/organizations/${organizationId}/members/${memberId}/sourcecontrol/metrics?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) {
      throw new Error(`Failed to get member metrics: ${response.statusText}`)
    }

    return await response.json()
  }

  // Source Control Account API functions
  static async getOrganizationSourceControlAccounts(organizationId: string): Promise<SourceControlAccount[]> {
    const response = await this.get(`/organizations/${organizationId}/source-control-accounts`)
    return response.source_control_accounts
  }

  // Member Pull Requests API functions
  static async getMemberPullRequests(organizationId: string, memberId: string, params: GetMemberPullRequestsParams): Promise<PullRequest[]> {
    const token = this.accessToken
    const queryParams = new URLSearchParams()
    
    if (params.startDate) {
      queryParams.append('startDate', params.startDate)
    }
    if (params.endDate) {
      queryParams.append('endDate', params.endDate)
    }

    const response = await fetch(`${this.API_BASE_URL}/organizations/${organizationId}/members/${memberId}/pull-requests?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) {
      throw new Error(`Failed to get member pull requests: ${response.statusText}`)
    }

    const data = await response.json()
    return data.pull_requests
  }
}

