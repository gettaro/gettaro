import { useState, useEffect } from 'react'
import Api from '../api/api'
import { Title, CreateTitleRequest } from '../types/title'
import { useAuth } from '../hooks/useAuth'
import { useOrganizationStore } from '../stores/organization'

export default function Titles() {
  const { user } = useAuth()
  const { currentOrganization } = useOrganizationStore()
  const [titles, setTitles] = useState<Title[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [newTitleName, setNewTitleName] = useState('')
  const [isCreating, setIsCreating] = useState(false)

  useEffect(() => {
    if (currentOrganization) {
      loadTitles()
    }
  }, [currentOrganization])

  const loadTitles = async () => {
    if (!currentOrganization) {
      setError('No organization selected')
      setLoading(false)
      return
    }

    try {
      setLoading(true)
      setError(null)
      const titlesData = await Api.getOrganizationTitles(currentOrganization.id)
      setTitles(titlesData)
    } catch (err) {
      setError('Failed to load titles')
      console.error('Error loading titles:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateTitle = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newTitleName.trim() || !currentOrganization) return

    try {
      setIsCreating(true)
      setError(null)
      const request: CreateTitleRequest = { name: newTitleName.trim() }
      await Api.createTitle(currentOrganization.id, request)
      setNewTitleName('')
      setIsCreateModalOpen(false)
      await loadTitles() // Reload the list
    } catch (err) {
      setError('Failed to create title')
      console.error('Error creating title:', err)
    } finally {
      setIsCreating(false)
    }
  }

  const handleDeleteTitle = async (titleId: string) => {
    if (!confirm('Are you sure you want to delete this title?') || !currentOrganization) return

    try {
      setError(null)
      await Api.deleteTitle(currentOrganization.id, titleId)
      await loadTitles() // Reload the list
    } catch (err) {
      setError('Failed to delete title')
      console.error('Error deleting title:', err)
    }
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
            <div className="text-muted-foreground">Loading titles...</div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">Titles</h1>
          <p className="text-muted-foreground">
            Manage job titles within your organization.
          </p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        <div className="bg-card border border-border rounded-lg">
          <div className="p-6 border-b border-border">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-semibold text-foreground">Organization Titles</h2>
              <button
                onClick={() => setIsCreateModalOpen(true)}
                className="bg-primary text-primary-foreground hover:bg-primary/90 px-4 py-2 rounded-md transition-colors"
              >
                Add Title
              </button>
            </div>
          </div>

          <div className="divide-y divide-border">
            {titles.length === 0 ? (
              <div className="p-6 text-center text-muted-foreground">
                No titles found. Create your first title to get started.
              </div>
            ) : (
              titles.map((title) => (
                <div key={title.id} className="p-6 flex items-center justify-between">
                  <div>
                    <h3 className="font-medium text-foreground">{title.name}</h3>
                    <p className="text-sm text-muted-foreground">
                      Created {new Date(title.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="flex items-center space-x-3">
                    <button
                      onClick={() => handleDeleteTitle(title.id)}
                      className="text-red-600 hover:text-red-800 transition-colors"
                      title="Delete title"
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
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

        {/* Create Title Modal */}
        {isCreateModalOpen && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md mx-4">
              <h3 className="text-lg font-semibold text-foreground mb-4">Create New Title</h3>
              <form onSubmit={handleCreateTitle}>
                <div className="mb-4">
                  <label htmlFor="titleName" className="block text-sm font-medium text-foreground mb-2">
                    Title Name
                  </label>
                  <input
                    type="text"
                    id="titleName"
                    value={newTitleName}
                    onChange={(e) => setNewTitleName(e.target.value)}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
                    placeholder="e.g., Senior Software Engineer"
                    required
                  />
                </div>
                <div className="flex justify-end space-x-3">
                  <button
                    type="button"
                    onClick={() => {
                      setIsCreateModalOpen(false)
                      setNewTitleName('')
                    }}
                    className="px-4 py-2 text-muted-foreground hover:text-foreground transition-colors"
                    disabled={isCreating}
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="bg-primary text-primary-foreground hover:bg-primary/90 px-4 py-2 rounded-md transition-colors disabled:opacity-50"
                    disabled={isCreating || !newTitleName.trim()}
                  >
                    {isCreating ? 'Creating...' : 'Create Title'}
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