import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createOrganization } from "../../api/organizations";
import { useAuth0 } from "@auth0/auth0-react";
import { OrganizationConflictError } from "../../api/organizations";

export function CreateOrganizationForm() {
  const navigate = useNavigate();
  const { getAccessTokenSilently } = useAuth0();
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      const organization = await createOrganization(name, slug, getAccessTokenSilently);
      
      if (!organization) {
        console.error("Organization is null or undefined");
        throw new Error("Organization creation failed");
      }

      if (!organization.slug) {
        console.error("Organization response missing slug:", organization);
        throw new Error("Invalid organization response: missing slug");
      }

      // Navigate to the organization's dashboard
      navigate(`/organizations/${organization.id}/dashboard`);
    } catch (err) {
      console.error("Error creating organization:", err);
      if (err instanceof OrganizationConflictError) {
        setError("An organization with this name already exists. Please choose a different name.");
      } else {
        setError("Failed to create organization. Please try again.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="name" className="block text-sm font-medium text-foreground mb-1">
          Organization Name
        </label>
        <input
          type="text"
          id="name"
          value={name}
          onChange={(e) => {
            setName(e.target.value);
            setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "-"));
          }}
          className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground"
          required
        />
      </div>
      <div>
        <label htmlFor="slug" className="block text-sm font-medium text-foreground mb-1">
          Organization Slug
        </label>
        <input
          type="text"
          id="slug"
          value={slug}
          onChange={(e) => setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "-"))}
          className="w-full px-3 py-2 border border-border rounded-md bg-background text-foreground"
          required
        />
      </div>
      {error && (
        <div className="text-sm text-destructive">
          {error}
        </div>
      )}
      <button
        type="submit"
        disabled={isLoading}
        className="w-full px-4 py-2 text-sm bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50"
      >
        {isLoading ? "Creating..." : "Create Organization"}
      </button>
    </form>
  );
} 