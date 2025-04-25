import { useAuth0 } from '@auth0/auth0-react'
import { User } from '@auth0/auth0-react'
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

  const login = async () => {
    await loginWithRedirect()
  }

  const logout = async () => {
    await auth0Logout({ logoutParams: { returnTo: window.location.origin } })
  }

  const getToken = async () => {
    try {
      const token = await getAccessTokenSilently()
      Api.setAccessToken(token)
      return token
    } catch (error) {
      console.error('Failed to get token:', error)
      Api.setAccessToken(null)
      throw error
    }
  }

  return {
    user: user || null,
    isAuthenticated,
    isLoading,
    getToken,
    login,
    logout,
  }
}
