import { useApiErrorStore } from '../stores/apiError'
import { Button } from './ui/button'
import { AlertCircle, X } from 'lucide-react'

export default function ApiErrorBanner() {
  const { isApiUnavailable, errorMessage, clearApiError } = useApiErrorStore()

  if (!isApiUnavailable) {
    return null
  }

  return (
    <div className="fixed top-0 left-0 right-0 z-50 bg-destructive text-destructive-foreground shadow-lg">
      <div className="container mx-auto px-4 py-3">
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-3 flex-1">
            <AlertCircle className="h-5 w-5 flex-shrink-0" />
            <p className="text-sm font-medium flex-1">{errorMessage}</p>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                clearApiError()
                // Reload the page to retry
                window.location.reload()
              }}
              className="text-destructive-foreground hover:bg-destructive-foreground/20"
            >
              Retry
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={clearApiError}
              className="text-destructive-foreground hover:bg-destructive-foreground/20 h-8 w-8"
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
