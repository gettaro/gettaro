import { useState } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { useNavigate } from "react-router-dom";
import Api from "../../api/api";
import { useAuth0 } from "@auth0/auth0-react";
import { OrganizationConflictError } from "../../api/errors/organizatinos";

export function CreateOrganizationForm() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: "",
    slug: "",
  });
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      const organization = await Api.createOrganization(formData.name, formData.slug);
      
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
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-2">
        <Label htmlFor="name">Organization Name</Label>
        <Input
          id="name"
          placeholder="Enter organization name"
          value={formData.name}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => 
            setFormData({ ...formData, name: e.target.value })
          }
          required
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="slug">Organization Slug</Label>
        <Input
          id="slug"
          placeholder="Enter organization slug"
          value={formData.slug}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => 
            setFormData({ ...formData, slug: e.target.value })
          }
          required
        />
      </div>
      {error && (
        <div className="text-sm text-destructive">
          {error}
        </div>
      )}
      <Button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Creating..." : "Create Organization"}
      </Button>
    </form>
  );
} 