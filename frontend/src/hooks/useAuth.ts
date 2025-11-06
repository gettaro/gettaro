import { useAuth0 } from '@auth0/auth0-react'
import { User } from '@auth0/auth0-react'
import { useCallback, useRef } from 'react'
import Api from '../api/api'

export interface AuthContext {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  getToken: () => Promise<string>
  login: () => Promise<void>
  logout: () => Promise<void>
}

export const useAuth = (): AuthContext => {
  const {
    getAccessTokenSilently,
    user,
    isAuthenticated,
    isLoading,
    loginWithRedirect,
    logout: auth0Logout,
  } = useAuth0<{
    getAccessTokenSilently: () => Promise<string>
  }>()

  const tokenRequestRef = useRef<Promise<string> | null>(null)

  const login = async () => {
    await loginWithRedirect()
  }

  const logout = async () => {
    await auth0Logout({ logoutParams: { returnTo: window.location.origin } })
    tokenRequestRef.current = null
  }

  const getToken = useCallback(async () => {
    // If there's already a token request in progress, return it
    if (tokenRequestRef.current) {
      return tokenRequestRef.current
    }

    // Create new token request
    const tokenPromise = (async () => {
      try {
        const token = await getAccessTokenSilently()
        Api.setAccessToken(token)
        tokenRequestRef.current = null // Clear ref when done
        return token
      } catch (error) {
        console.error('Failed to get token:', error)
        Api.setAccessToken(null)
        tokenRequestRef.current = null // Clear ref on error
        throw error
      }
    })()

    tokenRequestRef.current = tokenPromise
    return tokenPromise
  }, [getAccessTokenSilently])

  return {
    user: user || null,
    isAuthenticated,
    isLoading,
    getToken,
    login,
    logout,
  }
}
