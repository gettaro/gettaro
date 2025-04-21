import { useState, useCallback } from 'react';
import { useAuth } from './useAuth';

interface Team {
  id: string;
  name: string;
  description?: string;
  organizationId: string;
  members?: TeamMember[];
}

interface TeamMember {
  id: string;
  userId: string;
  teamId: string;
  role: string;
}

interface CreateTeamRequest {
  name: string;
  description?: string;
  organizationId: string;
}

export const useTeams = () => {
  const [teams, setTeams] = useState<Team[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const { token } = useAuth();

  const fetchTeams = useCallback(async (organizationId?: string) => {
    if (!token || !organizationId) return;

    setIsLoading(true);
    try {
      const url = new URL('http://localhost:8080/api/teams');
      url.searchParams.append('organizationId', organizationId);

      const response = await fetch(url.toString(), {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch teams');
      }

      const data = await response.json();
      setTeams(data.teams || []);
    } catch (error) {
      console.error('Error fetching teams:', error);
    } finally {
      setIsLoading(false);
    }
  }, [token]);

  const createTeam = useCallback(async (team: CreateTeamRequest) => {
    if (!token) return;

    try {
      const response = await fetch('http://localhost:8080/api/teams', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(team),
      });

      if (!response.ok) {
        throw new Error('Failed to create team');
      }

      const newTeam = await response.json();
      setTeams((prev) => [...prev, newTeam.team]);
    } catch (error) {
      console.error('Error creating team:', error);
      throw error;
    }
  }, [token]);

  const addTeamMember = useCallback(async (teamId: string, userId: string, role: string) => {
    if (!token) return;

    try {
      const response = await fetch(`http://localhost:8080/api/teams/${teamId}/members`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ userId, role }),
      });

      if (!response.ok) {
        throw new Error('Failed to add team member');
      }
    } catch (error) {
      console.error('Error adding team member:', error);
      throw error;
    }
  }, [token]);

  return {
    teams,
    isLoading,
    fetchTeams,
    createTeam,
    addTeamMember,
  };
}; 