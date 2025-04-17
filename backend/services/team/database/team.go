package database

import (
	"context"
	"errors"

	"ems.dev/backend/services/team/types"
	"gorm.io/gorm"
)

type TeamDB struct {
	db *gorm.DB
}

func NewTeamDB(db *gorm.DB) *TeamDB {
	return &TeamDB{
		db: db,
	}
}

// CreateTeam creates a new team
func (t *TeamDB) CreateTeam(ctx context.Context, team *types.Team) error {
	return t.db.WithContext(ctx).Create(team).Error
}

// GetTeam retrieves a team by ID
func (t *TeamDB) GetTeam(ctx context.Context, id string) (*types.Team, error) {
	var team types.Team
	err := t.db.WithContext(ctx).
		Preload("Organization").
		Preload("Members").
		First(&team, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}

// ListTeams retrieves teams based on search parameters
func (t *TeamDB) ListTeams(ctx context.Context, params types.TeamSearchParams) ([]types.Team, error) {
	var teams []types.Team
	query := t.db.WithContext(ctx).
		Preload("Organization").
		Preload("Members")

	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*params.Name+"%")
	}
	if params.OrganizationID != nil {
		query = query.Where("organization_id = ?", *params.OrganizationID)
	}

	err := query.Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// UpdateTeam updates a team
func (t *TeamDB) UpdateTeam(ctx context.Context, id string, team *types.Team) error {
	return t.db.WithContext(ctx).Model(&types.Team{}).
		Where("id = ?", id).
		Updates(team).Error
}

// DeleteTeam deletes a team
func (t *TeamDB) DeleteTeam(ctx context.Context, id string) error {
	return t.db.WithContext(ctx).Delete(&types.Team{}, "id = ?", id).Error
}

// AddTeamMember adds a user to a team
func (t *TeamDB) AddTeamMember(ctx context.Context, teamID string, member *types.TeamMember) error {
	return t.db.WithContext(ctx).Create(member).Error
}

// RemoveTeamMember removes a user from a team
func (t *TeamDB) RemoveTeamMember(ctx context.Context, teamID, userID string) error {
	return t.db.WithContext(ctx).
		Delete(&types.TeamMember{}, "team_id = ? AND user_id = ?", teamID, userID).Error
}

// GetTeamMember retrieves a team member
func (t *TeamDB) GetTeamMember(ctx context.Context, teamID, userID string) (*types.TeamMember, error) {
	var member types.TeamMember
	err := t.db.WithContext(ctx).
		First(&member, "team_id = ? AND user_id = ?", teamID, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}
