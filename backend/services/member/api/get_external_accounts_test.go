package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/member/types"
	"github.com/stretchr/testify/assert"
)

func TestGetExternalAccounts(t *testing.T) {
	ctx := context.Background()
	orgID := "org-1"

	tests := []struct {
		name          string
		params        *types.ExternalAccountParams
		mockAccounts  []types.ExternalAccount
		mockError     error
		expectedAccounts []types.ExternalAccount
		expectedError error
	}{
		{
			name:          "success - empty list",
			params:        &types.ExternalAccountParams{OrganizationID: orgID},
			mockAccounts:  []types.ExternalAccount{},
			expectedAccounts: []types.ExternalAccount{},
		},
		{
			name:   "success - with accounts",
			params: &types.ExternalAccountParams{OrganizationID: orgID},
			mockAccounts: []types.ExternalAccount{
				{ID: "account-1", OrganizationID: &orgID},
				{ID: "account-2", OrganizationID: &orgID},
			},
			expectedAccounts: []types.ExternalAccount{
				{ID: "account-1", OrganizationID: &orgID},
				{ID: "account-2", OrganizationID: &orgID},
			},
		},
		{
			name:          "database error",
			params:        &types.ExternalAccountParams{OrganizationID: orgID},
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

			mockDB.On("GetExternalAccounts", ctx, tt.params).Return(tt.mockAccounts, tt.mockError)

			result, err := api.GetExternalAccounts(ctx, tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedAccounts), len(result))
				if len(tt.expectedAccounts) > 0 {
					assert.Equal(t, tt.expectedAccounts[0].ID, result[0].ID)
				}
			}

			mockDB.AssertExpectations(t)
		})
	}
}
