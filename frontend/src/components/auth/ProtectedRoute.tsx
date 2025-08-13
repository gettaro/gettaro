import React from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'

interface ProtectedRouteProps {
  children: React.ReactNode
}



export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, getToken } = useAuth()
  console.log("hello world")
  if (isAuthenticated) {
    console.log("initialising token...")
    getToken()
  }

  const location = useLocation()

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
      </div>
    )
  }

  if (!isAuthenticated) {
    return <Navigate to="/" state={{ from: location }} replace />
  }

  return <>{children}</>
} 