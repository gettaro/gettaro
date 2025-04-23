import { CreateOrganizationForm } from "../components/forms/CreateOrganizationForm";
import Navigation from "../components/Navigation";

export default function NoOrganization() {
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
        <div className="container mx-auto px-4 py-16">
          <div className="max-w-4xl mx-auto text-center">
            <h2 className="text-5xl font-bold tracking-tight mb-6 bg-clip-text text-transparent bg-gradient-to-r from-primary to-primary/60">
              Welcome to EMS.dev! 🎉
            </h2>
            <p className="text-xl text-muted-foreground mb-12">
              EMS.dev streamlines your 1:1s by automatically tracking and aggregating your team's engineering metrics.
            </p>
            <div className="bg-card/50 backdrop-blur-sm rounded-lg border p-8 max-w-md mx-auto">
              <h3 className="text-lg font-semibold mb-4">Let's get started</h3>
              <p className="text-sm text-muted-foreground mb-6">
                Create your first organization to start having more effective conversations with your engineers.
              </p>
              <CreateOrganizationForm />
            </div>
          </div>
        </div>
      </main>
    </div>
  );
} 