import { create } from 'zustand'

export type ApiErrorType = 'network' | 'server' | null

interface ApiErrorState {
  isApiUnavailable: boolean
  errorType: ApiErrorType
  errorMessage: string | null
  setApiUnavailable: (errorType: ApiErrorType, message?: string) => void
  clearApiError: () => void
}

const DEFAULT_ERROR_MESSAGE = 'Unable to connect to the server. Please check your connection and try again.'

export const useApiErrorStore = create<ApiErrorState>((set) => ({
  isApiUnavailable: false,
  errorType: null,
  errorMessage: null,

  setApiUnavailable: (errorType: ApiErrorType, message?: string) => {
    set({
      isApiUnavailable: true,
      errorType,
      errorMessage: message || DEFAULT_ERROR_MESSAGE,
    })
  },

  clearApiError: () => {
    set({
      isApiUnavailable: false,
      errorType: null,
      errorMessage: null,
    })
  },
}))
