package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateExternalAccounts(t *testing.T) {
	ctx := context.Background()
	orgID := "org-1"

	tests := []struct {
		name          string
		accounts      []*types.ExternalAccount
		mockError     error
		expectedError error
	}{
		{
			name: "success",
			accounts: []*types.ExternalAccount{
				{ID: "account-1", OrganizationID: &orgID},
			},
			mockError: nil,
		},
		{
			name: "database error",
			accounts: []*types.ExternalAccount{
				{ID: "account-1", OrganizationID: &orgID},
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

			mockDB.On("CreateExternalAccounts", ctx, tt.accounts).Return(tt.mockError)

			err := api.CreateExternalAccounts(ctx, tt.accounts)

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
