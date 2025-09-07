import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { useOrganizationStore } from '../stores/organization'
import { Member } from '../types/member'
import { GetMemberMetricsResponse, GetMemberMetricsParams } from '../types/memberMetrics'
import { Title } from '../types/title'
import { SourceControlAccount, PullRequest, GetMemberPullRequestsParams, GetMemberPullRequestReviewsParams } from '../types/sourcecontrol'
import { MemberActivity, GetMemberActivityParams } from '../types/sourcecontrol'
import { OrgChartNode } from '../types/directs'
import Api from '../api/api'
import { formatMetricValue, formatTimeMetric } from '../utils/formatMetrics'
import MetricIcon from '../components/MetricIcon'

type TabType = 'overview' | 'source-control-metrics'

export default function MemberProfilePage() {
  const { memberId } = useParams<{ memberId: string }>()
  const { currentOrganization } = useOrganizationStore()
  const [member, setMember] = useState<Member | null>(null)
  const [title, setTitle] = useState<Title | null>(null)
  const [sourceControlAccount, setSourceControlAccount] = useState<SourceControlAccount | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<TabType>('overview')
  
  // Date filter state for Code Contributions tab
  const [dateParams, setDateParams] = useState<GetMemberActivityParams>(() => {
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 30) // 1 month ago
    
    return {
      startDate: startDate.toISOString().split('T')[0], // YYYY-MM-DD format
      endDate: endDate.toISOString().split('T')[0]
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
  }, [activeTab, currentOrganization?.id, memberId, dateParams.startDate, dateParams.endDate])

  // Load management tree when member is loaded and is a manager
  useEffect(() => {
    if (member && title?.isManager && currentOrganization?.id && memberId) {
      loadManagementTree()
    }
  }, [member, title, currentOrganization?.id, memberId])

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
        interval: 'monthly'
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

  const loadManagementTree = async () => {
    if (!currentOrganization?.id || !memberId || !member) return

    try {
      setManagementTreeLoading(true)
      setError(null)

      const response = await Api.getManagerTree(currentOrganization.id, memberId)
      console.log('Management tree response:', response)
      console.log('OrgChart:', response.orgChart)
      if (response.orgChart && response.orgChart.length > 0) {
        console.log('First node:', response.orgChart[0])
      }
      setManagementTree(response.orgChart)
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

  const renderPullRequestContent = (pr: PullRequest) => {
    return (
      <div>
        <h3 className="text-lg font-medium text-foreground mb-3">
          {pr.title}
        </h3>
        
        {/* PR Statistics */}
        <div className="flex flex-wrap gap-4 text-sm text-muted-foreground mb-3">
          {/* PR Status */}
          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
            pr.status === 'open' 
              ? 'bg-green-100 text-green-800' 
              : pr.status === 'closed' && pr.merged_at
              ? 'bg-purple-100 text-purple-800'
              : pr.status === 'closed'
              ? 'bg-red-100 text-red-800'
              : 'bg-gray-100 text-gray-800'
          }`}>
            {pr.status === 'open' ? 'Open' : 
             pr.status === 'closed' && pr.merged_at ? 'Merged' :
             pr.status === 'closed' ? 'Closed' : 
             pr.status}
          </span>
          
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            <span>+{pr.additions}</span>
          </span>
          
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
            </svg>
            <span>-{pr.deletions}</span>
          </span>
          
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
            <span>{pr.comments} comments</span>
          </span>
          
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
            <span>{pr.review_comments} review comments</span>
          </span>
          
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span>{pr.changed_files} files</span>
          </span>
          
          {pr.merged_at && (
            <span className="flex items-center space-x-1">
              <svg className="w-4 h-4 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{(() => {
                const created = new Date(pr.created_at)
                const merged = new Date(pr.merged_at)
                const diffMs = merged.getTime() - created.getTime()
                const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
                const diffHours = Math.floor((diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
                
                if (diffDays > 0) {
                  return `${diffDays}d ${diffHours}h`
                } else if (diffHours > 0) {
                  return `${diffHours}h`
                } else {
                  const diffMinutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60))
                  return `${diffMinutes}m`
                }
              })()} to merge</span>
            </span>
          )}
        </div>
        
        {/* Expandable PR Description */}
        {pr.description && (
          <div className="mb-3">
            <button
              onClick={() => toggleExpanded(`${pr.id}-description`)}
              className="flex items-center space-x-2 text-sm text-primary hover:text-primary/80 transition-colors"
            >
              <svg 
                className={`w-4 h-4 transition-transform ${expandedItems.has(`${pr.id}-description`) ? 'rotate-90' : ''}`}
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
              <span>
                {expandedItems.has(`${pr.id}-description`) ? 'Hide description' : 'Show description'}
              </span>
            </button>
            {expandedItems.has(`${pr.id}-description`) && (
              <div className="mt-2 p-3 bg-muted/30 rounded-md border border-border">
                <p className="text-muted-foreground text-sm whitespace-pre-wrap">
                  {pr.description}
                </p>
              </div>
            )}
          </div>
        )}
      </div>
    )
  }

  const renderActivityContent = (activity: MemberActivity) => {
    switch (activity.type) {
      case 'pull_request':
        return (
          <div>
            <div className="text-sm text-muted-foreground mb-2">
              <span className="font-medium">@{activity.authorUsername}</span> has created a new PR
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
            {activity.prMetrics && (
              <div className="flex flex-wrap gap-4 text-sm text-muted-foreground mb-3">
                {/* Time to merge */}
                {activity.prMetrics.time_to_merge_seconds !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span>Merge: {formatTimeMetric(activity.prMetrics.time_to_merge_seconds)}</span>
                  </span>
                )}
                
                {/* Time to first review */}
                {activity.prMetrics.time_to_first_non_bot_review_seconds !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span>First review: {formatTimeMetric(activity.prMetrics.time_to_first_non_bot_review_seconds)}</span>
                  </span>
                )}

                {/* Show opened duration for open PRs */}
                {activity.metadata?.state === 'open' && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span>Opened {Math.ceil((Date.now() - new Date(activity.createdAt).getTime()) / (1000 * 60 * 60 * 24))}d</span>
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
              <span className="font-medium">@{activity.authorUsername}</span> commented on{' '}
              <span className="font-medium">{activity.prTitle}</span> from{' '}
              <span className="font-medium">@{activity.prAuthorUsername}</span>
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
              <span className="font-medium">@{activity.authorUsername}</span> has reviewed{' '}
              <span className="font-medium">{activity.prTitle}</span> from{' '}
              <span className="font-medium">@{activity.prAuthorUsername}</span>
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
      if (foundMember.titleId) {
        try {
          const titles = await Api.getOrganizationTitles(currentOrganization.id)
          const foundTitle = titles.find(t => t.id === foundMember.titleId)
          setTitle(foundTitle || null)
        } catch (err) {
          console.error('Error loading title:', err)
        }
      }

      // Load source control account if member has one
      try {
        const sourceControlAccounts = await Api.getOrganizationSourceControlAccounts(currentOrganization.id)
        const foundAccount = sourceControlAccounts.find(acc => acc.memberId === memberId)
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

  const renderManagementTreeNode = (node: OrgChartNode, depth: number = 0) => {
    console.log('Rendering node:', node)
    console.log('Node member:', node.member)
    
    const indentClass = `ml-${depth * 4}`
    
    return (
      <div key={node.member.id} className={`${indentClass} mb-2`}>
        <div className="flex items-center space-x-3 p-3 bg-muted/30 rounded-lg border border-border">
          <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
            <span className="text-primary font-medium text-sm">
              {node.member.username.charAt(0).toUpperCase()}
            </span>
          </div>
          <div>
            <h4 className="font-medium text-foreground">{node.member.username}</h4>
            <p className="text-sm text-muted-foreground">{node.member.email}</p>
          </div>
          {depth > 0 && (
            <div className="ml-auto">
              <span className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded">
                Level {depth}
              </span>
            </div>
          )}
        </div>
        {node.directReports != null && node.directReports.length > 0 && (
          <div className="mt-2">
            {node.directReports.map(report => renderManagementTreeNode(report, depth + 1))}
          </div>
        )}
      </div>
    )
  }

  const renderTabContent = () => {
    switch (activeTab) {
      case 'overview':
        return (
          <div className="space-y-6">
            {/* Management Tree Section - Only show if member is a manager */}
            {title?.isManager && (
              <div className="bg-card border border-border rounded-lg p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold text-foreground">Management Tree</h3>
                  {managementTreeLoading && (
                    <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                      <div className="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
                      <span>Loading...</span>
                    </div>
                  )}
                </div>
                
                {managementTreeLoading ? (
                  <div className="text-center py-8">
                    <div className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-2"></div>
                    <p className="text-muted-foreground">Loading management tree...</p>
                  </div>
                ) : managementTree.length > 0 ? (
                  <div className="space-y-2">
                    <p className="text-sm text-muted-foreground mb-4">
                      People who report directly or indirectly to {member?.username}:
                    </p>
                    {managementTree.map(node => renderManagementTreeNode(node))}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <div className="w-12 h-12 bg-muted rounded-full flex items-center justify-center mx-auto mb-3">
                      <svg className="w-6 h-6 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                      </svg>
                    </div>
                    <p className="text-muted-foreground">No direct reports found</p>
                    <p className="text-sm text-muted-foreground mt-1">
                      This manager doesn't have any direct reports yet.
                    </p>
                  </div>
                )}
              </div>
            )}

            <div className="bg-card border border-border rounded-lg p-6">
              <h3 className="text-lg font-semibold text-foreground mb-4">Source Control Overview</h3>
              <p className="text-muted-foreground">
                This section will contain source control metrics. Coming soon...
              </p>
            </div>
          </div>
        )

      case 'source-control-metrics':
        return (
          <div className="space-y-6">
            {/* Date Filter Controls */}
            <div className="bg-card border border-border rounded-lg p-6">
              <div className="flex items-center justify-between mb-4">
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
              <div className="flex space-x-4">
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Start Date
                  </label>
                  <input
                    type="date"
                    value={dateParams.startDate || ''}
                    onChange={(e) => handleDateChange('startDate', e.target.value)}
                    className="px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
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
                    className="px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                  />
                </div>
              </div>
            </div>

            {/* Metrics Section */}
            <div className="bg-card border border-border rounded-lg p-6">
              <h3 className="text-lg font-semibold text-foreground mb-6">Metrics</h3>
              
              {metricsLoading ? (
                <div className="flex items-center justify-center h-32">
                  <div className="flex items-center space-x-2 text-muted-foreground">
                    <svg className="animate-spin h-6 w-6" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Loading metrics...</span>
                  </div>
                </div>
              ) : metrics ? (
                <div className="space-y-8">
                  {metrics.snapshotMetrics.map((category) => (
                    <div key={category.category.name} className="space-y-4">
                      <h4 className="text-md font-semibold text-foreground">{category.category.name}</h4>
                      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {category.metrics.map((metric) => (
                          <div key={metric.label} className="text-center p-4 bg-muted/30 rounded-lg border border-border">
                            <div className="flex justify-center mb-2">
                              {/* Use backend-provided icon */}
                              <MetricIcon 
                                iconIdentifier={metric.iconIdentifier || 'default'} 
                                iconColor={metric.iconColor || 'gray'} 
                              />
                            </div>
                            <div className="text-2xl font-bold text-foreground">
                              {formatMetricValue(metric.value, metric.unit)}
                            </div>
                            <div className="text-sm text-muted-foreground">{metric.label}</div>
                            {typeof metric.peersValue === 'number' && metric.peersValue > 0 && (
                              <div className="text-xs text-muted-foreground mt-1">
                                vs {formatMetricValue(metric.peersValue, metric.unit)} (peers)
                              </div>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="flex items-center justify-center h-32">
                  <div className="text-center">
                    <div className="text-muted-foreground mb-2">No metrics available</div>
                    <div className="text-sm text-muted-foreground">
                      Metrics will appear here once you have source control activity
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* Recent Pull Requests Table */}
            <div className="bg-card border border-border rounded-lg">
              <div className="p-6 border-b border-border">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-lg font-semibold text-foreground">Pull Requests</h3>
                    <p className="text-sm text-muted-foreground mt-1">
                      Pull requests created in the selected date range
                    </p>
                  </div>
                  <div className="flex items-center space-x-4">
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
                <div className="divide-y divide-border">
                  {pullRequests.length === 0 ? (
                    <div className="p-6 text-center">
                      <div className="text-muted-foreground">
                        No pull requests found for the selected date range.
                      </div>
                    </div>
                  ) : (
                    <div className="h-96 overflow-y-auto">
                      {pullRequests.map((pr) => (
                        <div key={pr.id} className="p-6">
                          <div className="flex items-start space-x-4">
                            <div className="flex-shrink-0 mt-1">
                              <svg className="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" />
                              </svg>
                            </div>
                            <div className="flex-1 min-w-0">
                              {renderPullRequestContent(pr)}
                              
                              <div className="flex items-center space-x-4 text-sm text-muted-foreground mt-3">
                                <span>
                                  {new Date(pr.created_at).toLocaleDateString('en-US', {
                                    year: 'numeric',
                                    month: 'long',
                                    day: 'numeric',
                                    hour: '2-digit',
                                    minute: '2-digit'
                                  })}
                                </span>
                                {pr.merged_at && (
                                  <span className="text-green-600">
                                    Merged {new Date(pr.merged_at).toLocaleString('en-US', {
                                      year: 'numeric',
                                      month: 'long',
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
            <div className="bg-card border border-border rounded-lg">
              <div className="p-6 border-b border-border">
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
                <div className="divide-y divide-border">
                  {pullRequestReviewsLoading ? (
                    <div className="p-6 text-center">
                      <div className="text-muted-foreground">Loading PR reviews...</div>
                    </div>
                  ) : pullRequestReviews.length === 0 ? (
                    <div className="p-6 text-center">
                      <div className="text-muted-foreground">
                        No PR reviews found for the selected date range.
                      </div>
                    </div>
                  ) : (
                    <div className="h-96 overflow-y-auto">
                      {pullRequestReviews.map((review) => (
                          <div key={review.id} className="p-6">
                            <div className="flex items-start space-x-4">
                              <div className="flex-shrink-0 mt-1">
                                {getActivityIcon(review.type)}
                              </div>
                              <div className="flex-1 min-w-0">
                                {renderActivityContent(review)}
                                
                                <div className="flex items-center space-x-4 text-sm text-muted-foreground mt-3">
                                  <span>
                                    {new Date(review.createdAt).toLocaleDateString('en-US', {
                                      year: 'numeric',
                                      month: 'long',
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
      <div className="max-w-6xl mx-auto">
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
        <div className="bg-card border border-border rounded-lg p-6 mb-8">
          <div className="flex items-start justify-between">
            <div className="flex items-center space-x-4">
              {/* Avatar Placeholder */}
              <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                <svg className="w-8 h-8 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              
              <div>
                <h1 className="text-3xl font-bold text-foreground mb-2">
                  {member.username}
                </h1>
                <div className="space-y-2 text-sm text-muted-foreground">
                  <div className="flex items-center space-x-2">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                    <span>{member.email}</span>
                  </div>
                  
                  {title && (
                    <div className="flex items-center space-x-2">
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.815-8.764-2.245m0 0A23.023 23.023 0 014 12c0-3.183.815-6.22 2.245-8.764m0 0A23.023 23.023 0 0112 4c3.183 0 6.22.815 8.764 2.245M12 4v8m0 0v8" />
                      </svg>
                      <span>{title.name}</span>
                    </div>
                  )}
                  
                  {sourceControlAccount && (
                    <div className="flex items-center space-x-2">
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                      </svg>
                      <span>{sourceControlAccount.username} ({sourceControlAccount.providerName})</span>
                    </div>
                  )}
                </div>
              </div>
            </div>
            
            {/* Role Badge */}
            <div className="flex items-center space-x-2">
              {member.isOwner && (
                <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                  Owner
                </span>
              )}
              <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                Member
              </span>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="bg-card border border-border rounded-lg">
          <div className="border-b border-border">
            <nav className="flex space-x-8 px-6">
              <button
                onClick={() => setActiveTab('overview')}
                className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'overview'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border'
                }`}
              >
                Overview
              </button>
              <button
                onClick={() => setActiveTab('source-control-metrics')}
                className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'source-control-metrics'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border'
                }`}
              >
                Code Contributions
              </button>
            </nav>
          </div>
          
          <div className="p-6">
            {renderTabContent()}
          </div>
        </div>
      </div>
    </div>
  )
} 