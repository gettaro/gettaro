import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getOrganizations } from "../api/organizations";
import { useAuth0 } from "@auth0/auth0-react";

export default function OrganizationRedirect() {
  const navigate = useNavigate();
  const { getAccessTokenSilently, isAuthenticated, isLoading: isAuthLoading } = useAuth0();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    async function checkOrganizations() {
      // Wait for auth to be ready
      if (isAuthLoading) {
        return;
      }
      
      // If not authenticated, don't proceed
      if (!isAuthenticated) {
        return;
      }

      try {
        console.log("Checking organizations...");
        const token = await getAccessTokenSilently();
        console.log("Got token, fetching organizations...");
        const organizations = await getOrganizations(() => Promise.resolve(token));
        console.log("Organizations:", organizations);
        
        if (!mounted) return;

        if (organizations.length === 0) {
          console.log("No organizations found, redirecting to /no-organization");
          navigate("/no-organization", { replace: true });
        } else {
          console.log("Found organizations, redirecting to first organization's dashboard");
          navigate(`/organizations/${organizations[0].slug}/dashboard`, { replace: true });
        }
      } catch (error) {
        console.error("Error checking organizations:", error);
        if (!mounted) return;
        setError("Failed to load organizations");
        navigate("/no-organization", { replace: true });
      } finally {
        if (mounted) {
          setIsLoading(false);
        }
      }
    }

    checkOrganizations();

    return () => {
      mounted = false;
    };
  }, [navigate, getAccessTokenSilently, isAuthenticated, isAuthLoading]);

  if (isLoading || isAuthLoading) {
    return nil;
  }

  if (error) {
    return <div className="text-red-500">{error}</div>;
  }

  return null;
} 