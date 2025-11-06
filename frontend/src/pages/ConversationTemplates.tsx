import { useState, useEffect } from 'react'
import { useOrganizationStore } from '../stores/organization'
import Api from '../api/api'
import { ConversationTemplate, TemplateField } from '../types/conversationTemplate'
import { Button } from '../components/ui/button'
import { Card } from '../components/ui/card'
import { Input } from '../components/ui/input'
import { Label } from '../components/ui/label'
import { useToast } from '../hooks/useToast'

export default function ConversationTemplates() {
  const { currentOrganization } = useOrganizationStore()
  const { toast } = useToast()
  const [templates, setTemplates] = useState<ConversationTemplate[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editingTemplate, setEditingTemplate] = useState<ConversationTemplate | null>(null)

  // Form state
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    is_active: true
  })
  const [templateFields, setTemplateFields] = useState<TemplateField[]>([])

  useEffect(() => {
    if (currentOrganization) {
      fetchTemplates()
    }
  }, [currentOrganization])

  const fetchTemplates = async () => {
    if (!currentOrganization) return

    try {
      setLoading(true)
      const response = await Api.getConversationTemplates(currentOrganization.id)
      setTemplates(response.conversation_templates || [])
    } catch (error) {
      console.error('Error fetching conversation templates:', error)
      toast({
        title: 'Error',
        description: 'Failed to fetch conversation templates',
        variant: 'destructive'
      })
      setTemplates([]) // Ensure templates is always an array
    } finally {
      setLoading(false)
    }
  }

  const handleCreateTemplate = async () => {
    if (!currentOrganization) return

    try {
      const response = await Api.createConversationTemplate(currentOrganization.id, {
        name: formData.name,
        description: formData.description || undefined,
        template_fields: templateFields,
        is_active: formData.is_active
      })

      setTemplates([response.conversation_template, ...(templates || [])])
      resetForm()
      toast({
        title: 'Success',
        description: 'Conversation template created successfully'
      })
    } catch (error) {
      console.error('Error creating conversation template:', error)
      toast({
        title: 'Error',
        description: 'Failed to create conversation template',
        variant: 'destructive'
      })
    }
  }

  const handleUpdateTemplate = async () => {
    if (!editingTemplate) return

    try {
      const response = await Api.updateConversationTemplate(editingTemplate.id, {
        name: formData.name,
        description: formData.description || undefined,
        template_fields: templateFields,
        is_active: formData.is_active
      })

      setTemplates((templates || []).map(t => t.id === editingTemplate.id ? response.conversation_template : t))
      resetForm()
      toast({
        title: 'Success',
        description: 'Conversation template updated successfully'
      })
    } catch (error) {
      console.error('Error updating conversation template:', error)
      toast({
        title: 'Error',
        description: 'Failed to update conversation template',
        variant: 'destructive'
      })
    }
  }

  const handleDeleteTemplate = async (templateId: string) => {
    try {
      await Api.deleteConversationTemplate(templateId)
      setTemplates((templates || []).filter(t => t.id !== templateId))
      toast({
        title: 'Success',
        description: 'Conversation template deleted successfully'
      })
    } catch (error) {
      console.error('Error deleting conversation template:', error)
      toast({
        title: 'Error',
        description: 'Failed to delete conversation template',
        variant: 'destructive'
      })
    }
  }

  const resetForm = () => {
    setFormData({ name: '', description: '', is_active: true })
    setTemplateFields([])
    setShowCreateForm(false)
    setEditingTemplate(null)
  }

  const startEdit = (template: ConversationTemplate) => {
    setEditingTemplate(template)
    setFormData({
      name: template.name,
      description: template.description || '',
      is_active: template.is_active
    })
    setTemplateFields([...template.template_fields])
    setShowCreateForm(true)
  }

  const addTemplateField = () => {
    const newField: TemplateField = {
      id: `field_${Date.now()}`,
      label: '',
      type: 'text',
      required: false,
      order: templateFields.length
    }
    setTemplateFields([...templateFields, newField])
  }

  const updateTemplateField = (index: number, field: Partial<TemplateField>) => {
    const updatedFields = [...templateFields]
    updatedFields[index] = { ...updatedFields[index], ...field }
    setTemplateFields(updatedFields)
  }

  const removeTemplateField = (index: number) => {
    setTemplateFields(templateFields.filter((_, i) => i !== index))
  }

  if (!currentOrganization) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-7xl mx-auto">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-foreground mb-2">Conversation Templates</h1>
          </div>
          <div className="bg-card rounded-lg p-6">
            <p className="text-muted-foreground">Please select an organization to manage conversation templates.</p>
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
              <h1 className="text-3xl font-bold text-foreground mb-2">Conversation Templates</h1>
              <p className="text-muted-foreground">
                Create and manage templates for performance conversations and 1:1s to standardize your team's feedback process.
              </p>
            </div>
            <Button onClick={() => setShowCreateForm(true)}>
              Create Template
            </Button>
          </div>
        </div>

        {showCreateForm && (
          <Card className="p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">
            {editingTemplate ? 'Edit Template' : 'Create New Template'}
          </h2>
          
          <div className="space-y-4">
            <div>
              <Label htmlFor="name">Template Name</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="e.g., Weekly 1:1 Conversation"
              />
            </div>

            <div>
              <Label htmlFor="description">Description</Label>
              <Input
                id="description"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="Optional description"
              />
            </div>

            <div className="flex items-center space-x-2">
              <input
                type="checkbox"
                id="isActive"
                checked={formData.is_active}
                onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
              />
              <Label htmlFor="isActive">Active</Label>
            </div>

            <div className="bg-muted/30 rounded-lg p-4 border-2 border-primary/20">
              <div className="flex justify-between items-center mb-4">
                <div>
                  <Label className="text-base font-semibold">Template Fields</Label>
                  <p className="text-xs text-muted-foreground mt-1">Configure the fields for this conversation template</p>
                </div>
                <Button type="button" onClick={addTemplateField}>
                  Add Field
                </Button>
              </div>

              {templateFields.length === 0 ? (
                <div className="text-center py-6 border-2 border-dashed border-border rounded-lg">
                  <p className="text-sm text-muted-foreground">No fields added yet. Click "Add Field" to get started.</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {templateFields.map((field, index) => (
                    <Card key={field.id} className="p-4 bg-card border-2 border-border">
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <Label>Field Label</Label>
                          <Input
                            value={field.label}
                            onChange={(e) => updateTemplateField(index, { label: e.target.value })}
                            placeholder="e.g., What went well this week?"
                          />
                        </div>
                        <div>
                          <Label>Field Type</Label>
                          <select
                            className="w-full p-2 border rounded bg-background"
                            value={field.type}
                            onChange={(e) => updateTemplateField(index, { type: e.target.value as any })}
                          >
                            <option value="text">Text</option>
                            <option value="textarea">Textarea</option>
                            <option value="select">Select</option>
                            <option value="checkbox">Checkbox</option>
                            <option value="rating">Rating</option>
                            <option value="date">Date</option>
                            <option value="number">Number</option>
                          </select>
                        </div>
                        <div>
                          <Label>Placeholder</Label>
                          <Input
                            value={field.placeholder || ''}
                            onChange={(e) => updateTemplateField(index, { placeholder: e.target.value })}
                            placeholder="Optional placeholder text"
                          />
                        </div>
                        <div className="flex items-center space-x-2">
                          <input
                            type="checkbox"
                            checked={field.required}
                            onChange={(e) => updateTemplateField(index, { required: e.target.checked })}
                          />
                          <Label>Required</Label>
                        </div>
                      </div>
                      <Button
                        type="button"
                        variant="destructive"
                        size="sm"
                        onClick={() => removeTemplateField(index)}
                        className="mt-3"
                      >
                        Remove Field
                      </Button>
                    </Card>
                  ))}
                </div>
              )}
            </div>

            <div className="flex space-x-2">
              <Button onClick={editingTemplate ? handleUpdateTemplate : handleCreateTemplate}>
                {editingTemplate ? 'Update Template' : 'Create Template'}
              </Button>
              <Button variant="outline" onClick={resetForm}>
                Cancel
              </Button>
            </div>
          </div>
        </Card>
      )}

        {loading ? (
          <div className="bg-card rounded-lg p-6 text-center">
            <p className="text-muted-foreground">Loading conversation templates...</p>
          </div>
        ) : (
          <div className="space-y-4">
            {!templates || templates.length === 0 ? (
              <Card className="p-6 text-center">
                <p className="text-muted-foreground">No conversation templates found.</p>
              </Card>
            ) : (
              (templates || []).filter(template => template).map((template) => (
                <Card key={template.id} className="p-6">
                  <div className="flex justify-between items-start">
                    <div>
                      <h3 className="text-lg font-semibold text-foreground">{template.name}</h3>
                      {template.description && (
                        <p className="text-muted-foreground mt-1">{template.description}</p>
                      )}
                      <div className="flex items-center space-x-4 mt-2">
                        <span className={`px-2 py-1 rounded text-sm ${
                          template.is_active 
                            ? 'bg-success/10 text-success dark:text-success' 
                            : 'bg-muted text-muted-foreground'
                        }`}>
                          {template.is_active ? 'Active' : 'Inactive'}
                        </span>
                        <span className="text-sm text-muted-foreground">
                          {template.template_fields.length} field{template.template_fields.length !== 1 ? 's' : ''}
                        </span>
                      </div>
                    </div>
                    <div className="flex space-x-2">
                      <Button
                        size="sm"
                        onClick={() => startEdit(template)}
                      >
                        Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDeleteTemplate(template.id)}
                      >
                        Delete
                      </Button>
                    </div>
                  </div>
                </Card>
              ))
            )}
          </div>
        )}
      </div>
    </div>
  )
}
