import { useState, useEffect, useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import { useOrganizationStore } from '../stores/organization'
import Api from '../api/api'
import { Team, TeamType } from '../types/team'
import { Member } from '../types/member'
import { OrganizationMetricsResponse } from '../types/organizationMetrics'
import MetricChart from '../components/MetricChart'
import MetricInfoButton from '../components/MetricInfoButton'

export default function Teams() {
  const navigate = useNavigate()
  const { currentOrganization } = useOrganizationStore()
  
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

  // Team type filter
  const [selectedTeamType, setSelectedTeamType] = useState<TeamType | 'all'>('all')

  // Data state
  const [teams, setTeams] = useState<Team[]>([])
  const [members, setMembers] = useState<Member[]>([])
  const [teamMetrics, setTeamMetrics] = useState<Map<string, OrganizationMetricsResponse>>(new Map())
  const [teamMetricsLoading, setTeamMetricsLoading] = useState<Set<string>>(new Set())
  const [teamGraphIndices, setTeamGraphIndices] = useState<Map<string, number>>(new Map())
  const [error, setError] = useState<string | null>(null)

  // Load teams
  useEffect(() => {
    const loadTeams = async () => {
      if (!currentOrganization?.id) return
      try {
        const response = await Api.listTeams(currentOrganization.id)
        setTeams(response.teams || [])
      } catch (err) {
        console.error('Error loading teams:', err)
      }
    }
    loadTeams()
  }, [currentOrganization?.id])

  // Load members
  useEffect(() => {
    const loadMembers = async () => {
      if (!currentOrganization?.id) return
      try {
        const membersData = await Api.getOrganizationMembers(currentOrganization.id)
        setMembers(membersData)
      } catch (err) {
        console.error('Error loading members:', err)
      }
    }
    loadMembers()
  }, [currentOrganization?.id])

  // Filter teams by type
  const filteredTeams = useMemo(() => {
    if (selectedTeamType === 'all') {
      return teams
    }
    const filtered = teams.filter(team => {
      return team.type === selectedTeamType
    })
    return filtered
  }, [teams, selectedTeamType])

  // Load metrics for all teams
  useEffect(() => {
    if (currentOrganization?.id && filteredTeams.length > 0) {
      const loadAllTeamMetrics = async () => {
        setTeamMetricsLoading(new Set(filteredTeams.map(t => t.id)))
        
        try {
          const allTeamIds = filteredTeams.map(team => team.id)
          const metricsData = await Api.getOrganizationMetrics(currentOrganization.id, {
            startDate: dateParams.startDate,
            endDate: dateParams.endDate,
            interval: dateParams.interval,
            teamIds: allTeamIds
          })
          
          const newTeamMetrics = new Map<string, OrganizationMetricsResponse>()
          if (metricsData.teams_breakdown && metricsData.teams_breakdown.length > 0) {
            metricsData.teams_breakdown.forEach(teamBreakdown => {
              newTeamMetrics.set(teamBreakdown.team_id, {
                snapshot_metrics: teamBreakdown.snapshot_metrics || [],
                graph_metrics: teamBreakdown.graph_metrics || []
              })
            })
          }
          
          setTeamMetrics(newTeamMetrics)
        } catch (err) {
          console.error('Error loading team metrics:', err)
          setError('Failed to load team metrics')
        } finally {
          setTeamMetricsLoading(new Set())
        }
      }
      
      loadAllTeamMetrics()
    }
  }, [currentOrganization?.id, selectedTeamType, filteredTeams.map(t => t.id).join(','), dateParams.startDate, dateParams.endDate, dateParams.interval])

  const getMemberName = (memberId: string) => {
    const member = members.find(m => m.id === memberId)
    return member ? member.username : 'Unknown Member'
  }

  // Get all graphs for a team
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

  const handleDateChange = (field: 'startDate' | 'endDate', value: string) => {
    setDateParams(prev => ({
      ...prev,
      [field]: value
    }))
  }

  const handleIntervalChange = (value: 'daily' | 'weekly' | 'monthly') => {
    setDateParams(prev => ({
      ...prev,
      interval: value
    }))
  }

  if (!currentOrganization) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-7xl mx-auto">
          <p className="text-muted-foreground">Please select an organization</p>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">Teams</h1>
        </div>

        {/* Date Range Picker for Teams */}
        <div className="bg-card rounded-lg p-4 mb-6 border border-border">
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
        {error && (
          <div className="bg-destructive/10 border border-destructive rounded-lg p-4 mb-6">
            <p className="text-destructive">{error}</p>
          </div>
        )}

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
                        No metrics available for this team
                      </p>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}

