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
    <div className="container mx-auto py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Integrations</h1>
        {currentOrganization?.id && (
          <NewIntegrationModal 
            organizationId={currentOrganization.id} 
            onSuccess={loadIntegrations} 
          />
        )}
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {integrations.map((integration) => (
          <Card key={integration.id}>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <span className="capitalize">{integration.providerName}</span>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleDelete(integration.id)}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </CardTitle>
              <CardDescription>
                Created {new Date(integration.createdAt).toLocaleDateString()}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <div>
                  <span className="font-medium">Type:</span> {integration.providerType}
                </div>
                {integration.lastSyncedAt && (
                  <div>
                    <span className="font-medium">Last Synced:</span>{' '}
                    {new Date(integration.lastSyncedAt).toLocaleDateString()}
                  </div>
                )}
                {integration.metadata && (
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
  )
} 