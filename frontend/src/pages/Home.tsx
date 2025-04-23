import { Link } from "react-router-dom";
import { useAuth0 } from "@auth0/auth0-react";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function Home() {
  const { isAuthenticated, isLoading } = useAuth0();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      navigate("/dashboard", { replace: true });
    }
  }, [isAuthenticated, isLoading, navigate]);

  if (isLoading) {
    return null;
  }

  return (
    <div className="container mx-auto px-4 py-16">
      <div className="max-w-4xl mx-auto text-center">
        <h2 className="text-5xl font-bold tracking-tight mb-6 bg-clip-text text-transparent bg-gradient-to-r from-primary to-primary/60">
          Welcome to EMS.dev! 🎉
        </h2>
        <p className="text-xl text-muted-foreground mb-12">
          EMS.dev streamlines your 1:1s by automatically tracking and aggregating your team's engineering metrics.
        </p>
        <div className="bg-card/50 backdrop-blur-sm rounded-lg border p-8 max-w-md mx-auto">
          <h3 className="text-lg font-semibold mb-4">Get Started</h3>
          <p className="text-sm text-muted-foreground mb-6">
            Sign in to start having more effective conversations with your engineers.
          </p>
          <Link
            to="/dashboard"
            className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
          >
            Sign In
          </Link>
        </div>
      </div>
    </div>
  );
} 