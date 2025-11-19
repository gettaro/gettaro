package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"ems.dev/backend/services/title/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAssignMemberTitle(t *testing.T) {
	tests := []struct {
		name          string
		memberTitle   types.MemberTitle
		mockError     error
		expectedError error
	}{
		{
			name: "successful assignment",
			memberTitle: types.MemberTitle{
				MemberID:       "member-1",
				TitleID:        "title-1",
				OrganizationID: "org-1",
			},
		},
		{
			name: "database error",
			memberTitle: types.MemberTitle{
				MemberID:       "member-1",
				TitleID:        "title-1",
				OrganizationID: "org-1",
			},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name: "duplicate assignment",
			memberTitle: types.MemberTitle{
				MemberID:       "member-1",
				TitleID:        "title-1",
				OrganizationID: "org-1",
			},
			mockError:     errors.New("duplicate assignment"),
			expectedError: errors.New("duplicate assignment"),
		},
		{
			name: "invalid member id",
			memberTitle: types.MemberTitle{
				MemberID:       "",
				TitleID:        "title-1",
				OrganizationID: "org-1",
			},
			mockError:     errors.New("invalid member id"),
			expectedError: errors.New("invalid member id"),
		},
		{
			name: "invalid title id",
			memberTitle: types.MemberTitle{
				MemberID:       "member-1",
				TitleID:        "",
				OrganizationID: "org-1",
			},
			mockError:     errors.New("invalid title id"),
			expectedError: errors.New("invalid title id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockTitleDB)
			api := &Api{db: mockDB}

			mockDB.On("AssignMemberTitle", tt.memberTitle).Return(tt.mockError).Run(func(args mock.Arguments) {
				mt := args.Get(0).(types.MemberTitle)
				mt.ID = "member-title-1"
				mt.CreatedAt = time.Now()
				mt.UpdatedAt = time.Now()
			})

			err := api.AssignMemberTitle(context.Background(), tt.memberTitle)

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
