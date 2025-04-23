import { Outlet } from "react-router-dom";
import Navigation from "../components/Navigation";

export default function RootLayout() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-background/80">
      <header className="bg-card/50 backdrop-blur-sm border-b">
        <div className="container">
          <div className="flex items-center justify-between py-4">
            <h1 className="text-3xl font-bold text-foreground">EMS.dev</h1>
            <Navigation />
          </div>
        </div>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
} 