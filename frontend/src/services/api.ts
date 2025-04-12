import axios from 'axios';
import { useAuth0 } from '@auth0/auth0-react';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
});

// Add a request interceptor to add the auth token to requests
api.interceptors.request.use(async (config) => {
  const { getAccessTokenSilently } = useAuth0();
  const token = await getAccessTokenSilently();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Teams
export const getTeams = () => api.get('/teams');
export const getTeam = (id: string) => api.get(`/teams/${id}`);
export const createTeam = (data: any) => api.post('/teams', data);
export const updateTeam = (id: string, data: any) => api.put(`/teams/${id}`, data);
export const deleteTeam = (id: string) => api.delete(`/teams/${id}`);

// Projects
export const getProjects = () => api.get('/projects');
export const getProject = (id: string) => api.get(`/projects/${id}`);
export const createProject = (data: any) => api.post('/projects', data);
export const updateProject = (id: string, data: any) => api.put(`/projects/${id}`, data);
export const deleteProject = (id: string) => api.delete(`/projects/${id}`);

// Tasks
export const getTasks = () => api.get('/tasks');
export const getTask = (id: string) => api.get(`/tasks/${id}`);
export const createTask = (data: any) => api.post('/tasks', data);
export const updateTask = (id: string, data: any) => api.put(`/tasks/${id}`, data);
export const deleteTask = (id: string) => api.delete(`/tasks/${id}`);

// Work Logs
export const getWorkLogs = () => api.get('/work-logs');
export const createWorkLog = (data: any) => api.post('/work-logs', data);

// Metrics
export const getTeamMetrics = (params?: { team_id?: string; start_date?: string; end_date?: string }) =>
  api.get('/metrics/teams', { params });
export const getProjectMetrics = (params?: { project_id?: string; start_date?: string; end_date?: string }) =>
  api.get('/metrics/projects', { params });
export const getUserMetrics = (params?: { user_id?: string; start_date?: string; end_date?: string }) =>
  api.get('/metrics/users', { params });

// GenAI Usage
export const getGenAIUsage = (params?: { user_id?: string; start_date?: string; end_date?: string }) =>
  api.get('/genai-usage', { params });
export const createGenAIUsage = (data: any) => api.post('/genai-usage', data); 