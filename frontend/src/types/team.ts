export interface Team {
  id: string
  name: string
  description: string
  organization_id: string
  created_at: string
  updated_at: string
  members: TeamMember[]
}

export interface TeamMember {
  id: string
  member_id: string
  role: string
  created_at: string
  updated_at: string
}

export interface CreateTeamRequest {
  name: string
  description: string
  organization_id: string
}

export interface UpdateTeamRequest {
  name?: string
  description?: string
}

export interface AddTeamMemberRequest {
  member_id: string
  role: string
}

export interface ListTeamsResponse {
  teams: Team[]
}

export interface GetTeamResponse {
  team: Team
}

export interface CreateTeamResponse {
  team: Team
}

export interface UpdateTeamResponse {
  team: Team
}


