export interface Member {
  id: string
  user_id: string
  email: string
  username: string
  organization_id: string
  is_owner: boolean
  title_id?: string
  manager_id?: string
  created_at: string
  updated_at: string
}

export interface ListOrganizationMembersResponse {
  members: Member[]
}

export interface AddMemberRequest {
  email: string
  username: string
  title_id: string
  external_account_id: string // Renamed from source_control_account_id
  manager_id?: string
}

export interface UpdateMemberRequest {
  username: string
  title_id: string
  external_account_id: string // Renamed from source_control_account_id
  manager_id?: string
} 