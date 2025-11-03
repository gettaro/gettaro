import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { useOrganizationStore } from '../stores/organization'
import { Member } from '../types/member'
import { GetMemberMetricsResponse, GetMemberMetricsParams } from '../types/memberMetrics'
import { Title } from '../types/title'
import { ExternalAccount, PullRequest, GetMemberPullRequestsParams, GetMemberPullRequestReviewsParams } from '../types/sourcecontrol'
import { MemberActivity } from '../types/sourcecontrol'
import { OrgChartNode } from '../types/directs'
import Api from '../api/api'
import { formatMetricValue, formatTimeMetric } from '../utils/formatMetrics'
import { ConversationsTab } from '../components/ConversationsTab'
import { ConversationSidebar } from '../components/ConversationSidebar'
import { ConversationWithDetails } from '../types/conversation'
import { AIChat } from '../components/AIChat'
import { ChatContext } from '../types/ai'
import MetricChart from '../components/MetricChart'
import PullRequestItem from '../components/PullRequestItem'
import { AICodeAssistantUsageStats, GetMemberAICodeAssistantUsageParams } from '../types/aicodeassistant'

type TabType = 'overview' | 'source-control-metrics' | 'management-tree' | 'conversations' | 'ai-chat' | 'ai-code-assistant-usage'

