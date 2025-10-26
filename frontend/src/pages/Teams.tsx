import { useState, useEffect } from 'react'
import Api from '../api/api'
import { Team, CreateTeamRequest, UpdateTeamRequest, AddTeamMemberRequest } from '../types/team'
import { Member } from '../types/member'
import { useOrganizationStore } from '../stores/organization'
import { useToast } from '../hooks/useToast'

export default function Teams() {
  const { currentOrganization } = useOrganizationStore()
  const { toast } = useToast()
  const [teams, setTeams] = useState<Team[]>([])
  const [members, setMembers] = useState<Member[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [isEditModalOpen, setIsEditModalOpen] = useState(false)
  const [isAddMemberModalOpen, setIsAddMemberModalOpen] = useState(false)
  const [isCreating, setIsCreating] = useState(false)
  const [isUpdating, setIsUpdating] = useState(false)
  const [isAddingMember, setIsAddingMember] = useState(false)
  const [selectedTeam, setSelectedTeam] = useState<Team | null>(null)
  const [activeTab, setActiveTab] = useState<'teams' | 'members'>('teams')

  const [createFormData, setCreateFormData] = useState<CreateTeamRequest>({
    name: '',
    description: '',
    organization_id: ''
  })

  const [updateFormData, setUpdateFormData] = useState<UpdateTeamRequest>({
    name: '',
    description: ''
  })

  const [addMemberFormData, setAddMemberFormData] = useState<AddTeamMemberRequest>({
    member_id: ''
  })

  useEffect(() => {
    if (currentOrganization) {
      loadTeams()
      loadMembers()
    }
  }, [currentOrganization])

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

  const loadMembers = async () => {
    if (!currentOrganization) return

    try {
      const response = await Api.getOrganizationMembers(currentOrganization.id)
      setMembers(response)
    } catch (err) {
      console.error('Error loading members:', err)
    }
  }

  const handleCreateTeam = async () => {
    if (!currentOrganization) return

    try {
      setIsCreating(true)
      const teamData = { ...createFormData, organization_id: currentOrganization.id }
      await Api.createTeam(currentOrganization.id, teamData)
      
      toast({
        title: 'Success',
        description: 'Team created successfully',
      })
      
      setIsCreateModalOpen(false)
      setCreateFormData({ name: '', description: '', organization_id: '' })
      loadTeams()
    } catch (err) {
      console.error('Error creating team:', err)
      toast({
        title: 'Error',
        description: 'Failed to create team',
        variant: 'destructive',
      })
    } finally {
      setIsCreating(false)
    }
  }

  const handleUpdateTeam = async () => {
    if (!selectedTeam) return

    try {
      setIsUpdating(true)
      await Api.updateTeam(currentOrganization.id, selectedTeam.id, updateFormData)
      
      toast({
        title: 'Success',
        description: 'Team updated successfully',
      })
      
      setIsEditModalOpen(false)
      setSelectedTeam(null)
      setUpdateFormData({ name: '', description: '' })
      loadTeams()
    } catch (err) {
      console.error('Error updating team:', err)
      toast({
        title: 'Error',
        description: 'Failed to update team',
        variant: 'destructive',
      })
    } finally {
      setIsUpdating(false)
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

  const handleAddMember = async () => {
    if (!selectedTeam) return

    try {
      setIsAddingMember(true)
      await Api.addTeamMember(currentOrganization.id, selectedTeam.id, addMemberFormData)
      
      toast({
        title: 'Success',
        description: 'Member added to team successfully',
      })
      
      setIsAddMemberModalOpen(false)
      setSelectedTeam(null)
      setAddMemberFormData({ member_id: '' })
      loadTeams()
    } catch (err) {
      console.error('Error adding member to team:', err)
      toast({
        title: 'Error',
        description: 'Failed to add member to team',
        variant: 'destructive',
      })
    } finally {
      setIsAddingMember(false)
    }
  }

  const handleRemoveMember = async (teamId: string, memberId: string) => {
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

  const openEditModal = (team: Team) => {
    setSelectedTeam(team)
    setUpdateFormData({
      name: team.name,
      description: team.description
    })
    setIsEditModalOpen(true)
  }

  const openAddMemberModal = (team: Team) => {
    setSelectedTeam(team)
    setAddMemberFormData({ member_id: '' })
    setIsAddMemberModalOpen(true)
  }

  const getMemberName = (memberId: string) => {
    const member = members.find(m => m.id === memberId)
    return member ? member.username : 'Unknown Member'
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">Loading teams...</div>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="text-center">
            <div className="text-destructive mb-4">{error}</div>
            <button
              onClick={loadTeams}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">Teams</h1>
          <p className="text-muted-foreground">
            Manage teams and their members.
          </p>
        </div>

        {/* Tab Navigation */}
        <div className="mb-6">
          <nav className="flex space-x-8 border-b border-border">
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
            <button
              onClick={() => setActiveTab('members')}
              className={`py-2 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === 'members'
                  ? 'border-primary text-primary'
                  : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              All Members
            </button>
          </nav>
        </div>

        {/* Tab Content */}
        {activeTab === 'teams' ? (
          <div>
            {/* Teams Header */}
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-xl font-semibold">Teams ({teams.length})</h2>
              <button
                onClick={() => setIsCreateModalOpen(true)}
                className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
              >
                Create Team
              </button>
            </div>

            {/* Teams List */}
            {teams.length === 0 ? (
              <div className="text-center py-12">
                <div className="text-muted-foreground mb-4">No teams found</div>
                <button
                  onClick={() => setIsCreateModalOpen(true)}
                  className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
                >
                  Create your first team
                </button>
              </div>
            ) : (
              <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {teams.map((team) => (
                  <div key={team.id} className="border border-border rounded-lg p-6">
                    <div className="flex justify-between items-start mb-4">
                      <div>
                        <h3 className="text-lg font-semibold">{team.name}</h3>
                        <p className="text-sm text-muted-foreground mt-1">
                          {team.description || 'No description'}
                        </p>
                      </div>
                      <div className="flex space-x-2">
                        <button
                          onClick={() => openEditModal(team)}
                          className="text-sm text-primary hover:text-primary/80"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDeleteTeam(team)}
                          className="text-sm text-destructive hover:text-destructive/80"
                        >
                          Delete
                        </button>
                      </div>
                    </div>

                    <div className="mb-4">
                      <div className="text-sm font-medium mb-2">
                        Members ({team.members.length})
                      </div>
                      {team.members.length === 0 ? (
                        <div className="text-sm text-muted-foreground">No members</div>
                      ) : (
                        <div className="space-y-1">
                          {team.members.map((member) => (
                            <div key={member.id} className="flex justify-between items-center text-sm">
                              <span>{getMemberName(member.member_id)}</span>
                              <div className="flex items-center space-x-2">
                                <button
                                  onClick={() => handleRemoveMember(team.id, member.member_id)}
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

                    <button
                      onClick={() => openAddMemberModal(team)}
                      className="w-full text-sm text-primary hover:text-primary/80 border border-primary rounded-md py-2"
                    >
                      Add Member
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        ) : (
          <div>
            {/* Members Header */}
            <div className="mb-6">
              <h2 className="text-xl font-semibold">All Members ({members.length})</h2>
            </div>

            {/* Members List */}
            {members.length === 0 ? (
              <div className="text-center py-12">
                <div className="text-muted-foreground">No members found</div>
              </div>
            ) : (
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {members.map((member) => (
                  <div key={member.id} className="border border-border rounded-lg p-4">
                    <div className="font-medium">{member.username}</div>
                    <div className="text-sm text-muted-foreground">{member.email}</div>
                    {member.title_id && (
                      <div className="text-sm text-muted-foreground mt-1">
                        Title ID: {member.title_id}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Create Team Modal */}
        {isCreateModalOpen && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-background border border-border rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-semibold mb-4">Create Team</h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Name</label>
                  <input
                    type="text"
                    value={createFormData.name}
                    onChange={(e) => setCreateFormData({ ...createFormData, name: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="Team name"
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-2">Description</label>
                  <textarea
                    value={createFormData.description}
                    onChange={(e) => setCreateFormData({ ...createFormData, description: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="Team description"
                    rows={3}
                  />
                </div>
              </div>

              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => setIsCreateModalOpen(false)}
                  className="px-4 py-2 text-muted-foreground hover:text-foreground"
                >
                  Cancel
                </button>
                <button
                  onClick={handleCreateTeam}
                  disabled={isCreating || !createFormData.name}
                  className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
                >
                  {isCreating ? 'Creating...' : 'Create'}
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Edit Team Modal */}
        {isEditModalOpen && selectedTeam && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-background border border-border rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-semibold mb-4">Edit Team</h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Name</label>
                  <input
                    type="text"
                    value={updateFormData.name}
                    onChange={(e) => setUpdateFormData({ ...updateFormData, name: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="Team name"
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-2">Description</label>
                  <textarea
                    value={updateFormData.description}
                    onChange={(e) => setUpdateFormData({ ...updateFormData, description: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                    placeholder="Team description"
                    rows={3}
                  />
                </div>
              </div>

              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => setIsEditModalOpen(false)}
                  className="px-4 py-2 text-muted-foreground hover:text-foreground"
                >
                  Cancel
                </button>
                <button
                  onClick={handleUpdateTeam}
                  disabled={isUpdating || !updateFormData.name}
                  className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
                >
                  {isUpdating ? 'Updating...' : 'Update'}
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Add Member Modal */}
        {isAddMemberModalOpen && selectedTeam && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-background border border-border rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-semibold mb-4">Add Member to {selectedTeam.name}</h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Member</label>
                  <select
                    value={addMemberFormData.member_id}
                    onChange={(e) => setAddMemberFormData({ ...addMemberFormData, member_id: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md"
                  >
                    <option value="">Select a member</option>
                    {members
                      .filter(member => !selectedTeam.members.some(tm => tm.member_id === member.id))
                      .map((member) => (
                        <option key={member.id} value={member.id}>
                          {member.username} ({member.email})
                        </option>
                      ))}
                  </select>
                </div>
              </div>

              <div className="flex justify-end space-x-3 mt-6">
                <button
                  onClick={() => setIsAddMemberModalOpen(false)}
                  className="px-4 py-2 text-muted-foreground hover:text-foreground"
                >
                  Cancel
                </button>
                <button
                  onClick={handleAddMember}
                  disabled={isAddingMember || !addMemberFormData.member_id}
                  className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
                >
                  {isAddingMember ? 'Adding...' : 'Add Member'}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
