import React, { useState, useEffect } from 'react';
import { Button } from './ui/button';
import { Card } from './ui/card';
import { Input } from './ui/input';
import { DateInput } from './ui/date-input';
import { Label } from './ui/label';
import { useToast } from '../hooks/useToast';
import Api from '../api/api';
import { ConversationWithDetails, UpdateConversationRequest, TemplateField } from '../types/conversation';

interface ConversationModalProps {
  conversation: ConversationWithDetails | null;
  isOpen: boolean;
  onClose: () => void;
  onUpdate: () => void;
}

export const ConversationModal: React.FC<ConversationModalProps> = ({
  conversation,
  isOpen,
  onClose,
  onUpdate
}) => {
  const [isEditing, setIsEditing] = useState(false);
  const [conversationDate, setConversationDate] = useState('');
  const [content, setContent] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    if (conversation) {
      setConversationDate(conversation.conversation_date || '');
      setContent(conversation.content || {});
      setIsEditing(false);
    }
  }, [conversation]);

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleCancel = () => {
    if (conversation) {
      setConversationDate(conversation.conversation_date || '');
      setContent(conversation.content || {});
    }
    setIsEditing(false);
  };

  const handleSave = async () => {
    if (!conversation) return;

    setLoading(true);
    try {
      const updateData: UpdateConversationRequest = {
        conversation_date: conversationDate || undefined,
        content: content
      };

      await Api.updateConversation(conversation.id, updateData);
      toast({
        title: 'Success',
        description: 'Conversation updated successfully',
      });
      setIsEditing(false);
      onUpdate();
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

  const handleFieldChange = (fieldId: string, value: any) => {
    setContent(prev => ({
      ...prev,
      [fieldId]: value
    }));
  };

  const renderField = (field: any) => {
    const value = content[field.id] || '';

    if (!isEditing) {
      return (
        <div className="p-3 bg-gray-50 rounded-md">
          <p className="text-sm text-gray-600">
            {value ? (
              typeof value === 'string' ? value : JSON.stringify(value)
            ) : (
              <span className="text-gray-400 italic">Not filled</span>
            )}
          </p>
        </div>
      );
    }

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
        return (
          <Input
            id={field.id}
            value={value}
            onChange={(e) => handleFieldChange(field.id, e.target.value)}
            placeholder={field.placeholder}
            required={field.required}
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

  if (!isOpen || !conversation) {
    return null;
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <Card className="w-full max-w-4xl max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <div>
              <h2 className="text-xl font-semibold">
                {conversation.template?.name || 'Untitled Conversation'}
              </h2>
              <p className="text-sm text-gray-600">
                {conversation.template?.description}
              </p>
            </div>
            <div className="flex items-center space-x-2">
              {conversation.status === 'draft' && !isEditing && (
                <Button variant="outline" onClick={handleEdit}>
                  Edit
                </Button>
              )}
              <Button variant="ghost" onClick={onClose}>
                ✕
              </Button>
            </div>
          </div>

          <div className="space-y-6">
            {/* Conversation Date */}
            <div>
              <Label htmlFor="conversationDate">Conversation Date</Label>
              {isEditing ? (
                <DateInput
                  id="conversationDate"
                  value={conversationDate}
                  onChange={(e) => setConversationDate(e.target.value)}
                />
              ) : (
                <div className="p-3 bg-gray-50 rounded-md">
                  <p className="text-sm text-gray-600">
                    {conversationDate ? new Date(conversationDate).toLocaleDateString() : 'Not scheduled'}
                  </p>
                </div>
              )}
            </div>

            {/* Status */}
            <div>
              <Label>Status</Label>
              <div className="p-3 bg-gray-50 rounded-md">
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                  conversation.status === 'completed'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-yellow-100 text-yellow-800'
                }`}>
                  {conversation.status}
                </span>
              </div>
            </div>

            {/* Template Fields */}
            {(() => {
              // Extract template fields from content metadata or use template fields
              const savedTemplateFields = conversation.content?._template_fields as TemplateField[] | undefined;
              const fields = savedTemplateFields || conversation.template?.template_fields;
              
              if (fields && fields.length > 0) {
                return (
                  <div className="space-y-4">
                    <h3 className="font-medium">Conversation Details:</h3>
                    {fields
                      .sort((a, b) => a.order - b.order)
                      .map((field) => (
                        <div key={field.id}>
                          <Label htmlFor={field.id}>
                            {field.label}
                            {field.required && <span className="text-red-500 ml-1">*</span>}
                          </Label>
                          {renderField(field)}
                        </div>
                      ))}
                  </div>
                );
              } else {
                return (
                  <div className="space-y-4">
                    <h3 className="font-medium">Conversation Content:</h3>
                    <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-md">
                      <p className="text-sm text-yellow-800">
                        No template fields defined for this conversation.
                      </p>
                    </div>
                  </div>
                );
              }
            })()}


            {/* Action Buttons */}
            {isEditing && (
              <div className="flex justify-end space-x-3 pt-4 border-t">
                <Button variant="outline" onClick={handleCancel} disabled={loading}>
                  Cancel
                </Button>
                <Button onClick={handleSave} disabled={loading}>
                  {loading ? 'Saving...' : 'Save Changes'}
                </Button>
              </div>
            )}
          </div>
        </div>
      </Card>
    </div>
  );
};
