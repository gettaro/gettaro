export interface Member {
  id: string
  userId: string
  email: string
  username: string
  organizationId: string
  isOwner: boolean
  titleId?: string
  createdAt: string
  updatedAt: string
}

export interface ListOrganizationMembersResponse {
  members: Member[]
}

export interface AddMemberRequest {
  email: string
  username: string
  titleId: string
  sourceControlAccountId: string
}

export interface UpdateMemberRequest {
  username: string
  titleId: string
  sourceControlAccountId: string
} 