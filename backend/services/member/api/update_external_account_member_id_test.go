package api

import (
	"context"
	"errors"
	"testing"

	liberrors "ems.dev/backend/libraries/errors"
	"ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateExternalAccountMemberID(t *testing.T) {
	ctx := context.Background()
	orgID := "org-1"
	accountID := "account-1"
	memberID := "member-1"

	tests := []struct {
		name            string
		organizationID string
		accountID       string
		memberID        *string
		mockAccount     *types.ExternalAccount
		mockError       error
		mockUpdateError error
		expectedAccount *types.ExternalAccount
		expectedError   error
	}{
		{
			name:            "success - set member ID",
			organizationID:  orgID,
			accountID:       accountID,
			memberID:        stringPtr(memberID),
			mockAccount:     &types.ExternalAccount{ID: accountID, OrganizationID: &orgID},
			expectedAccount: &types.ExternalAccount{ID: accountID, OrganizationID: &orgID, MemberID: stringPtr(memberID)},
		},
		{
			name:            "success - clear member ID",
			organizationID:  orgID,
			accountID:       accountID,
			memberID:        nil,
			mockAccount:     &types.ExternalAccount{ID: accountID, OrganizationID: &orgID, MemberID: stringPtr(memberID)},
			expectedAccount: &types.ExternalAccount{ID: accountID, OrganizationID: &orgID, MemberID: nil},
		},
		{
			name:          "error - account not found",
			organizationID: orgID,
			accountID:     accountID,
			memberID:      stringPtr(memberID),
			mockAccount:   nil,
			expectedError: liberrors.NewNotFoundError("external account not found"),
		},
		{
			name:          "error - wrong organization",
			organizationID: orgID,
			accountID:     accountID,
			memberID:      stringPtr(memberID),
			mockAccount:   &types.ExternalAccount{ID: accountID, OrganizationID: stringPtr("wrong-org")},
			expectedError: liberrors.NewNotFoundError("external account not found in this organization"),
		},
		{
			name:          "error - database error on get",
			organizationID: orgID,
			accountID:     accountID,
			memberID:      stringPtr(memberID),
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:            "error - database error on update",
			organizationID:  orgID,
			accountID:       accountID,
			memberID:        stringPtr(memberID),
			mockAccount:     &types.ExternalAccount{ID: accountID, OrganizationID: &orgID},
			mockUpdateError: errors.New("database error"),
			expectedError:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockMemberDB)
			mockUserAPI := new(MockUserAPI)
			mockDirectsAPI := new(MockDirectReportsAPI)
			mockSourceControlAPI := new(MockSourceControlAPI)
			mockTitleAPI := new(MockTitleAPI)

			api := NewApi(mockDB, mockUserAPI, mockSourceControlAPI, mockTitleAPI, mockDirectsAPI)

			mockDB.On("GetExternalAccount", ctx, tt.accountID).Return(tt.mockAccount, tt.mockError)

			if tt.mockError == nil && tt.mockAccount != nil {
				if tt.mockAccount.OrganizationID == nil || *tt.mockAccount.OrganizationID != tt.organizationID {
					// Error case - wrong org, no update call
				} else {
					mockDB.On("UpdateExternalAccount", ctx, mock.AnythingOfType("*types.ExternalAccount")).Return(tt.mockUpdateError)
				}
			}

			result, err := api.UpdateExternalAccountMemberID(ctx, tt.organizationID, tt.accountID, tt.memberID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedAccount.MemberID, result.MemberID)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
