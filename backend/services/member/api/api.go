package api

import (
	"context"

	directsapi "ems.dev/backend/services/directs/api"
	memberdb "ems.dev/backend/services/member/database"
	"ems.dev/backend/services/member/types"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
	titleapi "ems.dev/backend/services/title/api"
	userapi "ems.dev/backend/services/user/api"
)

// MemberAPI defines the interface for member operations
type MemberAPI interface {
	AddOrganizationMember(ctx context.Context, req types.AddMemberRequest, member *types.OrganizationMember) (*types.OrganizationMember, error)
	RemoveOrganizationMember(ctx context.Context, orgID string, userID string) error
	GetOrganizationMembers(ctx context.Context, orgID string, params *types.OrganizationMemberParams) ([]types.OrganizationMember, error)
	GetOrganizationMemberByID(ctx context.Context, memberID string) (*types.OrganizationMember, error)
	IsOrganizationOwner(ctx context.Context, orgID string, userID string) (bool, error)
	UpdateOrganizationMember(ctx context.Context, orgID string, memberID string, req types.UpdateMemberRequest) error
	CalculateSourceControlMemberMetrics(ctx context.Context, organizationID string, memberID string, params sourcecontroltypes.MemberMetricsParams) (*sourcecontroltypes.MetricsResponse, error)

	// External Accounts
	GetExternalAccounts(ctx context.Context, params *types.ExternalAccountParams) ([]types.ExternalAccount, error)
	CreateExternalAccounts(ctx context.Context, accounts []*types.ExternalAccount) error
	GetExternalAccount(ctx context.Context, id string) (*types.ExternalAccount, error)
	UpdateExternalAccount(ctx context.Context, account *types.ExternalAccount) error
	// UpdateExternalAccountMemberID updates the member_id association for an external account
	// Validates that the account belongs to the specified organization
	UpdateExternalAccountMemberID(ctx context.Context, organizationID string, accountID string, memberID *string) (*types.ExternalAccount, error)
}

type Api struct {
	db               memberdb.DB
	userApi          userapi.UserAPI
	sourceControlApi sourcecontrolapi.SourceControlAPI
	titleApi         titleapi.TitleAPI
	directsApi       directsapi.DirectReportsAPI
}

func NewApi(memberDb memberdb.DB, userApi userapi.UserAPI, sourceControlApi sourcecontrolapi.SourceControlAPI, titleApi titleapi.TitleAPI, directsApi directsapi.DirectReportsAPI) *Api {
	return &Api{
		db:               memberDb,
		userApi:          userApi,
		sourceControlApi: sourceControlApi,
		titleApi:         titleApi,
		directsApi:       directsApi,
	}
}
