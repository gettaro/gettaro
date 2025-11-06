import { useEffect, useState } from 'react'
import { IntegrationConfig } from '../types/integration'
import Api from '../api/api'
import { Button } from '../components/ui/button'
import { Card, CardContent } from '../components/ui/card'
import { Trash2, Github, Zap, Calendar, Code2 } from 'lucide-react'
import { toast } from 'sonner'
import NewIntegrationModal from '../components/CreateIntegrationModal'
import { useOrganizationStore } from '../stores/organization'

export default function Integrations() {
  const [integrations, setIntegrations] = useState<IntegrationConfig[]>([])
  const [loading, setLoading] = useState(true)
  const { currentOrganization } = useOrganizationStore()

  useEffect(() => {
    if (!currentOrganization?.id) return
    loadIntegrations()
  }, [currentOrganization])

  async function loadIntegrations() {
    if (!currentOrganization?.id) return
    try {
      const data = await Api.getOrganizationIntegrations(currentOrganization.id)
      setIntegrations(data)
    } catch (error) {
      toast.error('Failed to load integrations')
    } finally {
      setLoading(false)
    }
  }

  async function handleDelete(integrationId: string) {
    if (!currentOrganization?.id) return
    try {
      await Api.deleteIntegrationConfig(currentOrganization.id, integrationId)
      toast.success('Integration deleted')
      loadIntegrations()
    } catch (error) {
      toast.error('Failed to delete integration')
    }
  }

  const getProviderIcon = (providerName: string) => {
    switch (providerName.toLowerCase()) {
      case 'github':
        return <Github className="w-6 h-6" />
      case 'cursor':
        return <Zap className="w-6 h-6" />
      default:
        return <Code2 className="w-6 h-6" />
    }
  }

  const getProviderColor = (providerName: string) => {
    switch (providerName.toLowerCase()) {
      case 'github':
        return {
          bg: 'bg-gray-900 dark:bg-gray-100',
          text: 'text-white dark:text-gray-900',
          border: 'border-gray-800 dark:border-gray-200',
          accent: 'text-gray-600 dark:text-gray-400'
        }
      case 'cursor':
        return {
          bg: 'bg-purple-600 dark:bg-purple-400',
          text: 'text-white',
          border: 'border-purple-500 dark:border-purple-300',
          accent: 'text-purple-400 dark:text-purple-300'
        }
      default:
        return {
          bg: 'bg-primary/10',
          text: 'text-primary',
          border: 'border-primary/20',
          accent: 'text-muted-foreground'
        }
    }
  }

  const renderGitHubCard = (integration: IntegrationConfig) => {
    const colors = getProviderColor('github')
    const repositories = integration.metadata?.repositories 
      ? (typeof integration.metadata.repositories === 'string' 
          ? integration.metadata.repositories.split(',').map(r => r.trim())
          : [])
      : []

    return (
      <Card key={integration.id} className="border-border hover:border-primary/50 transition-colors">
        <CardContent className="p-6">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center space-x-3">
              <div className={`${colors.bg} ${colors.text} p-3 rounded-lg`}>
                {getProviderIcon('github')}
              </div>
              <div>
                <h3 className="text-lg font-semibold text-foreground">GitHub</h3>
                <p className="text-sm text-muted-foreground">Source Control</p>
              </div>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => handleDelete(integration.id)}
              className="text-muted-foreground hover:text-destructive"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>

          <div className="space-y-3">
            {repositories.length > 0 && (
              <div>
                <p className="text-xs font-medium text-muted-foreground mb-2">Repositories</p>
                <div className="flex flex-wrap gap-2">
                  {repositories.slice(0, 3).map((repo, idx) => (
                    <span
                      key={idx}
                      className="px-2 py-1 bg-muted/50 rounded-md text-xs text-foreground font-mono"
                    >
                      {repo}
                    </span>
                  ))}
                  {repositories.length > 3 && (
                    <span className="px-2 py-1 bg-muted/50 rounded-md text-xs text-muted-foreground">
                      +{repositories.length - 3} more
                    </span>
                  )}
                </div>
              </div>
            )}

            <div className="flex items-center space-x-4 text-xs text-muted-foreground">
              {integration.last_synced_at && (
                <div className="flex items-center space-x-1">
                  <Calendar className="w-3 h-3" />
                  <span>Synced {new Date(integration.last_synced_at).toLocaleDateString()}</span>
                </div>
              )}
              <div className="flex items-center space-x-1">
                <span>Created {new Date(integration.created_at).toLocaleDateString()}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    )
  }

  const renderCursorCard = (integration: IntegrationConfig) => {
    const colors = getProviderColor('cursor')

    return (
      <Card key={integration.id} className="border-border hover:border-primary/50 transition-colors">
        <CardContent className="p-6">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center space-x-3">
              <div className={`${colors.bg} ${colors.text} p-3 rounded-lg`}>
                {getProviderIcon('cursor')}
              </div>
              <div>
                <h3 className="text-lg font-semibold text-foreground">Cursor</h3>
                <p className="text-sm text-muted-foreground">AI Code Assistant</p>
              </div>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => handleDelete(integration.id)}
              className="text-muted-foreground hover:text-destructive"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>

          <div className="space-y-3">
            <div className="flex items-center space-x-4 text-xs text-muted-foreground">
              {integration.last_synced_at && (
                <div className="flex items-center space-x-1">
                  <Calendar className="w-3 h-3" />
                  <span>Synced {new Date(integration.last_synced_at).toLocaleDateString()}</span>
                </div>
              )}
              <div className="flex items-center space-x-1">
                <span>Created {new Date(integration.created_at).toLocaleDateString()}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    )
  }

  const renderDefaultCard = (integration: IntegrationConfig) => {
    const colors = getProviderColor(integration.provider_name)

    return (
      <Card key={integration.id} className="border-border hover:border-primary/50 transition-colors">
        <CardContent className="p-6">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center space-x-3">
              <div className={`${colors.bg} ${colors.text} p-3 rounded-lg`}>
                {getProviderIcon(integration.provider_name)}
              </div>
              <div>
                <h3 className="text-lg font-semibold text-foreground capitalize">
                  {integration.provider_name}
                </h3>
                <p className="text-sm text-muted-foreground">{integration.provider_type}</p>
              </div>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => handleDelete(integration.id)}
              className="text-muted-foreground hover:text-destructive"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>

          <div className="space-y-3">
            <div className="flex items-center space-x-4 text-xs text-muted-foreground">
              {integration.last_synced_at && (
                <div className="flex items-center space-x-1">
                  <Calendar className="w-3 h-3" />
                  <span>Synced {new Date(integration.last_synced_at).toLocaleDateString()}</span>
                </div>
              )}
              <div className="flex items-center space-x-1">
                <span>Created {new Date(integration.created_at).toLocaleDateString()}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    )
  }

  const renderIntegrationCard = (integration: IntegrationConfig) => {
    switch (integration.provider_name.toLowerCase()) {
      case 'github':
        return renderGitHubCard(integration)
      case 'cursor':
        return renderCursorCard(integration)
      default:
        return renderDefaultCard(integration)
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-7xl mx-auto">
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading integrations...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-foreground mb-2">Integrations</h1>
              <p className="text-muted-foreground">
                Connect and manage third-party services and APIs to enhance your organization's capabilities.
              </p>
            </div>
            {currentOrganization?.id && (
              <NewIntegrationModal 
                organizationId={currentOrganization.id} 
                onSuccess={loadIntegrations} 
              />
            )}
          </div>
        </div>

        {integrations.length === 0 ? (
          <div className="text-center py-12 bg-card rounded-lg border border-border">
            <Code2 className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
            <p className="text-muted-foreground mb-4">No integrations found</p>
            {currentOrganization?.id && (
              <NewIntegrationModal 
                organizationId={currentOrganization.id} 
                onSuccess={loadIntegrations} 
              />
            )}
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {integrations.map(renderIntegrationCard)}
          </div>
        )}
      </div>
    </div>
  )
} 