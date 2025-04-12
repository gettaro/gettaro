import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Auth0ProviderWrapper from './providers/Auth0Provider';
import MainLayout from './layouts/MainLayout';
import Dashboard from './pages/Dashboard';
import Teams from './pages/Teams';
import Projects from './pages/Projects';
import Tasks from './pages/Tasks';
import Metrics from './pages/Metrics';

function App() {
  return (
    <Auth0ProviderWrapper>
      <BrowserRouter>
        <Routes>
          <Route element={<MainLayout />}>
            <Route path="/" element={<Dashboard />} />
            <Route path="/teams" element={<Teams />} />
            <Route path="/projects" element={<Projects />} />
            <Route path="/tasks" element={<Tasks />} />
            <Route path="/metrics" element={<Metrics />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </Auth0ProviderWrapper>
  );
}

export default App; 