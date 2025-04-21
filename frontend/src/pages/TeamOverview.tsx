import React, { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useTeams } from '../hooks/useTeams';

export default function TeamOverview() {
  const { id } = useParams<{ id: string }>();
  const { teams, fetchTeams } = useTeams();

  useEffect(() => {
    if (id) {
      fetchTeams();
    }
  }, [id, fetchTeams]);

  const team = teams.find((t) => t.id === id);

  if (!team) {
    return <div>Team not found</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-4">{team.name}</h1>
        {team.description && <p className="text-gray-600">{team.description}</p>}
      </div>

      <div>
        <h2 className="text-2xl font-semibold mb-4">Team Members</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {team.members?.map((member) => (
            <div
              key={member.id}
              className="p-4 border rounded hover:shadow-md transition-shadow"
            >
              <h3 className="text-xl font-semibold">User ID: {member.userId}</h3>
              <p className="text-gray-600">Role: {member.role}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
} 