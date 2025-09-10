// Conversation types for the frontend

export interface Conversation {
  id: string;
  organization_id: string;
  template_id?: string;
  title: string;
  manager_member_id: string;
  direct_member_id: string;
  conversation_date?: string;
  status: 'draft' | 'completed';
  content?: ConversationContent;
  created_at: string;
  updated_at: string;
}

export interface ConversationContent {
  [key: string]: any; // Flexible JSON structure for template data
}

export interface ConversationWithDetails extends Conversation {
  template?: ConversationTemplate;
  manager?: OrganizationMember;
  direct_report?: OrganizationMember;
}

export interface ConversationTemplate {
  id: string;
  organization_id: string;
  name: string;
  description?: string;
  template_fields: TemplateField[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface TemplateField {
  id: string;
  label: string;
  type: string;
  required: boolean;
  placeholder?: string;
  order: number;
}

export interface OrganizationMember {
  id: string;
  user_id: string;
  email: string;
  username: string;
  organization_id: string;
  is_owner: boolean;
  title_id?: string;
  manager_id?: string;
  created_at: string;
  updated_at: string;
}

// API Request/Response types
export interface CreateConversationRequest {
  template_id?: string;
  title: string;
  direct_member_id: string;
  conversation_date?: string; // ISO date string
  content?: ConversationContent;
}

export interface UpdateConversationRequest {
  conversation_date?: string; // ISO date string
  status?: 'draft' | 'completed';
  content?: ConversationContent;
}

export interface ListConversationsQuery {
  manager_member_id?: string;
  direct_member_id?: string;
  template_id?: string;
  status?: string;
  limit?: number;
  offset?: number;
}

export interface ListConversationsResponse {
  conversations: Conversation[];
}

export interface GetConversationResponse {
  conversation: Conversation;
}

export interface GetConversationWithDetailsResponse {
  conversation: ConversationWithDetails;
}

export interface ConversationStatsResponse {
  stats: Record<string, number>;
}
