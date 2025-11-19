package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteIntegrationConfig(t *testing.T) {
	validKey := make([]byte, 32)
	for i := range validKey {
		validKey[i] = byte(i)
	}

	tests := []struct {
		name          string
		id            string
		mockError     error
		expectedError error
	}{
		{
			name: "success",
			id:   "config-1",
		},
		{
			name:          "error - config not found",
			id:            "non-existent",
			mockError:     errors.New("record not found"),
			expectedError: errors.New("record not found"),
		},
		{
			name:          "error - database error",
			id:            "config-1",
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewApi(mockDB, validKey)

			mockDB.On("DeleteIntegrationConfig", tt.id).Return(tt.mockError)

			err := api.DeleteIntegrationConfig(context.Background(), tt.id)

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
