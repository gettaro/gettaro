package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteTitle(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockError     error
		expectedError error
	}{
		{
			name: "successful deletion",
			id:   "title-1",
		},
		{
			name:          "database error",
			id:            "title-1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:          "title not found",
			id:            "non-existent",
			mockError:     errors.New("record not found"),
			expectedError: errors.New("record not found"),
		},
		{
			name:          "invalid id",
			id:            "",
			mockError:     errors.New("invalid id"),
			expectedError: errors.New("invalid id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTitleDB)
			api := &Api{db: mockDB}

			mockDB.On("DeleteTitle", tt.id).Return(tt.mockError)

			err := api.DeleteTitle(context.Background(), tt.id)

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
