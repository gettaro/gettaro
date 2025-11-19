import { useState, useEffect, useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import { useOrganizationStore } from '../stores/organization'
import Api from '../api/api'
import { Team, TeamType } from '../types/team'
import { Member } from '../types/member'
import { Title } from '../types/title'
import { OrganizationMetricsResponse } from '../types/organizationMetrics'
import MetricChart from '../components/MetricChart'
import MetricInfoButton from '../components/MetricInfoButton'
import { DateInput } from '../components/ui/date-input'

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
  const [titles, setTitles] = useState<Title[]>([])
  const [teamMetrics, setTeamMetrics] = useState<Map<string, OrganizationMetricsResponse>>(new Map())
  const [teamMetricsLoading, setTeamMetricsLoading] = useState<Set<string>>(new Set())
  const [teamGraphIndices, setTeamGraphIndices] = useState<Map<string, number>>(new Map())
  const [expandedMembers, setExpandedMembers] = useState<Set<string>>(new Set())
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

  // Load titles
  useEffect(() => {
    const loadTitles = async () => {
      if (!currentOrganization?.id) return
      try {
        const titlesData = await Api.getOrganizationTitles(currentOrganization.id)
        setTitles(titlesData)
      } catch (err) {
        console.error('Error loading titles:', err)
      }
    }
    loadTitles()
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

  const toggleMembersExpansion = (teamId: string) => {
    setExpandedMembers(prev => {
      const newSet = new Set(prev)
      if (newSet.has(teamId)) {
        newSet.delete(teamId)
      } else {
        newSet.add(teamId)
      }
      return newSet
    })
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
                  <DateInput
                    value={dateParams.startDate || ''}
                    onChange={(e) => handleDateChange('startDate', e.target.value)}
                    className="px-3 py-2 border border-border/50 rounded focus:outline-none focus:ring-1 focus:ring-primary text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    End Date
                  </label>
                  <DateInput
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
                      className="text-left w-full flex items-center justify-between group"
                    >
                      <div className="flex-1">
                        <h3 className="text-xl font-semibold mb-2 group-hover:text-primary transition-colors">{team.name}</h3>
                        {team.description && (
                          <p className="text-sm text-muted-foreground">{team.description}</p>
                        )}
                      </div>
                      <svg
                        className="w-5 h-5 text-muted-foreground group-hover:text-primary transition-colors flex-shrink-0 ml-4"
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

                  {/* Team Metrics */}
                  <div className="border-t border-border pt-4 mb-4">
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

                  {/* Team Members - Expandable */}
                  <div className="border-t border-border pt-4">
                    <button
                      onClick={() => toggleMembersExpansion(team.id)}
                      className="w-full flex items-center justify-between text-left"
                    >
                      <h4 className="text-sm font-medium text-muted-foreground">
                        Members ({team.members.length})
                      </h4>
                      <svg
                        className={`w-4 h-4 text-muted-foreground transition-transform ${
                          expandedMembers.has(team.id) ? 'rotate-180' : ''
                        }`}
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M19 9l-7 7-7-7"
                        />
                      </svg>
                    </button>
                    {expandedMembers.has(team.id) && (
                      <div className="mt-4">
                        {team.members.length === 0 ? (
                          <p className="text-xs text-muted-foreground">No members</p>
                        ) : (
                          <div className="space-y-2">
                            {team.members.map((teamMember) => {
                              const member = members.find(m => m.id === teamMember.member_id)
                              if (!member) return null
                              
                              const memberTitle = member.title_id 
                                ? titles.find(t => t.id === member.title_id)?.name 
                                : null
                              
                              return (
                                <button
                                  key={teamMember.id}
                                  onClick={() => navigate(`/members/${member.id}/profile`)}
                                  className="w-full flex items-center space-x-3 p-2 bg-muted/20 rounded border border-border/30 hover:bg-muted/30 transition-colors text-left"
                                >
                                  {/* Avatar */}
                                  <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center flex-shrink-0">
                                    <span className="text-primary font-medium text-sm">
                                      {member.username.charAt(0).toUpperCase()}
                                    </span>
                                  </div>
                                  
                                  {/* Member Info */}
                                  <div className="flex-1 min-w-0">
                                    <h4 className="font-medium text-foreground text-sm">{member.username}</h4>
                                    {memberTitle && (
                                      <p className="text-xs text-muted-foreground">{memberTitle}</p>
                                    )}
                                  </div>
                                </button>
                              )
                            })}
                          </div>
                        )}
                      </div>
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

