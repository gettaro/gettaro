package team

import (
	"time"

	"ems.dev/backend/services/team/types"
)

type TeamResponse struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	OrganizationID string           `json:"organizationId"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Members        []MemberResponse `json:"members"`
}

type MemberResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
			UserID:    m.UserID,
			Role:      m.Role,
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
