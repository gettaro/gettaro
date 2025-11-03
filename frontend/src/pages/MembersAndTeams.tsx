import { useState, useEffect, useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import Api from '../api/api'
import { Team, CreateTeamRequest, UpdateTeamRequest, AddTeamMemberRequest, TeamType } from '../types/team'
import { Member, AddMemberRequest, UpdateMemberRequest } from '../types/member'
import { Title } from '../types/title'
import { ExternalAccount } from '../types/sourcecontrol'
import { useOrganizationStore } from '../stores/organization'
import { useToast } from '../hooks/useToast'

export default function MembersAndTeams() {
  const navigate = useNavigate()
  const { currentOrganization } = useOrganizationStore()
  const { toast } = useToast()
  
  // Teams state
  const [teams, setTeams] = useState<Team[]>([])
  const [isCreateTeamModalOpen, setIsCreateTeamModalOpen] = useState(false)
  const [isEditTeamModalOpen, setIsEditTeamModalOpen] = useState(false)
  const [isAddMemberToTeamModalOpen, setIsAddMemberToTeamModalOpen] = useState(false)
  const [isCreatingTeam, setIsCreatingTeam] = useState(false)
  const [isUpdatingTeam, setIsUpdatingTeam] = useState(false)
  const [isAddingMemberToTeam, setIsAddingMemberToTeam] = useState(false)
  const [selectedTeam, setSelectedTeam] = useState<Team | null>(null)

  // Members state
  const [members, setMembers] = useState<Member[]>([])
  const [titles, setTitles] = useState<Title[]>([])
  const [sourceControlAccounts, setSourceControlAccounts] = useState<ExternalAccount[]>([])
  const [isAddMemberModalOpen, setIsAddMemberModalOpen] = useState(false)
  const [isAddingMember, setIsAddingMember] = useState(false)
  const [isUpdateMemberModalOpen, setIsUpdateMemberModalOpen] = useState(false)
  const [isUpdatingMember, setIsUpdatingMember] = useState(false)
  const [selectedMember, setSelectedMember] = useState<Member | null>(null)

  // General state
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<'members' | 'teams'>('members')
  const [activeTeamTypeTab, setActiveTeamTypeTab] = useState<TeamType | 'all' | null>(null)

  // Form data
  const [createTeamFormData, setCreateTeamFormData] = useState<CreateTeamRequest>({
    name: '',
    description: '',
    type: undefined,
    pr_prefix: undefined,
    organization_id: ''
  })

  const [updateTeamFormData, setUpdateTeamFormData] = useState<UpdateTeamRequest>({
    name: '',
    description: '',
    type: undefined,
    pr_prefix: undefined
  })

  // Group teams by type
  const teamsByType = useMemo(() => {
    const grouped: Record<string, Team[]> = {
      all: teams,
      squad: [],
      chapter: [],
      tribe: [],
      guild: []
    }

    teams.forEach(team => {
      const type = team.type || 'squad' // Default to squad if null
      if (grouped[type]) {
        grouped[type].push(team)
      }
    })

    return grouped
  }, [teams])

  // Get available team type tabs (only show tabs that have teams)
  const availableTeamTypes = useMemo(() => {
    const types: Array<{ type: TeamType | 'all' | null; label: string; count: number }> = [
      { type: 'all', label: 'All Teams', count: teams.length }
    ]

    const typeLabels: Record<TeamType, string> = {
      squad: 'Squads',
      chapter: 'Chapters',
      tribe: 'Tribes',
      guild: 'Guilds'
    }

    ;(['squad', 'chapter', 'tribe', 'guild'] as TeamType[]).forEach(type => {
      const count = teamsByType[type].length
      if (count > 0) {
        types.push({ type, label: typeLabels[type], count })
      }
    })

    return types
  }, [teams, teamsByType])

  // Set default active team type tab
  useEffect(() => {
    if (activeTab === 'teams' && activeTeamTypeTab === null && availableTeamTypes.length > 0) {
      setActiveTeamTypeTab(availableTeamTypes[0].type)
    }
  }, [activeTab, activeTeamTypeTab, availableTeamTypes])

  // Get filtered teams based on active tab
  const filteredTeams = useMemo(() => {
    if (activeTeamTypeTab === 'all' || activeTeamTypeTab === null) {
      return teams
    }
    return teamsByType[activeTeamTypeTab] || []
  }, [activeTeamTypeTab, teams, teamsByType])

  const [addMemberToTeamFormData, setAddMemberToTeamFormData] = useState<AddTeamMemberRequest>({
    member_id: ''
  })

  const [addMemberFormData, setAddMemberFormData] = useState<AddMemberRequest>({
    email: '',
    username: '',
    title_id: '',
    external_account_id: '',
    manager_id: undefined
  })

  const [updateMemberFormData, setUpdateMemberFormData] = useState<UpdateMemberRequest>({
    username: '',
    title_id: '',
    external_account_id: '',
    manager_id: undefined
  })

  useEffect(() => {
    if (currentOrganization) {
      loadTeams()
      loadMembers()
      loadTitles()
      loadSourceControlAccounts()
    }
  }, [currentOrganization])


  // Teams functions
  const loadTeams = async () => {
    if (!currentOrganization) return

    try {
      setLoading(true)
      const response = await Api.listTeams(currentOrganization.id)
      setTeams(response.teams)
    } catch (err) {
      console.error('Error loading teams:', err)
      setError('Failed to load teams')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateTeam = async () => {
    if (!currentOrganization) return

    try {
      setIsCreatingTeam(true)
      const teamData = { ...createTeamFormData, organization_id: currentOrganization.id }
      await Api.createTeam(currentOrganization.id, teamData)
      
      toast({
        title: 'Success',
        description: 'Team created successfully',
      })
      
      setIsCreateTeamModalOpen(false)
      setCreateTeamFormData({ name: '', description: '', type: undefined, pr_prefix: undefined, organization_id: '' })
      loadTeams()
    } catch (err) {
      console.error('Error creating team:', err)
      toast({
        title: 'Error',
        description: 'Failed to create team',
        variant: 'destructive',
      })
    } finally {
      setIsCreatingTeam(false)
    }
  }

  const handleUpdateTeam = async () => {
    if (!selectedTeam) return

    try {
      setIsUpdatingTeam(true)
      await Api.updateTeam(currentOrganization.id, selectedTeam.id, updateTeamFormData)
      
      toast({
        title: 'Success',
        description: 'Team updated successfully',
      })
      
      setIsEditTeamModalOpen(false)
      setSelectedTeam(null)
      setUpdateTeamFormData({ name: '', description: '', type: undefined, pr_prefix: undefined })
      loadTeams()
    } catch (err) {
      console.error('Error updating team:', err)
      toast({
        title: 'Error',
        description: 'Failed to update team',
        variant: 'destructive',
      })
    } finally {
      setIsUpdatingTeam(false)
    }
  }

  const handleDeleteTeam = async (team: Team) => {
    if (!confirm(`Are you sure you want to delete the team "${team.name}"?`)) {
      return
    }

    try {
      await Api.deleteTeam(currentOrganization.id, team.id)
      
      toast({
        title: 'Success',
        description: 'Team deleted successfully',
      })
      
      loadTeams()
    } catch (err) {
      console.error('Error deleting team:', err)
      toast({
        title: 'Error',
        description: 'Failed to delete team',
        variant: 'destructive',
      })
    }
  }

  const handleAddMemberToTeam = async () => {
    if (!selectedTeam) return

    try {
      setIsAddingMemberToTeam(true)
      await Api.addTeamMember(currentOrganization.id, selectedTeam.id, addMemberToTeamFormData)
      
      toast({
        title: 'Success',
        description: 'Member added to team successfully',
      })
      
      setIsAddMemberToTeamModalOpen(false)
      setSelectedTeam(null)
      setAddMemberToTeamFormData({ member_id: '' })
      loadTeams()
    } catch (err) {
      console.error('Error adding member to team:', err)
      toast({
        title: 'Error',
        description: 'Failed to add member to team',
        variant: 'destructive',
      })
    } finally {
      setIsAddingMemberToTeam(false)
    }
  }

  const handleRemoveMemberFromTeam = async (teamId: string, memberId: string) => {
    if (!confirm('Are you sure you want to remove this member from the team?')) {
      return
    }

    try {
      await Api.removeTeamMember(currentOrganization.id, teamId, memberId)
      
      toast({
        title: 'Success',
        description: 'Member removed from team successfully',
      })
      
      loadTeams()
    } catch (err) {
      console.error('Error removing member from team:', err)
      toast({
        title: 'Error',
        description: 'Failed to remove member from team',
        variant: 'destructive',
      })
    }
  }

  const openEditTeamModal = (team: Team) => {
    setSelectedTeam(team)
    setUpdateTeamFormData({ name: team.name, description: team.description, type: team.type, pr_prefix: team.pr_prefix })
    setIsEditTeamModalOpen(true)
  }

  const openAddMemberToTeamModal = (team: Team) => {
    setSelectedTeam(team)
    setAddMemberToTeamFormData({ member_id: '' })
    setIsAddMemberToTeamModalOpen(true)
  }

  const getMemberName = (memberId: string) => {
    const member = members.find(m => m.id === memberId)
    return member ? `${member.username} (${member.email})` : 'Unknown Member'
  }

  // Members functions
  const loadMembers = async () => {
    if (!currentOrganization) {
      setError('No organization selected')
      setLoading(false)
      return
    }

    try {
      setLoading(true)
      setError(null)
      const membersData = await Api.getOrganizationMembers(currentOrganization.id)
      setMembers(membersData)
    } catch (err) {
      setError('Failed to load members')
      console.error('Error loading members:', err)
    } finally {
      setLoading(false)
    }
  }

  const loadTitles = async () => {
    if (!currentOrganization) return

    try {
      const titlesData = await Api.getOrganizationTitles(currentOrganization.id)
      setTitles(titlesData)
    } catch (err) {
      console.error('Error loading titles:', err)
    }
  }

  const loadSourceControlAccounts = async () => {
    if (!currentOrganization) return

    try {
      const accountsData = await Api.getOrganizationSourceControlAccounts(currentOrganization.id)
      setSourceControlAccounts(accountsData)
    } catch (err) {
      console.error('Error loading source control accounts:', err)
    }
  }

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!currentOrganization) return

    try {
      setIsAddingMember(true)
      setError(null)
      const createdMember = await Api.addOrganizationMember(currentOrganization.id, addMemberFormData)
      setAddMemberFormData({
        email: '',
        username: '',
        title_id: '',
        external_account_id: '',
        manager_id: undefined
      })
      setIsAddMemberModalOpen(false)
      await loadMembers()
      
      // Navigate to the member profile page
      if (createdMember?.id) {
        navigate(`/members/${createdMember.id}/profile`)
      }
    } catch (err) {
      setError('Failed to add member')
      console.error('Error adding member:', err)
    } finally {
      setIsAddingMember(false)
    }
  }

  const handleUpdateMember = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!currentOrganization || !selectedMember) return

    try {
      setIsUpdatingMember(true)
      setError(null)
      await Api.updateOrganizationMember(currentOrganization.id, selectedMember.id, updateMemberFormData)
      setUpdateMemberFormData({
        username: '',
        title_id: '',
        external_account_id: '',
        manager_id: undefined
      })
      setSelectedMember(null)
      setIsUpdateMemberModalOpen(false)
      await loadMembers()
    } catch (err) {
      setError('Failed to update member')
      console.error('Error updating member:', err)
    } finally {
      setIsUpdatingMember(false)
    }
  }

  const handleEditMember = (member: Member) => {
    setSelectedMember(member)
    setUpdateMemberFormData({
      username: member.username,
      title_id: member.title_id || '',
      external_account_id: '',
      manager_id: member.manager_id
    })
    setIsUpdateMemberModalOpen(true)
  }

  const handleDeleteMember = async (member: Member) => {
    if (!confirm(`Are you sure you want to delete ${member.username}?`)) {
      return
    }

    try {
      await Api.removeOrganizationMember(currentOrganization.id, member.id)
      await loadMembers()
    } catch (err) {
      console.error('Error deleting member:', err)
    }
  }

  const getTitleName = (titleId: string | null) => {
    if (!titleId) return 'No Title'
    const title = titles.find(t => t.id === titleId)
    return title ? title.name : 'Unknown Title'
  }

  const getManagerName = (managerId: string | null) => {
    if (!managerId) return 'No Manager'
    const manager = members.find(m => m.id === managerId)
    return manager ? manager.username : 'Unknown Manager'
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
              <p className="text-muted-foreground">Loading...</p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">Members & Teams</h1>
          <p className="text-muted-foreground">
            Manage team members and teams within your organization.
          </p>
        </div>

        {/* Tab Navigation */}
        <div className="mb-6">
          <nav className="flex space-x-8 border-b border-border">
            <button
              onClick={() => setActiveTab('members')}
              className={`py-2 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === 'members'
                  ? 'border-primary text-primary'
                  : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              Members
            </button>
            <button
              onClick={() => setActiveTab('teams')}
              className={`py-2 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === 'teams'
                  ? 'border-primary text-primary'
                  : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              Teams
            </button>
          </nav>
        </div>

        {/* Tab Content */}
        {activeTab === 'members' ? (
          <div>
            {/* Members Header */}
            <div className="flex justify-between items-center mb-6">
              <div>
                <h2 className="text-2xl font-semibold text-foreground">Team Members</h2>
                <p className="text-muted-foreground">Manage your organization's team members</p>
              </div>
              <button
                onClick={() => setIsAddMemberModalOpen(true)}
                className="bg-primary text-primary-foreground hover:bg-primary/90 px-4 py-2 rounded-md transition-colors"
              >
                Add Member
              </button>
            </div>

            {/* Error Message */}
            {error && (
              <div className="mb-4 p-4 bg-destructive/10 border border-destructive/20 rounded-md">
                <p className="text-destructive">{error}</p>
              </div>
            )}

            {/* Members List */}
            <div className="bg-card border border-border rounded-lg">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="border-b border-border">
                    <tr>
                      <th className="text-left p-4 font-medium text-foreground">Name</th>
                      <th className="text-left p-4 font-medium text-foreground">Email</th>
                      <th className="text-left p-4 font-medium text-foreground">Title</th>
                      <th className="text-left p-4 font-medium text-foreground">Manager</th>
                      <th className="text-left p-4 font-medium text-foreground">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {members.map((member) => (
                      <tr 
                        key={member.id} 
                        className="border-b border-border last:border-b-0 hover:bg-muted/30 cursor-pointer transition-colors"
                        onClick={() => navigate(`/members/${member.id}/profile`)}
                      >
                        <td className="p-4">
                          <div className="flex items-center space-x-3">
                            <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                              <span className="text-sm font-medium text-primary">
                                {member.username.charAt(0).toUpperCase()}
                              </span>
                            </div>
                            <span className="font-medium text-foreground">{member.username}</span>
                          </div>
                        </td>
                        <td className="p-4 text-muted-foreground">{member.email}</td>
                        <td className="p-4 text-muted-foreground">{getTitleName(member.title_id)}</td>
                        <td className="p-4 text-muted-foreground">{getManagerName(member.manager_id)}</td>
                        <td className="p-4">
                          <div className="flex items-center gap-2" onClick={(e) => e.stopPropagation()}>
                            <button
                              onClick={() => handleEditMember(member)}
                              className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-white bg-primary hover:bg-primary/90 rounded-md transition-colors"
                              title="Edit member"
                            >
                              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                              </svg>
                              Edit
                            </button>
                            <button
                              onClick={() => handleDeleteMember(member)}
                              className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-white bg-destructive hover:bg-destructive/90 rounded-md transition-colors"
                              title="Delete member"
                            >
                              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                              </svg>
                              Delete
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        ) : (
          <div>
            {/* Teams Header */}
            <div className="flex justify-between items-center mb-6">
              <div>
                <h2 className="text-2xl font-semibold text-foreground">Teams</h2>
                <p className="text-muted-foreground">Create and manage teams within your organization</p>
              </div>
              <button
                onClick={() => setIsCreateTeamModalOpen(true)}
                className="bg-primary text-primary-foreground hover:bg-primary/90 px-4 py-2 rounded-md transition-colors"
              >
                Create Team
              </button>
            </div>

            {/* Team Type Tabs */}
            {availableTeamTypes.length > 1 && (
              <div className="mb-6">
                <nav className="flex space-x-4 border-b border-border">
                  {availableTeamTypes.map(({ type, label, count }) => (
                    <button
                      key={type || 'all'}
                      onClick={() => setActiveTeamTypeTab(type)}
                      className={`py-2 px-1 border-b-2 font-medium text-sm transition-colors ${
                        activeTeamTypeTab === type
                          ? 'border-primary text-primary'
                          : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border'
                      }`}
                    >
                      {label} ({count})
                    </button>
                  ))}
                </nav>
              </div>
            )}

            {/* Teams List */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredTeams.map((team) => {
                const getTeamTypeIcon = (type?: TeamType) => {
                  switch (type) {
                    case 'squad':
                      return (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                        </svg>
                      )
                    case 'chapter':
                      return (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.582 4.786 8 4.285v13C9.582 17.786 10.832 18.477 12 19m0-13C13.168 5.477 14.418 4.786 16 4.285v13c-1.582.501-2.832 1.192-4 1.585m0 0V19" />
                        </svg>
                      )
                    case 'tribe':
                      return (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      )
                    case 'guild':
                      return (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h1" />
                        </svg>
                      )
                    default:
                      return (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                        </svg>
                      )
                  }
                }

                const getTeamTypeColor = (type?: TeamType) => {
                  switch (type) {
                    case 'squad':
                      return 'bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20'
                    case 'chapter':
                      return 'bg-purple-500/10 text-purple-600 dark:text-purple-400 border-purple-500/20'
                    case 'tribe':
                      return 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20'
                    case 'guild':
                      return 'bg-orange-500/10 text-orange-600 dark:text-orange-400 border-orange-500/20'
                    default:
                      return 'bg-muted text-muted-foreground border-border'
                  }
                }

                return (
                  <div key={team.id} className="bg-card border border-border rounded-lg overflow-hidden hover:border-primary/50 transition-all duration-200 hover:shadow-lg">
                    {/* Header */}
                    <div className="p-5 bg-gradient-to-br from-muted/50 to-transparent border-b border-border">
                      <div className="flex items-start justify-between mb-3">
                        <div className="flex items-center gap-3 flex-1 min-w-0">
                          <div className={`p-2 rounded-lg ${getTeamTypeColor(team.type)} border`}>
                            {getTeamTypeIcon(team.type)}
                          </div>
                          <div className="flex-1 min-w-0">
                            <h3 className="text-lg font-semibold text-foreground truncate">{team.name}</h3>
                            <div className="flex items-center gap-2 mt-1 flex-wrap">
                              {team.type && (
                                <span className={`inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded ${getTeamTypeColor(team.type)} border`}>
                                  {team.type.charAt(0).toUpperCase() + team.type.slice(1)}
                                </span>
                              )}
                              {team.pr_prefix && (
                                <span className="inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded bg-muted text-muted-foreground border border-border">
                                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" />
                                  </svg>
                                  {team.pr_prefix}
                                </span>
                              )}
                            </div>
                          </div>
                        </div>
                        <div className="flex items-center gap-1">
                          <button
                            onClick={() => openEditTeamModal(team)}
                            className="p-1.5 rounded hover:bg-muted transition-colors text-muted-foreground hover:text-foreground"
                            title="Edit team"
                          >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                            </svg>
                          </button>
                          <button
                            onClick={() => handleDeleteTeam(team)}
                            className="p-1.5 rounded hover:bg-destructive/10 transition-colors text-muted-foreground hover:text-destructive"
                            title="Delete team"
                          >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                          </button>
                        </div>
                      </div>
                      {team.description && (
                        <p className="text-sm text-muted-foreground line-clamp-2 mt-2">{team.description}</p>
                      )}
                    </div>

                    {/* Members Section */}
                    <div className="p-5">
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-2">
                          <svg className="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                          </svg>
                          <h4 className="text-sm font-medium text-foreground">
                            {team.members.length} {team.members.length === 1 ? 'Member' : 'Members'}
                          </h4>
                        </div>
                        <button
                          onClick={() => openAddMemberToTeamModal(team)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-primary hover:bg-primary/10 rounded transition-colors"
                          title="Add member"
                        >
                          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                          </svg>
                          Add
                        </button>
                      </div>

                      {team.members.length === 0 ? (
                        <div className="flex flex-col items-center justify-center py-6 text-center">
                          <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center mb-2">
                            <svg className="w-6 h-6 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                            </svg>
                          </div>
                          <p className="text-sm text-muted-foreground">No members yet</p>
                          <p className="text-xs text-muted-foreground mt-1">Add members to get started</p>
                        </div>
                      ) : (
                        <div className="space-y-2">
                          {team.members.slice(0, 4).map((member) => {
                            const memberData = members.find(m => m.id === member.member_id)
                            const initial = memberData?.username?.charAt(0).toUpperCase() || '?'
                            return (
                              <div key={member.id} className="flex items-center justify-between p-2 rounded-lg hover:bg-muted/50 transition-colors group">
                                <div className="flex items-center gap-2 flex-1 min-w-0">
                                  <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
                                    <span className="text-xs font-medium text-primary">{initial}</span>
                                  </div>
                                  <div className="flex-1 min-w-0">
                                    <p className="text-sm font-medium text-foreground truncate">
                                      {memberData?.username || 'Unknown'}
                                    </p>
                                    <p className="text-xs text-muted-foreground truncate">
                                      {memberData?.email || 'No email'}
                                    </p>
                                  </div>
                                </div>
                                <button
                                  onClick={() => handleRemoveMemberFromTeam(team.id, member.member_id)}
                                  className="opacity-0 group-hover:opacity-100 p-1 rounded hover:bg-destructive/10 text-destructive transition-all"
                                  title="Remove member"
                                >
                                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                  </svg>
                                </button>
                              </div>
                            )
                          })}
                          {team.members.length > 4 && (
                            <div className="pt-2 text-center">
                              <p className="text-xs text-muted-foreground">
                                +{team.members.length - 4} more {team.members.length - 4 === 1 ? 'member' : 'members'}
                              </p>
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        )}

        {/* Add Member Modal */}
        {isAddMemberModalOpen && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md mx-4">
              <h3 className="text-lg font-semibold text-foreground mb-4">Add Team Member</h3>
              <form onSubmit={handleAddMember}>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium mb-2">Email</label>
                    <input
                      type="email"
                      value={addMemberFormData.email}
                      onChange={(e) => setAddMemberFormData({ ...addMemberFormData, email: e.target.value })}
                      className="w-full px-3 py-2 border border-border rounded-md"
                      placeholder="member@example.com"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-2">Username</label>
                    <input
                      type="text"
                      value={addMemberFormData.username}
                      onChange={(e) => setAddMemberFormData({ ...addMemberFormData, username: e.target.value })}
                      className="w-full px-3 py-2 border border-border rounded-md"
                      placeholder="username"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-2">Title</label>
                    <select
                      value={addMemberFormData.title_id}
                      onChange={(e) => setAddMemberFormData({ ...addMemberFormData, title_id: e.target.value })}
                      className="w-full px-3 py-2 border border-border rounded-md"
                    >
                      <option value="">Select a title</option>
                      {titles.map((title) => (
                        <option key={title.id} value={title.id}>
                          {title.name}
                        </option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-2">Manager</label>
                    <select
                      value={addMemberFormData.manager_id || ''}
                      onChange={(e) => setAddMemberFormData({ ...addMemberFormData, manager_id: e.target.value || undefined })}
                      className="w-full px-3 py-2 border border-border rounded-md"
                    >
                      <option value="">No manager</option>
                      {members.map((member) => (
                        <option key={member.id} value={member.id}>
                          {member.username} ({member.email})
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div className="flex justify-end space-x-3 mt-6">
                  <button
                    type="button"
                    onClick={() => setIsAddMemberModalOpen(false)}
                    className="px-4 py-2 text-muted-foreground hover:text-foreground"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={isAddingMember}
                    className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
                  >
                    {isAddingMember ? 'Adding...' : 'Add Member'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* Update Member Modal */}
        {isUpdateMemberModalOpen && selectedMember && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md mx-4">
              <h3 className="text-lg font-semibold text-foreground mb-4">Update Team Member</h3>
              <form onSubmit={handleUpdateMember}>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium mb-2">Username</label>
                    <input
                      type="text"
                      value={updateMemberFormData.username}
                      onChange={(e) => setUpdateMemberFormData({ ...updateMemberFormData, username: e.target.value })}
                      className="w-full px-3 py-2 border border-border rounded-md"
                      placeholder="username"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-2">Title</label>
                    <select
                      value={updateMemberFormData.title_id}
                      onChange={(e) => setUpdateMemberFormData({ ...updateMemberFormData, title_id: e.target.value })}
                      className="w-full px-3 py-2 border border-border rounded-md"
                    >
                      <option value="">Select a title</option>
                      {titles.map((title) => (
                        <option key={title.id} value={title.id}>
                          {title.name}
                        </option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium mb-2">Manager</label>
                    <select
                      value={updateMemberFormData.manager_id || ''}
                      onChange={(e) => setUpdateMemberFormData({ ...updateMemberFormData, manager_id: e.target.value || undefined })}
                      className="w-full px-3 py-2 border border-border rounded-md"
                    >
                      <option value="">No manager</option>
                      {members.filter(m => m.id !== selectedMember.id).map((member) => (
                        <option key={member.id} value={member.id}>
                          {member.username} ({member.email})
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div className="flex justify-end space-x-3 mt-6">
                  <button
                    type="button"
                    onClick={() => setIsUpdateMemberModalOpen(false)}
                    className="px-4 py-2 text-muted-foreground hover:text-foreground"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={isUpdatingMember}
                    className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
                  >
                    {isUpdatingMember ? 'Updating...' : 'Update Member'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* Create Team Modal */}
        {isCreateTeamModalOpen && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md mx-4">
              <h3 className="text-lg font-semibold text-foreground mb-4">Create Team</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Team Name</label>
                  <input
                    type="text"
                    value={createTeamFormData.name}
                    onChange={(e) => setCreateTeamFormData({ ...createTeamFormData, name: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="Enter team name"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Description</label>
                  <textarea
                    value={createTeamFormData.description}
                    onChange={(e) => setCreateTeamFormData({ ...createTeamFormData, description: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="Enter team description"
                    rows={3}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Team Type</label>
                  <select
                    value={createTeamFormData.type || ''}
                    onChange={(e) => setCreateTeamFormData({ ...createTeamFormData, type: e.target.value as TeamType || undefined })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                  >
                    <option value="">Select a type (optional)</option>
                    <option value="squad">Squad</option>
                    <option value="chapter">Chapter</option>
                    <option value="tribe">Tribe</option>
                    <option value="guild">Guild</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">PR Prefix</label>
                  <input
                    type="text"
                    value={createTeamFormData.pr_prefix || ''}
                    onChange={(e) => setCreateTeamFormData({ ...createTeamFormData, pr_prefix: e.target.value || undefined })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="e.g., FE, BE, API"
                  />
                  <p className="text-xs text-muted-foreground mt-1">Optional prefix for pull requests (e.g., FE, BE, API)</p>
                </div>
              </div>

              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => setIsCreateTeamModalOpen(false)}
                  className="px-4 py-2 text-muted-foreground hover:text-foreground"
                >
                  Cancel
                </button>
                <button
                  onClick={handleCreateTeam}
                  disabled={isCreatingTeam || !createTeamFormData.name}
                  className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
                >
                  {isCreatingTeam ? 'Creating...' : 'Create Team'}
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Edit Team Modal */}
        {isEditTeamModalOpen && selectedTeam && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md mx-4">
              <h3 className="text-lg font-semibold text-foreground mb-4">Edit Team</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Team Name</label>
                  <input
                    type="text"
                    value={updateTeamFormData.name}
                    onChange={(e) => setUpdateTeamFormData({ ...updateTeamFormData, name: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="Enter team name"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Description</label>
                  <textarea
                    value={updateTeamFormData.description}
                    onChange={(e) => setUpdateTeamFormData({ ...updateTeamFormData, description: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="Enter team description"
                    rows={3}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Team Type</label>
                  <select
                    value={updateTeamFormData.type || ''}
                    onChange={(e) => setUpdateTeamFormData({ ...updateTeamFormData, type: e.target.value as TeamType || undefined })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                  >
                    <option value="">Select a type (optional)</option>
                    <option value="squad">Squad</option>
                    <option value="chapter">Chapter</option>
                    <option value="tribe">Tribe</option>
                    <option value="guild">Guild</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">PR Prefix</label>
                  <input
                    type="text"
                    value={updateTeamFormData.pr_prefix || ''}
                    onChange={(e) => setUpdateTeamFormData({ ...updateTeamFormData, pr_prefix: e.target.value || undefined })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="e.g., FE, BE, API"
                  />
                  <p className="text-xs text-muted-foreground mt-1">Optional prefix for pull requests (e.g., FE, BE, API)</p>
                </div>
              </div>

              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => setIsEditTeamModalOpen(false)}
                  className="px-4 py-2 text-muted-foreground hover:text-foreground"
                >
                  Cancel
                </button>
                <button
                  onClick={handleUpdateTeam}
                  disabled={isUpdatingTeam || !updateTeamFormData.name}
                  className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
                >
                  {isUpdatingTeam ? 'Updating...' : 'Update Team'}
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Add Member to Team Modal */}
        {isAddMemberToTeamModalOpen && selectedTeam && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md mx-4">
              <h3 className="text-lg font-semibold text-foreground mb-4">Add Member to Team</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Member</label>
                  <select
                    value={addMemberToTeamFormData.member_id}
                    onChange={(e) => setAddMemberToTeamFormData({ ...addMemberToTeamFormData, member_id: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                  >
                    <option value="">Select a member</option>
                    {members.map((member) => (
                      <option key={member.id} value={member.id}>
                        {member.username} ({member.email})
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => setIsAddMemberToTeamModalOpen(false)}
                  className="px-4 py-2 text-muted-foreground hover:text-foreground"
                >
                  Cancel
                </button>
                <button
                  onClick={handleAddMemberToTeam}
                  disabled={isAddingMemberToTeam || !addMemberToTeamFormData.member_id}
                  className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
                >
                  {isAddingMemberToTeam ? 'Adding...' : 'Add Member'}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
