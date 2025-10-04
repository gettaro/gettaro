import React, { useState, useEffect, useRef } from 'react';
import { Button } from './ui/button';
import { useToast } from '../hooks/useToast';
import Api from '../api/api';
import { ConversationWithDetails, TemplateField } from '../types/conversation';

interface ConversationsTabProps {
  organizationId: string;
  memberId: string;
  memberName: string;
  onViewConversation: (conversationId: string) => void;
  onCreateConversation: () => void;
}

export const ConversationsTab: React.FC<ConversationsTabProps> = ({ 
  organizationId, 
  memberId, 
  memberName,
  onViewConversation,
  onCreateConversation
}) => {
  const [conversations, setConversations] = useState<ConversationWithDetails[]>([]);
  const [loading, setLoading] = useState(true);
  const [expandedConversations, setExpandedConversations] = useState<Set<string>>(new Set());
  const { toast } = useToast();
  const fetchRef = useRef<boolean>(false);
  const loadedRef = useRef<string | null>(null);

  const fetchConversations = async (forceRefresh = false) => {
    const cacheKey = `${organizationId}-${memberId}`;
    
    if (fetchRef.current) {
      console.log('ConversationsTab: Fetch already in progress, skipping');
      return;
    }
    
    if (!forceRefresh && loadedRef.current === cacheKey) {
      console.log('ConversationsTab: Data already loaded for this member, skipping');
      return;
    }
    
    try {
      fetchRef.current = true;
      setLoading(true);
      console.log('ConversationsTab: Starting fetch for', { organizationId, memberId });
      const response = await Api.getConversations(organizationId, {
        direct_member_id: memberId
      });
      console.log('ConversationsTab: Fetch completed, got', response.conversations?.length || 0, 'conversations');
      setConversations(response.conversations);
      loadedRef.current = cacheKey;
    } catch (error) {
      console.error('Error fetching conversations:', error);
      toast({
        title: 'Error',
        description: 'Failed to fetch conversations',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
      fetchRef.current = false;
    }
  };

  useEffect(() => {
    if (organizationId && memberId) {
      console.log('ConversationsTab: Fetching conversations for', { organizationId, memberId });
      fetchConversations();
    }
  }, [organizationId, memberId]);


  const handleUpdateConversation = async (conversationId: string, status: 'draft' | 'completed') => {
    try {
      await Api.updateConversation(conversationId, { status });
      toast({
        title: 'Success',
        description: 'Conversation updated successfully',
      });
      fetchConversations(true);
    } catch (error) {
      console.error('Error updating conversation:', error);
      toast({
        title: 'Error',
        description: 'Failed to update conversation',
        variant: 'destructive',
      });
    }
  };

  const handleDeleteConversation = async (conversationId: string) => {
    if (!confirm('Are you sure you want to delete this conversation?')) {
      return;
    }

    try {
      await Api.deleteConversation(conversationId);
      toast({
        title: 'Success',
        description: 'Conversation deleted successfully',
      });
      fetchConversations(true);
    } catch (error) {
      console.error('Error deleting conversation:', error);
      toast({
        title: 'Error',
        description: 'Failed to delete conversation',
        variant: 'destructive',
      });
    }
  };

  const toggleConversationExpansion = (conversationId: string) => {
    setExpandedConversations(prev => {
      const newSet = new Set(prev);
      if (newSet.has(conversationId)) {
        newSet.delete(conversationId);
      } else {
        newSet.add(conversationId);
      }
      return newSet;
    });
  };

  const handleEditConversation = async (conversationId: string) => {
    onViewConversation(conversationId);
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Not scheduled';
    return new Date(dateString).toLocaleDateString();
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'draft':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const renderField = (field: any, content: Record<string, any>) => {
    const value = content[field.id] || '';
    
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
  };

  const renderExpandedContent = (conversation: ConversationWithDetails) => {
    // Extract template fields from content metadata or use template fields
    const savedTemplateFields = conversation.content?._template_fields as TemplateField[] | undefined;
    const fields = savedTemplateFields || conversation.template?.template_fields;
    
    if (!fields || fields.length === 0) {
      return (
        <div className="p-3 bg-muted/20 rounded border border-border/30">
          <p className="text-sm text-muted-foreground">No template fields defined for this conversation.</p>
        </div>
      );
    }

    return (
      <div className="space-y-3">
        <h4 className="font-medium text-sm text-muted-foreground">Conversation Details:</h4>
        {fields
          .sort((a, b) => a.order - b.order)
          .map((field) => {
            const notes = conversation.content?.[`${field.id}_notes`] || '';
            return (
              <div key={field.id} className="space-y-2">
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    {field.label}
                    {field.required && <span className="text-red-500 ml-1">*</span>}
                  </label>
                  {renderField(field, conversation.content || {})}
                </div>
                
                {/* Display notes if they exist */}
                {notes && (
                  <div>
                    <label className="block text-xs font-medium text-muted-foreground mb-1">
                      Notes:
                    </label>
                    <div className="p-2 bg-muted/10 rounded border border-border/20">
                      <p className="text-xs text-muted-foreground">{notes}</p>
                    </div>
                  </div>
                )}
              </div>
            );
          })}
      </div>
    );
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-gray-500">Loading conversations...</div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Conversations Table */}
      <div className="bg-card rounded-lg">
        <div className="p-4 border-b border-border/50">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-semibold text-foreground">Conversations with {memberName}</h3>
              <p className="text-sm text-muted-foreground mt-1">
                One-on-one conversations and feedback sessions
              </p>
            </div>
            <div className="flex items-center space-x-3">
              {loading && (
                <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                  <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <span>Loading...</span>
                </div>
              )}
              <Button onClick={onCreateConversation}>
                New Conversation
              </Button>
            </div>
          </div>
        </div>

        <div className="divide-y divide-border/50">
          {conversations.length === 0 ? (
            <div className="p-6 text-center">
              <div className="w-10 h-10 bg-muted rounded-full flex items-center justify-center mx-auto mb-3">
                <svg className="w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
              </div>
              <div className="text-muted-foreground mb-3 text-sm">No conversations yet</div>
              <Button onClick={onCreateConversation} size="sm">
                Start First Conversation
              </Button>
            </div>
          ) : (
            <div className="max-h-80 overflow-y-auto">
              {conversations.map((conversation) => {
                const isExpanded = expandedConversations.has(conversation.id);
                return (
                  <div key={conversation.id} className="p-4 hover:bg-muted/20 transition-colors">
                    <div className="flex items-start space-x-3">
                      <div className="flex-shrink-0 mt-1">
                        <div className="w-6 h-6 bg-primary/10 rounded-full flex items-center justify-center">
                          <svg className="w-3 h-3 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                          </svg>
                        </div>
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between">
                          <div className="flex-1 min-w-0">
                            <h4 className="text-base font-semibold text-foreground mb-1">
                              {conversation.title || 'Untitled Conversation'}
                            </h4>
                            {conversation.template?.description && (
                              <p className="text-sm text-muted-foreground mb-2 line-clamp-2">
                                {conversation.template.description}
                              </p>
                            )}
                            <div className="flex items-center space-x-3 text-xs text-muted-foreground">
                              <span className="flex items-center space-x-1">
                                <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                </svg>
                                <span>{formatDate(conversation.conversation_date)}</span>
                              </span>
                              <span className="flex items-center space-x-1">
                                <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                                <span>Created {new Date(conversation.created_at).toLocaleDateString()}</span>
                              </span>
                            </div>
                          </div>
                          <div className="flex items-center space-x-2 ml-3">
                            <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${getStatusColor(conversation.status)}`}>
                              {conversation.status}
                            </span>
                          </div>
                        </div>
                        
                        <div className="flex items-center justify-between mt-3">
                          <div className="flex items-center space-x-2">
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => toggleConversationExpansion(conversation.id)}
                              className="text-xs h-7"
                            >
                              {isExpanded ? 'Hide Details' : 'View Details'}
                            </Button>
                            {conversation.status === 'draft' && (
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => handleEditConversation(conversation.id)}
                                className="text-xs h-7"
                              >
                                Edit
                              </Button>
                            )}
                            {conversation.status === 'draft' && (
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => handleUpdateConversation(conversation.id, 'completed')}
                                className="text-xs h-7"
                              >
                                Mark Complete
                              </Button>
                            )}
                            {conversation.status === 'completed' && (
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => handleUpdateConversation(conversation.id, 'draft')}
                                className="text-xs h-7"
                              >
                                Reopen
                              </Button>
                            )}
                          </div>
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => handleDeleteConversation(conversation.id)}
                            className="text-destructive hover:text-destructive hover:bg-destructive/10 h-7 w-7 p-0"
                          >
                            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                          </Button>
                        </div>

                        {/* Expanded Content */}
                        {isExpanded && (
                          <div className="mt-4 pt-3 border-t border-border/50">
                            {renderExpandedContent(conversation)}
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>

    </div>
  );
};
