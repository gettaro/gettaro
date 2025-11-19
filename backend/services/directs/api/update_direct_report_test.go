package api

import (
	"context"
	"errors"
	"testing"

	"ems.dev/backend/services/directs/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDirectReport(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		id            string
		params        types.UpdateDirectReportParams
		mockError     error
		expectedError error
	}{
		{
			name: "successful update",
			id:   "dr-1",
			params: types.UpdateDirectReportParams{
				Depth: intPtr(2),
			},
		},
		{
			name: "update with nil depth",
			id:   "dr-1",
			params: types.UpdateDirectReportParams{
				Depth: nil,
			},
		},
		{
			name:          "not found",
			id:            "dr-nonexistent",
			params:        types.UpdateDirectReportParams{Depth: intPtr(2)},
			mockError:     errors.New("record not found"),
			expectedError: errors.New("record not found"),
		},
		{
			name:          "database error",
			id:            "dr-1",
			params:        types.UpdateDirectReportParams{Depth: intPtr(2)},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			api := NewDirectReportsAPI(mockDB)

			mockDB.On("UpdateDirectReport", ctx, tt.id, tt.params).Return(tt.mockError)

			err := api.UpdateDirectReport(ctx, tt.id, tt.params)

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
