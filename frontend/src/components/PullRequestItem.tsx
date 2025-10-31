import { useState } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { PullRequest } from '../types/sourcecontrol'

interface PullRequestItemProps {
  pr: PullRequest
  showRepository?: boolean
  showAuthor?: boolean
}

export default function PullRequestItem({ pr, showRepository = false, showAuthor = false }: PullRequestItemProps) {
  const [isDescriptionExpanded, setIsDescriptionExpanded] = useState(false)

  const toggleDescription = () => {
    setIsDescriptionExpanded(!isDescriptionExpanded)
  }

  const getStatusBadge = () => {
    if (pr.status === 'open') {
      return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'
    } else if (pr.status === 'closed' && pr.merged_at) {
      return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300'
    } else if (pr.status === 'closed') {
      return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300'
    } else {
      return 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300'
    }
  }

  const getStatusLabel = () => {
    if (pr.status === 'open') return 'Open'
    if (pr.status === 'closed' && pr.merged_at) return 'Merged'
    if (pr.status === 'closed') return 'Closed'
    return pr.status
  }

  const formatTimeToMerge = () => {
    if (!pr.merged_at) return null
    
    const created = new Date(pr.created_at)
    const merged = new Date(pr.merged_at)
    const diffMs = merged.getTime() - created.getTime()
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
    const diffHours = Math.floor((diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    
    if (diffDays > 0) {
      return `${diffDays}d ${diffHours}h`
    } else if (diffHours > 0) {
      return `${diffHours}h`
    } else {
      const diffMinutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60))
      return `${diffMinutes}m`
    }
  }

  return (
    <div>
      <h3 className="text-lg font-medium text-foreground mb-3">
        {pr.url ? (
          <a
            href={pr.url}
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-primary transition-colors"
          >
            {pr.title}
          </a>
        ) : (
          pr.title
        )}
      </h3>
      
      {/* PR Statistics */}
      <div className="flex flex-wrap gap-4 text-sm text-muted-foreground mb-3">
        {/* PR Status */}
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusBadge()}`}>
          {getStatusLabel()}
        </span>
        
        {showRepository && pr.repository_name && (
          <span className="px-2 py-1 rounded text-xs bg-muted text-muted-foreground">
            {pr.repository_name}
          </span>
        )}
        
        {showAuthor && pr.author && (
          <span className="flex items-center space-x-1">
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
              />
            </svg>
            <span>{pr.author.username || pr.author.email || 'Unknown'}</span>
          </span>
        )}
        
        {pr.additions !== undefined && (
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            <span>+{pr.additions}</span>
          </span>
        )}
        
        {pr.deletions !== undefined && (
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
            </svg>
            <span>-{pr.deletions}</span>
          </span>
        )}
        
        {pr.comments !== undefined && pr.comments > 0 && (
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
            <span>{pr.comments} comment{pr.comments !== 1 ? 's' : ''}</span>
          </span>
        )}
        
        {pr.review_comments !== undefined && pr.review_comments > 0 && (
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-purple-600 dark:text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
            <span>{pr.review_comments} review comment{pr.review_comments !== 1 ? 's' : ''}</span>
          </span>
        )}
        
        {pr.changed_files !== undefined && (
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-orange-600 dark:text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span>{pr.changed_files} file{pr.changed_files !== 1 ? 's' : ''}</span>
          </span>
        )}
        
        {pr.merged_at && (
          <span className="flex items-center space-x-1">
            <svg className="w-4 h-4 text-indigo-600 dark:text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{formatTimeToMerge()} to merge</span>
          </span>
        )}
      </div>
      
      {/* Expandable PR Description */}
      {pr.description && (
        <div className="mb-3">
          <button
            onClick={toggleDescription}
            className="flex items-center space-x-2 text-sm text-primary hover:text-primary/80 transition-colors"
          >
            <svg 
              className={`w-4 h-4 transition-transform ${isDescriptionExpanded ? 'rotate-90' : ''}`}
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
            <span>
              {isDescriptionExpanded ? 'Hide description' : 'Show description'}
            </span>
          </button>
          {isDescriptionExpanded && (
            <div className="mt-2 p-3 bg-muted/30 rounded-md border border-border">
              <div className="prose prose-sm max-w-none dark:prose-invert">
                <ReactMarkdown
                  remarkPlugins={[remarkGfm]}
                  components={{
                    p: ({ children }) => <p className="mb-2 last:mb-0 text-sm text-muted-foreground">{children}</p>,
                    ul: ({ children }) => <ul className="list-disc list-inside mb-2 space-y-1 text-sm text-muted-foreground">{children}</ul>,
                    ol: ({ children }) => <ol className="list-decimal list-inside mb-2 space-y-1 text-sm text-muted-foreground">{children}</ol>,
                    li: ({ children }) => <li className="text-sm">{children}</li>,
                    strong: ({ children }) => <strong className="font-semibold text-foreground">{children}</strong>,
                    em: ({ children }) => <em className="italic">{children}</em>,
                    code: ({ children }) => <code className="bg-background px-1.5 py-0.5 rounded text-xs font-mono border border-border">{children}</code>,
                    pre: ({ children }) => <pre className="bg-background p-2 rounded text-xs font-mono overflow-x-auto mb-2 border border-border">{children}</pre>,
                    blockquote: ({ children }) => <blockquote className="border-l-4 border-primary/30 pl-4 italic mb-2 text-muted-foreground">{children}</blockquote>,
                    h1: ({ children }) => <h1 className="text-base font-bold mb-2 text-foreground">{children}</h1>,
                    h2: ({ children }) => <h2 className="text-sm font-bold mb-2 text-foreground">{children}</h2>,
                    h3: ({ children }) => <h3 className="text-sm font-semibold mb-1 text-foreground">{children}</h3>,
                    a: ({ children, href }) => (
                      <a href={href} target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 underline">
                        {children}
                      </a>
                    ),
                    hr: () => <hr className="my-2 border-border" />,
                    table: ({ children }) => <div className="overflow-x-auto"><table className="min-w-full border-collapse border border-border">{children}</table></div>,
                    thead: ({ children }) => <thead className="bg-muted/50">{children}</thead>,
                    tbody: ({ children }) => <tbody>{children}</tbody>,
                    tr: ({ children }) => <tr className="border-b border-border">{children}</tr>,
                    th: ({ children }) => <th className="border border-border px-2 py-1 text-left text-xs font-semibold">{children}</th>,
                    td: ({ children }) => <td className="border border-border px-2 py-1 text-xs">{children}</td>,
                  }}
                >
                  {pr.description}
                </ReactMarkdown>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

