import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useOrganizationStore } from '../stores/organization'
import { Team } from '../types/team'
import { Member } from '../types/member'
import { PullRequest } from '../types/sourcecontrol'
import { OrganizationMetricsResponse } from '../types/organizationMetrics'
import { GetMemberAICodeAssistantMetricsParams, GetMemberAICodeAssistantMetricsResponse } from '../types/aicodeassistant'
import Api from '../api/api'
import MetricChart from '../components/MetricChart'
import MetricInfoButton from '../components/MetricInfoButton'
import PullRequestItem from '../components/PullRequestItem'

type TabType = 'overview' | 'code-contributions' | 'ai-code-assistant'

export default function TeamProfilePage() {
  const { teamId } = useParams<{ teamId: string }>()
  const navigate = useNavigate()
  const { currentOrganization } = useOrganizationStore()
  
  const [team, setTeam] = useState<Team | null>(null)
  const [members, setMembers] = useState<Member[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<TabType>('overview')
  
  // Date filter state
  const [dateParams, setDateParams] = useState(() => {
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 30) // 1 month ago
    
    return {
      startDate: startDate.toISOString().split('T')[0],
      endDate: endDate.toISOString().split('T')[0],
      interval: 'weekly' as 'daily' | 'weekly' | 'monthly'
    }
  })
  
  const [openPRs, setOpenPRs] = useState<PullRequest[]>([])
  const [openPRsLoading, setOpenPRsLoading] = useState(false)
  const [allPRs, setAllPRs] = useState<PullRequest[]>([])
  const [allPRsLoading, setAllPRsLoading] = useState(false)
  const [metrics, setMetrics] = useState<OrganizationMetricsResponse | null>(null)
  const [metricsLoading, setMetricsLoading] = useState(false)
  const [aiMetrics, setAiMetrics] = useState<GetMemberAICodeAssistantMetricsResponse | null>(null)
  const [aiMetricsLoading, setAiMetricsLoading] = useState(false)
  const [expandedTables, setExpandedTables] = useState<Set<string>>(new Set(['pull-requests']))
  const [currentGraphIndex, setCurrentGraphIndex] = useState(0)
  const [aiCurrentGraphIndex, setAiCurrentGraphIndex] = useState(0)

  // Load team data
  useEffect(() => {
    if (!currentOrganization?.id || !teamId) return

    const loadTeam = async () => {
      try {
        setLoading(true)
        const response = await Api.getTeam(currentOrganization.id, teamId)
        setTeam(response.team)
      } catch (err) {
        console.error('Error loading team:', err)
        setError('Failed to load team')
      } finally {
        setLoading(false)
      }
    }

    loadTeam()
  }, [currentOrganization?.id, teamId])

  // Load team members
  useEffect(() => {
    if (!team || !currentOrganization?.id) return

    const loadMembers = async () => {
      try {
        const allMembers = await Api.getOrganizationMembers(currentOrganization.id)
        const teamMemberIds = team.members.map(tm => tm.member_id)
        const teamMembers = allMembers.filter(m => teamMemberIds.includes(m.id))
        setMembers(teamMembers)
      } catch (err) {
        console.error('Error loading team members:', err)
      }
    }

    loadMembers()
  }, [team, currentOrganization?.id])

  // Load open PRs for team members
  useEffect(() => {
    if (!currentOrganization?.id || !team) return

    const loadOpenPRs = async () => {
      setOpenPRsLoading(true)
      try {
        // Always use team prefix for filtering
        if (team.pr_prefix) {
          const prs = await Api.getOrganizationPullRequests(currentOrganization.id, {
            prefix: team.pr_prefix,
            status: 'open'
          })
          setOpenPRs(prs)
        } else {
          // No prefix, return empty array
          setOpenPRs([])
        }
      } catch (err) {
        console.error('Error loading open PRs:', err)
        setOpenPRs([])
      } finally {
        setOpenPRsLoading(false)
      }
    }

    loadOpenPRs()
  }, [currentOrganization?.id, team])

  // Load all PRs for team members in date range
  useEffect(() => {
    if (!currentOrganization?.id || !team || activeTab !== 'code-contributions') {
      return
    }

    const loadAllPRs = async () => {
      setAllPRsLoading(true)
      try {
        // Always use team prefix for filtering
        if (team.pr_prefix) {
          const prs = await Api.getOrganizationPullRequests(currentOrganization.id, {
            prefix: team.pr_prefix,
            startDate: dateParams.startDate,
            endDate: dateParams.endDate
          })
          setAllPRs(prs)
        } else {
          // No prefix, return empty array
          setAllPRs([])
        }
      } catch (err) {
        console.error('Error loading PRs:', err)
        setAllPRs([])
      } finally {
        setAllPRsLoading(false)
      }
    }

    loadAllPRs()
  }, [currentOrganization?.id, team, dateParams.startDate, dateParams.endDate, activeTab])

  // Load metrics
  useEffect(() => {
    if (!currentOrganization?.id || !team || activeTab !== 'code-contributions') return

    const loadMetrics = async () => {
      setMetricsLoading(true)
      try {
        const metricsData = await Api.getOrganizationMetrics(currentOrganization.id, {
          startDate: dateParams.startDate,
          endDate: dateParams.endDate,
          interval: dateParams.interval,
          teamIds: [team.id]
        })
        setMetrics(metricsData)
      } catch (err) {
        console.error('Error loading metrics:', err)
      } finally {
        setMetricsLoading(false)
      }
    }

    loadMetrics()
  }, [currentOrganization?.id, team, dateParams.startDate, dateParams.endDate, dateParams.interval, activeTab])

  // Load AI code assistant metrics
  useEffect(() => {
    if (!currentOrganization?.id || !team || activeTab !== 'ai-code-assistant') return

    const loadAiMetrics = async () => {
      setAiMetricsLoading(true)
      try {
        const params: GetMemberAICodeAssistantMetricsParams = {
          startDate: dateParams.startDate,
          endDate: dateParams.endDate,
          interval: dateParams.interval
        }
        const aiMetricsData = await Api.getTeamAICodeAssistantMetrics(currentOrganization.id, teamId!, params)
        setAiMetrics(aiMetricsData)
      } catch (err) {
        console.error('Error loading AI metrics:', err)
      } finally {
        setAiMetricsLoading(false)
      }
    }

    loadAiMetrics()
  }, [currentOrganization?.id, team, dateParams.startDate, dateParams.endDate, dateParams.interval, activeTab, teamId])

  // Reset graph index when metrics change
  useEffect(() => {
    if (metrics) {
      const allGraphs = getAllGraphs(metrics)
      if (allGraphs.length > 0 && currentGraphIndex >= allGraphs.length) {
        setCurrentGraphIndex(0)
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [metrics])

  const handleDateChange = (field: 'startDate' | 'endDate', value: string) => {
    setDateParams(prev => ({
      ...prev,
      [field]: value
    }))
  }

  const handleIntervalChange = (interval: 'daily' | 'weekly' | 'monthly') => {
    setDateParams(prev => ({
      ...prev,
      interval
    }))
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

  // Get all graphs (flattened from all categories)
  const getAllGraphs = (metricsData: OrganizationMetricsResponse) => {
    if (!metricsData.graph_metrics || metricsData.graph_metrics.length === 0) {
      return []
    }

    const allGraphs: Array<{ metric: any; category: string; description: string }> = []

    metricsData.graph_metrics.forEach((category) => {
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

  // Get all graphs for AI metrics
  const getAllAiGraphs = (metricsData: GetMemberAICodeAssistantMetricsResponse) => {
    if (!metricsData.graph_metrics || metricsData.graph_metrics.length === 0) {
      return []
    }

    const allGraphs: Array<{ 
      metric: any
      category: string
    }> = []

    metricsData.graph_metrics.forEach((category) => {
      const metricsWithData = category.metrics.filter((metric) => {
        if (!metric.time_series || metric.time_series.length === 0) {
          return false
        }
        return metric.time_series.some(entry => 
          entry.data && entry.data.length > 0
        )
      })

      metricsWithData.forEach((metric) => {
        allGraphs.push({
          metric: {
            ...metric,
            type: metric.type || 'line',
            unit: metric.unit || 'count'
          },
          category: category.category
        })
      })
    })

    return allGraphs
  }

  const navigateAiGraph = (direction: 'prev' | 'next', totalGraphs: number) => {
    let newIndex: number
    
    if (direction === 'prev') {
      newIndex = aiCurrentGraphIndex > 0 ? aiCurrentGraphIndex - 1 : totalGraphs - 1
    } else {
      newIndex = aiCurrentGraphIndex < totalGraphs - 1 ? aiCurrentGraphIndex + 1 : 0
    }
    
    setAiCurrentGraphIndex(newIndex)
  }

  // Reset AI graph index when AI metrics change
  useEffect(() => {
    if (aiMetrics) {
      const allGraphs = getAllAiGraphs(aiMetrics)
      if (allGraphs.length > 0 && aiCurrentGraphIndex >= allGraphs.length) {
        setAiCurrentGraphIndex(0)
      }
    }
  }, [aiMetrics, aiCurrentGraphIndex])


  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center h-64">
          <div className="text-muted-foreground">Loading team profile...</div>
        </div>
      </div>
    )
  }

  if (error || !team) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="bg-destructive/10 text-destructive px-4 py-3 rounded">
          {error || 'Team not found'}
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-6">
          <button
            onClick={() => navigate('/dashboard')}
            className="text-muted-foreground hover:text-foreground mb-4 text-sm flex items-center gap-1"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            Back to Dashboard
          </button>
          
          <div className="flex items-start justify-between">
            <div>
              <h1 className="text-3xl font-bold mb-2">{team.name}</h1>
              {team.description && (
                <p className="text-muted-foreground">{team.description}</p>
              )}
              <div className="mt-4 flex items-center gap-4">
                <div>
                  <span className="text-sm text-muted-foreground">Members: </span>
                  <span className="text-sm font-medium">{members.length}</span>
                </div>
              </div>
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
                onClick={() => setActiveTab('code-contributions')}
                className={`py-3 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'code-contributions'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border/50'
                }`}
              >
                Code Contributions
              </button>
              <button
                onClick={() => setActiveTab('ai-code-assistant')}
                className={`py-3 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'ai-code-assistant'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border/50'
                }`}
              >
                AI Code Assistant
              </button>
            </nav>
          </div>
          
          <div className="p-4">
            {activeTab === 'overview' && (
              <div className="space-y-4">
                {/* Team Members */}
                <div className="bg-muted/30 rounded-lg p-4">
                  <h3 className="text-lg font-semibold mb-3">Team Members</h3>
                  {members.length === 0 ? (
                    <p className="text-sm text-muted-foreground">No members in this team</p>
                  ) : (
                    <div className="flex flex-wrap gap-2">
                      {members.map((member) => (
                        <button
                          key={member.id}
                          onClick={() => navigate(`/members/${member.id}/profile`)}
                          className="px-3 py-1 bg-background rounded border border-border text-sm hover:bg-primary/10 hover:border-primary/50 cursor-pointer transition-colors text-left"
                        >
                          {member.username} ({member.email})
                        </button>
                      ))}
                    </div>
                  )}
                </div>

                {/* Open Pull Requests */}
                <div className="bg-card rounded-lg">
                  <div className="p-4 border-b border-border/50">
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="text-lg font-semibold text-foreground">Open Pull Requests</h3>
                        <p className="text-sm text-muted-foreground mt-1">
                          Open pull requests from team members
                        </p>
                      </div>
                      <button
                        onClick={() => toggleTableExpanded('open-prs')}
                        className="flex items-center space-x-2 text-sm text-primary hover:text-primary/80 transition-colors"
                      >
                        <svg 
                          className={`w-4 h-4 transition-transform ${expandedTables.has('open-prs') ? 'rotate-90' : ''}`}
                          fill="none" 
                          stroke="currentColor" 
                          viewBox="0 0 24 24"
                        >
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                        <span>{expandedTables.has('open-prs') ? 'Collapse' : 'Expand'}</span>
                      </button>
                    </div>
                  </div>
                  {expandedTables.has('open-prs') && (
                    <div className="divide-y divide-border/50">
                      {openPRsLoading ? (
                        <div className="p-4 text-center">
                          <div className="text-muted-foreground text-sm">Loading...</div>
                        </div>
                      ) : openPRs.length === 0 ? (
                        <div className="p-4 text-center">
                          <div className="text-muted-foreground text-sm">No open pull requests</div>
                        </div>
                      ) : (
                        <div className="max-h-80 overflow-y-auto">
                          {openPRs.map((pr) => (
                            <div key={pr.id} className="p-4">
                              <div className="flex items-start space-x-3">
                                <div className="flex-shrink-0 mt-1">
                                  <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" />
                                  </svg>
                                </div>
                                <div className="flex-1 min-w-0">
                                  <PullRequestItem pr={pr} showRepository={true} showAuthor={true} />
                                  <div className="flex items-center space-x-3 text-xs text-muted-foreground mt-2">
                                    <span>
                                      {new Date(pr.created_at).toLocaleDateString('en-US', {
                                        year: 'numeric',
                                        month: 'short',
                                        day: 'numeric'
                                      })}
                                    </span>
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
            )}

            {activeTab === 'code-contributions' && (
              <div className="space-y-4">
                {/* Date Filter Controls */}
                <div className="bg-card rounded-lg p-4">
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-lg font-semibold text-foreground">Filter by Date Range</h3>
                    {(metricsLoading || allPRsLoading) && (
                      <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                        <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        <span>Loading...</span>
                      </div>
                    )}
                  </div>
                  <div className="flex flex-wrap gap-4">
                    <div>
                      <label className="block text-sm font-medium text-foreground mb-1">
                        Start Date
                      </label>
                      <input
                        type="date"
                        value={dateParams.startDate || ''}
                        onChange={(e) => handleDateChange('startDate', e.target.value)}
                        className="px-3 py-2 border border-border/50 rounded focus:outline-none focus:ring-1 focus:ring-primary text-sm"
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
                        className="px-3 py-2 border border-border/50 rounded focus:outline-none focus:ring-1 focus:ring-primary text-sm"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-foreground mb-1">
                        Interval
                      </label>
                      <select
                        value={dateParams.interval}
                        onChange={(e) => handleIntervalChange(e.target.value as 'daily' | 'weekly' | 'monthly')}
                        className="px-3 py-2 border border-border/50 rounded focus:outline-none focus:ring-1 focus:ring-primary text-sm"
                      >
                        <option value="daily">Daily</option>
                        <option value="weekly">Weekly</option>
                        <option value="monthly">Monthly</option>
                      </select>
                    </div>
                  </div>
                </div>

                {/* Metrics Graphs */}
                {metrics && (() => {
                  const allGraphs = getAllGraphs(metrics)
                  
                  if (allGraphs.length === 0) {
                    return (
                      <div className="bg-card rounded-lg p-4">
                        <div className="text-center py-8">
                          <p className="text-muted-foreground text-sm">No metrics available</p>
                        </div>
                      </div>
                    )
                  }

                  const currentGraph = allGraphs[currentGraphIndex] || allGraphs[0]
                  
                  return (
                    <div className="bg-card rounded-lg p-4">
                      <h3 className="text-lg font-semibold mb-4">Metrics</h3>
                      <div className="relative">
                        {/* Navigation Buttons */}
                        <div className="flex items-center justify-between mb-3">
                          <button
                            onClick={() => navigateGraph('prev', allGraphs.length)}
                            className="p-2 rounded hover:bg-muted/50 transition-colors"
                            aria-label="Previous graph"
                          >
                            <svg
                              className="w-5 h-5 text-muted-foreground"
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
                            <span className="text-xs text-muted-foreground">
                              {currentGraphIndex + 1} of {allGraphs.length}
                            </span>
                          </div>
                          
                          <button
                            onClick={() => navigateGraph('next', allGraphs.length)}
                            className="p-2 rounded hover:bg-muted/50 transition-colors"
                            aria-label="Next graph"
                          >
                            <svg
                              className="w-5 h-5 text-muted-foreground"
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

                        {/* Current Graph */}
                        <div className="bg-muted/30 rounded-lg p-3">
                          <div className="flex items-center gap-2 mb-2">
                            <div className="flex-1">
                              <p className="text-xs text-muted-foreground mb-1">{currentGraph.category}</p>
                              <div className="flex items-center gap-2">
                                <h5 className="text-sm font-medium">{currentGraph.metric.label}</h5>
                                {currentGraph.description && (
                                  <MetricInfoButton description={currentGraph.description} />
                                )}
                              </div>
                            </div>
                          </div>
                          <MetricChart metric={currentGraph.metric} height={300} />
                        </div>
                      </div>
                    </div>
                  )
                })()}

                {/* Pull Requests Table */}
                <div className="bg-card rounded-lg">
                  <div className="p-4 border-b border-border/50">
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="text-lg font-semibold text-foreground">Pull Requests</h3>
                        <p className="text-sm text-muted-foreground mt-1">
                          Pull requests created in the selected date range
                          {!allPRsLoading && (
                            <span className="ml-2">
                              ({allPRs.length} {allPRs.length === 1 ? 'PR' : 'PRs'})
                            </span>
                          )}
                        </p>
                      </div>
                      <div className="flex items-center space-x-3">
                        {allPRsLoading && (
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
                          <span>{expandedTables.has('pull-requests') ? 'Collapse' : 'Expand'}</span>
                        </button>
                      </div>
                    </div>
                  </div>
                  {expandedTables.has('pull-requests') && (
                    <div className="divide-y divide-border/50">
                      {allPRsLoading ? (
                        <div className="p-4 text-center">
                          <div className="text-muted-foreground text-sm">Loading...</div>
                        </div>
                      ) : allPRs.length === 0 ? (
                        <div className="p-4 text-center">
                          <div className="text-muted-foreground text-sm">
                            No pull requests found for the selected date range.
                          </div>
                        </div>
                      ) : (
                        <div className="max-h-80 overflow-y-auto">
                          {allPRs.map((pr) => (
                            <div key={pr.id} className="p-4">
                              <div className="flex items-start space-x-3">
                                <div className="flex-shrink-0 mt-1">
                                  <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" />
                                  </svg>
                                </div>
                                <div className="flex-1 min-w-0">
                                  <PullRequestItem pr={pr} showRepository={true} showAuthor={true} />
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
            )}

            {activeTab === 'ai-code-assistant' && (
              <div className="space-y-4">
                {/* Date Filter Controls */}
                <div className="bg-card rounded-lg p-4">
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-lg font-semibold text-foreground">Filter by Date Range</h3>
                    {aiMetricsLoading && (
                      <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                        <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        <span>Loading...</span>
                      </div>
                    )}
                  </div>
                  <div className="flex flex-wrap gap-4">
                    <div>
                      <label className="block text-sm font-medium text-foreground mb-1">
                        Start Date
                      </label>
                      <input
                        type="date"
                        value={dateParams.startDate || ''}
                        onChange={(e) => handleDateChange('startDate', e.target.value)}
                        className="px-3 py-2 border border-border/50 rounded focus:outline-none focus:ring-1 focus:ring-primary text-sm"
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
                        className="px-3 py-2 border border-border/50 rounded focus:outline-none focus:ring-1 focus:ring-primary text-sm"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-foreground mb-1">
                        Interval
                      </label>
                      <select
                        value={dateParams.interval}
                        onChange={(e) => handleIntervalChange(e.target.value as 'daily' | 'weekly' | 'monthly')}
                        className="px-3 py-2 border border-border/50 rounded focus:outline-none focus:ring-1 focus:ring-primary text-sm"
                      >
                        <option value="daily">Daily</option>
                        <option value="weekly">Weekly</option>
                        <option value="monthly">Monthly</option>
                      </select>
                    </div>
                  </div>
                </div>

                {/* AI Code Assistant Metrics Graphs */}
                {aiMetricsLoading ? (
                  <div className="bg-card rounded-lg p-4">
                    <div className="flex justify-center items-center py-12">
                      <span className="text-muted-foreground">Loading AI metrics...</span>
                    </div>
                  </div>
                ) : aiMetrics ? (() => {
                  const allGraphs = getAllAiGraphs(aiMetrics)
                  
                  if (allGraphs.length === 0) {
                    return (
                      <div className="bg-card rounded-lg p-4">
                        <div className="text-center py-8">
                          <p className="text-muted-foreground text-sm">No AI metrics available</p>
                        </div>
                      </div>
                    )
                  }

                  const currentGraph = allGraphs[aiCurrentGraphIndex] || allGraphs[0]
                  
                  return (
                    <div className="bg-card rounded-lg p-4">
                      <h3 className="text-lg font-semibold mb-4">AI Code Assistant Metrics</h3>
                      <div className="relative">
                        {/* Navigation Buttons */}
                        <div className="flex items-center justify-between mb-3">
                          <button
                            onClick={() => navigateAiGraph('prev', allGraphs.length)}
                            className="p-2 rounded hover:bg-muted/50 transition-colors"
                            aria-label="Previous graph"
                          >
                            <svg
                              className="w-5 h-5 text-muted-foreground"
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
                            <span className="text-xs text-muted-foreground">
                              {aiCurrentGraphIndex + 1} of {allGraphs.length}
                            </span>
                          </div>
                          
                          <button
                            onClick={() => navigateAiGraph('next', allGraphs.length)}
                            className="p-2 rounded hover:bg-muted/50 transition-colors"
                            aria-label="Next graph"
                          >
                            <svg
                              className="w-5 h-5 text-muted-foreground"
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

                        {/* Current Graph */}
                        <div className="bg-muted/30 rounded-lg p-3">
                          <div className="flex items-center gap-2 mb-2">
                            <div className="flex-1">
                              <p className="text-xs text-muted-foreground mb-1">{currentGraph.category}</p>
                              <div className="flex items-center gap-2">
                                <h5 className="text-sm font-medium">{currentGraph.metric.label}</h5>
                              </div>
                            </div>
                          </div>
                          <MetricChart metric={currentGraph.metric} height={300} />
                        </div>
                      </div>
                    </div>
                  )
                })() : (
                  <div className="bg-card rounded-lg p-4">
                    <div className="text-center py-8">
                      <p className="text-muted-foreground text-sm">No AI metrics available</p>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

