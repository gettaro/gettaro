export interface Title {
  id: string
  name: string
  organization_id: string
  is_manager: boolean
  created_at: string
  updated_at: string
}

export interface CreateTitleRequest {
  name: string
  is_manager: boolean
}

export interface UpdateTitleRequest {
  name: string
  is_manager: boolean
}

export interface AssignUserTitleRequest {
  user_id: string
  title_id: string
} 