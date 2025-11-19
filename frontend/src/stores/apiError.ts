import { create } from 'zustand'

export type ApiErrorType = 'network' | 'server' | null

interface ApiErrorState {
  errorType: ApiErrorType
  errorMessage: string | null
  isApiUnavailable: boolean
  setApiError: (type: ApiErrorType, message: string | null) => void
  clearApiError: () => void
}

export const useApiErrorStore = create<ApiErrorState>((set) => ({
  errorType: null,
  errorMessage: null,
  isApiUnavailable: false,
  
  setApiError: (type, message) => {
    set({
      errorType: type,
      errorMessage: message,
      isApiUnavailable: type !== null,
    })
  },
  
  clearApiError: () => {
    set({
      errorType: null,
      errorMessage: null,
      isApiUnavailable: false,
    })
  },
}))
