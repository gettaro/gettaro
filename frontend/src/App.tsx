import { RouterProvider } from "react-router-dom";
import { router } from "./router";
import { Auth0Provider } from '@auth0/auth0-react'
import { Toaster } from "sonner";
import { useEffect } from "react";

// Initialize theme on app load
function ThemeInitializer() {
  useEffect(() => {
    // Check localStorage first
    const savedTheme = localStorage.getItem('theme')
    const root = document.documentElement
    
    if (savedTheme === 'dark') {
      root.classList.add('dark')
    } else if (savedTheme === 'light') {
      root.classList.remove('dark')
    } else {
      // Check system preference if no saved theme
      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        root.classList.add('dark')
      } else {
        root.classList.remove('dark')
      }
    }
  }, [])
  
  return null
}

export default function App() {
  return (
      <Auth0Provider
        domain={import.meta.env.VITE_AUTH0_DOMAIN}
        clientId={import.meta.env.VITE_AUTH0_CLIENT_ID}
        authorizationParams={{
          redirect_uri: window.location.origin,
          audience: import.meta.env.VITE_AUTH0_AUDIENCE,
          scope: 'openid profile email offline_access',
        }}
        cacheLocation="localstorage"
        useRefreshTokens={true}
      >
        <ThemeInitializer />
        <RouterProvider router={router} />
        <Toaster />
      </Auth0Provider>
  )
} 