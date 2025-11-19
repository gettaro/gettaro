import { Organization } from '../types/organization'
import { OrganizationConflictError } from './errors/organizations'
import { CreateIntegrationConfigRequest, IntegrationConfig, UpdateIntegrationConfigRequest } from '../types/integration'
import { Title, CreateTitleRequest, UpdateTitleRequest } from '../types/title'
import { Member, AddMemberRequest, UpdateMemberRequest } from '../types/member'
import { ExternalAccount, PullRequest, GetMemberPullRequestsParams, GetMemberPullRequestReviewsParams, MemberActivity } from '../types/sourcecontrol'
import { GetManagerTreeResponse } from '../types/directs'
import { GetMemberMetricsParams, GetMemberMetricsResponse } from '../types/memberMetrics'
import { OrganizationMetricsResponse } from '../types/organizationMetrics'
import { 
  CreateConversationTemplateRequest, 
  UpdateConversationTemplateRequest, 
  ListConversationTemplatesQuery,
  ListConversationTemplatesResponse,
  GetConversationTemplateResponse,
  CreateConversationTemplateResponse,
  UpdateConversationTemplateResponse
} from '../types/conversationTemplate'
import {
  CreateConversationRequest,
  UpdateConversationRequest,
  ListConversationsQuery,
  ListConversationsResponse,
  GetConversationResponse,
  GetConversationWithDetailsResponse,
  ConversationStatsResponse
} from '../types/conversation'
import {
  GetMemberAICodeAssistantUsageParams,
  GetMemberAICodeAssistantUsageResponse,
  GetMemberAICodeAssistantMetricsParams,
  GetMemberAICodeAssistantMetricsResponse
} from '../types/aicodeassistant'
import {
  AIQueryRequest,
  AIQueryResponse,
  AIQueryHistoryItem,
  AIQueryStats
} from '../types/ai'
import {
  CreateTeamRequest,
  UpdateTeamRequest,
  AddTeamMemberRequest,
  ListTeamsResponse,
  GetTeamResponse,
  CreateTeamResponse,
  UpdateTeamResponse
} from '../types/team'
import { useApiErrorStore } from '../stores/apiError'

export default class Api {
  private static API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'
  private static accessToken: string | null = null
  private static REQUEST_TIMEOUT = 30000 // 30 seconds

  static setAccessToken(token: string | null) {
    console.log("setting access token", token ? `${token.substring(0, 20)}...` : 'null')
    this.accessToken = token
  }

  static hasAccessToken(): boolean {
    return this.accessToken !== null && this.accessToken !== ''
  }

  /**
   * Detects if an error indicates API unavailability (network or server errors)
   * Distinguishes from auth errors (401, 403) and validation errors (400, 404)
   */
  private static isApiUnavailableError(error: any, status?: number): boolean {
    // Check status codes first
    if (status !== undefined) {
      // Server errors indicate API unavailability
      if (status >= 500 && status <= 504) {
        return true
      }
      // Auth and validation errors are NOT API unavailability
      if (status === 401 || status === 403 || status === 400 || status === 404) {
        return false
      }
    }

    // Network errors indicate API unavailability
    if (error instanceof TypeError) {
      // TypeError typically means network failure, connection refused, etc.
      const message = error.message.toLowerCase()
      if (
        message.includes('failed to fetch') ||
        message.includes('networkerror') ||
        message.includes('network error') ||
        message.includes('connection') ||
        message.includes('timeout') ||
        message.includes('aborted')
      ) {
        return true
      }
    }

    // Check error message for network-related keywords
    const errorMessage = error?.message?.toLowerCase() || ''
    if (
      errorMessage.includes('network') ||
      errorMessage.includes('connection') ||
      errorMessage.includes('timeout') ||
      errorMessage.includes('refused') ||
      errorMessage.includes('fetch failed')
    ) {
      return true
    }

    return false
  }

