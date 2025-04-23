import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getOrganizations } from "../api/organizations";
import { useAuth0 } from "@auth0/auth0-react";
import { CreateOrganizationForm } from "../components/forms/CreateOrganizationForm";

export default function NoOrganization() {
  const navigate = useNavigate();
  const { getAccessTokenSilently, isAuthenticated, isLoading: isAuthLoading } = useAuth0();
  const [isChecking, setIsChecking] = useState(true);
  const [hasOrganizations, setHasOrganizations] = useState(false);

  useEffect(() => {
    let mounted = true;

    async function checkOrganizations() {
      if (isAuthLoading || !isAuthenticated) return;

      try {
        const token = await getAccessTokenSilently();
        const organizations = await getOrganizations(() => Promise.resolve(token));
        
        if (!mounted) return;

        if (organizations.length > 0) {
          // Sort by createdAt and get the most recent
          const sortedOrgs = [...organizations].sort((a, b) => 
            new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
          );
          setHasOrganizations(true);
          navigate(`/organizations/${sortedOrgs[0].id}/dashboard`, { replace: true });
        } else {
          setHasOrganizations(false);
        }
      } catch (error) {
        console.error("Error checking organizations:", error);
        setHasOrganizations(false);
      } finally {
        if (mounted) {
          setIsChecking(false);
        }
      }
    }

    checkOrganizations();

    return () => {
      mounted = false;
    };
  }, [navigate, getAccessTokenSilently, isAuthenticated, isAuthLoading]);

  if (isChecking || isAuthLoading) {
    return <></>;
  }

  if (hasOrganizations) {
    return <div>Redirecting to organization dashboard...</div>;
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
          <h3 className="text-lg font-semibold mb-4">Let's get started</h3>
          <p className="text-sm text-muted-foreground mb-6">
            Create your first organization to start having more effective conversations with your engineers.
          </p>
          <CreateOrganizationForm />
        </div>
      </div>
    </div>
  );
} 