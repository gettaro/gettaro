package api

import (
	"context"

	orgapi "ems.dev/backend/services/organization/api"
	teamdb "ems.dev/backend/services/team/database"
	"ems.dev/backend/services/team/types"
)

// TeamAPI defines the interface for team-related operations
type TeamAPI interface {
	// CreateTeam creates a new team
	CreateTeam(ctx context.Context, team *types.Team) error

	// ListTeams returns all teams
	ListTeams(ctx context.Context, params types.TeamSearchParams) ([]types.Team, error)

	// GetTeam returns a team by ID
	GetTeam(ctx context.Context, id string) (*types.Team, error)

	// UpdateTeam updates a team
	UpdateTeam(ctx context.Context, id string, team *types.Team) error

	// DeleteTeam deletes a team by ID
	DeleteTeam(ctx context.Context, id string) error

	// AddTeamMember adds a member to a team
	AddTeamMember(ctx context.Context, teamID string, member *types.TeamMember) error

	// RemoveTeamMember removes a member from a team
	RemoveTeamMember(ctx context.Context, teamID string, userID string) error
}

type Api struct {
	db     teamdb.DB
	orgApi orgapi.OrganizationAPI
}

func NewApi(teamDb teamdb.DB, orgApi orgapi.OrganizationAPI) *Api {
	return &Api{
		db:     teamDb,
		orgApi: orgApi,
	}
}
