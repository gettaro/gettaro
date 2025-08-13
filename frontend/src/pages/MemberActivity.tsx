import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import Api from '../api/api'
import { MemberActivity, GetMemberActivityParams } from '../types/memberActivity'
import { Member } from '../types/member'
import { useOrganizationStore } from '../stores/organization'
import { useAuth } from '../hooks/useAuth'

export default function MemberActivityPage() {
  const { id: organizationId, memberId } = useParams<{ id: string; memberId: string }>()
  const { currentOrganization } = useOrganizationStore()
  const { isAuthenticated, isLoading: authLoading, getToken } = useAuth()
  const [activities, setActivities] = useState<MemberActivity[]>([])
  const [member, setMember] = useState<Member | null>(null)
  const [loading, setLoading] = useState(true)
  const [filterLoading, setFilterLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [dateParams, setDateParams] = useState<GetMemberActivityParams>({})
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set())

  useEffect(() => {
    if (isAuthenticated && !authLoading && currentOrganization && organizationId && memberId) {
      // Ensure we have a token before making API calls
      initializePage()
    }
  }, [isAuthenticated, authLoading, currentOrganization, organizationId, memberId])

  // Separate effect for date parameter changes
  useEffect(() => {
    if (isAuthenticated && !authLoading && currentOrganization && organizationId && memberId && (dateParams.startDate || dateParams.endDate)) {
      loadActivities()
    }
  }, [dateParams])

  const initializePage = async () => {
    try {
      // Ensure we have a fresh token
      await getToken()
      await loadMember()
      await loadActivities()
    } catch (err) {
      console.error('Error initializing page:', err)
      setError('Failed to initialize page')
    }
  }

  const loadMember = async () => {
    if (!organizationId || !memberId) return

    try {
      const members = await Api.getOrganizationMembers(organizationId)
      const foundMember = members.find(m => m.id === memberId)
      if (foundMember) {
        setMember(foundMember)
      }
    } catch (err) {
      console.error('Error loading member:', err)
    }
  }

  const loadActivities = async () => {
    if (!organizationId || !memberId) return

    try {
      if (dateParams.startDate || dateParams.endDate) {
        setFilterLoading(true)
      } else {
        setLoading(true)
      }
      setError(null)
      const activitiesData = await Api.getMemberActivity(organizationId, memberId, dateParams)
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
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
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

  const getActivityTypeLabel = (type: string) => {
    switch (type) {
      case 'pull_request':
        return 'Pull Request'
      case 'pr_comment':
        return 'Comment'
      case 'pr_review':
        return 'Review'
      default:
        return type
    }
  }

  if (!organizationId || !memberId) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="text-red-600">Invalid organization or member ID</div>
        </div>
      </div>
    )
  }

  if (authLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
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
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-red-600">Please log in to view this page</div>
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

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
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
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-red-600">{error}</div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto">
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
            <div className="divide-y divide-border">
              {activities.map((activity) => (
                <div key={activity.id} className={`p-6 ${activity.type !== 'pull_request' ? 'ml-8 border-l-2 border-border/30' : ''}`}>
                  <div className="flex items-start space-x-4">
                    <div className="flex-shrink-0 mt-1">
                      {getActivityIcon(activity.type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center space-x-2 mb-2">
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary/10 text-primary">
                          {getActivityTypeLabel(activity.type)}
                        </span>
                        {activity.repository && (
                          <span className="text-sm text-muted-foreground">
                            in {activity.repository}
                          </span>
                        )}
                        {activity.authorUsername && (
                          <span className="text-sm text-muted-foreground">
                            by @{activity.authorUsername}
                          </span>
                        )}
                      </div>
                      
                      <h3 className="text-lg font-medium text-foreground mb-2">
                        {activity.title}
                      </h3>
                      
                      {/* Show PR statistics for pull request activities */}
                      {activity.type === 'pull_request' && activity.metadata && (
                        <div className="mb-3 flex flex-wrap gap-4 text-sm text-muted-foreground">
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
                          {activity.metadata.comments !== undefined && (
                            <span className="flex items-center space-x-1">
                              <svg className="w-4 h-4 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                              </svg>
                              <span>{activity.metadata.comments} comments</span>
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
                        </div>
                      )}
                      
                      {/* Expandable Content - PR Description or Comment */}
                      {activity.description && (
                        <div className="mb-3">
                          <button
                            onClick={() => toggleExpanded(activity.type === 'pull_request' ? `${activity.id}-description` : `${activity.id}-comment`)}
                            className="flex items-center space-x-2 text-sm text-primary hover:text-primary/80 transition-colors"
                          >
                            <svg 
                              className={`w-4 h-4 transition-transform ${expandedItems.has(activity.type === 'pull_request' ? `${activity.id}-description` : `${activity.id}-comment`) ? 'rotate-90' : ''}`}
                              fill="none" 
                              stroke="currentColor" 
                              viewBox="0 0 24 24"
                            >
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                            </svg>
                            <span>
                              {expandedItems.has(activity.type === 'pull_request' ? `${activity.id}-description` : `${activity.id}-comment`) 
                                ? (activity.type === 'pull_request' ? 'Hide description' : 'Hide comment')
                                : (activity.type === 'pull_request' ? 'Show description' : 'Show comment')
                              }
                            </span>
                          </button>
                          {expandedItems.has(activity.type === 'pull_request' ? `${activity.id}-description` : `${activity.id}-comment`) && (
                            <div className="mt-2 p-3 bg-muted/30 rounded-md border border-border">
                              <p className="text-muted-foreground text-sm whitespace-pre-wrap">
                                {activity.description}
                              </p>
                            </div>
                          )}
                        </div>
                      )}
                      
                      <div className="flex items-center space-x-4 text-sm text-muted-foreground">
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
          )}
        </div>
      </div>
    </div>
  )
} 