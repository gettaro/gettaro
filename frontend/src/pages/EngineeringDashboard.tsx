import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useOrganizationStore } from '../stores/organization'
import Api from '../api/api'
import { PullRequest } from '../types/sourcecontrol'
import { OrganizationMetricsResponse } from '../types/organizationMetrics'
import { Team } from '../types/team'
import { Member } from '../types/member'
import MetricChart from '../components/MetricChart'
import MetricInfoButton from '../components/MetricInfoButton'
import PullRequestItem from '../components/PullRequestItem'
import { formatMetricValue } from '../utils/formatMetrics'

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

  // Load teams
  useEffect(() => {
    if (!currentOrganization?.id) return
    
    const loadTeams = async () => {
      try {
        const response = await Api.get(`/organizations/${currentOrganization.id}/teams`)
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

  const handleTeamBreakdownToggle = (checked: boolean) => {
    setShowTeamBreakdown(checked)
    if (checked && teams.length > 0) {
      // When enabling team breakdown, select all teams by default
      setSelectedTeams(teams.map(team => team.id))
    }
  }

  // Load metrics for a specific team
  const loadTeamMetrics = async (teamId: string) => {
    if (!currentOrganization?.id) return

    setTeamMetricsLoading(prev => new Set(prev).add(teamId))
    try {
      const metricsData = await Api.getOrganizationMetrics(currentOrganization.id, {
        startDate: dateParams.startDate,
        endDate: dateParams.endDate,
        interval: dateParams.interval,
        teamIds: [teamId]
      })
      setTeamMetrics(prev => {
        const newMap = new Map(prev)
        newMap.set(teamId, metricsData)
        return newMap
      })
    } catch (err) {
      console.error('Error loading team metrics:', err)
    } finally {
      setTeamMetricsLoading(prev => {
        const newSet = new Set(prev)
        newSet.delete(teamId)
        return newSet
      })
    }
  }

  // Load metrics for all teams when Teams tab is active
  useEffect(() => {
    if (activeTab === 'teams' && currentOrganization?.id && teams.length > 0) {
      teams.forEach(team => {
        loadTeamMetrics(team.id)
      })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab, currentOrganization?.id, teams.map(t => t.id).join(','), dateParams.startDate, dateParams.endDate, dateParams.interval])

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
              <h3 className="text-sm font-medium mb-2">Select Teams:</h3>
              <div className="flex flex-wrap gap-2">
                {teams.map(team => (
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
            </div>
          )}

          {metricsLoading ? (
            <div className="flex justify-center items-center py-12">
              <span>Loading metrics...</span>
            </div>
          ) : metrics ? (
            <div className="space-y-8">
              {/* Graph Metrics - Single graph per metric with team lines when breakdown enabled */}
              {metrics.graph_metrics && metrics.graph_metrics.length > 0 && (
                <div>
                  <h3 className="text-lg font-semibold mb-4">Trends</h3>
                  <div className="space-y-6">
                    {metrics.graph_metrics.map((category) => {
                      // Filter metrics that have data
                      const metricsWithData = category.metrics.filter((metric) => {
                        // Check if metric has time series data
                        if (!metric.time_series || metric.time_series.length === 0) {
                          return false
                        }
                        
                        // Check if team breakdown should be checked
                        if (showTeamBreakdown && metrics.teams_breakdown && metrics.teams_breakdown.length > 0) {
                          const teamMetricsData = metrics.teams_breakdown
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
                          
                          // Check if any team has data (0 is a valid value)
                          return teamMetricsData.some(team => 
                            team.metric.time_series && 
                            team.metric.time_series.length > 0 &&
                            team.metric.time_series.some(entry => 
                              entry.data && entry.data.length > 0
                            )
                          )
                        }
                        
                        // Check cumulative metric has data (0 is a valid value)
                        return metric.time_series.some(entry => 
                          entry.data && entry.data.length > 0
                        )
                      })
                      
                      if (metricsWithData.length === 0) {
                        return null
                      }
                      
                      return (
                        <div key={category.category.name} className="space-y-4">
                          <h4 className="font-medium text-muted-foreground">{category.category.name}</h4>
                          {metricsWithData.map((metric) => {
                            // Find description from snapshot metrics if available
                            const snapshotMetric = metrics.snapshot_metrics
                              ?.flatMap(cat => cat.metrics)
                              .find(m => m.label === metric.label)
                            const description = snapshotMetric?.description || ''
                            
                            // If team breakdown is enabled, collect metrics from all teams for this metric
                            let teamMetricsData: { teamName: string; metric: typeof metric }[] | undefined
                            if (showTeamBreakdown && metrics.teams_breakdown && metrics.teams_breakdown.length > 0) {
                              teamMetricsData = metrics.teams_breakdown
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
                            
                            const chartElement = teamMetricsData && teamMetricsData.length > 0 ? (
                              <MetricChart teamMetrics={teamMetricsData} height={300} />
                            ) : (
                              <MetricChart metric={metric} height={300} />
                            )
                            
                            // Only render if chart component returns something (has data)
                            if (!chartElement) {
                              return null
                            }
                            
                            return (
                              <div key={metric.label} className="bg-muted/30 rounded-lg p-4">
                                <div className="flex items-center gap-2 mb-3">
                                  <h5 className="font-medium">{metric.label}</h5>
                                  {description && <MetricInfoButton description={description} />}
                                </div>
                                {chartElement}
                              </div>
                            )
                          })}
                        </div>
                      )
                    }).filter((element) => element !== null)}
                  </div>
                </div>
              )}
            </div>
          ) : (
            <p className="text-muted-foreground py-8 text-center">No metrics available</p>
          )}
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
            {teams.length === 0 ? (
              <div className="bg-card rounded-lg p-6 text-center">
                <p className="text-muted-foreground">No teams found</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {teams.map((team) => {
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
