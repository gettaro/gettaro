package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestGetExternalAccount(t *testing.T) {
	ctx := context.Background()
	accountID := "account-1"
	orgID := "org-1"

	tests := []struct {
		name          string
		accountID     string
		mockAccount   *types.ExternalAccount
		mockError     error
		expectedAccount *types.ExternalAccount
		expectedError error
	}{
		{
			name:     "success",
			accountID: accountID,
			mockAccount: &types.ExternalAccount{
				ID:             accountID,
				OrganizationID: &orgID,
			},
			expectedAccount: &types.ExternalAccount{
				ID:             accountID,
				OrganizationID: &orgID,
			},
		},
		{
			name:          "not found",
			accountID:     accountID,
			mockAccount:   nil,
			mockError:     nil,
			expectedAccount: nil,
		},
		{
			name:          "database error",
			accountID:     accountID,
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
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

			result, err := api.GetExternalAccount(ctx, tt.accountID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectedAccount == nil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Equal(t, tt.expectedAccount.ID, result.ID)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