  /**
   * Wraps fetch calls with error handling and timeout
   * Detects API unavailability and updates global error state
   */
  private static async fetchWithErrorHandling(
    url: string,
    options: RequestInit = {}
  ): Promise<Response> {
    const errorStore = useApiErrorStore.getState()

    // Create abort controller for timeout
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), this.REQUEST_TIMEOUT)

    try {
      const response = await fetch(url, {
        ...options,
        signal: controller.signal,
      })

      clearTimeout(timeoutId)

      // Clear error state on successful request
      if (response.ok) {
        errorStore.clearApiError()
        return response
      }

      // Check if this is an API unavailability error
      if (this.isApiUnavailableError(null, response.status)) {
        errorStore.setApiError(
          'server',
          'Unable to connect to the server. Please check your connection and try again.'
        )
        throw new Error(`Server error: ${response.status} ${response.statusText}`)
      }

      // Auth errors are handled separately, don't set API unavailable
      if (response.status === 401 || response.status === 403) {
        this.accessToken = null
        throw new Error('Authentication failed. Please log in again.')
      }

      // For other errors (400, 404, 409, etc.), return the response
      // so calling code can handle them appropriately (e.g., check status === 409)
      // Clear error state since these are not API unavailability errors
      errorStore.clearApiError()
      return response
    } catch (error: any) {
      clearTimeout(timeoutId)

      // Check if this is a network error or timeout
      if (this.isApiUnavailableError(error)) {
        const errorType = error.name === 'AbortError' ? 'network' : 'network'
        errorStore.setApiError(
          errorType,
          'Unable to connect to the server. Please check your connection and try again.'
        )
        throw new Error('Network error: Unable to connect to the server')
      }

      // Re-throw other errors (auth, validation, etc.)
      throw error
    }
  }

  private static async get(path: string): Promise<any> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}${path}`, {
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    if (!response.ok) {
      throw new Error(`Failed to fetch data: ${response.status} ${response.statusText}`)
    }

    return response.json()
  }

  private static async post(path: string, data: any): Promise<any> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}${path}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      throw new Error(`Failed to create data: ${response.status} ${response.statusText}`)
    }

    // Handle empty responses (like 201 Created with no body)
    const contentType = response.headers.get('content-type')
    if (contentType && contentType.includes('application/json')) {
      return response.json()
    }
    
    // Return null for empty responses
    return null
  }

  private static async put(path: string, data: any): Promise<any> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}${path}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      throw new Error(`Failed to update data: ${response.status} ${response.statusText}`)
    }

    // Handle empty responses (like 200 OK with no body)
    const contentType = response.headers.get('content-type')
    if (contentType && contentType.includes('application/json')) {
      return response.json()
    }
    
    // Return null for empty responses
    return null
  }

  private static async delete(path: string): Promise<void> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}${path}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    if (!response.ok) {
      throw new Error(`Failed to delete data: ${response.status} ${response.statusText}`)
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
    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${id}`, {
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    const data = await response.json()
    return data.organization
  }

  static async createOrganization(name: string, slug: string): Promise<Organization> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }
    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations`, {
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

    const data = await response.json()
    return data.organization
  }

  static async getOrganizationIntegrations(organizationId: string): Promise<IntegrationConfig[]> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/integrations`, {
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })

    const data = await response.json()
    return data.integrations
  }

  static async createIntegrationConfig(organizationId: string, config: CreateIntegrationConfigRequest): Promise<IntegrationConfig> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/integrations`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(config),
    })

    const data = await response.json()
    return data.integration
  }

  static async updateIntegrationConfig(organizationId: string, integrationId: string, config: UpdateIntegrationConfigRequest): Promise<IntegrationConfig> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/integrations/${integrationId}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(config),
    })

    const data = await response.json()
    return data.integration
  }

  static async deleteIntegrationConfig(organizationId: string, integrationId: string): Promise<void> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/integrations/${integrationId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
      },
    })
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

  static async addOrganizationMember(organizationId: string, request: AddMemberRequest): Promise<Member | null> {
    const response = await this.post(`/organizations/${organizationId}/members`, request)
    return response.member || null
  }

  static async updateOrganizationMember(organizationId: string, memberId: string, request: UpdateMemberRequest): Promise<void> {
    await this.put(`/organizations/${organizationId}/members/${memberId}`, request)
  }

  static async deleteOrganizationMember(organizationId: string, memberId: string): Promise<void> {
    await this.delete(`/organizations/${organizationId}/members/${memberId}`)
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

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/members/${memberId}/sourcecontrol/metrics?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    return await response.json()
  }

  // External Account API functions
  static async getOrganizationSourceControlAccounts(organizationId: string): Promise<ExternalAccount[]> {
    const response = await this.get(`/organizations/${organizationId}/external-accounts?account_type=sourcecontrol`)
    return response.external_accounts || []
  }

  static async getOrganizationExternalAccounts(organizationId: string, accountType?: string): Promise<ExternalAccount[]> {
    const queryParam = accountType ? `?account_type=${accountType}` : ''
    const response = await this.get(`/organizations/${organizationId}/external-accounts${queryParam}`)
    return response.external_accounts || []
  }

  static async updateExternalAccount(organizationId: string, accountId: string, account: ExternalAccount): Promise<void> {
    // Update only the member_id field
    await this.put(`/organizations/${organizationId}/external-accounts/${accountId}`, {
      member_id: account.member_id
    })
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

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/members/${memberId}/pull-requests?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    const data = await response.json()
    return data.pull_requests
  }

  // Member Pull Request Reviews API functions
  static async getMemberPullRequestReviews(organizationId: string, memberId: string, params: GetMemberPullRequestReviewsParams): Promise<MemberActivity[]> {
    const token = this.accessToken
    const queryParams = new URLSearchParams()
    
    if (params.startDate) {
      queryParams.append('startDate', params.startDate)
    }
    if (params.endDate) {
      queryParams.append('endDate', params.endDate)
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/members/${memberId}/pull-request-reviews?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    const data = await response.json()
    return data.reviews
  }

  // Member AI Code Assistant Usage API functions
  static async getMemberAICodeAssistantMetrics(organizationId: string, memberId: string, params?: GetMemberAICodeAssistantMetricsParams): Promise<GetMemberAICodeAssistantMetricsResponse> {
    const token = this.accessToken
    const queryParams = new URLSearchParams()
    
    if (params?.startDate) {
      queryParams.append('startDate', params.startDate)
    }
    if (params?.endDate) {
      queryParams.append('endDate', params.endDate)
    }
    if (params?.interval) {
      queryParams.append('interval', params.interval)
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/members/${memberId}/ai-code-assistant/metrics?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    return await response.json()
  }

  static async getOrganizationAICodeAssistantMetrics(organizationId: string, params?: GetMemberAICodeAssistantMetricsParams): Promise<GetMemberAICodeAssistantMetricsResponse> {
    const token = this.accessToken
    const queryParams = new URLSearchParams()
    
    if (params?.startDate) {
      queryParams.append('startDate', params.startDate)
    }
    if (params?.endDate) {
      queryParams.append('endDate', params.endDate)
    }
    if (params?.interval) {
      queryParams.append('interval', params.interval)
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/ai-code-assistant/metrics?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    return await response.json()
  }

  static async getTeamAICodeAssistantMetrics(organizationId: string, teamId: string, params?: GetMemberAICodeAssistantMetricsParams): Promise<GetMemberAICodeAssistantMetricsResponse> {
    const token = this.accessToken
    const queryParams = new URLSearchParams()
    
    if (params?.startDate) {
      queryParams.append('startDate', params.startDate)
    }
    if (params?.endDate) {
      queryParams.append('endDate', params.endDate)
    }
    if (params?.interval) {
      queryParams.append('interval', params.interval)
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/teams/${teamId}/ai-code-assistant/metrics?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    return await response.json()
  }

  // Organization Pull Requests API functions
  static async getOrganizationPullRequests(organizationId: string, params?: {
    userIds?: string[]
    repositoryName?: string
    prefix?: string
    startDate?: string
    endDate?: string
    status?: string
  }): Promise<PullRequest[]> {
    const token = this.accessToken
    const queryParams = new URLSearchParams()
    
    if (params?.userIds && params.userIds.length > 0) {
      params.userIds.forEach(id => queryParams.append('userIds', id))
    }
    if (params?.repositoryName) {
      queryParams.append('repositoryName', params.repositoryName)
    }
    if (params?.prefix) {
      queryParams.append('prefix', params.prefix)
    }
    if (params?.startDate) {
      queryParams.append('startDate', params.startDate)
    }
    if (params?.endDate) {
      queryParams.append('endDate', params.endDate)
    }
    if (params?.status) {
      queryParams.append('status', params.status)
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/pull-requests?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    const data = await response.json()
    return data.pull_requests
  }

  // Organization Metrics API functions
  static async getOrganizationMetrics(organizationId: string, params: {
    startDate?: string
    endDate?: string
    interval?: string
    teamIds?: string[]
  }): Promise<OrganizationMetricsResponse> {
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
    if (params.teamIds && params.teamIds.length > 0) {
      params.teamIds.forEach(id => queryParams.append('teamIds', id))
    }

    const response = await this.fetchWithErrorHandling(`${this.API_BASE_URL}/organizations/${organizationId}/sourcecontrol/metrics?${queryParams}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    return response.json()
  }

  static async getManagerTree(organizationId: string, managerId: string): Promise<GetManagerTreeResponse> {
    const response = await this.get(`/organizations/${organizationId}/managers/${managerId}/directs/tree`)
    return response
  }

  // Conversation Template API functions
  static async getConversationTemplates(organizationId: string, query?: ListConversationTemplatesQuery): Promise<ListConversationTemplatesResponse> {
    const queryParams = new URLSearchParams()
    if (query?.is_active !== undefined) {
      queryParams.append('is_active', query.is_active.toString())
    }
    
    const queryString = queryParams.toString()
    const url = queryString ? `/organizations/${organizationId}/conversation-templates?${queryString}` : `/organizations/${organizationId}/conversation-templates`
    
    const response = await this.get(url)
    return response
  }

  static async getConversationTemplate(templateId: string): Promise<GetConversationTemplateResponse> {
    const response = await this.get(`/conversation-templates/${templateId}`)
    return response
  }

  static async createConversationTemplate(organizationId: string, template: CreateConversationTemplateRequest): Promise<CreateConversationTemplateResponse> {
    const response = await this.post(`/organizations/${organizationId}/conversation-templates`, template)
    return response
  }

  static async updateConversationTemplate(templateId: string, template: UpdateConversationTemplateRequest): Promise<UpdateConversationTemplateResponse> {
    const response = await this.put(`/conversation-templates/${templateId}`, template)
    return response
  }

  static async deleteConversationTemplate(templateId: string): Promise<void> {
    await this.delete(`/conversation-templates/${templateId}`)
  }

  // Conversation API functions
  static async getConversations(organizationId: string, query?: ListConversationsQuery): Promise<ListConversationsResponse> {
    const queryParams = new URLSearchParams()
    if (query?.manager_member_id) queryParams.append('manager_member_id', query.manager_member_id)
    if (query?.direct_member_id) queryParams.append('direct_member_id', query.direct_member_id)
    if (query?.template_id) queryParams.append('template_id', query.template_id)
    if (query?.status) queryParams.append('status', query.status)
    if (query?.limit) queryParams.append('limit', query.limit.toString())
    if (query?.offset) queryParams.append('offset', query.offset.toString())
    
    const queryString = queryParams.toString()
    const url = queryString ? `/organizations/${organizationId}/conversations?${queryString}` : `/organizations/${organizationId}/conversations`
    
    const response = await this.get(url)
    return response
  }

  static async getConversation(conversationId: string): Promise<GetConversationResponse> {
    const response = await this.get(`/conversations/${conversationId}`)
    return response
  }

  static async getConversationWithDetails(conversationId: string): Promise<GetConversationWithDetailsResponse> {
    const response = await this.get(`/conversations/${conversationId}/details`)
    return response
  }

  static async createConversation(organizationId: string, conversation: CreateConversationRequest): Promise<GetConversationResponse> {
    const response = await this.post(`/organizations/${organizationId}/conversations`, conversation)
    return response
  }

  static async updateConversation(conversationId: string, conversation: UpdateConversationRequest): Promise<void> {
    await this.put(`/conversations/${conversationId}`, conversation)
  }

  static async deleteConversation(conversationId: string): Promise<void> {
    await this.delete(`/conversations/${conversationId}`)
  }

  static async getConversationStats(organizationId: string, managerMemberId?: string): Promise<ConversationStatsResponse> {
    const queryParams = managerMemberId ? `?manager_member_id=${managerMemberId}` : ''
    const response = await this.get(`/organizations/${organizationId}/conversations/stats${queryParams}`)
    return response
  }

  // AI API functions
  static async queryAI(organizationId: string, query: AIQueryRequest): Promise<AIQueryResponse> {
    const response = await this.post(`/organizations/${organizationId}/ai/query`, query)
    return response.ai_response
  }

  static async getAIQueryHistory(organizationId: string, limit?: number): Promise<AIQueryHistoryItem[]> {
    const queryParams = limit ? `?limit=${limit}` : ''
    const response = await this.get(`/organizations/${organizationId}/ai/history${queryParams}`)
    return response.queries
  }

  static async getAIQueryStats(organizationId: string, days?: number): Promise<AIQueryStats> {
    const queryParams = days ? `?days=${days}` : ''
    const response = await this.get(`/organizations/${organizationId}/ai/stats${queryParams}`)
    return response.stats
  }

  // Team Management
  static async listTeams(organizationId: string, params?: { name?: string }): Promise<ListTeamsResponse> {
    const queryParams = new URLSearchParams()
    if (params?.name) {
      queryParams.append('name', params.name)
    }
    const queryString = queryParams.toString()
    return this.get(`/organizations/${organizationId}/teams${queryString ? `?${queryString}` : ''}`)
  }

  static async getTeam(organizationId: string, teamId: string): Promise<GetTeamResponse> {
    return this.get(`/organizations/${organizationId}/teams/${teamId}`)
  }

  static async createTeam(organizationId: string, team: CreateTeamRequest): Promise<CreateTeamResponse> {
    return this.post(`/organizations/${organizationId}/teams`, team)
  }

  static async updateTeam(organizationId: string, teamId: string, team: UpdateTeamRequest): Promise<UpdateTeamResponse> {
    return this.put(`/organizations/${organizationId}/teams/${teamId}`, team)
  }

  static async deleteTeam(organizationId: string, teamId: string): Promise<void> {
    return this.delete(`/organizations/${organizationId}/teams/${teamId}`)
  }

  static async addTeamMember(organizationId: string, teamId: string, member: AddTeamMemberRequest): Promise<void> {
    return this.post(`/organizations/${organizationId}/teams/${teamId}/members`, member)
  }

  static async removeTeamMember(organizationId: string, teamId: string, memberId: string): Promise<void> {
    return this.delete(`/organizations/${organizationId}/teams/${teamId}/members/${memberId}`)
  }
}

