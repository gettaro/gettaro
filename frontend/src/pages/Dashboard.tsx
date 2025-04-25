import { useAuth0 } from "@auth0/auth0-react";
import { useNavigate } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getOrganizations, createOrganization } from "../api/organizations";
import { CreateOrganizationForm } from "../components/forms/CreateOrganizationForm";
import { useToast } from "../hooks/use-toast";
import { useEffect } from "react";

export default function Dashboard() {
  const { isAuthenticated, isLoading: isAuthLoading, getAccessTokenSilently } = useAuth0();
  const navigate = useNavigate();
  const { toast } = useToast();
  const queryClient = useQueryClient();

  // Fetch organizations
  const { data: organizations, isLoading: isOrgsLoading, error: orgsError } = useQuery({
    queryKey: ["organizations"],
    queryFn: async () => {
      return await getOrganizations(getAccessTokenSilently);
    },
    enabled: isAuthenticated,
  });

  // Create organization mutation
  const createOrganizationMutation = useMutation({
    mutationFn: async (data: { name: string; slug: string }) => {
      return await createOrganization(data.name, data.slug, getAccessTokenSilently);
    },
    onSuccess: (data) => {
      // Invalidate and refetch organizations query
      queryClient.invalidateQueries({ queryKey: ["organizations"] });
      toast({
        title: "Success",
        description: "Organization created successfully",
      });
      // Navigate to the new organization's dashboard
      navigate(`/organizations/${data.id}/dashboard`);
    },
    onError: (error) => {
      toast({
        title: "Error",
        description: "Failed to create organization",
        variant: "destructive",
      });
    },
  });

  // Redirect to home if not authenticated
  useEffect(() => {
    if (!isAuthLoading && !isAuthenticated) {
      navigate("/", { replace: true });
    }
  }, [isAuthenticated, isAuthLoading, navigate]);

  // Show loading state while checking auth or fetching organizations
  if (isAuthLoading || isOrgsLoading) {
    return null;
  }

  // Show error state if organizations fetch failed
  if (orgsError) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-2xl font-bold mb-4">Error</h2>
          <p className="text-muted-foreground">
            Failed to load organizations. Please try again later.
          </p>
        </div>
      </div>
    );
  }

  // If user has organizations, redirect to the most recent one
  if (organizations && organizations.length > 0) {
    const mostRecentOrg = organizations[0];
    navigate(`/organizations/${mostRecentOrg.id}/dashboard`, { replace: true });
    return null;
  }

  // Show create organization form if user has no organizations
  return (
    <div className="container mx-auto px-4 py-16">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-2xl font-bold mb-8">Create Your First Organization</h2>
        <CreateOrganizationForm
          onSubmit={(data: { name: string; slug: string }) => createOrganizationMutation.mutate(data)}
          isLoading={createOrganizationMutation.isPending}
        />
      </div>
    </div>
  );
} 