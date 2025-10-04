import React, { useState, useEffect } from 'react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { useToast } from '../hooks/useToast';
import Api from '../api/api';
import { ConversationWithDetails, UpdateConversationRequest, CreateConversationRequest, TemplateField, ConversationTemplate } from '../types/conversation';
import RatingSlider from './RatingSlider';

interface ConversationSidebarProps {
  conversation: ConversationWithDetails | null;
  isOpen: boolean;
  onClose: () => void;
  onUpdate: (conversation: ConversationWithDetails) => void;
  onCreate?: (conversation: ConversationWithDetails) => void;
  mode?: 'edit' | 'create';
  organizationId?: string;
  memberId?: string;
}

export const ConversationSidebar: React.FC<ConversationSidebarProps> = ({
  conversation,
  isOpen,
  onClose,
  onUpdate,
  onCreate,
  mode = 'edit',
  organizationId,
  memberId
}) => {
  const [isEditing, setIsEditing] = useState(false);
  const [conversationDate, setConversationDate] = useState('');
  const [content, setContent] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState(false);
  const [templates, setTemplates] = useState<ConversationTemplate[]>([]);
  const [selectedTemplate, setSelectedTemplate] = useState<ConversationTemplate | null>(null);
  const [title, setTitle] = useState('');
  const [validationAttempted, setValidationAttempted] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    if (conversation) {
      setIsEditing(false);
      // Format date for HTML date input (YYYY-MM-DD)
      const dateValue = conversation.conversation_date 
        ? new Date(conversation.conversation_date).toISOString().split('T')[0]
        : '';
      setConversationDate(dateValue);
      setContent(conversation.content || {});
      setTitle(conversation.title || '');
    } else if (mode === 'create') {
      setIsEditing(true);
      setConversationDate('');
      setContent({});
      setTitle('');
      setSelectedTemplate(null);
      setValidationAttempted(false); // Reset validation state when creating
    }
  }, [conversation, mode]);

  useEffect(() => {
    if (mode === 'create' && organizationId) {
      fetchTemplates();
    }
  }, [mode, organizationId]);

  const fetchTemplates = async () => {
    if (!organizationId) return;
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

  const validateRequiredFields = (): { isValid: boolean; errors: string[] } => {
    const errors: string[] = [];
    
    // Check title (always required)
    if (!title.trim()) {
      errors.push('Title is required');
    }
    
    // Get template fields for validation
    const savedTemplateFields = conversation?.content?._template_fields as TemplateField[] | undefined;
    const fields = savedTemplateFields || conversation?.template?.template_fields || selectedTemplate?.template_fields;
    
    if (fields) {
      fields.forEach(field => {
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

  const handleSave = async () => {
    // Set validation attempted flag
    setValidationAttempted(true);
    
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

    if (mode === 'create') {
      await handleCreate();
    } else {
      await handleUpdate();
    }
  };

  const handleCreate = async () => {
    if (!organizationId || !memberId || !selectedTemplate) return;

    setLoading(true);
    try {
      const createData: CreateConversationRequest = {
        template_id: selectedTemplate.id,
        title: title,
        direct_member_id: memberId,
        conversation_date: conversationDate ? new Date(conversationDate).toISOString().split('T')[0] : undefined,
        content: content
      };

      const response = await Api.createConversation(organizationId, createData);
      
      if (onCreate) {
        // Fetch the full conversation details
        const fullResponse = await Api.getConversationWithDetails(response.conversation.id);
        onCreate(fullResponse.conversation);
      }
      
      toast({
        title: 'Success',
        description: 'Conversation created successfully',
      });
    } catch (error) {
      console.error('Error creating conversation:', error);
      toast({
        title: 'Error',
        description: 'Failed to create conversation',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleUpdate = async () => {
    if (!conversation) return;

    setLoading(true);
    try {
      const updateData: UpdateConversationRequest = {
        conversation_date: conversationDate ? new Date(conversationDate).toISOString().split('T')[0] : undefined,
        content: content
      };

      await Api.updateConversation(conversation.id, updateData);
      
      // Fetch updated conversation
      const response = await Api.getConversationWithDetails(conversation.id);
      onUpdate(response.conversation);
      
      setIsEditing(false);
      toast({
        title: 'Success',
        description: 'Conversation updated successfully',
      });
    } catch (error) {
      console.error('Error updating conversation:', error);
      toast({
        title: 'Error',
        description: 'Failed to update conversation',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    if (conversation) {
      // Format date for HTML date input (YYYY-MM-DD)
      const dateValue = conversation.conversation_date 
        ? new Date(conversation.conversation_date).toISOString().split('T')[0]
        : '';
      setConversationDate(dateValue);
      setContent(conversation.content || {});
      setTitle(conversation.title || '');
    } else if (mode === 'create') {
      setConversationDate('');
      setContent({});
      setTitle('');
      setSelectedTemplate(null);
    }
    setIsEditing(false);
  };

  const handleTemplateSelect = (template: ConversationTemplate) => {
    setSelectedTemplate(template);
    setTitle(template.name);
    // Initialize content with empty values for each field
    const initialContent: Record<string, any> = {};
    template.template_fields.forEach(field => {
      initialContent[field.id] = '';
    });
    setContent(initialContent);
  };

  const handleFieldChange = (fieldId: string, value: any) => {
    setContent(prev => ({
      ...prev,
      [fieldId]: value
    }));
  };

  const renderField = (field: TemplateField) => {
    const value = content[field.id] || '';

    if (isEditing) {
      switch (field.type) {
        case 'text':
          return (
            <Input
              id={field.id}
              value={value}
              onChange={(e) => handleFieldChange(field.id, e.target.value)}
              placeholder={field.placeholder}
              required={field.required}
              className="text-sm"
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
              className="w-full px-3 py-2 border border-border/50 rounded bg-background text-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm"
              rows={3}
            />
          );
        case 'select':
          return (
            <Input
              id={field.id}
              value={value}
              onChange={(e) => handleFieldChange(field.id, e.target.value)}
              placeholder={field.placeholder}
              required={field.required}
              className="text-sm"
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
              className="text-sm"
            />
          );
      }
    } else {
      if (field.type === 'rating') {
        const ratingValue = typeof value === 'number' ? value : (typeof value === 'string' ? parseInt(value) : 0);
        const getRatingLabel = (rating: number) => {
          switch (rating) {
            case 1: return 'Poor'
            case 2: return 'Fair'
            case 3: return 'Good'
            case 4: return 'Very Good'
            case 5: return 'Excellent'
            default: return `${rating}/5`
          }
        };
        
        return (
          <div className="p-2 bg-muted/20 rounded border border-border/30">
            <div className="flex items-center space-x-2">
              <div className="flex space-x-1">
                {[1, 2, 3, 4, 5].map((star) => (
                  <span
                    key={star}
                    className={`text-lg ${
                      star <= ratingValue ? 'text-yellow-400' : 'text-muted-foreground/30'
                    }`}
                  >
                    ★
                  </span>
                ))}
              </div>
              <span className="text-sm text-muted-foreground">
                {getRatingLabel(ratingValue)} ({ratingValue}/5)
              </span>
            </div>
          </div>
        );
      }
      
      return (
        <div className="p-2 bg-muted/20 rounded border border-border/30">
          <p className="text-sm text-muted-foreground">
            {value ? (
              typeof value === 'string' ? value : JSON.stringify(value)
            ) : (
              <span className="text-muted-foreground/60 italic">Not filled</span>
            )}
          </p>
        </div>
      );
    }
  };

  if (!isOpen) {
    return null;
  }

  // For edit mode, we need a conversation
  if (mode === 'edit' && !conversation) {
    return null;
  }

  // Extract template fields from content metadata or use template fields
  const savedTemplateFields = conversation?.content?._template_fields as TemplateField[] | undefined;
  const fields = savedTemplateFields || conversation?.template?.template_fields || selectedTemplate?.template_fields;

  return (
    <div className="fixed inset-y-0 right-0 w-[28rem] bg-card border-l border-border/50 shadow-lg z-50 flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-border/50">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-foreground">
              {mode === 'create' ? 'Create New Conversation' : (conversation?.title || 'Untitled Conversation')}
            </h2>
            <p className="text-sm text-muted-foreground">
              {mode === 'create' ? 'Select a template to get started' : (conversation?.template?.name || 'No template')}
            </p>
          </div>
          <Button variant="ghost" onClick={onClose} className="h-8 w-8 p-0">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </Button>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {/* Template Selection (Create Mode) */}
        {mode === 'create' && (
          <div>
            <Label htmlFor="template" className="text-sm font-medium">
              Conversation Template
            </Label>
            <select
              id="template"
              value={selectedTemplate?.id || ''}
              onChange={(e) => {
                const template = templates.find(t => t.id === e.target.value);
                if (template) handleTemplateSelect(template);
              }}
              className="w-full px-3 py-2 border border-border/50 rounded bg-background text-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm mt-1"
              required
            >
              <option value="">Select a template...</option>
              {templates.map((template) => (
                <option key={template.id} value={template.id}>
                  {template.name}
                </option>
              ))}
            </select>
          </div>
        )}

        {/* Title */}
        <div>
          <Label htmlFor="title" className={`text-sm font-medium ${validationAttempted && !title.trim() ? 'text-red-500' : 'text-foreground'}`}>
            Conversation Title <span className="text-red-500 ml-1">*</span>
          </Label>
          {isEditing || mode === 'create' ? (
            <div className={validationAttempted && !title.trim() ? 'ring-1 ring-red-500 rounded' : ''}>
              <Input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Enter conversation title..."
                className="text-sm mt-1"
                required
              />
            </div>
          ) : (
            <div className="p-2 bg-muted/20 rounded border border-border/30 mt-1">
              <p className="text-sm text-muted-foreground">
                {conversation?.title || 'Untitled Conversation'}
              </p>
            </div>
          )}
          {validationAttempted && !title.trim() && (
            <p className="text-xs text-red-500 mt-1">Title is required</p>
          )}
        </div>

        {/* Conversation Date */}
        <div>
          <Label htmlFor="conversationDate" className="text-sm font-medium">
            Conversation Date
          </Label>
          {isEditing || mode === 'create' ? (
            <Input
              id="conversationDate"
              type="date"
              value={conversationDate}
              onChange={(e) => setConversationDate(e.target.value)}
              className="text-sm mt-1"
            />
          ) : (
            <div className="p-2 bg-muted/20 rounded border border-border/30 mt-1">
              <p className="text-sm text-muted-foreground">
                {conversationDate ? new Date(conversationDate).toLocaleDateString() : 'Not scheduled'}
              </p>
            </div>
          )}
        </div>

        {/* Status */}
        {mode === 'edit' && (
          <div>
            <Label className="text-sm font-medium">Status</Label>
            <div className="mt-1">
              <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
                conversation?.status === 'completed' 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-yellow-100 text-yellow-800'
              }`}>
                {conversation?.status}
              </span>
            </div>
          </div>
        )}

        {/* Template Fields */}
        {fields && fields.length > 0 ? (
          <div className="space-y-3">
            <h3 className="font-medium text-sm text-foreground">Conversation Details</h3>
            {fields
              .sort((a, b) => a.order - b.order)
              .map((field) => {
                const value = content[field.id];
                const isEmpty = !value || (typeof value === 'string' && !value.trim()) || value === '';
                const isRequiredAndEmpty = field.required && isEmpty;
                const showValidation = validationAttempted && isRequiredAndEmpty;
                
                return (
                  <div key={field.id}>
                    <Label htmlFor={field.id} className={`text-sm font-medium ${showValidation ? 'text-red-500' : 'text-foreground'}`}>
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
        ) : mode === 'create' ? (
          <div className="p-3 bg-muted/20 rounded border border-border/30">
            <p className="text-sm text-muted-foreground">Select a template to see conversation fields.</p>
          </div>
        ) : (
          <div className="p-3 bg-muted/20 rounded border border-border/30">
            <p className="text-sm text-muted-foreground">No template fields defined for this conversation.</p>
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="p-4 border-t border-border/50">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            {mode === 'create' ? (
              <>
                <Button
                  onClick={handleSave}
                  size="sm"
                  disabled={loading || !selectedTemplate}
                  className="text-xs"
                >
                  {loading ? 'Creating...' : 'Create Conversation'}
                </Button>
                <Button
                  onClick={handleCancel}
                  size="sm"
                  variant="outline"
                  className="text-xs"
                >
                  Cancel
                </Button>
              </>
            ) : !isEditing ? (
              <Button
                onClick={() => {
                  setIsEditing(true);
                  setValidationAttempted(false); // Reset validation state when starting to edit
                }}
                size="sm"
                className="text-xs"
              >
                Edit
              </Button>
            ) : (
              <>
                <Button
                  onClick={handleSave}
                  size="sm"
                  disabled={loading}
                  className="text-xs"
                >
                  {loading ? 'Saving...' : 'Save'}
                </Button>
                <Button
                  onClick={handleCancel}
                  size="sm"
                  variant="outline"
                  className="text-xs"
                >
                  Cancel
                </Button>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
