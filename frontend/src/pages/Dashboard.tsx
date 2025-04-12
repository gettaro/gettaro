import { useEffect, useState } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import { getTeams, getProjects, getTasks, getTeamMetrics } from '../services/api';

export default function Dashboard() {
  const { isAuthenticated } = useAuth0();
  const [teams, setTeams] = useState([]);
  const [projects, setProjects] = useState([]);
  const [tasks, setTasks] = useState([]);
  const [metrics, setMetrics] = useState([]);

  useEffect(() => {
    if (isAuthenticated) {
      Promise.all([
        getTeams(),
        getProjects(),
        getTasks(),
        getTeamMetrics(),
      ]).then(([teamsRes, projectsRes, tasksRes, metricsRes]) => {
        setTeams(teamsRes.data);
        setProjects(projectsRes.data);
        setTasks(tasksRes.data);
        setMetrics(metricsRes.data);
      });
    }
  }, [isAuthenticated]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-gray-900">Dashboard</h1>
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div className="card">
          <h3 className="text-lg font-medium text-gray-900">Teams</h3>
          <p className="mt-2 text-3xl font-semibold text-primary-600">{teams.length}</p>
        </div>

        <div className="card">
          <h3 className="text-lg font-medium text-gray-900">Projects</h3>
          <p className="mt-2 text-3xl font-semibold text-primary-600">{projects.length}</p>
        </div>

        <div className="card">
          <h3 className="text-lg font-medium text-gray-900">Tasks</h3>
          <p className="mt-2 text-3xl font-semibold text-primary-600">{tasks.length}</p>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-medium text-gray-900">Recent Activity</h3>
        <div className="mt-4">
          {metrics.length > 0 ? (
            <ul className="divide-y divide-gray-200">
              {metrics.slice(0, 5).map((metric: any) => (
                <li key={metric.id} className="py-4">
                  <div className="flex items-center space-x-4">
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm font-medium text-gray-900">
                        {metric.metric}
                      </p>
                      <p className="truncate text-sm text-gray-500">
                        {new Date(metric.timestamp).toLocaleDateString()}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm font-medium text-primary-600">
                        {metric.value}
                      </p>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-gray-500">No recent activity</p>
          )}
        </div>
      </div>
    </div>
  );
} 