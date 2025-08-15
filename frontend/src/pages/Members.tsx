import { useState, useEffect } from 'react'
import Api from '../api/api'
import { Member, AddMemberRequest, UpdateMemberRequest } from '../types/member'
import { Title } from '../types/title'
import { SourceControlAccount } from '../types/sourcecontrol'
import { useOrganizationStore } from '../stores/organization'

export default function Members() {
  const { currentOrganization } = useOrganizationStore()
  const [members, setMembers] = useState<Member[]>([])
  const [titles, setTitles] = useState<Title[]>([])
  const [sourceControlAccounts, setSourceControlAccounts] = useState<SourceControlAccount[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isAddModalOpen, setIsAddModalOpen] = useState(false)
  const [isAdding, setIsAdding] = useState(false)
  const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false)
  const [isUpdating, setIsUpdating] = useState(false)
  const [selectedMember, setSelectedMember] = useState<Member | null>(null)
  const [formData, setFormData] = useState<AddMemberRequest>({
    email: '',
    username: '',
    titleId: '',
    sourceControlAccountId: ''
  })
  const [updateFormData, setUpdateFormData] = useState<UpdateMemberRequest>({
    username: '',
    titleId: '',
    sourceControlAccountId: ''
  })

  useEffect(() => {
    if (currentOrganization) {
      loadMembers()
      loadTitles()
      loadSourceControlAccounts()
    }
  }, [currentOrganization])

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
      console.log('Loading source control accounts for organization:', currentOrganization.id)
      const accountsData = await Api.getOrganizationSourceControlAccounts(currentOrganization.id)
      console.log('Source control accounts loaded:', accountsData)
      setSourceControlAccounts(accountsData)
    } catch (err) {
      console.error('Error loading source control accounts:', err)
    }
  }

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!currentOrganization) return

    try {
      setIsAdding(true)
      setError(null)
      await Api.addOrganizationMember(currentOrganization.id, formData)
      setFormData({
        email: '',
        username: '',
        titleId: '',
        sourceControlAccountId: ''
      })
      setIsAddModalOpen(false)
      await loadMembers() // Reload the list
    } catch (err) {
      setError('Failed to add member')
      console.error('Error adding member:', err)
    } finally {
      setIsAdding(false)
    }
  }

  const handleUpdateMember = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!currentOrganization || !selectedMember) return

    try {
      setIsUpdating(true)
      setError(null)
      await Api.updateOrganizationMember(currentOrganization.id, selectedMember.id, updateFormData)
      setUpdateFormData({
        username: '',
        titleId: '',
        sourceControlAccountId: ''
      })
      setSelectedMember(null)
      setIsUpdateModalOpen(false)
      await loadMembers() // Reload the list
    } catch (err) {
      setError('Failed to update member')
      console.error('Error updating member:', err)
    } finally {
      setIsUpdating(false)
    }
  }

  const handleEditMember = (member: Member) => {
    setSelectedMember(member)
    setUpdateFormData({
      username: member.username,
      titleId: member.titleId || '', // Pre-populate with current title
      sourceControlAccountId: sourceControlAccounts.find(acc => acc.memberId === member.id)?.id || '' // Pre-populate with current source control account
    })
    setIsUpdateModalOpen(true)
  }

  const handleDeleteMember = async (member: Member) => {
    if (!currentOrganization) return

    if (!confirm(`Are you sure you want to remove ${member.username} from the organization?`)) {
      return
    }

    try {
      setError(null)
      await Api.deleteOrganizationMember(currentOrganization.id, member.id)
      await loadMembers() // Reload the list
    } catch (err) {
      setError('Failed to delete member')
      console.error('Error deleting member:', err)
    }
  }

  const getRoleBadge = (isOwner: boolean) => {
    const baseClasses = "px-2 py-1 rounded-full text-xs font-medium"
    return isOwner 
      ? `${baseClasses} bg-red-100 text-red-800`
      : `${baseClasses} bg-gray-100 text-gray-800`
  }

  if (!currentOrganization) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">No organization selected</div>
          </div>
        </div>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">Loading members...</div>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <div className="text-red-600">{error}</div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">Members</h1>
          <p className="text-muted-foreground">
            Manage team members and their permissions.
          </p>
        </div>

        <div className="bg-card border border-border rounded-lg">
          <div className="p-6 border-b border-border">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-semibold text-foreground">Team Members</h2>
              <button 
                onClick={() => setIsAddModalOpen(true)}
                className="bg-primary text-primary-foreground hover:bg-primary/90 px-4 py-2 rounded-md transition-colors"
              >
                Add Member
              </button>
            </div>
          </div>

          <div className="divide-y divide-border">
            {members.map((member) => (
              <div key={member.id} className="p-6 flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                    <span className="text-primary font-medium">
                      {member.username.charAt(0).toUpperCase()}
                    </span>
                  </div>
                  <div>
                    <h3 className="font-medium text-foreground">{member.username}</h3>
                    <p className="text-sm text-muted-foreground">{member.email}</p>
                    <div className="mt-2 space-y-1">
                      {member.titleId && (
                        <div className="flex items-center space-x-2 text-sm text-foreground">
                          <svg className="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2-2v2m8 0V6a2 2 0 012 2v6a2 2 0 01-2 2H8a2 2 0 01-2-2V8a2 2 0 012-2z" />
                          </svg>
                          <span>{titles.find(t => t.id === member.titleId)?.name || 'Unknown Title'}</span>
                        </div>
                      )}
                      {sourceControlAccounts.find(acc => acc.memberId === member.id) && (
                        <div className="flex items-center space-x-2 text-sm text-foreground">
                          <svg className="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                          </svg>
                          <span>
                            {sourceControlAccounts.find(acc => acc.memberId === member.id)?.username}
                            {sourceControlAccounts.find(acc => acc.memberId === member.id)?.providerName && 
                              ` (${sourceControlAccounts.find(acc => acc.memberId === member.id)?.providerName})`
                            }
                          </span>
                        </div>
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground mt-2">
                      Joined {new Date(member.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-3">
                  <span className={getRoleBadge(member.isOwner)}>
                    {member.isOwner ? 'Owner' : 'Member'}
                  </span>
                  {!member.isOwner && (
                    <>
                      <button 
                        onClick={() => handleEditMember(member)}
                        className="text-muted-foreground hover:text-foreground transition-colors"
                        title="Edit member"
                      >
                        <svg
                          className="w-5 h-5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                          xmlns="http://www.w3.org/2000/svg"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                          />
                        </svg>
                      </button>
                      <a
                        href={`/members/${member.id}/profile`}
                        className="text-green-500 hover:text-green-700 transition-colors"
                        title="View profile"
                      >
                        <svg
                          className="w-5 h-5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                          xmlns="http://www.w3.org/2000/svg"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                          />
                        </svg>
                      </a>
                      <a
                        href={`/members/${member.id}/activity`}
                        className="text-blue-500 hover:text-blue-700 transition-colors"
                        title="View activity"
                      >
                        <svg
                          className="w-5 h-5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                          xmlns="http://www.w3.org/2000/svg"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                          />
                        </svg>
                      </a>
                      <button 
                        onClick={() => handleDeleteMember(member)}
                        className="text-red-500 hover:text-red-700 transition-colors"
                        title="Delete member"
                      >
                        <svg
                          className="w-5 h-5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                          xmlns="http://www.w3.org/2000/svg"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                          />
                        </svg>
                      </button>
                    </>
                  )}
                  {/* Profile and Activity buttons for all members (including owners) */}
                  <a
                    href={`/members/${member.id}/profile`}
                    className="text-green-500 hover:text-green-700 transition-colors"
                    title="View profile"
                  >
                    <svg
                      className="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                      />
                    </svg>
                  </a>
                  <a
                    href={`/members/${member.id}/activity`}
                    className="text-blue-500 hover:text-blue-700 transition-colors"
                    title="View activity"
                  >
                    <svg
                      className="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                      />
                    </svg>
                  </a>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Add Member Modal */}
        {isAddModalOpen && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md mx-4">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-semibold text-foreground">Add Member</h2>
                <button
                  onClick={() => setIsAddModalOpen(false)}
                  className="text-muted-foreground hover:text-foreground"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <form onSubmit={handleAddMember} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Email *
                  </label>
                  <input
                    type="email"
                    required
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                    placeholder="member@example.com"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Username *
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.username}
                    onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                    placeholder="username"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Title *
                  </label>
                  <select
                    required
                    value={formData.titleId}
                    onChange={(e) => setFormData({ ...formData, titleId: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
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
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Source Control Account *
                  </label>
                  <select
                    required
                    value={formData.sourceControlAccountId}
                    onChange={(e) => setFormData({ ...formData, sourceControlAccountId: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                  >
                    <option value="">Select a source control account</option>
                    {sourceControlAccounts.length > 0 ? (
                      sourceControlAccounts.map((account) => {
                        console.log('Account:', account);
                        return (
                          <option key={account.id} value={account.id}>
                            {account.username} {account.providerName ? `(${account.providerName})` : ''}
                          </option>
                        );
                      })
                    ) : (
                      <option value="" disabled>No source control accounts available</option>
                    )}
                  </select>
                  {sourceControlAccounts.length === 0 && (
                    <p className="text-xs text-muted-foreground mt-1">
                      No source control accounts found for this organization
                    </p>
                  )}
                </div>

                <div className="flex space-x-3 pt-4">
                  <button
                    type="button"
                    onClick={() => setIsAddModalOpen(false)}
                    className="flex-1 px-4 py-2 border border-border rounded-md text-foreground hover:bg-muted transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={isAdding}
                    className="flex-1 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50 transition-colors"
                  >
                    {isAdding ? 'Adding...' : 'Add Member'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* Update Member Modal */}
        {isUpdateModalOpen && selectedMember && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md mx-4">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-semibold text-foreground">Update Member</h2>
                <button
                  onClick={() => {
                    setIsUpdateModalOpen(false)
                    setSelectedMember(null)
                    setUpdateFormData({ username: '', titleId: '', sourceControlAccountId: '' })
                  }}
                  className="text-muted-foreground hover:text-foreground"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <form onSubmit={handleUpdateMember} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Email
                  </label>
                  <input
                    type="email"
                    disabled
                    value={selectedMember.email}
                    className="w-full px-3 py-2 border border-border rounded-md bg-muted text-muted-foreground cursor-not-allowed"
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    Email cannot be updated
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Username *
                  </label>
                  <input
                    type="text"
                    required
                    value={updateFormData.username}
                    onChange={(e) => setUpdateFormData({ ...updateFormData, username: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                    placeholder="username"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Title *
                  </label>
                  <select
                    required
                    value={updateFormData.titleId}
                    onChange={(e) => setUpdateFormData({ ...updateFormData, titleId: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
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
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Source Control Account *
                  </label>
                  <select
                    required
                    value={updateFormData.sourceControlAccountId}
                    onChange={(e) => setUpdateFormData({ ...updateFormData, sourceControlAccountId: e.target.value })}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                  >
                    <option value="">Select a source control account</option>
                    {sourceControlAccounts.length > 0 ? (
                      sourceControlAccounts.map((account) => (
                        <option key={account.id} value={account.id}>
                          {account.username} {account.providerName ? `(${account.providerName})` : ''}
                        </option>
                      ))
                    ) : (
                      <option value="" disabled>No source control accounts available</option>
                    )}
                  </select>
                  {sourceControlAccounts.length === 0 && (
                    <p className="text-xs text-muted-foreground mt-1">
                      No source control accounts found for this organization
                    </p>
                  )}
                </div>

                <div className="flex space-x-3 pt-4">
                  <button
                    type="button"
                    onClick={() => {
                      setIsUpdateModalOpen(false)
                      setSelectedMember(null)
                      setUpdateFormData({ username: '', titleId: '', sourceControlAccountId: '' })
                    }}
                    className="flex-1 px-4 py-2 border border-border rounded-md text-foreground hover:bg-muted transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={isUpdating}
                    className="flex-1 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50 transition-colors"
                  >
                    {isUpdating ? 'Updating...' : 'Update Member'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </div>
  )
} 