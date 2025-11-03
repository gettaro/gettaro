import { useState, useEffect, useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import { useOrganizationStore } from '../stores/organization'
import Api from '../api/api'
import { PullRequest } from '../types/sourcecontrol'
import { OrganizationMetricsResponse } from '../types/organizationMetrics'
import { Team, TeamType } from '../types/team'
import { Member } from '../types/member'
import MetricChart from '../components/MetricChart'
import MetricInfoButton from '../components/MetricInfoButton'
import PullRequestItem from '../components/PullRequestItem'

type TabType = 'engineering-productivity' | 'tech-health' | 'teams'

export default function EngineeringDashboard() {
  const navigate = useNavigate()
  const { currentOrganization } = useOrganizationStore()
  
  // Date filter state
  const [dateParams, setDateParams] = useState(() => {
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 30) // 1 month ago
    
    return {
      startDate: startDate.toISOString().split('T')[0], // YYYY-MM-DD format
      endDate: endDate.toISOString().split('T')[0],
      interval: 'weekly' as 'daily' | 'weekly' | 'monthly'
    }
  })

  // Data state
  const [openPRs, setOpenPRs] = useState<PullRequest[]>([])
  const [openPRsLoading, setOpenPRsLoading] = useState(false)
  const [metrics, setMetrics] = useState<OrganizationMetricsResponse | null>(null)
  const [metricsLoading, setMetricsLoading] = useState(false)
  const [aiMetrics, setAiMetrics] = useState<GetMemberAICodeAssistantMetricsResponse | null>(null)
  const [aiMetricsLoading, setAiMetricsLoading] = useState(false)
  const [aiCurrentGraphIndex, setAiCurrentGraphIndex] = useState(0)
  const [teams, setTeams] = useState<Team[]>([])
  const [selectedTeams, setSelectedTeams] = useState<string[]>([])
  const [showTeamBreakdown, setShowTeamBreakdown] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [expandedRepos, setExpandedRepos] = useState<Set<string>>(new Set())
  const [showRepositories, setShowRepositories] = useState(false)
  const [activeTab, setActiveTab] = useState<TabType>('engineering-productivity')
  const [members, setMembers] = useState<Member[]>([])
  const [teamMetrics, setTeamMetrics] = useState<Map<string, OrganizationMetricsResponse>>(new Map())
  const [teamMetricsLoading, setTeamMetricsLoading] = useState<Set<string>>(new Set())
  const [teamGraphIndices, setTeamGraphIndices] = useState<Map<string, number>>(new Map())
  const [currentGraphIndex, setCurrentGraphIndex] = useState(0)
  const [selectedTeamType, setSelectedTeamType] = useState<TeamType | 'all'>('squad')

  // Load teams
  useEffect(() => {
    if (!currentOrganization?.id) return
    
    const loadTeams = async () => {
      try {
        const response = await Api.listTeams(currentOrganization.id)
        setTeams(response.teams || [])
      } catch (err) {
        console.error('Error loading teams:', err)
        // Teams are optional, don't set error state
      }
    }
    
    loadTeams()
  }, [currentOrganization?.id])

  // Load members
  useEffect(() => {
    if (!currentOrganization?.id) return
    
    const loadMembers = async () => {
      try {
        const membersData = await Api.getOrganizationMembers(currentOrganization.id)
        setMembers(membersData)
      } catch (err) {
        console.error('Error loading members:', err)
      }
    }
    
    loadMembers()
  }, [currentOrganization?.id])

  // Load open PRs
  const loadOpenPRs = async () => {
    if (!currentOrganization?.id) return

    setOpenPRsLoading(true)
    setError(null)

    try {
      const prs = await Api.getOrganizationPullRequests(currentOrganization.id, {
        status: 'open'
      })
      setOpenPRs(prs)
    } catch (err) {
      console.error('Error loading open PRs:', err)
      setError('Failed to load open PRs')
    } finally {
      setOpenPRsLoading(false)
    }
  }

  // Load metrics
  const loadMetrics = async () => {
    if (!currentOrganization?.id) return

    setMetricsLoading(true)
    setError(null)

    try {
      const teamIds = showTeamBreakdown && selectedTeams.length > 0 ? selectedTeams : undefined
      const metricsData = await Api.getOrganizationMetrics(currentOrganization.id, {
        startDate: dateParams.startDate,
        endDate: dateParams.endDate,
        interval: dateParams.interval,
        teamIds: teamIds
      })
      setMetrics(metricsData)
    } catch (err) {
      console.error('Error loading metrics:', err)
      setError('Failed to load metrics')
    } finally {
      setMetricsLoading(false)
    }
  }

  // Load AI code assistant metrics
  const loadAiMetrics = async () => {
    if (!currentOrganization?.id) return

    setAiMetricsLoading(true)
    setError(null)

    try {
      const params: GetMemberAICodeAssistantMetricsParams = {
        startDate: dateParams.startDate,
        endDate: dateParams.endDate,
        interval: dateParams.interval
      }
      const aiMetricsData = await Api.getOrganizationAICodeAssistantMetrics(currentOrganization.id, params)
      setAiMetrics(aiMetricsData)
    } catch (err) {
      console.error('Error loading AI metrics:', err)
      // Don't set error for AI metrics, they're optional
    } finally {
      setAiMetricsLoading(false)
    }
  }

  // Load open PRs when organization changes
  useEffect(() => {
    if (currentOrganization?.id) {
      loadOpenPRs()
    }
  }, [currentOrganization?.id])

  // Load metrics when params change
  useEffect(() => {
    if (currentOrganization?.id) {
      loadMetrics()
      loadAiMetrics()
    }
  }, [currentOrganization?.id, dateParams.startDate, dateParams.endDate, dateParams.interval, showTeamBreakdown, selectedTeams])

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

  const toggleTeamSelection = (teamId: string) => {
    setSelectedTeams(prev => 
      prev.includes(teamId) 
        ? prev.filter(id => id !== teamId)
        : [...prev, teamId]
    )
  }

  // Filter teams by type
  const filteredTeams = useMemo(() => {
    if (selectedTeamType === 'all') {
      return teams
    }
    // Filter teams that have the selected type (exclude teams without a type or with a different type)
    const filtered = teams.filter(team => {
      // Only include teams that have a type set and match the selected type
      return team.type === selectedTeamType
    })
    return filtered
  }, [teams, selectedTeamType])

  const handleTeamBreakdownToggle = (checked: boolean) => {
    setShowTeamBreakdown(checked)
    if (checked && filteredTeams.length > 0) {
      // When enabling team breakdown, select all filtered teams by default
      setSelectedTeams(filteredTeams.map(team => team.id))
    } else if (!checked) {
      // Clear selection when disabling breakdown
      setSelectedTeams([])
    }
  }

  // Update selected teams when team type filter changes
  useEffect(() => {
    if (showTeamBreakdown) {
      // When team type filter changes, update selected teams to only include teams of the selected type
      const newSelectedTeams = selectedTeams.filter(teamId => 
        filteredTeams.some(team => team.id === teamId)
      )
      // Also add all filtered teams if none are selected
      if (newSelectedTeams.length === 0 && filteredTeams.length > 0) {
        setSelectedTeams(filteredTeams.map(team => team.id))
      } else if (newSelectedTeams.length !== selectedTeams.length) {
        setSelectedTeams(newSelectedTeams)
      }
    }
  }, [selectedTeamType, showTeamBreakdown, filteredTeams, selectedTeams])

  // Load metrics for all teams when Teams tab is active
  useEffect(() => {
    if (activeTab === 'teams' && currentOrganization?.id && filteredTeams.length > 0) {
      const loadAllTeamMetrics = async () => {
        // Mark all teams as loading
        setTeamMetricsLoading(new Set(filteredTeams.map(t => t.id)))
        
        try {
          // Make one request with all team IDs
          const allTeamIds = filteredTeams.map(team => team.id)
          const metricsData = await Api.getOrganizationMetrics(currentOrganization.id, {
            startDate: dateParams.startDate,
            endDate: dateParams.endDate,
            interval: dateParams.interval,
            teamIds: allTeamIds
          })
          
          // Extract team-specific metrics from teams_breakdown
          const newTeamMetrics = new Map<string, OrganizationMetricsResponse>()
          if (metricsData.teams_breakdown && metricsData.teams_breakdown.length > 0) {
            metricsData.teams_breakdown.forEach(teamBreakdown => {
              // Create a response with team-specific metrics
              newTeamMetrics.set(teamBreakdown.team_id, {
                snapshot_metrics: teamBreakdown.snapshot_metrics || [],
                graph_metrics: teamBreakdown.graph_metrics || []
              })
            })
          }
          
          setTeamMetrics(newTeamMetrics)
        } catch (err) {
          console.error('Error loading team metrics:', err)
        } finally {
          setTeamMetricsLoading(new Set())
        }
      }
      
      loadAllTeamMetrics()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab, currentOrganization?.id, selectedTeamType, filteredTeams.map(t => t.id).join(','), dateParams.startDate, dateParams.endDate, dateParams.interval])

  const getMemberName = (memberId: string) => {
    const member = members.find(m => m.id === memberId)
    return member ? member.username : 'Unknown Member'
  }

  // Get all graphs for a team (flattened from all categories)
  const getAllTeamGraphs = (teamMetricsData: OrganizationMetricsResponse) => {
    if (!teamMetricsData.graph_metrics || teamMetricsData.graph_metrics.length === 0) {
      return []
    }

    const allGraphs: Array<{ metric: any; category: string; description: string }> = []

    teamMetricsData.graph_metrics.forEach((category) => {
      const metricsWithData = category.metrics.filter((metric) => {
        if (!metric.time_series || metric.time_series.length === 0) {
          return false
        }
        return metric.time_series.some(entry => 
          entry.data && entry.data.length > 0
        )
      })

      metricsWithData.forEach((metric) => {
        const snapshotMetric = teamMetricsData.snapshot_metrics
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

  const getTeamGraphIndex = (teamId: string) => {
    return teamGraphIndices.get(teamId) || 0
  }

  const setTeamGraphIndex = (teamId: string, index: number) => {
    setTeamGraphIndices(prev => {
      const newMap = new Map(prev)
      newMap.set(teamId, index)
      return newMap
    })
  }

  const navigateTeamGraph = (teamId: string, direction: 'prev' | 'next', totalGraphs: number) => {
    const currentIndex = getTeamGraphIndex(teamId)
    let newIndex: number
    
    if (direction === 'prev') {
      newIndex = currentIndex > 0 ? currentIndex - 1 : totalGraphs - 1
    } else {
      newIndex = currentIndex < totalGraphs - 1 ? currentIndex + 1 : 0
    }
    
    setTeamGraphIndex(teamId, newIndex)
  }

  // Get all graphs for main organization metrics (handles team breakdown)
  const getAllMainGraphs = (metricsData: OrganizationMetricsResponse) => {
    if (!metricsData.graph_metrics || metricsData.graph_metrics.length === 0) {
      return []
    }

    const allGraphs: Array<{ 
      metric: any
      category: string
      description: string
      teamMetrics?: { teamName: string; metric: any }[]
    }> = []

    metricsData.graph_metrics.forEach((category) => {
      // Filter metrics that have data
      const metricsWithData = category.metrics.filter((metric) => {
        if (!metric.time_series || metric.time_series.length === 0) {
          return false
        }
        
        // Check if team breakdown should be checked
        if (showTeamBreakdown && metricsData.teams_breakdown && metricsData.teams_breakdown.length > 0) {
          const teamMetricsData = metricsData.teams_breakdown
            .map(teamData => {
              const teamMetric = teamData.graph_metrics
                ?.flatMap(cat => cat.metrics)
                .find(m => m.label === metric.label)
              return teamMetric ? {
                teamName: teamData.team_name,
                metric: teamMetric
              } : null
            })
            .filter((item): item is { teamName: string; metric: typeof metric } => item !== null)
          
          // Check if any selected team has data
          const selectedTeamMetrics = teamMetricsData.filter(team => 
            selectedTeams.includes(teams.find(t => t.name === team.teamName)?.id || '')
          )
          
          return selectedTeamMetrics.some(team => 
            team.metric.time_series && 
            team.metric.time_series.length > 0 &&
            team.metric.time_series.some(entry => 
              entry.data && entry.data.length > 0
            )
          )
        }
        
        // Check cumulative metric has data
        return metric.time_series.some(entry => 
          entry.data && entry.data.length > 0
        )
      })

      metricsWithData.forEach((metric) => {
        const snapshotMetric = metricsData.snapshot_metrics
          ?.flatMap(cat => cat.metrics)
          .find(m => m.label === metric.label)
        const description = snapshotMetric?.description || ''
        
        // If team breakdown is enabled, collect metrics from selected teams
        let teamMetricsData: { teamName: string; metric: typeof metric }[] | undefined
        if (showTeamBreakdown && metricsData.teams_breakdown && metricsData.teams_breakdown.length > 0 && selectedTeams.length > 0) {
          const selectedTeamIds = new Set(selectedTeams)
          teamMetricsData = metricsData.teams_breakdown
            .filter(teamData => selectedTeamIds.has(teamData.team_id))
            .map(teamData => {
              const teamMetric = teamData.graph_metrics
                ?.flatMap(cat => cat.metrics)
                .find(m => m.label === metric.label)
              return teamMetric ? {
                teamName: teamData.team_name,
                metric: teamMetric
              } : null
            })
            .filter((item): item is { teamName: string; metric: typeof metric } => item !== null)
        }
        
        allGraphs.push({
          metric,
          category: category.category.name,
          description,
          teamMetrics: teamMetricsData
        })
      })
    })

    return allGraphs
  }

  const navigateMainGraph = (direction: 'prev' | 'next', totalGraphs: number) => {
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

  // Reset graph index when metrics or team breakdown changes
  useEffect(() => {
    if (metrics) {
      const allGraphs = getAllMainGraphs(metrics)
      if (allGraphs.length > 0 && currentGraphIndex >= allGraphs.length) {
        setCurrentGraphIndex(0)
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [metrics, showTeamBreakdown, selectedTeams])

  // Reset AI graph index when AI metrics change
  useEffect(() => {
    if (aiMetrics) {
      const allGraphs = getAllAiGraphs(aiMetrics)
      if (allGraphs.length > 0 && aiCurrentGraphIndex >= allGraphs.length) {
        setAiCurrentGraphIndex(0)
      }
    }
  }, [aiMetrics, aiCurrentGraphIndex])

  const toggleRepoExpansion = (repoName: string) => {
    setExpandedRepos(prev => {
      const newSet = new Set(prev)
      if (newSet.has(repoName)) {
        newSet.delete(repoName)
      } else {
        newSet.add(repoName)
      }
      return newSet
    })
  }

  // Group open PRs by repository
  const prsByRepository = openPRs.reduce((acc, pr) => {
    // Use repository_name if available, otherwise extract from URL
    const repoName = pr.repository_name || (pr.url ? extractRepoFromUrl(pr.url) : 'Unknown')
    if (!acc[repoName]) {
      acc[repoName] = []
    }
    acc[repoName].push(pr)
    return acc
  }, {} as Record<string, PullRequest[]>)

  function extractRepoFromUrl(url: string): string {
    try {
      const urlObj = new URL(url)
      const parts = urlObj.pathname.split('/').filter(Boolean)
      if (parts.length >= 2) {
        return `${parts[0]}/${parts[1]}`
      }
      return 'Unknown'
    } catch {
      return 'Unknown'
    }
  }

  if (!currentOrganization) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p>Please select an organization</p>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">Engineering Dashboard</h1>

        {/* Tab Navigation */}
        <div className="border-b border-border mb-6">
          <nav className="flex space-x-8">
            <button
              onClick={() => setActiveTab('engineering-productivity')}
              className={`py-4 px-1 border-b-2 transition-colors ${
                activeTab === 'engineering-productivity'
                  ? 'border-primary text-primary font-medium'
                  : 'border-transparent text-muted-foreground hover:text-foreground'
              }`}
            >
              Engineering Productivity
            </button>
            <button
              onClick={() => setActiveTab('tech-health')}
              className={`py-4 px-1 border-b-2 transition-colors ${
                activeTab === 'tech-health'
                  ? 'border-primary text-primary font-medium'
                  : 'border-transparent text-muted-foreground hover:text-foreground'
              }`}
            >
              Tech Health
            </button>
            <button
              onClick={() => setActiveTab('teams')}
              className={`py-4 px-1 border-b-2 transition-colors ${
                activeTab === 'teams'
                  ? 'border-primary text-primary font-medium'
                  : 'border-transparent text-muted-foreground hover:text-foreground'
              }`}
            >
              Teams
            </button>
          </nav>
        </div>

        {error && (
          <div className="bg-destructive/10 text-destructive px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        {/* Engineering Productivity Tab Content */}
        {activeTab === 'engineering-productivity' && (
          <div className="bg-card rounded-lg p-6 mb-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold">Coding Contribution Metrics</h2>
            <div className="flex items-center gap-4">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={showTeamBreakdown}
                  onChange={(e) => handleTeamBreakdownToggle(e.target.checked)}
                  className="rounded"
                />
                <span className="text-sm">Break down by team</span>
              </label>
            </div>
          </div>

          {/* Date Range Picker */}
          <div className="bg-muted/30 rounded-lg p-4 mb-6">
            <div className="flex flex-wrap items-end gap-4">
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

          {showTeamBreakdown && teams.length > 0 && (
            <div className="mb-4 p-4 bg-muted/30 rounded">
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-sm font-medium">Select Teams:</h3>
                <div className="flex items-center gap-2">
                  <label className="text-sm text-muted-foreground">Filter by type:</label>
                  <select
                    value={selectedTeamType}
                    onChange={(e) => setSelectedTeamType(e.target.value as TeamType | 'all')}
                    className="px-3 py-1 border border-border/50 rounded focus:outline-none focus:ring-1 focus:ring-primary text-sm bg-background"
                  >
                    <option value="all">All Types</option>
                    <option value="squad">Squads</option>
                    <option value="chapter">Chapters</option>
                    <option value="tribe">Tribes</option>
                    <option value="guild">Guilds</option>
                  </select>
                </div>
              </div>
              {filteredTeams.length > 0 ? (
                <>
                  <div className="flex flex-wrap gap-2">
                    {filteredTeams.map(team => (
                      <label key={team.id} className="flex items-center gap-2 cursor-pointer px-3 py-1 bg-background rounded border border-border hover:bg-muted">
                        <input
                          type="checkbox"
                          checked={selectedTeams.includes(team.id)}
                          onChange={() => toggleTeamSelection(team.id)}
                          className="rounded"
                        />
                        <span className="text-sm">{team.name}</span>
                      </label>
                    ))}
                  </div>
                  {selectedTeams.length > 0 && (
                    <p className="text-xs text-muted-foreground mt-2">
                      Showing metrics for {selectedTeams.length} team{selectedTeams.length !== 1 ? 's' : ''}
                    </p>
                  )}
                </>
              ) : (
                <p className="text-sm text-muted-foreground">No teams of this type available</p>
              )}
            </div>
          )}

          {metricsLoading ? (
            <div className="flex justify-center items-center py-12">
              <span>Loading metrics...</span>
            </div>
          ) : metrics ? (() => {
            const allGraphs = getAllMainGraphs(metrics)
            
            if (allGraphs.length === 0) {
              return (
                <p className="text-muted-foreground py-8 text-center">No metrics available</p>
              )
            }
            
            const currentGraph = allGraphs[currentGraphIndex] || allGraphs[0]
            const chartElement = currentGraph.teamMetrics && currentGraph.teamMetrics.length > 0 ? (
              <MetricChart teamMetrics={currentGraph.teamMetrics} height={300} />
            ) : (
              <MetricChart metric={currentGraph.metric} height={300} />
            )
            
            if (!chartElement) {
              return (
                <p className="text-muted-foreground py-8 text-center">No metrics available</p>
              )
            }
            
            return (
              <div>
                <h3 className="text-lg font-semibold mb-4">Trends</h3>
                <div className="bg-muted/30 rounded-lg p-4">
                  <div className="relative">
                    {/* Navigation Buttons */}
                    <div className="flex items-center justify-between mb-4">
                      <button
                        onClick={() => navigateMainGraph('prev', allGraphs.length)}
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
                        onClick={() => navigateMainGraph('next', allGraphs.length)}
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
                    
                    {/* Current Graph */}
                    <div className="mb-2">
                      <p className="text-xs text-muted-foreground mb-1">{currentGraph.category}</p>
                      <div className="flex items-center gap-2 mb-3">
                        <h5 className="font-medium">{currentGraph.metric.label}</h5>
                        {currentGraph.description && (
                          <MetricInfoButton description={currentGraph.description} />
                        )}
                      </div>
                      {chartElement}
                    </div>
                  </div>
                </div>
              </div>
            )
          })() : (
            <p className="text-muted-foreground py-8 text-center">No metrics available</p>
          )}

          {/* AI Code Assistant Metrics */}
          <div className="mt-8">
            <h2 className="text-xl font-semibold mb-4">AI Code Assistant Metrics</h2>
            
            {aiMetricsLoading ? (
              <div className="flex justify-center items-center py-12">
                <span>Loading AI metrics...</span>
              </div>
            ) : aiMetrics ? (() => {
              const allGraphs = getAllAiGraphs(aiMetrics)
              
              if (allGraphs.length === 0) {
                return (
                  <p className="text-muted-foreground py-8 text-center">No AI metrics available</p>
                )
              }
              
              const currentGraph = allGraphs[aiCurrentGraphIndex] || allGraphs[0]
              const chartElement = currentGraph.metric ? (
                <MetricChart metric={currentGraph.metric} height={300} />
              ) : null
              
              if (!chartElement) {
                return (
                  <p className="text-muted-foreground py-8 text-center">No AI metrics available</p>
                )
              }
              
              return (
                <div>
                  <div className="bg-muted/30 rounded-lg p-4">
                    <div className="relative">
                      {/* Navigation Buttons */}
                      <div className="flex items-center justify-between mb-4">
                        <button
                          onClick={() => navigateAiGraph('prev', allGraphs.length)}
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
                            {aiCurrentGraphIndex + 1} of {allGraphs.length}
                          </span>
                        </div>
                        
                        <button
                          onClick={() => navigateAiGraph('next', allGraphs.length)}
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
                      
                      {/* Current Graph */}
                      <div className="mb-2">
                        <p className="text-xs text-muted-foreground mb-1">{currentGraph.category}</p>
                        <div className="flex items-center gap-2 mb-3">
                          <h5 className="font-medium">{currentGraph.metric.label}</h5>
                        </div>
                        {chartElement}
                      </div>
                    </div>
                  </div>
                </div>
              )
            })() : (
              <p className="text-muted-foreground py-8 text-center">No AI metrics available</p>
            )}
          </div>
          </div>
        )}

        {/* Tech Health Tab Content */}
        {activeTab === 'tech-health' && (
          <div className="bg-card rounded-lg p-6 mb-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold">Open Pull Requests</h2>
              {openPRsLoading && <span className="text-sm text-muted-foreground">Loading...</span>}
            </div>
            
            {!openPRsLoading && (
              <>
                <div className="mb-4">
                  <span className="text-3xl font-bold">{openPRs.length}</span>
                  <span className="text-muted-foreground ml-2">total open PRs</span>
                </div>
                
                {Object.keys(prsByRepository).length > 0 ? (
                  <div className="space-y-2">
                    <button
                      onClick={() => setShowRepositories(!showRepositories)}
                      className="w-full flex items-center justify-between p-2 hover:bg-muted/30 rounded transition-colors"
                    >
                      <div className="flex items-center gap-2">
                        <svg
                          className={`w-4 h-4 transition-transform ${showRepositories ? 'rotate-90' : ''}`}
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                        <h3 className="text-sm font-medium text-muted-foreground">Repositories</h3>
                        <span className="text-xs text-muted-foreground">
                          ({Object.keys(prsByRepository).length} {Object.keys(prsByRepository).length === 1 ? 'repository' : 'repositories'})
                        </span>
                      </div>
                    </button>
                    {showRepositories && (
                      <div className="space-y-2 pl-6">
                        {Object.entries(prsByRepository)
                          .sort((a, b) => b[1].length - a[1].length)
                          .map(([repo, prs]) => (
                            <div key={repo} className="bg-muted/30 rounded">
                              <button
                                onClick={() => toggleRepoExpansion(repo)}
                                className="w-full flex items-center justify-between p-3 hover:bg-muted/50 transition-colors"
                              >
                                <div className="flex items-center gap-2">
                                  <svg
                                    className={`w-4 h-4 transition-transform ${expandedRepos.has(repo) ? 'rotate-90' : ''}`}
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                  >
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                  </svg>
                                  <span className="font-medium">{repo}</span>
                                </div>
                                <span className="text-lg font-semibold">{prs.length}</span>
                              </button>
                              {expandedRepos.has(repo) && (
                                <div className="px-3 pb-3 space-y-4 border-t border-border/50 pt-3">
                                  {prs.map((pr) => (
                                    <div key={pr.id} className="bg-background rounded p-3">
                                      <PullRequestItem pr={pr} showRepository={false} showAuthor={true} />
                                      <div className="flex items-center gap-3 text-xs text-muted-foreground mt-2 pt-2 border-t border-border/30">
                                        <span>{new Date(pr.created_at).toLocaleDateString('en-US', {
                                          year: 'numeric',
                                          month: 'short',
                                          day: 'numeric'
                                        })}</span>
                                      </div>
                                    </div>
                                  ))}
                                </div>
                              )}
                            </div>
                          ))}
                      </div>
                    )}
                  </div>
                ) : (
                  <p className="text-muted-foreground">No open pull requests</p>
                )}
              </>
            )}
          </div>
        )}

        {/* Teams Tab Content */}
        {activeTab === 'teams' && (
          <>
            {/* Date Range Picker for Teams */}
            <div className="bg-muted/30 rounded-lg p-4 mb-6">
              <div className="flex flex-wrap items-end gap-4">
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Team Type
                  </label>
                  <select
                    value={selectedTeamType}
                    onChange={(e) => {
                      const newType = e.target.value as TeamType | 'all'
                      setSelectedTeamType(newType)
                    }}
                    className="px-3 py-2 border border-border/50 rounded focus:outline-none focus:ring-1 focus:ring-primary text-sm"
                  >
                    <option value="all">All Types</option>
                    <option value="squad">Squad</option>
                    <option value="chapter">Chapter</option>
                    <option value="tribe">Tribe</option>
                    <option value="guild">Guild</option>
                  </select>
                </div>
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

            {/* Team Cards */}
            {filteredTeams.length === 0 ? (
              <div className="bg-card rounded-lg p-6 text-center">
                <p className="text-muted-foreground">
                  {teams.length === 0 
                    ? 'No teams found' 
                    : `No teams found for type: ${selectedTeamType === 'all' ? 'All Types' : selectedTeamType}`
                  }
                </p>
              </div>
            ) : (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {filteredTeams.map((team) => {
                  const teamMetricsData = teamMetrics.get(team.id)
                  const isLoading = teamMetricsLoading.has(team.id)
                  
                  return (
                    <div key={team.id} className="bg-card rounded-lg p-6 border border-border hover:border-primary/50 transition-colors">
                      {/* Team Header */}
                      <div className="mb-4">
                        <button
                          onClick={() => navigate(`/teams/${team.id}/profile`)}
                          className="text-left w-full"
                        >
                          <h3 className="text-xl font-semibold mb-2 hover:text-primary transition-colors">{team.name}</h3>
                          {team.description && (
                            <p className="text-sm text-muted-foreground">{team.description}</p>
                          )}
                        </button>
                      </div>

                      {/* Team Members */}
                      <div className="mb-6">
                        <h4 className="text-sm font-medium text-muted-foreground mb-2">
                          Members ({team.members.length})
                        </h4>
                        {team.members.length === 0 ? (
                          <p className="text-xs text-muted-foreground">No members</p>
                        ) : (
                          <div className="flex flex-wrap gap-2">
                            {team.members.map((teamMember) => (
                              <span
                                key={teamMember.id}
                                className="px-2 py-1 text-xs bg-muted/50 rounded border border-border"
                              >
                                {getMemberName(teamMember.member_id)}
                              </span>
                            ))}
                          </div>
                        )}
                      </div>

                      {/* Team Metrics */}
                      <div className="border-t border-border pt-4">
                        <h4 className="text-sm font-medium mb-4">Code Contribution Metrics</h4>
                        {isLoading ? (
                          <div className="flex justify-center items-center py-8">
                            <span className="text-sm text-muted-foreground">Loading metrics...</span>
                          </div>
                        ) : teamMetricsData ? (() => {
                          const allGraphs = getAllTeamGraphs(teamMetricsData)
                          
                          if (allGraphs.length === 0) {
                            return (
                              <p className="text-xs text-muted-foreground text-center py-4">
                                No metrics available for this team
                              </p>
                            )
                          }

                          const currentIndex = getTeamGraphIndex(team.id)
                          const currentGraph = allGraphs[currentIndex]
                          
                          return (
                            <div className="relative">
                              {/* Navigation Buttons */}
                              <div className="flex items-center justify-between mb-3">
                                <button
                                  onClick={() => navigateTeamGraph(team.id, 'prev', allGraphs.length)}
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
                                    {currentIndex + 1} of {allGraphs.length}
                                  </span>
                                </div>
                                
                                <button
                                  onClick={() => navigateTeamGraph(team.id, 'next', allGraphs.length)}
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
                                      <h6 className="text-sm font-medium">{currentGraph.metric.label}</h6>
                                      {currentGraph.description && (
                                        <MetricInfoButton description={currentGraph.description} />
                                      )}
                                    </div>
                                  </div>
                                </div>
                                <MetricChart metric={currentGraph.metric} height={200} />
                              </div>
                            </div>
                          )
                        })() : (
                          <p className="text-xs text-muted-foreground text-center py-4">
                            No metrics available
                          </p>
                        )}
                      </div>
                    </div>
                  )
                })}
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
