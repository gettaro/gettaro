import { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './hooks/useAuth';
import { useUser } from './hooks/useUser';
import { useOrganizations } from './hooks/useOrganizations';
import OrganizationDashboard from './pages/OrganizationDashboard';
import TeamOverview from './pages/TeamOverview';
import Login from './pages/Login';
import Loading from './components/Loading';

function App() {
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const { user, isLoading: isUserLoading, fetchUser } = useUser();
  const { organizations, isLoading: isOrgsLoading, fetchOrganizations, setCurrentOrganization } = useOrganizations();

  useEffect(() => {
    if (isAuthenticated) {
      fetchUser();
    }
  }, [isAuthenticated, fetchUser]);

  useEffect(() => {
    if (user) {
      fetchOrganizations();
    }
  }, [user, fetchOrganizations]);

  useEffect(() => {
    if (organizations.length > 0) {
      setCurrentOrganization(organizations[0]);
    }
  }, [organizations, setCurrentOrganization]);

  if (isAuthLoading || isUserLoading || isOrgsLoading) {
    return <Loading />;
  }

  if (!isAuthenticated) {
    return <Login />;
  }

  return (
    <Router>
      <Routes>
        <Route path="/" element={<OrganizationDashboard />} />
        <Route path="/teams/:id" element={<TeamOverview />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

export default App; 