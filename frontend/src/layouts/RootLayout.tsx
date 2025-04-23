import { Outlet, useNavigate, useLocation } from "react-router-dom";
import { useEffect, useState } from "react";
import { useUserOrganizations } from "../hooks/use-user-organizations";

export default function RootLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { data: organizations, isLoading } = useUserOrganizations();
  const [shouldRender, setShouldRender] = useState(false);

  useEffect(() => {
    if (isLoading) {
      setShouldRender(false);
      return;
    }

    // Don't redirect if we're already on the no-organization page
    if (location.pathname === "/no-organization") {
      setShouldRender(true);
      return;
    }

    if (!organizations || organizations.length === 0) {
      navigate("/no-organization", { replace: true });
      setShouldRender(false);
    } else {
      setShouldRender(true);
    }
  }, [organizations, isLoading, navigate, location.pathname]);

  if (!shouldRender) {
    return null;
  }

  return <Outlet />;
} 