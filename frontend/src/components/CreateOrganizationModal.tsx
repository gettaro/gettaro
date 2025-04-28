import React, { useState } from 'react'
import { OrganizationConflictError } from '../api/errors/organizations'

interface CreateOrganizationModalProps {
  isOpen: boolean
  onClose: () => void
  onCreate: (name: string, slug: string) => Promise<void>
}

export default function CreateOrganizationModal({
  isOpen,
  onClose,
  onCreate,
}: CreateOrganizationModalProps) {
  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setIsLoading(true)

    try {
      await onCreate(name, slug)
      setName('')
      setSlug('')
      onClose()
      // The organization object will be available in the parent component
      // and the navigation will be handled there
    } catch (err) {
      if (err instanceof OrganizationConflictError) {
        setError('An organization with this name already exists. Please choose a different name.')
      } else {
        setError(err instanceof Error ? err.message : 'Failed to create organization')
      }
    } finally {
      setIsLoading(false)
    }
  }

  if (!isOpen) return null

  return (
    <div className="fixed top-64 left-0 right-0 bottom-0 bg-background/80 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-card p-6 rounded-lg shadow-lg w-full max-w-md mx-4 relative">
        <button
          onClick={onClose}
          className="absolute top-2 right-2 text-muted-foreground hover:text-foreground"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M18 6 6 18" />
            <path d="m6 6 12 12" />
          </svg>
        </button>
        <h2 className="text-2xl font-bold mb-4">Create Organization</h2>
        <form onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-foreground mb-1">
                Organization Name
              </label>
              <input
                type="text"
                id="name"
                value={name}
                onChange={(e) => {
                  setName(e.target.value)
                  setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, '-'))
                }}
                className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground"
                required
              />
            </div>
            <div>
              <label htmlFor="slug" className="block text-sm font-medium text-foreground mb-1">
                Organization Slug
              </label>
              <input
                type="text"
                id="slug"
                value={slug}
                onChange={(e) => setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, '-'))}
                className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground"
                required
              />
            </div>
            {error && (
              <div className="text-sm text-destructive">
                {error}
              </div>
            )}
          </div>
          <div className="mt-6 flex justify-end space-x-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-foreground hover:text-primary"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="px-4 py-2 text-sm bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
            >
              {isLoading ? 'Creating...' : 'Create Organization'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
} 