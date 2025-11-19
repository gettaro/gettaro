import { useApiErrorStore } from '../stores/apiError'
import { Button } from './ui/button'

export default function ApiErrorBanner() {
  const { isApiUnavailable, errorMessage, clearApiError } = useApiErrorStore()

  if (!isApiUnavailable) {
    return null
  }

  const handleRetry = () => {
    // Reload the page to retry
    window.location.reload()
  }

  const handleDismiss = () => {
    // Clear error state (user can dismiss but error will reappear on next API call)
    clearApiError()
  }

  return (
    <div
      className="fixed top-0 left-0 right-0 z-50 bg-destructive text-destructive-foreground shadow-lg"
      role="alert"
      aria-live="assertive"
    >
      <div className="container mx-auto px-4 py-3">
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-3 flex-1 min-w-0">
            <svg
              className="w-5 h-5 flex-shrink-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
            <p className="text-sm font-medium flex-1 truncate">
              {errorMessage || 'Unable to connect to the server. Please check your connection and try again.'}
            </p>
          </div>
          <div className="flex items-center gap-2 flex-shrink-0">
            <Button
              variant="secondary"
              size="sm"
              onClick={handleRetry}
              className="bg-destructive-foreground/20 hover:bg-destructive-foreground/30 text-destructive-foreground"
            >
              Retry
            </Button>
            <button
              onClick={handleDismiss}
              className="p-1 rounded-md hover:bg-destructive-foreground/20 transition-colors"
              aria-label="Dismiss error message"
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
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
