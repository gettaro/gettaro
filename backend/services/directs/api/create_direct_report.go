package api

import (
	"context"

	"ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/directs/types"
)

// CreateDirectReport creates a new direct report relationship
func (a *DirectReportsAPIImpl) CreateDirectReport(ctx context.Context, params types.CreateDirectReportParams) (*types.DirectReport, error) {
	// Validate that all required fields are not empty
	if params.ManagerMemberID == "" {
		return nil, errors.NewBadRequestError("manager member ID cannot be empty")
	}
	if params.ReportMemberID == "" {
		return nil, errors.NewBadRequestError("report member ID cannot be empty")
	}
	if params.OrganizationID == "" {
		return nil, errors.NewBadRequestError("organization ID cannot be empty")
	}

	directReport := &types.DirectReport{
		ManagerMemberID: params.ManagerMemberID,
		ReportMemberID:  params.ReportMemberID,
		OrganizationID:  params.OrganizationID,
		Depth:           params.Depth,
	}

	err := a.db.CreateDirectReport(ctx, directReport)
	if err != nil {
		return nil, err
	}

	return a.db.GetDirectReport(ctx, directReport.ID)
}
