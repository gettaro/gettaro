import { useState, useEffect } from "react";
import Api from "../api/api";
import { useAuth0 } from "@auth0/auth0-react";
import { Organization } from "../types/organization";

export function useUserOrganizations() {
  const { getAccessTokenSilently } = useAuth0();
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    async function fetchOrganizations() {
      try {
        const data = await Api.getOrganizations(getAccessTokenSilently);
        setOrganizations(data);
      } catch (err) {
        setError(err instanceof Error ? err : new Error("Failed to fetch organizations"));
      } finally {
        setIsLoading(false);
      }
    }

    fetchOrganizations();
  }, [getAccessTokenSilently]);

  return { organizations, isLoading, error };
} 