export default function MemberProfilePage() {
  const { memberId } = useParams<{ memberId: string }>()
  const { currentOrganization } = useOrganizationStore()
  const [member, setMember] = useState<Member | null>(null)
  const [title, setTitle] = useState<Title | null>(null)
  const [sourceControlAccount, setSourceControlAccount] = useState<ExternalAccount | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<TabType>('overview')
  
  // Date filter state for Code Contributions tab
  const [dateParams, setDateParams] = useState<GetMemberMetricsParams>(() => {
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 30) // 1 month ago
    
    return {
      startDate: startDate.toISOString().split('T')[0], // YYYY-MM-DD format
      endDate: endDate.toISOString().split('T')[0],
      interval: 'weekly' as 'daily' | 'weekly' | 'monthly'
    }
  })
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set())
  const [expandedTables, setExpandedTables] = useState<Set<string>>(new Set())
  const [metrics, setMetrics] = useState<GetMemberMetricsResponse | null>(null)
  const [metricsLoading, setMetricsLoading] = useState(false)
  const [pullRequests, setPullRequests] = useState<PullRequest[]>([])
  const [pullRequestsLoading, setPullRequestsLoading] = useState(false)
  const [pullRequestReviews, setPullRequestReviews] = useState<MemberActivity[]>([])
  const [pullRequestReviewsLoading, setPullRequestReviewsLoading] = useState(false)
  const [managementTree, setManagementTree] = useState<OrgChartNode[]>([])
  const [managementTreeLoading, setManagementTreeLoading] = useState(false)
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set())
  
  // Metrics view state
  const [metricsViewMode, setMetricsViewMode] = useState<'snapshot' | 'graph'>('snapshot')
  const [currentGraphIndex, setCurrentGraphIndex] = useState(0)
  
  // Conversation sidebar state
  const [showConversationSidebar, setShowConversationSidebar] = useState(false)
  const [selectedConversation, setSelectedConversation] = useState<ConversationWithDetails | null>(null)
  const [conversationSidebarMode, setConversationSidebarMode] = useState<'edit' | 'create'>('edit')
  
  // AI Code Assistant Usage state
  const [aiUsageStats, setAiUsageStats] = useState<AICodeAssistantUsageStats | null>(null)
  const [aiUsageStatsLoading, setAiUsageStatsLoading] = useState(false)
  const [aiUsageDateParams, setAiUsageDateParams] = useState<GetMemberAICodeAssistantUsageParams>(() => {
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 30) // 1 month ago
    
    return {
      startDate: startDate.toISOString().split('T')[0],
      endDate: endDate.toISOString().split('T')[0],
    }
  })

  useEffect(() => {
    if (currentOrganization && memberId) {
      initializePage()
    }
  }, [currentOrganization, memberId])

  // Load data when Code Contributions tab is selected or date parameters change
  useEffect(() => {
    if (activeTab === 'source-control-metrics' && currentOrganization?.id && memberId) {
      loadCodeContributionsData()
    }
  }, [activeTab, currentOrganization?.id, memberId, dateParams.startDate, dateParams.endDate, dateParams.interval])

  // Load management tree when management tree tab is selected
  useEffect(() => {
    if (activeTab === 'management-tree' && member && title?.is_manager && currentOrganization?.id && memberId) {
      loadManagementTree()
    }
  }, [activeTab, member, title, currentOrganization?.id, memberId])

  // Load AI Code Assistant usage stats when tab is selected
  useEffect(() => {
    if (activeTab === 'ai-code-assistant-usage' && currentOrganization?.id && memberId) {
      loadAICodeAssistantUsageStats()
    }
  }, [activeTab, currentOrganization?.id, memberId, aiUsageDateParams.startDate, aiUsageDateParams.endDate])

  const initializePage = async () => {
    try {
      await loadMemberData()
    } catch (err) {
      console.error('Error initializing page:', err)
      setError('Failed to initialize page')
    }
  }


  const loadCodeContributionsData = async () => {
    if (!currentOrganization?.id || !memberId) return

    // Set all loading states to true
    setMetricsLoading(true)
    setPullRequestsLoading(true)
    setPullRequestReviewsLoading(true)
    setError(null)

    try {
      // Prepare common parameters
      const metricsParams: GetMemberMetricsParams = {
        startDate: dateParams.startDate,
        endDate: dateParams.endDate,
        interval: dateParams.interval || 'weekly'
      }
      
      const prParams: GetMemberPullRequestsParams = {
        startDate: dateParams.startDate,
        endDate: dateParams.endDate,
      }
      
      const reviewsParams: GetMemberPullRequestReviewsParams = {
        startDate: dateParams.startDate,
        endDate: dateParams.endDate,
      }

      // Make all API calls in parallel for better performance
      const [metricsData, prs, reviews] = await Promise.all([
        Api.getMemberMetrics(currentOrganization.id, memberId, metricsParams),
        Api.getMemberPullRequests(currentOrganization.id, memberId, prParams),
        Api.getMemberPullRequestReviews(currentOrganization.id, memberId, reviewsParams)
      ])

      // Update all states
      setMetrics(metricsData)
      setPullRequests(prs)
      setPullRequestReviews(reviews)
    } catch (err) {
      console.error('Error loading code contributions data:', err)
      setError('Failed to load code contributions data')
    } finally {
      // Set all loading states to false
      setMetricsLoading(false)
      setPullRequestsLoading(false)
      setPullRequestReviewsLoading(false)
    }
  }

  const loadAICodeAssistantUsageStats = async () => {
    if (!currentOrganization?.id || !memberId) return

    setAiUsageStatsLoading(true)
    try {
      const response = await Api.getMemberAICodeAssistantUsageStats(
        currentOrganization.id,
        memberId,
        aiUsageDateParams
      )
      setAiUsageStats(response.stats)
    } catch (err) {
      console.error('Error loading AI code assistant usage stats:', err)
      setError('Failed to load AI code assistant usage stats')
    } finally {
      setAiUsageStatsLoading(false)
    }
  }

  const handleAiUsageDateChange = (field: 'startDate' | 'endDate', value: string) => {
    setAiUsageDateParams(prev => ({
      ...prev,
      [field]: value || undefined
    }))
  }

  const loadManagementTree = async () => {
    if (!currentOrganization?.id || !memberId || !member) return

    try {
      setManagementTreeLoading(true)
      setError(null)

      const response = await Api.getManagerTree(currentOrganization.id, memberId)
      setManagementTree(response.org_chart)
      
      // Auto-expand the top-level manager (the current member)
      if (response.org_chart && response.org_chart.length > 0) {
        setExpandedNodes(new Set([memberId]))
      }
    } catch (err) {
      console.error('Error loading management tree:', err)
      setError('Failed to load management tree')
    } finally {
      setManagementTreeLoading(false)
    }
  }

  const handleDateChange = (field: 'startDate' | 'endDate', value: string) => {
    setDateParams(prev => ({
      ...prev,
      [field]: value || undefined
    }))
  }

  const handleIntervalChange = (interval: 'daily' | 'weekly' | 'monthly') => {
    setDateParams(prev => ({
      ...prev,
      interval
    }))
  }

  // Get all graphs from metrics
  const getAllGraphs = (metricsData: GetMemberMetricsResponse) => {
    if (!metricsData.graph_metrics || metricsData.graph_metrics.length === 0) {
      return []
    }

    const allGraphs: Array<{ 
      metric: any
      category: string
      description: string
    }> = []

    metricsData.graph_metrics.forEach((category) => {
      // Filter metrics that have data
      const metricsWithData = category.metrics.filter((metric) => {
        if (!metric.time_series || metric.time_series.length === 0) {
          return false
        }
        return metric.time_series.some(entry => 
          entry.data && entry.data.length > 0
        )
      })

      metricsWithData.forEach((metric) => {
        const snapshotMetric = metricsData.snapshot_metrics
          ?.flatMap(cat => cat.metrics)
          .find(m => m.label === metric.label)
        const description = snapshotMetric?.description || ''
        
        allGraphs.push({
          metric,
          category: category.category.name,
          description
        })
      })
    })

    return allGraphs
  }

  const navigateGraph = (direction: 'prev' | 'next', totalGraphs: number) => {
    let newIndex: number
    
    if (direction === 'prev') {
      newIndex = currentGraphIndex > 0 ? currentGraphIndex - 1 : totalGraphs - 1
    } else {
      newIndex = currentGraphIndex < totalGraphs - 1 ? currentGraphIndex + 1 : 0
    }
    
    setCurrentGraphIndex(newIndex)
  }

  // Reset graph index when metrics change
  useEffect(() => {
    if (metrics && metricsViewMode === 'graph') {
      const allGraphs = getAllGraphs(metrics)
      if (allGraphs.length > 0 && currentGraphIndex >= allGraphs.length) {
        setCurrentGraphIndex(0)
      }
    }
  }, [metrics, metricsViewMode])

  const toggleExpanded = (itemId: string) => {
    setExpandedItems(prev => {
      const newSet = new Set(prev)
      if (newSet.has(itemId)) {
        newSet.delete(itemId)
      } else {
        newSet.add(itemId)
      }
      return newSet
    })
  }

  const toggleTableExpanded = (tableId: string) => {
    setExpandedTables(prev => {
      const newSet = new Set(prev)
      if (newSet.has(tableId)) {
        newSet.delete(tableId)
      } else {
        newSet.add(tableId)
      }
      return newSet
    })
  }

  const getActivityIcon = (type: string) => {
    switch (type) {
      case 'pull_request':
        return (
          <svg className="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" />
          </svg>
        )
      case 'pr_comment':
        return (
          <svg className="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        )
      case 'pr_review':
        return (
          <svg className="w-5 h-5 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        )
      default:
        return (
          <svg className="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        )
    }
  }

  const renderActivityContent = (activity: MemberActivity) => {
    switch (activity.type) {
      case 'pull_request':
        return (
          <div>
            <div className="text-sm text-muted-foreground mb-2">
              <span className="font-medium">@{activity.author_username}</span> has created a new PR
            </div>
            <h3 className="text-lg font-medium text-foreground mb-3">
              {activity.title}
            </h3>
            
            {/* PR Statistics */}
            {activity.metadata && (
              <div className="flex flex-wrap gap-4 text-sm text-muted-foreground mb-3">
                {/* PR Status */}
                {activity.metadata.state && (
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    activity.metadata.state === 'open' 
                      ? 'bg-green-100 text-green-800' 
                      : activity.metadata.state === 'closed' 
                      ? 'bg-red-100 text-red-800'
                      : 'bg-gray-100 text-gray-800'
                  }`}>
                    {activity.metadata.state === 'open' ? 'Open' : 
                     activity.metadata.state === 'closed' ? 'Closed' : 
                     activity.metadata.state}
                  </span>
                )}
                
                {activity.metadata.additions !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                    <span>+{activity.metadata.additions}</span>
                  </span>
                )}
                {activity.metadata.deletions !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
                    </svg>
                    <span>-{activity.metadata.deletions}</span>
                  </span>
                )}
                {activity.metadata.commits !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <span>{activity.metadata.commits} commits</span>
                  </span>
                )}
                {activity.metadata.changed_files !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <span>{activity.metadata.changed_files} files</span>
                  </span>
                )}
                {activity.metadata.review_comments !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                    </svg>
                    <span>{activity.metadata.review_comments} comments</span>
                  </span>
                )}
              </div>
            )}

            {/* PR Metrics */}
            {activity.pr_metrics && (
              <div className="flex flex-wrap gap-4 text-sm text-muted-foreground mb-3">
                {/* Time to merge */}
                {activity.pr_metrics.time_to_merge_seconds !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span>Merge: {formatTimeMetric(activity.pr_metrics.time_to_merge_seconds)}</span>
                  </span>
                )}
                
                {/* Time to first review */}
                {activity.pr_metrics.time_to_first_non_bot_review_seconds !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span>First review: {formatTimeMetric(activity.pr_metrics.time_to_first_non_bot_review_seconds)}</span>
                  </span>
                )}

                {/* Show opened duration for open PRs */}
                {activity.metadata?.state === 'open' && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span>Opened {Math.ceil((Date.now() - new Date(activity.created_at).getTime()) / (1000 * 60 * 60 * 24))}d</span>
                  </span>
                )}
              </div>
            )}
            
            {/* Expandable PR Description */}
            {activity.description && (
              <div className="mb-3">
                <button
                  onClick={() => toggleExpanded(`${activity.id}-description`)}
                  className="flex items-center space-x-2 text-sm text-primary hover:text-primary/80 transition-colors"
                >
                  <svg 
                    className={`w-4 h-4 transition-transform ${expandedItems.has(`${activity.id}-description`) ? 'rotate-90' : ''}`}
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                  <span>
                    {expandedItems.has(`${activity.id}-description`) ? 'Hide description' : 'Show description'}
                  </span>
                </button>
                {expandedItems.has(`${activity.id}-description`) && (
                  <div className="mt-2 p-3 bg-muted/30 rounded-md border border-border">
                    <p className="text-muted-foreground text-sm whitespace-pre-wrap">
                      {activity.description}
                    </p>
                  </div>
                )}
              </div>
            )}
          </div>
        )

      case 'pr_comment':
        return (
          <div>
            <div className="text-sm text-muted-foreground mb-2">
              <span className="font-medium">@{activity.author_username}</span> commented on{' '}
              <span className="font-medium">{activity.pr_title}</span> from{' '}
              <span className="font-medium">@{activity.pr_author_username}</span>
            </div>
            
            {/* Expandable Comment */}
            {activity.description && (
              <div className="mb-3">
                <button
                  onClick={() => toggleExpanded(`${activity.id}-comment`)}
                  className="flex items-center space-x-2 text-sm text-primary hover:text-primary/80 transition-colors"
                >
                  <svg 
                    className={`w-4 h-4 transition-transform ${expandedItems.has(`${activity.id}-comment`) ? 'rotate-90' : ''}`}
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                  <span>
                    {expandedItems.has(`${activity.id}-comment`) ? 'Hide comment' : 'Show comment'}
                  </span>
                </button>
                {expandedItems.has(`${activity.id}-comment`) && (
                  <div className="mt-2 p-3 bg-muted/30 rounded-md border border-border">
                    <p className="text-muted-foreground text-sm whitespace-pre-wrap">
                      {activity.description}
                    </p>
                  </div>
                )}
              </div>
            )}
          </div>
        )

      case 'pr_review':
        return (
          <div>
            <div className="text-sm text-muted-foreground mb-2">
              <span className="font-medium">@{activity.author_username}</span> has reviewed{' '}
              <span className="font-medium">{activity.pr_title}</span> from{' '}
              <span className="font-medium">@{activity.pr_author_username}</span>
            </div>
            
            {/* Expandable Review */}
            {activity.description && (
              <div className="mb-3">
                <button
                  onClick={() => toggleExpanded(`${activity.id}-review`)}
                  className="flex items-center space-x-2 text-sm text-primary hover:text-primary/80 transition-colors"
                >
                  <svg 
                    className={`w-4 h-4 transition-transform ${expandedItems.has(`${activity.id}-review`) ? 'rotate-90' : ''}`}
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                  <span>
                    {expandedItems.has(`${activity.id}-review`) ? 'Hide review' : 'Show review'}
                  </span>
                </button>
                {expandedItems.has(`${activity.id}-review`) && (
                  <div className="mt-2 p-3 bg-muted/30 rounded-md border border-border">
                    <p className="text-muted-foreground text-sm whitespace-pre-wrap">
                      {activity.description}
                    </p>
                  </div>
                )}
              </div>
            )}
          </div>
        )

      default:
        return (
          <div>
            <h3 className="text-lg font-medium text-foreground mb-2">
              {activity.title}
            </h3>
          </div>
        )
    }
  }

  const loadMemberData = async () => {
    if (!currentOrganization?.id || !memberId) return

    try {
      setLoading(true)
      setError(null)

      // Load member details
      const members = await Api.getOrganizationMembers(currentOrganization.id)
      const foundMember = members.find(m => m.id === memberId)
      if (!foundMember) {
        setError('Member not found')
        return
      }
      setMember(foundMember)

      // Load title if member has one
      if (foundMember.title_id) {
        try {
          const titles = await Api.getOrganizationTitles(currentOrganization.id)
          const foundTitle = titles.find(t => t.id === foundMember.title_id)
          setTitle(foundTitle || null)
        } catch (err) {
          console.error('Error loading title:', err)
        }
      }

      // Load source control account if member has one
      try {
        const sourceControlAccounts = await Api.getOrganizationSourceControlAccounts(currentOrganization.id)
        const foundAccount = sourceControlAccounts.find(acc => acc.member_id === memberId)
        setSourceControlAccount(foundAccount || null)
      } catch (err) {
        console.error('Error loading source control account:', err)
      }

    } catch (err) {
      console.error('Error loading member data:', err)
      setError('Failed to load member data')
    } finally {
      setLoading(false)
    }
  }

  const toggleNodeExpansion = (nodeId: string) => {
    setExpandedNodes(prev => {
      const newSet = new Set(prev)
      if (newSet.has(nodeId)) {
        newSet.delete(nodeId)
      } else {
        newSet.add(nodeId)
      }
      return newSet
    })
  }

  const expandAllNodes = () => {
    const getAllNodeIds = (nodes: OrgChartNode[]): string[] => {
      const ids: string[] = []
      nodes.forEach(node => {
        ids.push(node.member.id)
        if (node.direct_reports && node.direct_reports.length > 0) {
          ids.push(...getAllNodeIds(node.direct_reports))
        }
      })
      return ids
    }
    
    const allIds = getAllNodeIds(managementTree)
    setExpandedNodes(new Set(allIds))
  }

  const collapseAllNodes = () => {
    setExpandedNodes(new Set())
  }

  const handleViewConversation = async (conversationId: string) => {
    try {
      const response = await Api.getConversationWithDetails(conversationId)
      setSelectedConversation(response.conversation)
      setConversationSidebarMode('edit')
      setShowConversationSidebar(true)
    } catch (error) {
      console.error('Error fetching conversation details:', error)
    }
  }

  const handleCreateConversation = () => {
    setSelectedConversation(null)
    setConversationSidebarMode('create')
    setShowConversationSidebar(true)
  }

  const handleCloseConversationSidebar = () => {
    setShowConversationSidebar(false)
    setSelectedConversation(null)
    setConversationSidebarMode('edit')
  }

  const handleConversationUpdate = (updatedConversation: ConversationWithDetails) => {
    setSelectedConversation(updatedConversation)
    // You might want to refresh the conversations list here
  }

  const handleConversationCreate = (newConversation: ConversationWithDetails) => {
    setSelectedConversation(newConversation)
    setConversationSidebarMode('edit')
    // You might want to refresh the conversations list here
  }

  const renderManagementTreeNode = (node: OrgChartNode, depth: number = 0) => {
    const indentClass = `ml-${depth * 3}`
    const isExpanded = expandedNodes.has(node.member.id)
    const hasDirectReports = node.direct_reports && node.direct_reports.length > 0
    const isManager = hasDirectReports
    
    return (
      <div key={node.member.id} className={`${indentClass} mb-1`}>
        <div className="flex items-center space-x-2 p-2 bg-muted/20 rounded border border-border/30 hover:bg-muted/30 transition-colors">
          {/* Expand/Collapse Button */}
          {isManager && (
            <button
              onClick={() => toggleNodeExpansion(node.member.id)}
              className="w-5 h-5 flex items-center justify-center rounded hover:bg-muted/50 transition-colors"
            >
              {isExpanded ? (
                <svg className="w-3 h-3 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              ) : (
                <svg className="w-3 h-3 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              )}
            </button>
          )}
          
          {/* Avatar */}
          <div className="w-6 h-6 bg-primary/10 rounded-full flex items-center justify-center">
            <span className="text-primary font-medium text-xs">
              {node.member.username.charAt(0).toUpperCase()}
            </span>
          </div>
          
          {/* Member Info */}
          <div className="flex-1">
            <h4 className="font-medium text-foreground text-sm">{node.member.username}</h4>
            {node.member.title && (
              <p className="text-xs text-muted-foreground">{node.member.title}</p>
            )}
            {isManager && (
              <p className="text-xs text-muted-foreground">
                {node.direct_reports.length} direct report{node.direct_reports.length !== 1 ? 's' : ''}
              </p>
            )}
          </div>
          
          {/* Level Indicator */}
          {depth > 0 && (
            <div className="ml-auto">
              <span className="text-xs text-muted-foreground bg-muted/50 px-1.5 py-0.5 rounded text-xs">
                Level {depth}
              </span>
            </div>
          )}
        </div>
        
        {/* Direct Reports (only show if expanded) */}
        {isManager && isExpanded && (
          <div className="mt-1 ml-4 border-l border-muted/50 pl-3">
            {(node.direct_reports || []).map(report => renderManagementTreeNode(report, depth + 1))}
          </div>
        )}
      </div>
    )
  }

  const renderTabContent = () => {
    switch (activeTab) {
      case 'overview':
        return (
          <div className="space-y-4">
            <div className="bg-card rounded-lg p-4">
              <h3 className="text-lg font-semibold text-foreground mb-3">Source Control Overview</h3>
              <p className="text-muted-foreground text-sm">
                This section will contain source control metrics. Coming soon...
              </p>
            </div>
          </div>
        )

      case 'source-control-metrics':
        return (
          <div className="space-y-4">
            {/* Date Filter Controls */}
            <div className="bg-card rounded-lg p-4">
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-lg font-semibold text-foreground">Filter by Date Range</h3>
                {(metricsLoading || pullRequestsLoading || pullRequestReviewsLoading) && (
                  <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                    <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Loading...</span>
                  </div>
                )}
              </div>
              <div className="flex space-x-3">
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Start Date
                  </label>
                  <input
                    type="date"
                    value={dateParams.startDate || ''}
                    onChange={(e) => handleDateChange('startDate', e.target.value)}
                    className="px-3 py-2 border border-border rounded bg-card text-foreground focus:outline-none focus:ring-2 focus:ring-primary text-sm"
                    style={{ colorScheme: 'dark light' }}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    End Date
                  </label>
                  <input
                    type="date"
                    value={dateParams.endDate || ''}
                    onChange={(e) => handleDateChange('endDate', e.target.value)}
                    className="px-3 py-2 border border-border rounded bg-card text-foreground focus:outline-none focus:ring-2 focus:ring-primary text-sm"
                    style={{ colorScheme: 'dark light' }}
                  />
                </div>
                {metricsViewMode === 'graph' && (
                  <div>
                    <label className="block text-sm font-medium text-foreground mb-1">
                      Interval
                    </label>
                    <select
                      value={dateParams.interval || 'weekly'}
                      onChange={(e) => handleIntervalChange(e.target.value as 'daily' | 'weekly' | 'monthly')}
                      className="px-3 py-2 border border-border rounded bg-card text-foreground focus:outline-none focus:ring-2 focus:ring-primary text-sm"
                    >
                      <option value="daily">Daily</option>
                      <option value="weekly">Weekly</option>
                      <option value="monthly">Monthly</option>
                    </select>
                  </div>
                )}
              </div>
            </div>

            {/* Metrics Section */}
            <div className="bg-card rounded-lg p-4">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-foreground">Metrics</h3>
                
                {/* View Mode Toggle */}
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-muted-foreground">View:</span>
                  <div className="flex bg-muted/20 rounded-lg p-1">
                    <button
                      onClick={() => setMetricsViewMode('snapshot')}
                      className={`px-3 py-1 text-xs font-medium rounded-md transition-colors ${
                        metricsViewMode === 'snapshot'
                          ? 'bg-background text-foreground shadow-sm'
                          : 'text-muted-foreground hover:text-foreground'
                      }`}
                    >
                      Snapshot
                    </button>
                    <button
                      onClick={() => setMetricsViewMode('graph')}
                      className={`px-3 py-1 text-xs font-medium rounded-md transition-colors ${
                        metricsViewMode === 'graph'
                          ? 'bg-background text-foreground shadow-sm'
                          : 'text-muted-foreground hover:text-foreground'
                      }`}
                    >
                      Graph
                    </button>
                  </div>
                </div>
              </div>
              
              {metricsLoading ? (
                <div className="flex items-center justify-center h-24">
                  <div className="flex items-center space-x-2 text-muted-foreground">
                    <svg className="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Loading metrics...</span>
                  </div>
                </div>
              ) : metrics ? (
                <div className="space-y-6">
                  {metricsViewMode === 'snapshot' ? (
                    // Snapshot View
                    (metrics.snapshot_metrics || []).map((category) => (
                      <div key={category.category.name} className="space-y-3">
                        <h4 className="text-md font-semibold text-foreground">{category.category.name}</h4>
                        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
                          {(category.metrics || []).map((metric) => (
                            <div key={metric.label} className="text-center p-3 bg-muted/20 rounded border border-border/30">
                              <div className="text-xl font-bold text-foreground">
                                {formatMetricValue(metric.value, metric.unit)}
                              </div>
                              <div className="text-xs text-muted-foreground">{metric.label}</div>
                              {typeof metric.peers_value === 'number' && metric.peers_value > 0 && (
                                <div className="text-xs text-muted-foreground mt-1">
                                  vs {formatMetricValue(metric.peers_value, metric.unit)} (peers)
                                </div>
                              )}
                            </div>
                          ))}
                        </div>
                      </div>
                    ))
                  ) : (
                    // Graph View with Slider
                    (() => {
                      const allGraphs = getAllGraphs(metrics)
                      
                      if (allGraphs.length === 0) {
                        return (
                          <div className="flex items-center justify-center h-24">
                            <div className="text-center">
                              <div className="text-muted-foreground mb-1 text-sm">No graph data available</div>
                              <div className="text-xs text-muted-foreground">
                                Graph data will appear here once you have time-series metrics
                              </div>
                            </div>
                          </div>
                        )
                      }

                      const currentGraph = allGraphs[currentGraphIndex] || allGraphs[0]
                      const chartElement = currentGraph.metric ? (
                        <MetricChart metric={currentGraph.metric} height={300} />
                      ) : null

                      if (!chartElement) {
                        return (
                          <div className="flex items-center justify-center h-24">
                            <div className="text-center">
                              <div className="text-muted-foreground mb-1 text-sm">No graph data available</div>
                            </div>
                          </div>
                        )
                      }

                      return (
                        <div>
                          <div className="bg-muted/30 rounded-lg p-4">
                            <div className="relative">
                              {/* Navigation Buttons */}
                              <div className="flex items-center justify-between mb-4">
                                <button
                                  onClick={() => navigateGraph('prev', allGraphs.length)}
                                  className="p-2 rounded hover:bg-muted/50 transition-colors"
                                  aria-label="Previous graph"
                                >
                                  <svg
                                    className="w-5 h-5"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                  >
                                    <path
                                      strokeLinecap="round"
                                      strokeLinejoin="round"
                                      strokeWidth={2}
                                      d="M15 19l-7-7 7-7"
                                    />
                                  </svg>
                                </button>
                                
                                <div className="flex items-center gap-2">
                                  <span className="text-sm text-muted-foreground">
                                    {currentGraphIndex + 1} of {allGraphs.length}
                                  </span>
                                </div>
                                
                                <button
                                  onClick={() => navigateGraph('next', allGraphs.length)}
                                  className="p-2 rounded hover:bg-muted/50 transition-colors"
                                  aria-label="Next graph"
                                >
                                  <svg
                                    className="w-5 h-5"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                  >
                                    <path
                                      strokeLinecap="round"
                                      strokeLinejoin="round"
                                      strokeWidth={2}
                                      d="M9 5l7 7-7 7"
                                    />
                                  </svg>
                                </button>
                              </div>

                              {/* Graph Content */}
                              <div className="bg-muted/10 rounded-lg p-4">
                                <div className="flex items-center justify-between mb-3">
                                  <div>
                                    <h5 className="font-medium text-foreground">{currentGraph.metric.label}</h5>
                                    {currentGraph.description && (
                                      <p className="text-xs text-muted-foreground mt-1">{currentGraph.description}</p>
                                    )}
                                  </div>
                                  <span className="text-xs text-muted-foreground capitalize">{currentGraph.metric.type}</span>
                                </div>
                                <div className="bg-muted/5 rounded border border-border/30 p-4">
                                  {chartElement}
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      )
                    })()
                  )}
                </div>
              ) : (
                <div className="flex items-center justify-center h-24">
                  <div className="text-center">
                    <div className="text-muted-foreground mb-1 text-sm">No metrics available</div>
                    <div className="text-xs text-muted-foreground">
                      Metrics will appear here once you have source control activity
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* Recent Pull Requests Table */}
            <div className="bg-card rounded-lg">
              <div className="p-4 border-b border-border/50">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-lg font-semibold text-foreground">Pull Requests</h3>
                    <p className="text-sm text-muted-foreground mt-1">
                      Pull requests created in the selected date range
                    </p>
                  </div>
                  <div className="flex items-center space-x-3">
                    {pullRequestsLoading && (
                      <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                        <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        <span>Loading...</span>
                      </div>
                    )}
                    <button
                      onClick={() => toggleTableExpanded('pull-requests')}
                      className="flex items-center space-x-2 text-sm text-primary hover:text-primary/80 transition-colors"
                    >
                      <svg 
                        className={`w-4 h-4 transition-transform ${expandedTables.has('pull-requests') ? 'rotate-90' : ''}`}
                        fill="none" 
                        stroke="currentColor" 
                        viewBox="0 0 24 24"
                      >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                      </svg>
                      <span>
                        {expandedTables.has('pull-requests') ? 'Collapse' : 'Expand'}
                      </span>
                    </button>
                  </div>
                </div>
              </div>
              {expandedTables.has('pull-requests') && (
                <div className="divide-y divide-border/50">
                  {pullRequests.length === 0 ? (
                    <div className="p-4 text-center">
                      <div className="text-muted-foreground text-sm">
                        No pull requests found for the selected date range.
                      </div>
                    </div>
                  ) : (
                    <div className="max-h-80 overflow-y-auto">
                      {(pullRequests || []).map((pr) => (
                        <div key={pr.id} className="p-4">
                          <div className="flex items-start space-x-3">
                            <div className="flex-shrink-0 mt-1">
                              <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" />
                              </svg>
                            </div>
                            <div className="flex-1 min-w-0">
                              <PullRequestItem pr={pr} />
                              
                              <div className="flex items-center space-x-3 text-xs text-muted-foreground mt-2">
                                <span>
                                  {new Date(pr.created_at).toLocaleDateString('en-US', {
                                    year: 'numeric',
                                    month: 'short',
                                    day: 'numeric',
                                    hour: '2-digit',
                                    minute: '2-digit'
                                  })}
                                </span>
                                {pr.merged_at && (
                                  <span className="text-green-600">
                                    Merged {new Date(pr.merged_at).toLocaleString('en-US', {
                                      year: 'numeric',
                                      month: 'short',
                                      day: 'numeric',
                                      hour: '2-digit',
                                      minute: '2-digit'
                                    })}
                                  </span>
                                )}
                                {pr.url && (
                                  <a
                                    href={pr.url}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-primary hover:text-primary/80 transition-colors"
                                  >
                                    View on GitHub →
                                  </a>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* Recent PR Reviews Table */}
            <div className="bg-card rounded-lg">
              <div className="p-4 border-b border-border/50">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-lg font-semibold text-foreground">Pull Request Reviews</h3>
                    <p className="text-sm text-muted-foreground mt-1">
                      Pull request reviews submitted in the selected date range
                    </p>
                  </div>
                  <button
                    onClick={() => toggleTableExpanded('pr-reviews')}
                    className="flex items-center space-x-2 text-sm text-primary hover:text-primary/80 transition-colors"
                  >
                    <svg 
                      className={`w-4 h-4 transition-transform ${expandedTables.has('pr-reviews') ? 'rotate-90' : ''}`}
                      fill="none" 
                      stroke="currentColor" 
                      viewBox="0 0 24 24"
                    >
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                    <span>
                      {expandedTables.has('pr-reviews') ? 'Collapse' : 'Expand'}
                    </span>
                  </button>
                </div>
              </div>
              {expandedTables.has('pr-reviews') && (
                <div className="divide-y divide-border/50">
                  {pullRequestReviewsLoading ? (
                    <div className="p-4 text-center">
                      <div className="text-muted-foreground text-sm">Loading PR reviews...</div>
                    </div>
                  ) : pullRequestReviews.length === 0 ? (
                    <div className="p-4 text-center">
                      <div className="text-muted-foreground text-sm">
                        No PR reviews found for the selected date range.
                      </div>
                    </div>
                  ) : (
                    <div className="max-h-80 overflow-y-auto">
                      {(pullRequestReviews || []).map((review) => (
                          <div key={review.id} className="p-4">
                            <div className="flex items-start space-x-3">
                              <div className="flex-shrink-0 mt-1">
                                {getActivityIcon(review.type)}
                              </div>
                              <div className="flex-1 min-w-0">
                                {renderActivityContent(review)}
                                
                                <div className="flex items-center space-x-3 text-xs text-muted-foreground mt-2">
                                  <span>
                                    {new Date(review.created_at).toLocaleDateString('en-US', {
                                      year: 'numeric',
                                      month: 'short',
                                      day: 'numeric',
                                      hour: '2-digit',
                                      minute: '2-digit'
                                    })}
                                  </span>
                                  {review.url && (
                                    <a
                                      href={review.url}
                                      target="_blank"
                                      rel="noopener noreferrer"
                                      className="text-primary hover:text-primary/80 transition-colors"
                                    >
                                      View on GitHub →
                                    </a>
                                  )}
                                </div>
                              </div>
                            </div>
                          </div>
                        ))}
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        )

      case 'management-tree':
        return (
          <div className="space-y-4">
            <div className="bg-card rounded-lg p-4">
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-lg font-semibold text-foreground">Management Tree</h3>
                <div className="flex items-center space-x-2">
                  {managementTreeLoading ? (
                    <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                      <div className="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
                      <span>Loading...</span>
                    </div>
                  ) : managementTree.length > 0 ? (
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={expandAllNodes}
                        className="text-xs bg-muted hover:bg-muted/80 text-muted-foreground px-2 py-1 rounded transition-colors"
                      >
                        Expand All
                      </button>
                      <button
                        onClick={collapseAllNodes}
                        className="text-xs bg-muted hover:bg-muted/80 text-muted-foreground px-2 py-1 rounded transition-colors"
                      >
                        Collapse All
                      </button>
                    </div>
                  ) : null}
                </div>
              </div>
              
              {managementTreeLoading ? (
                <div className="text-center py-6">
                  <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-2"></div>
                  <p className="text-muted-foreground text-sm">Loading management tree...</p>
                </div>
              ) : managementTree.length > 0 ? (
                <div className="space-y-3">
                  <p className="text-sm text-muted-foreground mb-3">
                    Direct and indirect reports:
                  </p>
                  
                  {/* Direct Reports */}
                  {(managementTree || []).map(node => renderManagementTreeNode(node))}
                </div>
              ) : (
                <div className="text-center py-6">
                  <div className="w-10 h-10 bg-muted rounded-full flex items-center justify-center mx-auto mb-2">
                    <svg className="w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                  </div>
                  <p className="text-muted-foreground text-sm">No direct reports found</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    This manager doesn't have any direct reports yet.
                  </p>
                </div>
              )}
            </div>
          </div>
        )

      case 'conversations':
        return (
          <ConversationsTab
            organizationId={currentOrganization!.id}
            memberId={memberId!}
            memberName={member!.username}
            onViewConversation={handleViewConversation}
            onCreateConversation={handleCreateConversation}
          />
        )

      case 'ai-chat':
        if (!member || !currentOrganization) {
          return <div>Loading...</div>
        }
        
        const chatContext: ChatContext = {
          entityType: 'member',
          entityId: memberId!,
          entityName: member.username,
          organizationId: currentOrganization.id,
          context: 'overview'
        }
        
        return (
          <div className="space-y-4">
            <div className="bg-card rounded-lg p-4">
              <h3 className="text-lg font-semibold text-foreground mb-2">AI Assistant</h3>
              <p className="text-muted-foreground text-sm mb-4">
                Chat with our AI assistant to get insights about {member.username}. Ask questions about their performance, conversations, or any other aspects.
              </p>
              <AIChat context={chatContext} className="h-[600px]" />
            </div>
          </div>
        )

      case 'ai-code-assistant-usage':
        return (
          <div className="space-y-4">
            {/* Date Filter Controls */}
            <div className="bg-card rounded-lg p-4">
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-lg font-semibold text-foreground">Filter by Date Range</h3>
                {aiUsageStatsLoading && (
                  <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                    <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Loading...</span>
                  </div>
                )}
              </div>
              <div className="flex space-x-3">
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Start Date
                  </label>
                  <input
                    type="date"
                    value={aiUsageDateParams.startDate || ''}
                    onChange={(e) => handleAiUsageDateChange('startDate', e.target.value)}
                    className="px-3 py-2 border border-border rounded bg-card text-foreground focus:outline-none focus:ring-2 focus:ring-primary text-sm"
                    style={{ colorScheme: 'dark light' }}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    End Date
                  </label>
                  <input
                    type="date"
                    value={aiUsageDateParams.endDate || ''}
                    onChange={(e) => handleAiUsageDateChange('endDate', e.target.value)}
                    className="px-3 py-2 border border-border rounded bg-card text-foreground focus:outline-none focus:ring-2 focus:ring-primary text-sm"
                    style={{ colorScheme: 'dark light' }}
                  />
                </div>
              </div>
            </div>

            {/* Usage Stats */}
            <div className="bg-card rounded-lg p-4">
              <h3 className="text-lg font-semibold text-foreground mb-4">Usage Statistics</h3>
              {aiUsageStatsLoading ? (
                <div className="flex items-center justify-center h-24">
                  <div className="flex items-center space-x-2 text-muted-foreground">
                    <svg className="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Loading stats...</span>
                  </div>
                </div>
              ) : aiUsageStats !== null ? (
                <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
                  <div className="bg-muted/20 rounded-lg p-4 border border-border/30">
                    <div className="text-sm text-muted-foreground mb-1">Total Lines Accepted</div>
                    <div className="text-2xl font-bold text-foreground">
                      {aiUsageStats.total_lines_accepted.toLocaleString()}
                    </div>
                  </div>
                  <div className="bg-muted/20 rounded-lg p-4 border border-border/30">
                    <div className="text-sm text-muted-foreground mb-1">Total Suggestions</div>
                    <div className="text-2xl font-bold text-foreground">
                      {aiUsageStats.total_suggestions.toLocaleString()}
                    </div>
                  </div>
                  <div className="bg-muted/20 rounded-lg p-4 border border-border/30">
                    <div className="text-sm text-muted-foreground mb-1">Accept Rate</div>
                    <div className="text-2xl font-bold text-foreground">
                      {aiUsageStats.overall_accept_rate.toFixed(1)}%
                    </div>
                  </div>
                  <div className="bg-muted/20 rounded-lg p-4 border border-border/30">
                    <div className="text-sm text-muted-foreground mb-1">Active Sessions</div>
                    <div className="text-2xl font-bold text-foreground">
                      {aiUsageStats.active_users}
                    </div>
                  </div>
                </div>
              ) : (
                <div className="flex items-center justify-center h-24">
                  <div className="text-center">
                    <div className="text-muted-foreground mb-1 text-sm">No usage data available</div>
                    <div className="text-xs text-muted-foreground">
                      Usage statistics will appear here once AI code assistant data is synced
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        )

      default:
        return null
    }
  }

  if (!memberId) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="text-red-600">Invalid member ID</div>
        </div>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">Loading member profile...</div>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-red-600">{error}</div>
          </div>
        </div>
      </div>
    )
  }

  if (!member) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-red-600">Member not found</div>
          </div>
        </div>
      </div>
    )
  }

  if (!currentOrganization) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">No organization selected</div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className={`max-w-6xl mx-auto ${showConversationSidebar ? 'mr-[28rem]' : ''}`}>
        {/* Breadcrumb Navigation */}
        <nav className="flex items-center space-x-2 text-sm text-muted-foreground mb-6">
          <a
            href={`/settings/members`}
            className="hover:text-foreground transition-colors"
          >
            Members
          </a>
          <span>/</span>
          <span className="text-foreground">
            {member.username}
          </span>
        </nav>

        {/* Member Basic Information */}
        <div className="bg-card rounded-lg p-4 mb-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              {/* Avatar Placeholder */}
              <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center">
                <svg className="w-6 h-6 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              
              <div>
                <h1 className="text-2xl font-semibold text-foreground">
                  {member.username}
                </h1>
                <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                  <span className="flex items-center space-x-1">
                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                    <span>{member.email}</span>
                  </span>
                  
                  {title && (
                    <span className="flex items-center space-x-1">
                      <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.815-8.764-2.245m0 0A23.023 23.023 0 014 12c0-3.183.815-6.22 2.245-8.764m0 0A23.023 23.023 0 0112 4c3.183 0 6.22.815 8.764 2.245M12 4v8m0 0v8" />
                      </svg>
                      <span>{title.name}</span>
                    </span>
                  )}
                  
                  {sourceControlAccount && (
                    <span className="flex items-center space-x-1">
                      <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                      </svg>
                      <span>{sourceControlAccount.username}</span>
                    </span>
                  )}
                </div>
              </div>
            </div>
            
            {/* Role Badge */}
            <div className="flex items-center space-x-2">
              {member.is_owner && (
                <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-purple-100 text-purple-800">
                  Owner
                </span>
              )}
              <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-blue-100 text-blue-800">
                Member
              </span>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="bg-card rounded-lg">
          <div className="border-b border-border/50">
            <nav className="flex space-x-6 px-4">
              <button
                onClick={() => setActiveTab('overview')}
                className={`py-3 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'overview'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border/50'
                }`}
              >
                Overview
              </button>
              <button
                onClick={() => setActiveTab('source-control-metrics')}
                className={`py-3 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'source-control-metrics'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border/50'
                }`}
              >
                Code Contributions
              </button>
              {title?.is_manager && (
                <button
                  onClick={() => setActiveTab('management-tree')}
                  className={`py-3 px-1 border-b-2 font-medium text-sm transition-colors ${
                    activeTab === 'management-tree'
                      ? 'border-primary text-primary'
                      : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border/50'
                  }`}
                >
                  Management Tree
                </button>
              )}
              <button
                onClick={() => setActiveTab('conversations')}
                className={`py-3 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'conversations'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border/50'
                }`}
              >
                Conversations
              </button>
              <button
                onClick={() => setActiveTab('ai-chat')}
                className={`py-3 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'ai-chat'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border/50'
                }`}
              >
                AI Assistant
              </button>
              <button
                onClick={() => setActiveTab('ai-code-assistant-usage')}
                className={`py-3 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'ai-code-assistant-usage'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border/50'
                }`}
              >
                AI Code Assistant Usage
              </button>
            </nav>
          </div>
          
          <div className="p-4">
            {renderTabContent()}
          </div>
        </div>
      </div>

      {/* Conversation Sidebar */}
      <ConversationSidebar
        conversation={selectedConversation}
        isOpen={showConversationSidebar}
        onClose={handleCloseConversationSidebar}
        onUpdate={handleConversationUpdate}
        onCreate={handleConversationCreate}
        mode={conversationSidebarMode}
        organizationId={currentOrganization?.id}
        memberId={memberId}
      />
    </div>
  )
} 