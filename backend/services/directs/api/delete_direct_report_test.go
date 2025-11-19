package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteDirectReport(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		id            string
		mockError     error
		expectedError error
	}{
		{
			name: "successful deletion",
			id:   "dr-1",
		},
		{
			name:          "not found",
			id:            "dr-nonexistent",
			mockError:     errors.New("record not found"),
			expectedError: errors.New("record not found"),
		},
		{
			name:          "database error",
			id:            "dr-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			mockDB.On("DeleteDirectReport", ctx, tt.id).Return(tt.mockError)

			err := api.DeleteDirectReport(ctx, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
