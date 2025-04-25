import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import Api from "../api/api";
import { useAuth } from "../hooks/useAuth";
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card";
import { Organization } from "../types/organization";

export default function OrganizationDashboard() {
  const { id } = useParams<{ id: string }>();
  const { getAccessTokenSilently } = useAuth();
  const [organization, setOrganization] = useState<Organization | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchOrganization() {
      if (!id) return;
      
      try {
        const org = await Api.getOrganization(id);
        setOrganization(org);
      } catch (err) {
        setError("Failed to load organization");
        console.error("Error loading organization:", err);
      } finally {
        setIsLoading(false);
      }
    }

    fetchOrganization();
  }, [id, getAccessTokenSilently]);

  if (isLoading) {
    return null;
  }

  if (error) {
    return <div>Error loading organization</div>;
  }

  if (!organization) {
    return <div>Organization not found</div>;
  }

  return (
    <div className="container mx-auto py-8">
      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>{organization.name}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4">
              <div>
                <h3 className="text-lg font-medium">Organization Details</h3>
                <p className="text-muted-foreground">Slug: {organization.slug}</p>
                <p className="text-muted-foreground">Created: {new Date(organization.createdAt).toLocaleDateString()}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
} 