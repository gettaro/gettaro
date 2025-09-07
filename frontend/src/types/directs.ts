export interface Member {
  id: string
  email: string
  username: string
  organizationId: string
  isOwner: boolean
  titleId?: string
  createdAt: string
  updatedAt: string
}

export interface OrgChartNode {
  member: Member
  directReports: OrgChartNode[]
  depth: number
}

export interface GetManagerTreeResponse {
  orgChart: OrgChartNode[]
}
