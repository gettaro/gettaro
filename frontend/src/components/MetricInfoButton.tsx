import { useState } from 'react'

interface MetricInfoButtonProps {
  description: string
}

export default function MetricInfoButton({ description }: MetricInfoButtonProps) {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <div className="relative inline-flex">
      <button
        type="button"
        onClick={(e) => {
          e.stopPropagation()
          setIsOpen(!isOpen)
        }}
        className="ml-2 text-muted-foreground hover:text-foreground transition-colors"
        aria-label="Metric information"
      >
        <svg
          className="w-4 h-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
      </button>
      {isOpen && (
        <>
          <div
            className="fixed inset-0 z-40"
            onClick={() => setIsOpen(false)}
          />
          <div className="absolute left-0 bottom-full mb-2 z-50 w-64 p-3 bg-popover border border-border rounded-lg shadow-lg text-sm text-popover-foreground">
            <p className="whitespace-pre-wrap">{description}</p>
            <div className="absolute left-4 top-full w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-border" />
          </div>
        </>
      )}
    </div>
  )
}

