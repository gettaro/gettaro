import { useState, useCallback } from 'react';
import { useAuth } from './useAuth';

interface User {
  id: string;
  name: string;
  email: string;
  is_active: boolean;
  status: string;
  titleId: string;
  organizationId: string;
}

export const useUser = () => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const { token } = useAuth();

  const fetchUser = useCallback(async () => {
    if (!token) return;

    setIsLoading(true);
    try {
      const response = await fetch('http://localhost:8080/api/users/me', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch user');
      }

      const data = await response.json();
      setUser(data);
    } catch (error) {
      console.error('Error fetching user:', error);
    } finally {
      setIsLoading(false);
    }
  }, [token]);

  return {
    user,
    isLoading,
    fetchUser,
  };
}; 