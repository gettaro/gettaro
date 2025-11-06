import { useEffect, useState } from 'react'
import { IntegrationConfig } from '../types/integration'
import Api from '../api/api'
import { Button } from '../components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'
import { Trash2 } from 'lucide-react'
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

  if (loading) {
    return <div>Loading...</div>
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

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {integrations.map((integration) => (
            <Card key={integration.id}>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span className="capitalize">{integration.provider_name}</span>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDelete(integration.id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </CardTitle>
                <CardDescription>
                  Created {new Date(integration.created_at).toLocaleDateString()}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div>
                    <span className="font-medium">Type:</span> {integration.provider_type}
                  </div>
                  {integration.last_synced_at && (
                    <div>
                      <span className="font-medium">Last Synced:</span>{' '}
                      {new Date(integration.last_synced_at).toLocaleDateString()}
                    </div>
                  )}
                  {integration.metadata && integration.metadata.repositories && integration.provider_type === 'SourceControl' && (
                    <div>
                      <span className="font-medium">Repositories:</span>{' '}
                      {integration.metadata.repositories}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {integrations.length === 0 && (
          <div className="text-center py-12">
            <p className="text-muted-foreground">No integrations found</p>
          </div>
        )}
      </div>
    </div>
  )
} 