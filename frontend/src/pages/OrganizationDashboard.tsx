import React, { useState, useEffect } from 'react';
import { useOrganizations } from '../hooks/useOrganizations';
import { useTeams } from '../hooks/useTeams';
import { Link } from 'react-router-dom';
import Loading from '../components/Loading';

export default function OrganizationDashboard() {
  const { currentOrganization, addOrganizationMember } = useOrganizations();
  const { teams, isLoading, fetchTeams, createTeam } = useTeams();
  const [newTeamName, setNewTeamName] = useState('');
  const [inviteEmail, setInviteEmail] = useState('');

  useEffect(() => {
    if (currentOrganization) {
      fetchTeams(currentOrganization.id);
    }
  }, [currentOrganization, fetchTeams]);

  const handleCreateTeam = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!currentOrganization || !newTeamName) return;

    try {
      await createTeam({
        name: newTeamName,
        organizationId: currentOrganization.id,
      });
      setNewTeamName('');
    } catch (error) {
      console.error('Error creating team:', error);
    }
  };

  const handleInviteMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!currentOrganization || !inviteEmail) return;

    try {
      await addOrganizationMember(currentOrganization.id, inviteEmail);
      setInviteEmail('');
    } catch (error) {
      console.error('Error inviting member:', error);
    }
  };

  if (!currentOrganization) {
    return <div>No organization selected</div>;
  }

  if (isLoading) {
    return <Loading />;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">{currentOrganization.name}</h1>
        <div className="space-x-4">
          <form onSubmit={handleInviteMember} className="inline-flex items-center space-x-2">
            <input
              type="email"
              value={inviteEmail}
              onChange={(e) => setInviteEmail(e.target.value)}
              placeholder="Email to invite"
              className="px-4 py-2 border rounded"
              required
            />
            <button
              type="submit"
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Invite Member
            </button>
          </form>
        </div>
      </div>

      <div className="mb-8">
        <h2 className="text-2xl font-semibold mb-4">Teams</h2>
        <form onSubmit={handleCreateTeam} className="mb-4">
          <div className="flex space-x-2">
            <input
              type="text"
              value={newTeamName}
              onChange={(e) => setNewTeamName(e.target.value)}
              placeholder="New team name"
              className="flex-1 px-4 py-2 border rounded"
              required
            />
            <button
              type="submit"
              className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
            >
              Create Team
            </button>
          </div>
        </form>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {teams.map((team) => (
            <Link
              key={team.id}
              to={`/teams/${team.id}`}
              className="p-4 border rounded hover:shadow-md transition-shadow"
            >
              <h3 className="text-xl font-semibold">{team.name}</h3>
              <p className="text-gray-600">{team.members?.length || 0} members</p>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
} 