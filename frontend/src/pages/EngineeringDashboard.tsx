import { useState, useEffect } from 'react'
import { useOrganizationStore } from '../stores/organization'
import Api from '../api/api'
import { PullRequest } from '../types/sourcecontrol'
import { OrganizationMetricsResponse } from '../types/organizationMetrics'
import { Team } from '../types/team'
import MetricChart from '../components/MetricChart'
import MetricIcon from '../components/MetricIcon'
import { formatMetricValue } from '../utils/formatMetrics'

type TabType = 'engineering-productivity'

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
      interval: 'monthly' as 'daily' | 'weekly' | 'monthly'
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
              className="py-4 px-1 border-b-2 border-primary text-primary font-medium"
            >
              Engineering Productivity
            </button>
          </nav>
        </div>

        {/* Date Range Picker */}
        <div className="bg-card rounded-lg p-4 mb-6">
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

        {error && (
          <div className="bg-destructive/10 text-destructive px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        {/* Open PRs Widget */}
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
                <div className="space-y-3">
                  <h3 className="text-sm font-medium text-muted-foreground mb-2">By Repository</h3>
                  {Object.entries(prsByRepository)
                    .sort((a, b) => b[1].length - a[1].length)
                    .map(([repo, prs]) => (
                      <div key={repo} className="flex items-center justify-between p-3 bg-muted/30 rounded">
                        <div className="flex items-center gap-2">
                          <span className="font-medium">{repo}</span>
                        </div>
                        <span className="text-lg font-semibold">{prs.length}</span>
                      </div>
                    ))}
                </div>
              ) : (
                <p className="text-muted-foreground">No open pull requests</p>
              )}
            </>
          )}
        </div>

        {/* Metrics Section */}
        <div className="bg-card rounded-lg p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">Coding Contribution Metrics</h2>
            <div className="flex items-center gap-4">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={showTeamBreakdown}
                  onChange={(e) => setShowTeamBreakdown(e.target.checked)}
                  className="rounded"
                />
                <span className="text-sm">Break down by team</span>
              </label>
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
              {/* Snapshot Metrics */}
              {metrics.snapshot_metrics && metrics.snapshot_metrics.length > 0 && (
                <div>
                  <h3 className="text-lg font-semibold mb-4">Summary</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {metrics.snapshot_metrics.map((category) =>
                      category.metrics.map((metric) => (
                        <div key={metric.label} className="p-4 bg-muted/30 rounded-lg">
                          <div className="flex items-center gap-2 mb-2">
                            <MetricIcon identifier={metric.icon_identifier} color={metric.icon_color} />
                            <span className="font-medium text-sm">{metric.label}</span>
                          </div>
                          <div className="text-2xl font-bold">
                            {formatMetricValue(metric.value, metric.unit)}
                          </div>
                          <p className="text-xs text-muted-foreground mt-1">{metric.description}</p>
                        </div>
                      ))
                    )}
                  </div>
                </div>
              )}

              {/* Graph Metrics */}
              {metrics.graph_metrics && metrics.graph_metrics.length > 0 && (
                <div>
                  <h3 className="text-lg font-semibold mb-4">Trends</h3>
                  <div className="space-y-6">
                    {metrics.graph_metrics.map((category) => (
                      <div key={category.category.name} className="space-y-4">
                        <h4 className="font-medium text-muted-foreground">{category.category.name}</h4>
                        {category.metrics.map((metric) => (
                          <div key={metric.label} className="bg-muted/30 rounded-lg p-4">
                            <h5 className="font-medium mb-3">{metric.label}</h5>
                            <MetricChart metric={metric} height={300} />
                          </div>
                        ))}
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Team Breakdown */}
              {metrics.teams_breakdown && metrics.teams_breakdown.length > 0 && (
                <div>
                  <h3 className="text-lg font-semibold mb-4">Breakdown by Team</h3>
                  <div className="space-y-6">
                    {metrics.teams_breakdown.map((teamMetrics) => (
                      <div key={teamMetrics.team_id} className="bg-muted/30 rounded-lg p-6">
                        <h4 className="text-lg font-semibold mb-4">{teamMetrics.team_name}</h4>
                        
                        {/* Team Snapshot Metrics */}
                        {teamMetrics.snapshot_metrics && teamMetrics.snapshot_metrics.length > 0 && (
                          <div className="mb-6">
                            <h5 className="text-sm font-medium text-muted-foreground mb-3">Summary</h5>
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                              {teamMetrics.snapshot_metrics.map((category) =>
                                category.metrics.map((metric) => (
                                  <div key={metric.label} className="p-3 bg-background rounded-lg">
                                    <div className="flex items-center gap-2 mb-1">
                                      <MetricIcon identifier={metric.icon_identifier} color={metric.icon_color} />
                                      <span className="font-medium text-xs">{metric.label}</span>
                                    </div>
                                    <div className="text-xl font-bold">
                                      {formatMetricValue(metric.value, metric.unit)}
                                    </div>
                                  </div>
                                ))
                              )}
                            </div>
                          </div>
                        )}

                        {/* Team Graph Metrics */}
                        {teamMetrics.graph_metrics && teamMetrics.graph_metrics.length > 0 && (
                          <div>
                            <h5 className="text-sm font-medium text-muted-foreground mb-3">Trends</h5>
                            <div className="space-y-4">
                              {teamMetrics.graph_metrics.map((category) => (
                                <div key={category.category.name} className="space-y-3">
                                  <h6 className="text-xs font-medium text-muted-foreground">{category.category.name}</h6>
                                  {category.metrics.map((metric) => (
                                    <div key={metric.label} className="bg-background rounded-lg p-3">
                                      <h6 className="text-sm font-medium mb-2">{metric.label}</h6>
                                      <MetricChart metric={metric} height={200} />
                                    </div>
                                  ))}
                                </div>
                              ))}
                            </div>
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          ) : (
            <p className="text-muted-foreground py-8 text-center">No metrics available</p>
          )}
        </div>
      </div>
    </div>
  )
}
