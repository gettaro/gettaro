package api

import (
	"context"

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

// TeamDB defines the interface for team database operations
type TeamDB interface {
	CreateTeam(ctx context.Context, team *types.Team) error
	ListTeams(ctx context.Context, params types.TeamSearchParams) ([]types.Team, error)
	GetTeam(ctx context.Context, id string) (*types.Team, error)
	UpdateTeam(ctx context.Context, id string, team *types.Team) error
	DeleteTeam(ctx context.Context, id string) error
	AddTeamMember(ctx context.Context, teamID string, member *types.TeamMember) error
	RemoveTeamMember(ctx context.Context, teamID string, userID string) error
}

type Api struct {
	db TeamDB
}

func NewApi(teamDb TeamDB) *Api {
	return &Api{
		db: teamDb,
	}
}
