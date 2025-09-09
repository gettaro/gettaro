export interface TemplateField {
  id: string
  label: string
  type: 'text' | 'textarea' | 'select' | 'checkbox' | 'rating' | 'date' | 'number'
  required: boolean
  options?: string[]
  placeholder?: string
  order: number
}

export interface ConversationTemplate {
  id: string
  organization_id: string
  name: string
  description?: string
  template_fields: TemplateField[]
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateConversationTemplateRequest {
  name: string
  description?: string
  template_fields: TemplateField[]
  is_active?: boolean
}

export interface UpdateConversationTemplateRequest {
  name?: string
  description?: string
  template_fields?: TemplateField[]
  is_active?: boolean
}

export interface ListConversationTemplatesQuery {
  is_active?: boolean
}

export interface ListConversationTemplatesResponse {
  conversation_templates: ConversationTemplate[]
}

export interface GetConversationTemplateResponse {
  conversation_template: ConversationTemplate
}

export interface CreateConversationTemplateResponse {
  conversation_template: ConversationTemplate
}

export interface UpdateConversationTemplateResponse {
  conversation_template: ConversationTemplate
}
