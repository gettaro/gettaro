import { useAuth } from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { CreateOrganizationForm } from "../components/forms/CreateOrganizationForm";
import { useToast } from "../hooks/useToast";
import { useEffect } from "react";
import Api from "../api/api";

export default function Dashboard() {
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const navigate = useNavigate();
  const { toast } = useToast();
  const queryClient = useQueryClient();

 

  // Create organization mutation
  const createOrganizationMutation = useMutation({
    mutationFn: async (data: { name: string; slug: string }) => {
      return await Api.createOrganization(data.name, data.slug);
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