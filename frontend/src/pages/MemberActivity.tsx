import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import Api from '../api/api'
import { MemberActivity, GetMemberActivityParams } from '../types/memberActivity'
import { Member } from '../types/member'
import { useOrganizationStore } from '../stores/organization'
import { useAuth } from '../hooks/useAuth'

export default function MemberActivityPage() {
  const { memberId } = useParams<{ memberId: string }>()
  const { currentOrganization } = useOrganizationStore()
  const { isAuthenticated, isLoading: authLoading, getToken } = useAuth()
  const [activities, setActivities] = useState<MemberActivity[]>([])
  const [member, setMember] = useState<Member | null>(null)
  const [loading, setLoading] = useState(true)
  const [filterLoading, setFilterLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [dateParams, setDateParams] = useState<GetMemberActivityParams>(() => {
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 14) // 2 weeks ago
    
    return {
      startDate: startDate.toISOString().split('T')[0], // YYYY-MM-DD format
      endDate: endDate.toISOString().split('T')[0]
    }
  })
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set())
  const [hasMoreData, setHasMoreData] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)

  // Utility function to convert seconds to human readable format
  const formatTimeFromSeconds = (seconds: number): string => {
    if (seconds < 60) return `${seconds}s`
    if (seconds < 3600) return `${Math.round(seconds / 60)}m`
    if (seconds < 86400) return `${Math.round(seconds / 3600)}h`
    return `${Math.round(seconds / 86400)}d`
  }

  useEffect(() => {
    if (isAuthenticated && !authLoading && currentOrganization && memberId) {
      initializePage()
    }
  }, [isAuthenticated, authLoading, currentOrganization, memberId])

  useEffect(() => {
    if (isAuthenticated && !authLoading && currentOrganization && memberId && (dateParams.startDate || dateParams.endDate)) {
      loadActivities()
    }
  }, [dateParams])

  const initializePage = async () => {
    try {
      await getToken()
      await loadMember()
      await loadActivities()
    } catch (err) {
      console.error('Error initializing page:', err)
      setError('Failed to initialize page')
    }
  }

  const loadMember = async () => {
    if (!currentOrganization?.id || !memberId) return

    try {
      const members = await Api.getOrganizationMembers(currentOrganization.id)
      const foundMember = members.find(m => m.id === memberId)
      if (foundMember) {
        setMember(foundMember)
      }
    } catch (err) {
      console.error('Error loading member:', err)
    }
  }

  const loadActivities = async () => {
    if (!currentOrganization?.id || !memberId) return

    try {
      if (dateParams.startDate || dateParams.endDate) {
        setFilterLoading(true)
      } else {
        setLoading(true)
      }
      setError(null)
      const activitiesData = await Api.getMemberActivity(currentOrganization.id, memberId, dateParams)
      setActivities(activitiesData)
    } catch (err) {
      setError('Failed to load activities')
      console.error('Error loading activities:', err)
    } finally {
      setLoading(false)
      setFilterLoading(false)
    }
  }

  const handleDateChange = (field: 'startDate' | 'endDate', value: string) => {
    setDateParams(prev => ({
      ...prev,
      [field]: value || undefined
    }))
  }

  const loadMore = async () => {
    if (!currentOrganization?.id || !memberId || loadingMore) return

    try {
      setLoadingMore(true)
      
      // Calculate new start date (1 week earlier than current start date)
      const currentStartDate = new Date(dateParams.startDate!)
      const newStartDate = new Date(currentStartDate)
      newStartDate.setDate(newStartDate.getDate() - 7)
      
      const newDateParams = {
        startDate: newStartDate.toISOString().split('T')[0],
        endDate: dateParams.startDate! // End date becomes the old start date
      }
      
      // Load activities for the new date range
      const moreActivities = await Api.getMemberActivity(currentOrganization.id, memberId, newDateParams)
      
      if (moreActivities.length === 0) {
        setHasMoreData(false)
      } else {
        // Prepend new activities to existing ones
        setActivities(prev => [...moreActivities, ...prev])
        // Update date params to reflect the new range
        setDateParams(prev => ({
          ...prev,
          startDate: newDateParams.startDate
        }))
      }
    } catch (err) {
      console.error('Error loading more activities:', err)
      setError('Failed to load more activities')
    } finally {
      setLoadingMore(false)
    }
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
                    <span>Merge: {formatTimeFromSeconds(activity.prMetrics.time_to_merge_seconds)}</span>
                  </span>
                )}
                
                {/* Time to first review */}
                {activity.prMetrics.time_to_first_non_bot_review_seconds !== undefined && (
                  <span className="flex items-center space-x-1">
                    <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span>First review: {formatTimeFromSeconds(activity.prMetrics.time_to_first_non_bot_review_seconds)}</span>
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

  if (!currentOrganization) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">No organization selected</div>
          </div>
        </div>
      </div>
    )
  }

  if (authLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">Authenticating...</div>
          </div>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-red-600">Please log in to view this page</div>
          </div>
        </div>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">Loading activities...</div>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-red-600">{error}</div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          {/* Breadcrumb Navigation */}
          <nav className="flex items-center space-x-2 text-sm text-muted-foreground mb-4">
            <a
              href={`/settings/members`}
              className="hover:text-foreground transition-colors"
            >
              Members
            </a>
            <span>/</span>
            <span className="text-foreground">
              {member ? member.username : 'Member'} Activity
            </span>
          </nav>
          
          <h1 className="text-3xl font-bold text-foreground mb-2">
            {member ? `${member.username}'s Activity` : 'Member Activity'}
          </h1>
          <p className="text-muted-foreground">
            Source control activity timeline for this team member.
          </p>
        </div>

        {/* Date Filter Controls */}
        <div className="bg-card border border-border rounded-lg p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-foreground">Filter by Date Range</h2>
            {filterLoading && (
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

        {/* Activity Timeline */}
        <div className="bg-card border border-border rounded-lg">
          <div className="p-6 border-b border-border">
            <h2 className="text-xl font-semibold text-foreground">
              Activity Timeline
              {activities.length > 0 && (
                <span className="ml-2 text-sm text-muted-foreground">
                  ({activities.length} activities)
                </span>
              )}
            </h2>
          </div>

          {activities.length === 0 ? (
            <div className="p-6 text-center">
              <div className="text-muted-foreground">
                No activities found for the selected date range.
              </div>
            </div>
          ) : (
            <>
              {/* Scrollable Activity List */}
              <div className="h-96 overflow-y-auto">
                <div className="divide-y divide-border">
                  {activities.map((activity) => (
                    <div key={activity.id} className="p-6">
                      <div className="flex items-start space-x-4">
                        <div className="flex-shrink-0 mt-1">
                          {getActivityIcon(activity.type)}
                        </div>
                        <div className="flex-1 min-w-0">
                          {renderActivityContent(activity)}
                          
                          <div className="flex items-center space-x-4 text-sm text-muted-foreground mt-3">
                            <span>
                              {new Date(activity.createdAt).toLocaleDateString('en-US', {
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                                hour: '2-digit',
                                minute: '2-digit'
                              })}
                            </span>
                            {activity.url && (
                              <a
                                href={activity.url}
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
              </div>

              {/* Load More Button */}
              {hasMoreData && (
                <div className="p-6 border-t border-border">
                  <button
                    onClick={loadMore}
                    disabled={loadingMore}
                    className="w-full py-3 px-4 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center space-x-2"
                  >
                    {loadingMore ? (
                      <>
                        <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        <span>Loading...</span>
                      </>
                    ) : (
                      <span>Load More (Previous Week)</span>
                    )}
                  </button>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
} 