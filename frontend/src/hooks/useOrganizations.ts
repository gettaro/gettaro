import { useState, useCallback } from 'react';
import { useAuth } from './useAuth';

interface Organization {
  id: string;
  name: string;
  slug: string;
  isOwner: boolean;
  createdAt: string;
  updatedAt: string;
}

export const useOrganizations = () => {
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [currentOrganization, setCurrentOrganization] = useState<Organization | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const { token } = useAuth();

  const fetchOrganizations = useCallback(async () => {
    if (!token) return;

    setIsLoading(true);
    try {
      const response = await fetch('http://localhost:8080/api/organizations', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch organizations');
      }

      const data = await response.json();
      setOrganizations(data.organizations);
    } catch (error) {
      console.error('Error fetching organizations:', error);
    } finally {
      setIsLoading(false);
    }
  }, [token]);

  const addOrganizationMember = useCallback(async (orgId: string, email: string) => {
    if (!token) return;

    try {
      const response = await fetch(`http://localhost:8080/api/organizations/${orgId}/members`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });

      if (!response.ok) {
        throw new Error('Failed to add organization member');
      }
    } catch (error) {
      console.error('Error adding organization member:', error);
      throw error;
    }
  }, [token]);

  return {
    organizations,
    currentOrganization,
    isLoading,
    fetchOrganizations,
    setCurrentOrganization,
    addOrganizationMember,
  };
}; 