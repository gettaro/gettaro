import React from 'react'
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import Home from './pages/Home'
import Dashboard from './pages/Dashboard'

const App: React.FC = () => {
  return (
    <Router>
      <div className="min-h-screen bg-background">
        <header className="bg-card shadow-sm">
          <div className="container">
            <div className="flex items-center justify-between py-4">
              <h1 className="text-3xl font-bold text-foreground">EMS.dev</h1>
              <nav>
                <ul className="flex space-x-4">
                  <li>
                    <Link to="/" className="text-foreground hover:text-primary">Home</Link>
                  </li>
                  <li>
                    <Link to="/dashboard" className="text-foreground hover:text-primary">Dashboard</Link>
                  </li>
                </ul>
              </nav>
            </div>
          </div>
        </header>
        <main>
          <div className="container">
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/dashboard" element={<Dashboard />} />
            </Routes>
          </div>
        </main>
      </div>
    </Router>
  )
}

export default App 