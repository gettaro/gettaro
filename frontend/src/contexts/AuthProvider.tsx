import React, { useState, createContext, useContext, useEffect } from 'react'
import { useAuth0 } from '@auth0/auth0-react'
import { User } from '@auth0/auth0-react'
import Api from '../api/api'

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: () => Promise<void>
  logout: () => Promise<void>
  getToken: () => Promise<string>
}

const AuthContext = createContext<AuthContextType | null>(null)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
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

  const [token, setToken] = useState<string | null>(null)

  const login = async () => {
    await loginWithRedirect()
  }

  const logout = async () => {
    await auth0Logout({ logoutParams: { returnTo: window.location.origin } })
    setToken(null)
    Api.setAccessToken(null)
  }

  const getToken = async (): Promise<string> => {
    try {
      const accessToken = await getAccessTokenSilently()
      setToken(accessToken)
      Api.setAccessToken(accessToken)
      return accessToken
    } catch (error) {
      console.error('Failed to get token:', error)
      setToken(null)
      Api.setAccessToken(null)
      throw error
    }
  }

  // Automatically get token when user becomes authenticated
  useEffect(() => {
    if (isAuthenticated && !token) {
      getToken()
    }
  }, [isAuthenticated, token])

  const value: AuthContextType = {
    user: user || null,
    isAuthenticated,
    isLoading,
    login,
    logout,
    getToken,
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
