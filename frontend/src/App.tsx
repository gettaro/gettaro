import { RouterProvider } from "react-router-dom";
import { router } from "./router";
import { Auth0Provider } from '@auth0/auth0-react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from "sonner";
import { useEffect } from "react";
import { useAuth } from "./hooks/useAuth";
// Create a client
const queryClient = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Auth0Provider
        domain={import.meta.env.VITE_AUTH0_DOMAIN}
        clientId={import.meta.env.VITE_AUTH0_CLIENT_ID}
        authorizationParams={{
          redirect_uri: window.location.origin,
          audience: import.meta.env.VITE_AUTH0_AUDIENCE,
        }}
      >
        <RouterProvider router={router} />
        <Toaster />
      </Auth0Provider>
    </QueryClientProvider>
  )
} 