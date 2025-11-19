package api

import (
	"context"
	"errors"
	"testing"

	liberrors "ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/directs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateDirectReport(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		params         types.CreateDirectReportParams
		mockCreateErr  error
		mockGetResult  *types.DirectReport
		mockGetErr     error
		expectedResult *types.DirectReport
		expectedError  error
	}{
		{
			name: "successful creation",
			params: types.CreateDirectReportParams{
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
			},
			mockGetResult: &types.DirectReport{
				ID:              "dr-1",
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
			},
			expectedResult: &types.DirectReport{
				ID:              "dr-1",
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
			},
		},
		{
			name: "empty manager member ID",
			params: types.CreateDirectReportParams{
				ManagerMemberID: "",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
			},
			expectedError: liberrors.NewBadRequestError("manager member ID cannot be empty"),
		},
		{
			name: "empty report member ID",
			params: types.CreateDirectReportParams{
				ManagerMemberID: "manager-1",
				ReportMemberID:  "",
				OrganizationID:  "org-1",
				Depth:           1,
			},
			expectedError: liberrors.NewBadRequestError("report member ID cannot be empty"),
		},
		{
			name: "empty organization ID",
			params: types.CreateDirectReportParams{
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "",
				Depth:           1,
			},
			expectedError: liberrors.NewBadRequestError("organization ID cannot be empty"),
		},
		{
			name: "database create error",
			params: types.CreateDirectReportParams{
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
			},
			mockCreateErr: errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name: "database get error after create",
			params: types.CreateDirectReportParams{
				ManagerMemberID: "manager-1",
				ReportMemberID:  "report-1",
				OrganizationID:  "org-1",
				Depth:           1,
			},
			mockGetErr:    errors.New("get error"),
			expectedError: errors.New("get error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			// Only set up mocks if we're not expecting a validation error
			isValidationError := tt.expectedError != nil && 
				(tt.expectedError.Error() == "manager member ID cannot be empty" ||
				 tt.expectedError.Error() == "report member ID cannot be empty" ||
				 tt.expectedError.Error() == "organization ID cannot be empty")

			if !isValidationError {
				mockDB.On("CreateDirectReport", ctx, mock.AnythingOfType("*types.DirectReport")).Return(tt.mockCreateErr).Run(func(args mock.Arguments) {
					dr := args.Get(1).(*types.DirectReport)
					dr.ID = "dr-1"
				})

				if tt.mockCreateErr == nil {
					mockDB.On("GetDirectReport", ctx, "dr-1").Return(tt.mockGetResult, tt.mockGetErr)
				}
			}

			result, err := api.CreateDirectReport(ctx, tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if badReqErr, ok := tt.expectedError.(*liberrors.BadRequestError); ok {
					assert.IsType(t, &liberrors.BadRequestError{}, err)
					assert.Equal(t, badReqErr.Error(), err.Error())
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.ManagerMemberID, result.ManagerMemberID)
				assert.Equal(t, tt.expectedResult.ReportMemberID, result.ReportMemberID)
				assert.Equal(t, tt.expectedResult.OrganizationID, result.OrganizationID)
				assert.Equal(t, tt.expectedResult.Depth, result.Depth)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
