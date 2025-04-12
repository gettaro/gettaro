import { useEffect, useState } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import { getTeams, createTeam, updateTeam, deleteTeam } from '../services/api';

interface User {
  id: string;
  name: string;
  picture: string;
}

interface Team {
  id: string;
  name: string;
  description: string;
  users?: User[];
}

export default function Teams() {
  const { isAuthenticated } = useAuth0();
  const [teams, setTeams] = useState<Team[]>([]);
  const [isCreating, setIsCreating] = useState(false);
  const [newTeam, setNewTeam] = useState({ name: '', description: '' });

  useEffect(() => {
    if (isAuthenticated) {
      getTeams().then((res) => setTeams(res.data));
    }
  }, [isAuthenticated]);

  const handleCreateTeam = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await createTeam(newTeam);
      setTeams([...teams, res.data]);
      setIsCreating(false);
      setNewTeam({ name: '', description: '' });
    } catch (error) {
      console.error('Error creating team:', error);
    }
  };

  const handleDeleteTeam = async (id: string) => {
    try {
      await deleteTeam(id);
      setTeams(teams.filter((team) => team.id !== id));
    } catch (error) {
      console.error('Error deleting team:', error);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-gray-900">Teams</h1>
        <button
          onClick={() => setIsCreating(true)}
          className="btn btn-primary"
        >
          Create Team
        </button>
      </div>

      {isCreating && (
        <div className="card">
          <h3 className="text-lg font-medium text-gray-900">Create New Team</h3>
          <form onSubmit={handleCreateTeam} className="mt-4 space-y-4">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                Name
              </label>
              <input
                type="text"
                id="name"
                value={newTeam.name}
                onChange={(e) => setNewTeam({ ...newTeam, name: e.target.value })}
                className="input"
                required
              />
            </div>
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                Description
              </label>
              <textarea
                id="description"
                value={newTeam.description}
                onChange={(e) => setNewTeam({ ...newTeam, description: e.target.value })}
                className="input"
                rows={3}
              />
            </div>
            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={() => setIsCreating(false)}
                className="btn btn-secondary"
              >
                Cancel
              </button>
              <button type="submit" className="btn btn-primary">
                Create
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {teams.map((team) => (
          <div key={team.id} className="card">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-medium text-gray-900">{team.name}</h3>
              <button
                onClick={() => handleDeleteTeam(team.id)}
                className="text-red-600 hover:text-red-800"
              >
                Delete
              </button>
            </div>
            <p className="mt-2 text-sm text-gray-500">{team.description}</p>
            <div className="mt-4">
              <h4 className="text-sm font-medium text-gray-700">Members</h4>
              <ul className="mt-2 space-y-2">
                {team.users?.map((user) => (
                  <li key={user.id} className="flex items-center">
                    <img
                      src={user.picture}
                      alt={user.name}
                      className="h-6 w-6 rounded-full"
                    />
                    <span className="ml-2 text-sm text-gray-600">{user.name}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
} 