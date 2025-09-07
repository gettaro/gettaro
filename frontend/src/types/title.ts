export interface Title {
  id: string
  name: string
  organizationId: string
  isManager: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateTitleRequest {
  name: string
  isManager: boolean
}

export interface UpdateTitleRequest {
  name: string
  isManager: boolean
}

export interface AssignUserTitleRequest {
  userId: string
  titleId: string
} 