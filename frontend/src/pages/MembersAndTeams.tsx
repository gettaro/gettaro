import { useState, useEffect } from 'react'
import Api from '../api/api'
import { Team, CreateTeamRequest, UpdateTeamRequest, AddTeamMemberRequest } from '../types/team'
import { Member, AddMemberRequest, UpdateMemberRequest } from '../types/member'
import { Title } from '../types/title'
import { SourceControlAccount } from '../types/sourcecontrol'
import { useOrganizationStore } from '../stores/organization'
import { useToast } from '../hooks/useToast'

export default function MembersAndTeams() {
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
  const [sourceControlAccounts, setSourceControlAccounts] = useState<SourceControlAccount[]>([])
  const [isAddMemberModalOpen, setIsAddMemberModalOpen] = useState(false)
  const [isAddingMember, setIsAddingMember] = useState(false)
  const [isUpdateMemberModalOpen, setIsUpdateMemberModalOpen] = useState(false)
  const [isUpdatingMember, setIsUpdatingMember] = useState(false)
  const [selectedMember, setSelectedMember] = useState<Member | null>(null)

  // General state
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<'members' | 'teams'>('members')

  // Form data
  const [createTeamFormData, setCreateTeamFormData] = useState<CreateTeamRequest>({
    name: '',
    description: '',
    organization_id: ''
  })

  const [updateTeamFormData, setUpdateTeamFormData] = useState<UpdateTeamRequest>({
    name: '',
    description: ''
  })

  const [addMemberToTeamFormData, setAddMemberToTeamFormData] = useState<AddTeamMemberRequest>({
    member_id: ''
  })

  const [addMemberFormData, setAddMemberFormData] = useState<AddMemberRequest>({
    email: '',
    username: '',
    title_id: '',
    source_control_account_id: '',
    manager_id: undefined
  })

  const [updateMemberFormData, setUpdateMemberFormData] = useState<UpdateMemberRequest>({
    username: '',
    title_id: '',
    source_control_account_id: '',
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
      setCreateTeamFormData({ name: '', description: '', organization_id: '' })
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
      setUpdateTeamFormData({ name: '', description: '' })
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
    setUpdateTeamFormData({ name: team.name, description: team.description })
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
      await Api.addOrganizationMember(currentOrganization.id, addMemberFormData)
      setAddMemberFormData({
        email: '',
        username: '',
        title_id: '',
        source_control_account_id: '',
        manager_id: undefined
      })
      setIsAddMemberModalOpen(false)
      await loadMembers()
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
        source_control_account_id: '',
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
      source_control_account_id: sourceControlAccounts.find(acc => acc.member_id === member.id)?.id || '',
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
                      <tr key={member.id} className="border-b border-border last:border-b-0">
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
                          <div className="flex space-x-2">
                            <button
                              onClick={() => handleEditMember(member)}
                              className="text-primary hover:text-primary/80 text-sm"
                            >
                              Edit
                            </button>
                            <button
                              onClick={() => handleDeleteMember(member)}
                              className="text-destructive hover:text-destructive/80 text-sm"
                            >
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

            {/* Teams List */}
            <div className="grid gap-4">
              {teams.map((team) => (
                <div key={team.id} className="bg-card border border-border rounded-lg p-6">
                  <div className="flex justify-between items-start mb-4">
                    <div>
                      <h3 className="text-lg font-semibold text-foreground">{team.name}</h3>
                      {team.description && (
                        <p className="text-muted-foreground mt-1">{team.description}</p>
                      )}
                    </div>
                    <div className="flex space-x-2">
                      <button
                        onClick={() => openEditTeamModal(team)}
                        className="text-primary hover:text-primary/80 px-3 py-1 text-sm"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => openAddMemberToTeamModal(team)}
                        className="text-primary hover:text-primary/80 px-3 py-1 text-sm"
                      >
                        Add Member
                      </button>
                      <button
                        onClick={() => handleDeleteTeam(team)}
                        className="text-destructive hover:text-destructive/80 px-3 py-1 text-sm"
                      >
                        Delete
                      </button>
                    </div>
                  </div>

                  {/* Team Members */}
                  <div className="mt-4">
                    <h4 className="text-sm font-medium text-foreground mb-2">Members ({team.members.length})</h4>
                    {team.members.length === 0 ? (
                      <p className="text-muted-foreground text-sm">No members yet</p>
                    ) : (
                      <div className="space-y-1">
                        {team.members.map((member) => (
                          <div key={member.id} className="flex justify-between items-center text-sm">
                            <span>{getMemberName(member.member_id)}</span>
                            <div className="flex items-center space-x-2">
                              <button
                                onClick={() => handleRemoveMemberFromTeam(team.id, member.member_id)}
                                className="text-xs text-destructive hover:text-destructive/80"
                              >
                                Remove
                              </button>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              ))}
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
                    <label className="block text-sm font-medium mb-2">Source Control Account</label>
                    <select
                      value={addMemberFormData.source_control_account_id}
                      onChange={(e) => setAddMemberFormData({ ...addMemberFormData, source_control_account_id: e.target.value })}
                      className="w-full px-3 py-2 border border-border rounded-md"
                    >
                      <option value="">Select an account</option>
                      {sourceControlAccounts.map((account) => (
                        <option key={account.id} value={account.id}>
                          {account.username} ({account.provider_name})
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
                    <label className="block text-sm font-medium mb-2">Source Control Account</label>
                    <select
                      value={updateMemberFormData.source_control_account_id}
                      onChange={(e) => setUpdateMemberFormData({ ...updateMemberFormData, source_control_account_id: e.target.value })}
                      className="w-full px-3 py-2 border border-border rounded-md"
                    >
                      <option value="">Select an account</option>
                      {sourceControlAccounts.map((account) => (
                        <option key={account.id} value={account.id}>
                          {account.username} ({account.provider_name})
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
