import { useAuth } from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";
import { CreateOrganizationForm } from "../components/forms/CreateOrganizationForm";
import { useToast } from "../hooks/useToast";
import { useEffect } from "react";
import Api from "../api/api";

export default function Dashboard() {
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const navigate = useNavigate();

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
        <CreateOrganizationForm />
      </div>
    </div>
  );
} 