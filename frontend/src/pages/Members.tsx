import { useState, useEffect } from 'react'
import Api from '../api/api'
import { Member, AddMemberRequest } from '../types/member'
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
  const [formData, setFormData] = useState<AddMemberRequest>({
    email: '',
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

  const getProviderIcon = (providerName: string) => {
    switch (providerName.toLowerCase()) {
      case 'github':
        return (
          <svg className="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
        )
      default:
        return (
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
          </svg>
        )
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
                    <p className="text-xs text-muted-foreground">
                      Joined {new Date(member.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-3">
                  <span className={getRoleBadge(member.isOwner)}>
                    {member.isOwner ? 'Owner' : 'Member'}
                  </span>
                  <button className="text-muted-foreground hover:text-foreground transition-colors">
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
                        d="M12 5v.01M12 12v.01M12 19v.01M12 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
                      />
                    </svg>
                  </button>
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
      </div>
    </div>
  )
} 