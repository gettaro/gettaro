export interface Member {
  id: string
  email: string
  username: string
  organization_id: string
  is_owner: boolean
  title_id?: string
  title?: string
  created_at: string
  updated_at: string
}

export interface OrgChartNode {
  member: Member
  direct_reports: OrgChartNode[]
  depth: number
}

export interface GetManagerTreeResponse {
  org_chart: OrgChartNode[]
}
