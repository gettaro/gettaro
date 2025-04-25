import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import Api from "../api/api";
export default function OrganizationRedirect() {
  const navigate = useNavigate();
  const { getToken, isAuthenticated, isLoading: isAuthLoading } = useAuth();
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
        const organizations = await Api.getOrganizations();
        
        if (!mounted) return;

        if (organizations.length === 0) {
          navigate("/no-organization", { replace: true });
        } else {
          navigate(`/organizations/${organizations[0].slug}/dashboard`, { replace: true });
        }
      } catch (error) {
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
  }, [navigate, getToken, isAuthenticated, isAuthLoading]);

  if (isLoading || isAuthLoading) {
    return null;
  }

  if (error) {
    return <div className="text-red-500">{error}</div>;
  }

  return null;
} 