package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateExternalAccount(t *testing.T) {
	ctx := context.Background()
	accountID := "account-1"
	orgID := "org-1"

	tests := []struct {
		name          string
		account       *types.ExternalAccount
		mockError     error
		expectedError error
	}{
		{
			name: "success",
			account: &types.ExternalAccount{
				ID:             accountID,
				OrganizationID: &orgID,
			},
			mockError: nil,
		},
		{
			name: "database error",
			account: &types.ExternalAccount{
				ID:             accountID,
				OrganizationID: &orgID,
			},
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

			mockDB.On("UpdateExternalAccount", ctx, tt.account).Return(tt.mockError)

			err := api.UpdateExternalAccount(ctx, tt.account)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
