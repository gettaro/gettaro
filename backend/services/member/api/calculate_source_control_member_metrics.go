package api

import (
	"context"
	"encoding/json"

	"ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/member/types"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
	"gorm.io/datatypes"
)

// CalculateSourceControlMemberMetrics retrieves source control metrics for a specific member
func (a *Api) CalculateSourceControlMemberMetrics(ctx context.Context, organizationID string, memberID string, params sourcecontroltypes.MemberMetricsParams) (*sourcecontroltypes.MetricsResponse, error) {
	// Get external accounts for this member (filter by sourcecontrol type)
	sourceControlType := "sourcecontrol"
	externalAccounts, err := a.GetExternalAccounts(ctx, &types.ExternalAccountParams{
		OrganizationID: organizationID,
		MemberIDs:      []string{memberID},
		AccountType:    &sourceControlType,
	})
	if err != nil {
		return nil, err
	}

	sourceControlAccountIDs := []string{}
	for _, account := range externalAccounts {
		sourceControlAccountIDs = append(sourceControlAccountIDs, account.ID)
	}

	if len(sourceControlAccountIDs) == 0 {
		return nil, errors.NewNotFoundError("no source control accounts found for member")
	}

	member, err := a.GetOrganizationMemberByID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	orgMembers, err := a.GetOrganizationMembers(ctx, organizationID, &types.OrganizationMemberParams{
		TitleIDs: []string{*member.TitleID},
	})
	if err != nil {
		return nil, err
	}

	peerMemberIDs := []string{}
	for _, orgMember := range orgMembers {
		peerMemberIDs = append(peerMemberIDs, orgMember.ID)
	}

	peerExternalAccounts, err := a.GetExternalAccounts(ctx, &types.ExternalAccountParams{
		OrganizationID: organizationID,
		MemberIDs:      peerMemberIDs,
		AccountType:    &sourceControlType,
	})
	if err != nil {
		return nil, err
	}

	peerSourceControlAccountIDs := []string{}
	for _, account := range peerExternalAccounts {
		peerSourceControlAccountIDs = append(peerSourceControlAccountIDs, account.ID)
	}
	// Create the metric params with the source control account IDs
	metricParamsMap := map[string]interface{}{
		"organizationId":               organizationID,
		"sourceControlAccountIDs":      sourceControlAccountIDs,
		"peersSourceControlAccountIDs": peerSourceControlAccountIDs,
	}

	// Marshal to JSON bytes
	metricParamsJSON, err := json.Marshal(metricParamsMap)
	if err != nil {
		return nil, err
	}

	metricParams := sourcecontroltypes.MetricRuleParams{
		MetricParams: datatypes.JSON(metricParamsJSON),
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Interval:     params.Interval,
	}

	return a.sourceControlApi.CalculateMetrics(ctx, metricParams)
}
