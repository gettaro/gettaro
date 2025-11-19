import { Outlet } from "react-router-dom";
import Navigation from "../components/Navigation";
import ApiErrorBanner from "../components/ApiErrorBanner";
import { useEffect, useRef } from "react";
import { useAuth } from "../hooks/useAuth";
import { useApiErrorStore } from "../stores/apiError";

export default function RootLayout() {
  const { getToken, isAuthenticated } = useAuth()
  const { isApiUnavailable } = useApiErrorStore()
  const tokenInitializedRef = useRef(false)

  useEffect(() => {
    // Only initialize token once when authenticated
    if (isAuthenticated && !tokenInitializedRef.current) {
      tokenInitializedRef.current = true
      console.log("initialising token...")
      getToken().catch((error) => {
        console.error('Failed to get token on mount:', error)
        tokenInitializedRef.current = false // Reset on error so we can retry
      })
    } else if (!isAuthenticated) {
      // Reset when user logs out
      tokenInitializedRef.current = false
    }
  }, [isAuthenticated, getToken])

  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-background/80">
      <ApiErrorBanner />
      <header className={`bg-card/50 backdrop-blur-sm border-b ${isApiUnavailable ? 'mt-14' : ''}`}>
        <div className="container">
          <div className="flex items-center justify-between py-4">
            <img 
              src="/logo/taro_light.png" 
              alt="Taro" 
              className="h-12 block dark:hidden"
            />
            <img 
              src="/logo/taro_dark.png" 
              alt="Taro" 
              className="h-12 hidden dark:block"
            />
            <Navigation />
          </div>
        </div>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
} 