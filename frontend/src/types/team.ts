export type TeamType = 'squad' | 'chapter' | 'tribe' | 'guild'

export interface Team {
  id: string
  name: string
  description: string
  type?: TeamType
  organization_id: string
  created_at: string
  updated_at: string
  members: TeamMember[]
}

export interface TeamMember {
  id: string
  member_id: string
  created_at: string
  updated_at: string
}

export interface CreateTeamRequest {
  name: string
  description: string
  type?: TeamType
  organization_id: string
}

export interface UpdateTeamRequest {
  name?: string
  description?: string
  type?: TeamType
}

export interface AddTeamMemberRequest {
  member_id: string
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


