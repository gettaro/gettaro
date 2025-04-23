import { createBrowserRouter } from "react-router-dom";
import RootLayout from "./layouts/RootLayout";
import NoOrganization from "./pages/NoOrganization";
import Dashboard from "./pages/Dashboard";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <RootLayout />,
    children: [
      {
        path: "no-organization",
        element: <NoOrganization />,
      },
      {
        path: "/",
        element: <Dashboard />,
      },
    ],
  },
]); 