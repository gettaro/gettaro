package team

import (
	"time"

	"ems.dev/backend/services/team/types"
)

type TeamResponse struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	OrganizationID string           `json:"organization_id"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Members        []MemberResponse `json:"members"`
}

type MemberResponse struct {
	ID        string    `json:"id"`
	MemberID  string    `json:"member_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListTeamsResponse struct {
	Teams []TeamResponse `json:"teams"`
}

// GetTeamResponse converts a Team to a TeamResponse
func GetTeamResponse(t *types.Team) TeamResponse {
	members := make([]MemberResponse, len(t.Members))
	for i, m := range t.Members {
		members[i] = MemberResponse{
			ID:        m.ID,
			MemberID:  m.MemberID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return TeamResponse{
		ID:             t.ID,
		Name:           t.Name,
		Description:    t.Description,
		OrganizationID: t.OrganizationID,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
		Members:        members,
	}
}
