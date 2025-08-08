export interface Title {
  id: string
  name: string
  organizationId: string
  createdAt: string
  updatedAt: string
}

export interface CreateTitleRequest {
  name: string
}

export interface UpdateTitleRequest {
  name: string
}

export interface AssignUserTitleRequest {
  userId: string
  titleId: string
} 