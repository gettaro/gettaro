import React, { useState, useEffect } from 'react';
import { Button } from './ui/button';
import { Card } from './ui/card';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { useToast } from '../hooks/useToast';
import Api from '../api/api';
import { CreateConversationRequest, ConversationTemplate } from '../types/conversation';
import RatingSlider from './RatingSlider';

interface CreateConversationModalProps {
  organizationId: string;
  memberId: string;
  onClose: () => void;
  onSubmit: (data: CreateConversationRequest) => void;
}

export const CreateConversationModal: React.FC<CreateConversationModalProps> = ({
  organizationId,
  memberId,
  onClose,
  onSubmit
}) => {
  const [templates, setTemplates] = useState<ConversationTemplate[]>([]);
  const [selectedTemplate, setSelectedTemplate] = useState<ConversationTemplate | null>(null);
  const [title, setTitle] = useState('');
  const [conversationDate, setConversationDate] = useState('');
  const [content, setContent] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState(false);
  const [validationAttempted, setValidationAttempted] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    fetchTemplates();
  }, [organizationId]);

  const fetchTemplates = async () => {
    try {
      const response = await Api.getConversationTemplates(organizationId, { is_active: true });
      setTemplates(response.conversation_templates);
    } catch (error) {
      console.error('Error fetching templates:', error);
      toast({
        title: 'Error',
        description: 'Failed to fetch conversation templates',
        variant: 'destructive',
      });
    }
  };

  const handleTemplateSelect = (template: ConversationTemplate) => {
    setSelectedTemplate(template);
    // Set title to template name by default
    setTitle(template.name);
    // Initialize content with empty values for each field
    const initialContent: Record<string, any> = {};
    template.template_fields.forEach(field => {
      initialContent[field.id] = '';
    });
    setContent(initialContent);
    setValidationAttempted(false); // Reset validation state when template changes
  };

  const handleFieldChange = (fieldId: string, value: any) => {
    setContent(prev => ({
      ...prev,
      [fieldId]: value
    }));
  };

  const validateRequiredFields = (): { isValid: boolean; errors: string[] } => {
    const errors: string[] = [];
    
    // Check title (always required)
    if (!title.trim()) {
      errors.push('Title is required');
    }
    
    // Check template fields
    if (selectedTemplate?.template_fields) {
      selectedTemplate.template_fields.forEach(field => {
        if (field.required) {
          const value = content[field.id];
          if (!value || (typeof value === 'string' && !value.trim()) || value === '') {
            errors.push(`${field.label} is required`);
          }
        }
      });
    }
    
    return {
      isValid: errors.length === 0,
      errors
    };
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Set validation attempted flag
    setValidationAttempted(true);
    
    if (!selectedTemplate) {
      toast({
        title: 'Error',
        description: 'Please select a conversation template',
        variant: 'destructive',
      });
      return;
    }

    // Validate required fields
    const validation = validateRequiredFields();
    if (!validation.isValid) {
      toast({
        title: 'Validation Error',
        description: validation.errors.join(', '),
        variant: 'destructive',
      });
      return;
    }

    setLoading(true);
    try {
      const data: CreateConversationRequest = {
        template_id: selectedTemplate.id,
        title: title,
        direct_member_id: memberId,
        conversation_date: conversationDate || undefined,
        content: content
      };
      
      await onSubmit(data);
    } catch (error) {
      console.error('Error creating conversation:', error);
    } finally {
      setLoading(false);
    }
  };

  const renderField = (field: any) => {
    const value = content[field.id] || '';

    switch (field.type) {
      case 'text':
        return (
          <Input
            id={field.id}
            value={value}
            onChange={(e) => handleFieldChange(field.id, e.target.value)}
            placeholder={field.placeholder}
            required={field.required}
          />
        );
      case 'textarea':
        return (
          <textarea
            id={field.id}
            value={value}
            onChange={(e) => handleFieldChange(field.id, e.target.value)}
            placeholder={field.placeholder}
            required={field.required}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            rows={3}
          />
        );
      case 'select':
        // For now, treat select as text input
        // In a real implementation, you'd parse the options from the field definition
        return (
          <Input
            id={field.id}
            value={value}
            onChange={(e) => handleFieldChange(field.id, e.target.value)}
            placeholder={field.placeholder}
            required={field.required}
          />
        );
      case 'rating':
        return (
          <RatingSlider
            id={field.id}
            value={typeof value === 'number' ? value : 1}
            onChange={(rating) => handleFieldChange(field.id, rating)}
            min={1}
            max={5}
            step={1}
          />
        );
      default:
        return (
          <Input
            id={field.id}
            value={value}
            onChange={(e) => handleFieldChange(field.id, e.target.value)}
            placeholder={field.placeholder}
            required={field.required}
          />
        );
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <Card className="w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold">Create New Conversation</h2>
            <Button variant="ghost" onClick={onClose}>
              ✕
            </Button>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Template Selection */}
            <div>
              <Label htmlFor="template">Conversation Template</Label>
              <select
                id="template"
                value={selectedTemplate?.id || ''}
                onChange={(e) => {
                  const template = templates.find(t => t.id === e.target.value);
                  if (template) handleTemplateSelect(template);
                }}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              >
                <option value="">Select a template...</option>
                {templates?.map((template) => (
                  <option key={template.id} value={template.id}>
                    {template.name}
                  </option>
                ))}
              </select>
            </div>

            {/* Title */}
            <div>
              <Label htmlFor="title" className={validationAttempted && !title.trim() ? 'text-red-500' : ''}>
                Conversation Title <span className="text-red-500 ml-1">*</span>
              </Label>
              <div className={validationAttempted && !title.trim() ? 'ring-1 ring-red-500 rounded' : ''}>
                <Input
                  id="title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="Enter conversation title..."
                  required
                />
              </div>
              {validationAttempted && !title.trim() && (
                <p className="text-xs text-red-500 mt-1">Title is required</p>
              )}
            </div>

            {/* Conversation Date */}
            <div>
              <Label htmlFor="conversationDate">Conversation Date (Optional)</Label>
              <Input
                id="conversationDate"
                type="date"
                value={conversationDate}
                onChange={(e) => setConversationDate(e.target.value)}
              />
            </div>

            {/* Template Fields */}
            {selectedTemplate && (
              <div className="space-y-4">
                <h3 className="font-medium">Fill in the conversation details:</h3>
                {selectedTemplate.template_fields
                  ?.sort((a, b) => a.order - b.order)
                  .map((field) => {
                    const value = content[field.id];
                    const isEmpty = !value || (typeof value === 'string' && !value.trim()) || value === '';
                    const isRequiredAndEmpty = field.required && isEmpty;
                    const showValidation = validationAttempted && isRequiredAndEmpty;
                    
                    return (
                      <div key={field.id}>
                        <Label htmlFor={field.id} className={showValidation ? 'text-red-500' : ''}>
                          {field.label}
                          {field.required && <span className="text-red-500 ml-1">*</span>}
                        </Label>
                        <div className={showValidation ? 'ring-1 ring-red-500 rounded' : ''}>
                          {renderField(field)}
                        </div>
                        {showValidation && (
                          <p className="text-xs text-red-500 mt-1">This field is required</p>
                        )}
                      </div>
                    );
                  })}
              </div>
            )}

            {/* Submit Buttons */}
            <div className="flex justify-end space-x-3">
              <Button type="button" variant="outline" onClick={onClose}>
                Cancel
              </Button>
              <Button type="submit" disabled={loading}>
                {loading ? 'Creating...' : 'Create Conversation'}
              </Button>
            </div>
          </form>
        </div>
      </Card>
    </div>
  );
};
