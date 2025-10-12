import { createBrowserRouter } from "react-router-dom";
import RootLayout from "./layouts/RootLayout";
import Home from "./pages/Home";
import Dashboard from "./pages/Dashboard";
import Integrations from './pages/Integrations'
import Settings from './pages/Settings'
import Members from './pages/Members'
import Titles from './pages/Titles'
import MemberProfile from './pages/MemberProfile'
import ConversationTemplates from './pages/ConversationTemplates'
import Teams from './pages/Teams'

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
        path: "settings/conversation-templates",
        element: <ConversationTemplates />,
      },
      {
        path: "settings/teams",
        element: <Teams />,
      },
      {
        path: "members/:memberId/profile",
        element: <MemberProfile />,
      }
    ],
  },
]); 