import { Link } from 'react-router-dom';
import { useAuth0 } from '@auth0/auth0-react';

function Home() {
  const { isAuthenticated, isLoading } = useAuth0();

  if (isLoading) {
    return <div className="text-center">Loading...</div>;
  }

  return (
    <div className="text-center">
      <h1 className="text-4xl font-bold mb-8">Welcome to EMS.dev</h1>
      <p className="text-xl mb-8">
        Your all-in-one solution for email marketing and automation
      </p>
      {!isAuthenticated ? (
        <button
          onClick={() => loginWithRedirect()}
          className="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600"
        >
          Get Started
        </button>
      ) : (
        <Link
          to="/dashboard"
          className="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600"
        >
          Go to Dashboard
        </Link>
      )}
    </div>
  );
}

export default Home; 