import { useState, useEffect } from 'react';
import { useAuth0 } from '@auth0/auth0-react';

export const useAuth = () => {
  const { isAuthenticated, isLoading, getIdTokenClaims } = useAuth0();
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const getToken = async () => {
      if (isAuthenticated) {
        try {
          const idTokenClaims = await getIdTokenClaims();
          if (idTokenClaims) {
            setToken(idTokenClaims.__raw);
          }
        } catch (error) {
          console.error('Error getting ID token:', error);
        }
      }
    };

    getToken();
  }, [isAuthenticated, getIdTokenClaims]);

  return {
    isAuthenticated,
    isLoading,
    token,
  };
}; 