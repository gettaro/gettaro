// AI Query Request
export interface AIQueryRequest {
  entity_type: 'member' | 'team' | 'organization' | 'project';
  entity_id: string;
  query: string;
  context?: string;
  additional_data?: Record<string, any>;
}

// AI Query Response
export interface AIQueryResponse {
  answer: string;
  sources: string[];
  confidence: number;
  related_data?: Record<string, any>;
  suggestions?: string[];
}

// AI Query History Item
export interface AIQueryHistoryItem {
  id: string;
  organization_id: string;
  user_id: string;
  entity_type: string;
  entity_id: string;
  query: string;
  answer: string;
  context: string;
  confidence: number;
  sources: string[];
  created_at: string;
}

// AI Query Stats
export interface AIQueryStats {
  total_queries: number;
  queries_by_entity: Record<string, number>;
  queries_by_context: Record<string, number>;
  average_confidence: number;
  recent_queries: AIQueryHistoryItem[];
}

// Chat Message
export interface ChatMessage {
  id: string;
  type: 'user' | 'ai';
  content: string;
  timestamp: Date;
  confidence?: number;
  sources?: string[];
  suggestions?: string[];
}

// Chat Context
export interface ChatContext {
  entityType: 'member' | 'team' | 'organization' | 'project';
  entityId: string;
  entityName: string;
  organizationId: string;
  context?: string;
}
