import { createBrowserRouter } from "react-router-dom";
import RootLayout from "./layouts/RootLayout";
import Home from "./pages/Home";
import Dashboard from "./pages/Dashboard";
import Integrations from './pages/Integrations'
import Settings from './pages/Settings'
import Members from './pages/Members'
import Titles from './pages/Titles'
import MemberActivityPage from './pages/MemberActivity'
import MemberProfile from './pages/MemberProfile'

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
      },
      {
        path: "settings/titles",
        element: <Titles />,
      },
      {
        path: "members/:memberId/activity",
        element: <MemberActivityPage />,
      },
      {
        path: "members/:memberId/profile",
        element: <MemberProfile />,
      }
    ],
  },
]); 