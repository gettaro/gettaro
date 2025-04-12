import { useEffect, useState } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import { getProjects, createProject, updateProject, deleteProject } from '../services/api';

interface User {
  id: string;
  name: string;
  picture: string;
}

interface Team {
  id: string;
  name: string;
}

interface Project {
  id: string;
  name: string;
  description: string;
  team: Team;
  users?: User[];
  status: 'active' | 'completed' | 'archived';
}

export default function Projects() {
  const { isAuthenticated } = useAuth0();
  const [projects, setProjects] = useState<Project[]>([]);
  const [isCreating, setIsCreating] = useState(false);
  const [newProject, setNewProject] = useState({
    name: '',
    description: '',
    teamId: '',
    status: 'active' as const
  });

  useEffect(() => {
    if (isAuthenticated) {
      getProjects().then((res) => setProjects(res.data));
    }
  }, [isAuthenticated]);

  const handleCreateProject = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await createProject(newProject);
      setProjects([...projects, res.data]);
      setIsCreating(false);
      setNewProject({
        name: '',
        description: '',
        teamId: '',
        status: 'active'
      });
    } catch (error) {
      console.error('Error creating project:', error);
    }
  };

  const handleDeleteProject = async (id: string) => {
    try {
      await deleteProject(id);
      setProjects(projects.filter((project) => project.id !== id));
    } catch (error) {
      console.error('Error deleting project:', error);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-gray-900">Projects</h1>
        <button
          onClick={() => setIsCreating(true)}
          className="btn btn-primary"
        >
          Create Project
        </button>
      </div>

      {isCreating && (
        <div className="card">
          <h3 className="text-lg font-medium text-gray-900">Create New Project</h3>
          <form onSubmit={handleCreateProject} className="mt-4 space-y-4">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                Name
              </label>
              <input
                type="text"
                id="name"
                value={newProject.name}
                onChange={(e) => setNewProject({ ...newProject, name: e.target.value })}
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
                value={newProject.description}
                onChange={(e) => setNewProject({ ...newProject, description: e.target.value })}
                className="input"
                rows={3}
              />
            </div>
            <div>
              <label htmlFor="teamId" className="block text-sm font-medium text-gray-700">
                Team
              </label>
              <select
                id="teamId"
                value={newProject.teamId}
                onChange={(e) => setNewProject({ ...newProject, teamId: e.target.value })}
                className="input"
                required
              >
                <option value="">Select a team</option>
                {/* TODO: Add team options */}
              </select>
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
        {projects.map((project) => (
          <div key={project.id} className="card">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-medium text-gray-900">{project.name}</h3>
              <div className="flex items-center space-x-2">
                <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                  project.status === 'active' ? 'bg-green-100 text-green-800' :
                  project.status === 'completed' ? 'bg-blue-100 text-blue-800' :
                  'bg-gray-100 text-gray-800'
                }`}>
                  {project.status}
                </span>
                <button
                  onClick={() => handleDeleteProject(project.id)}
                  className="text-red-600 hover:text-red-800"
                >
                  Delete
                </button>
              </div>
            </div>
            <p className="mt-2 text-sm text-gray-500">{project.description}</p>
            <div className="mt-4">
              <h4 className="text-sm font-medium text-gray-700">Team</h4>
              <p className="mt-1 text-sm text-gray-600">{project.team.name}</p>
            </div>
            <div className="mt-4">
              <h4 className="text-sm font-medium text-gray-700">Members</h4>
              <ul className="mt-2 space-y-2">
                {project.users?.map((user) => (
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