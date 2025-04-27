import { useState } from 'react'
import { CreateIntegrationConfigRequest } from '../types/integration'
import Api from '../api/api'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import * as Select from '@radix-ui/react-select'
import { toast } from 'sonner'
import * as Dialog from '@radix-ui/react-dialog'
import { Plus } from 'lucide-react'

interface CreateIntegrationModalProps {
  organizationId: string
  onSuccess: () => void
}

export default function NewIntegrationModal({ organizationId, onSuccess }: CreateIntegrationModalProps) {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [repositories, setRepositories] = useState('')
  const [formData, setFormData] = useState<CreateIntegrationConfigRequest>({
    providerName: 'github',
    token: '',
    metadata: {},
  })

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    try {
      const submitData = { ...formData }
      if (formData.providerName === 'github') {
        submitData.metadata = { repositories }
      }
      await Api.createIntegrationConfig(organizationId, submitData)
      toast.success('Integration created')
      setOpen(false)
      onSuccess()
    } catch (error) {
      toast.error('Failed to create integration')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <Dialog.Trigger asChild>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Add Integration
        </Button>
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Content className="fixed left-[50%] top-[50%] max-h-[85vh] w-[90vw] max-w-[450px] translate-x-[-50%] translate-y-[-50%] rounded-md bg-white p-6 shadow-lg">
          <Dialog.Title className="text-lg font-medium">Add New Integration</Dialog.Title>
          <Dialog.Description className="mt-2 text-sm text-gray-500">
            Connect your organization with external services
          </Dialog.Description>
          <form onSubmit={handleSubmit} className="mt-4 space-y-4">
            <div className="space-y-2">
              <Label htmlFor="provider">Provider</Label>
              <Select.Root
                value={formData.providerName}
                onValueChange={(value: string) =>
                  setFormData({ ...formData, providerName: value as 'github' })
                }
              >
                <Select.Trigger className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
                  <Select.Value placeholder="Select a provider" />
                </Select.Trigger>
                <Select.Portal>
                  <Select.Content className="relative z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md">
                    <Select.Viewport className="p-1">
                      <Select.Item value="github" className="relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50">
                        <Select.ItemText>GitHub</Select.ItemText>
                      </Select.Item>
                    </Select.Viewport>
                  </Select.Content>
                </Select.Portal>
              </Select.Root>
            </div>

            <div className="space-y-2">
              <Label htmlFor="token">Access Token</Label>
              <Input
                id="token"
                type="password"
                value={formData.token}
                onChange={(e) =>
                  setFormData({ ...formData, token: e.target.value })
                }
                placeholder="Enter your access token"
                required
              />
            </div>

            {formData.providerName === 'github' && (
              <div className="space-y-2">
                <Label htmlFor="repositories">Repositories</Label>
                <Input
                  id="repositories"
                  value={repositories}
                  onChange={(e) => setRepositories(e.target.value)}
                  placeholder="owner/repo1,owner/repo2"
                  required
                />
                <p className="text-sm text-gray-500">
                  Enter repository names in format owner/repo, separated by commas
                </p>
              </div>
            )}

            <div className="flex justify-end space-x-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={loading}>
                {loading ? 'Creating...' : 'Create Integration'}
              </Button>
            </div>
          </form>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
} 