import { createBrowserRouter } from "react-router-dom";
import RootLayout from "./layouts/RootLayout";
import Home from "./pages/Home";
import Dashboard from "./pages/Dashboard";
import Integrations from './pages/Integrations'
import Settings from './pages/Settings'
import Members from './pages/Members'

export const router = createBrowserRouter([
  {
    path: "/",
    element: <RootLayout />,
    children: [
      {
        path: "/",
        element: <Home />,
      },
      {
        path: "dashboard",
        element: <Dashboard />,
      },
      {
        path: "settings",
        element: <Settings />,
      },
      {
        path: "settings/integrations",
        element: <Integrations />,
      },
      {
        path: "settings/members",
        element: <Members />,
      }
    ],
  },
]); 