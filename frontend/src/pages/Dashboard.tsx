import { useAuth } from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";
import { CreateOrganizationForm } from "../components/forms/CreateOrganizationForm";
import { useEffect } from "react";
import { useOrganizationStore } from "../stores/organization";
import { useIntegrationStore } from "../stores/integration";
import { Link } from "react-router-dom";

export default function Dashboard() {
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const { currentOrganization } = useOrganizationStore();
  const { integrations, fetchIntegrations } = useIntegrationStore();
  const navigate = useNavigate();

  // Redirect to home if not authenticated
  useEffect(() => {
    if (!isAuthLoading && !isAuthenticated) {
      navigate("/", { replace: true });
    }
  }, [isAuthenticated, isAuthLoading, navigate]);

  // Fetch integrations when organization changes
  useEffect(() => {
    if (currentOrganization?.id) {
      fetchIntegrations(currentOrganization.id);
    }
  }, [currentOrganization?.id, fetchIntegrations]);

  if (currentOrganization) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-2xl font-bold mb-8">Welcome to {currentOrganization.name}</h2>
          
          {integrations.length === 0 && (
            <div className="bg-yellow-50 border-l-4 border-yellow-400 p-4 mb-8">
              <div className="flex">
                <div className="flex-shrink-0">
                  <svg className="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                  </svg>
                </div>
                <div className="ml-3">
                  <p className="text-sm text-yellow-700">
                    No source control integrations are set up. <Link to="/settings/integrations" className="font-medium underline text-yellow-700 hover:text-yellow-600">Configure integrations</Link> to start tracking your repositories.
                  </p>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    );
  }

  // Show create organization form if user has no organizations
  return (
    <div className="container mx-auto px-4 py-16">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-2xl font-bold mb-8">Create Your First Organization</h2>
        <CreateOrganizationForm />
      </div>
    </div>
  );
} 