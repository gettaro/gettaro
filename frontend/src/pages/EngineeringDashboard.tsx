import { useState, useEffect } from 'react'
import { useOrganizationStore } from '../stores/organization'
import Api from '../api/api'
import { PullRequest } from '../types/sourcecontrol'
import { OrganizationMetricsResponse } from '../types/organizationMetrics'
import { Team } from '../types/team'
import MetricChart from '../components/MetricChart'
import MetricInfoButton from '../components/MetricInfoButton'
import { formatMetricValue } from '../utils/formatMetrics'

type TabType = 'engineering-productivity' | 'tech-health'

export default function EngineeringDashboard() {
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
                  className="px-3 py-2 border border-border/50 rounded bg-background text-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm"
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
                  className="px-3 py-2 border border-border/50 rounded bg-background text-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-foreground mb-1">
                  Interval
                </label>
                <select
                  value={dateParams.interval}
                  onChange={(e) => handleIntervalChange(e.target.value as 'daily' | 'weekly' | 'monthly')}
                  className="px-3 py-2 border border-border/50 rounded bg-background text-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm"
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
                                <div className="px-3 pb-3 space-y-2 border-t border-border/50 pt-3">
                                  {prs.map((pr) => (
                                    <a
                                      key={pr.id}
                                      href={pr.url}
                                      target="_blank"
                                      rel="noopener noreferrer"
                                      className="block p-2 bg-background rounded hover:bg-muted/50 transition-colors"
                                    >
                                      <div className="flex items-start justify-between gap-2">
                                        <div className="flex-1 min-w-0">
                                          <p className="font-medium text-sm truncate">{pr.title}</p>
                                          <p className="text-xs text-muted-foreground mt-1">
                                            {new Date(pr.created_at).toLocaleDateString()}
                                          </p>
                                        </div>
                                        <svg
                                          className="w-4 h-4 text-muted-foreground flex-shrink-0"
                                          fill="none"
                                          stroke="currentColor"
                                          viewBox="0 0 24 24"
                                        >
                                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                                        </svg>
                                      </div>
                                    </a>
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
      </div>
    </div>
  )
}
