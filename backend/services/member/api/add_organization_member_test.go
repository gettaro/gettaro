package api

import (
	"context"
	"errors"
	"testing"

	liberrors "ems.dev/backend/libraries/errors"
	directstypes "ems.dev/backend/services/directs/types"
	"ems.dev/backend/services/member/types"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestAddOrganizationMember(t *testing.T) {
	ctx := context.Background()
	orgID := "org-1"
	userID := "user-1"
	email := "test@example.com"
	titleID := "title-1"
	memberID := "member-1"

	tests := []struct {
		name                string
		req                 types.AddMemberRequest
		member              *types.OrganizationMember
		mockUser            *usertypes.User
		mockUserError       error
		mockCreateUser      *usertypes.User
		mockCreateUserError error
		mockExistingMember  *types.OrganizationMember
		mockExistingError    error
		mockAddError         error
		mockCreatedMember    *types.OrganizationMember
		mockCreatedError     error
		expectedMember      *types.OrganizationMember
		expectedError       error
	}{
		{
			name: "success - user exists",
			req: types.AddMemberRequest{
				Email:   email,
				TitleID: titleID,
			},
			member: &types.OrganizationMember{
				Email:          email,
				OrganizationID: orgID,
				Username:       "testuser",
			},
			mockUser: &usertypes.User{
				ID:    userID,
				Email: email,
			},
			mockExistingMember: nil,
			mockCreatedMember: &types.OrganizationMember{
				ID:             memberID,
				UserID:         userID,
				Email:          email,
				OrganizationID: orgID,
				TitleID:        &titleID,
			},
			expectedMember: &types.OrganizationMember{
				ID:             memberID,
				UserID:         userID,
				Email:          email,
				OrganizationID: orgID,
				TitleID:        &titleID,
			},
		},
		{
			name: "success - create new user",
			req: types.AddMemberRequest{
				Email:   email,
				TitleID: titleID,
			},
			member: &types.OrganizationMember{
				Email:          email,
				OrganizationID: orgID,
				Username:       "testuser",
			},
			mockUser:            nil,
			mockCreateUser:      &usertypes.User{ID: userID, Email: email},
			mockExistingMember:  nil,
			mockCreatedMember:   &types.OrganizationMember{ID: memberID, UserID: userID, Email: email, OrganizationID: orgID, TitleID: &titleID},
			expectedMember:      &types.OrganizationMember{ID: memberID, UserID: userID, Email: email, OrganizationID: orgID, TitleID: &titleID},
		},
		{
			name: "error - user lookup fails",
			req: types.AddMemberRequest{
				Email:   email,
				TitleID: titleID,
			},
			member: &types.OrganizationMember{
				Email:          email,
				OrganizationID: orgID,
			},
			mockUserError: errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name: "error - duplicate member",
			req: types.AddMemberRequest{
				Email:   email,
				TitleID: titleID,
			},
			member: &types.OrganizationMember{
				Email:          email,
				OrganizationID: orgID,
			},
			mockUser: &usertypes.User{
				ID:    userID,
				Email: email,
			},
			mockExistingMember: &types.OrganizationMember{ID: memberID},
			expectedError:      liberrors.NewConflictError("user already a member of organization"),
		},
		{
			name: "error - add member fails",
			req: types.AddMemberRequest{
				Email:   email,
				TitleID: titleID,
			},
			member: &types.OrganizationMember{
				Email:          email,
				OrganizationID: orgID,
			},
			mockUser:           &usertypes.User{ID: userID, Email: email},
			mockExistingMember: nil,
			mockAddError:        errors.New("database error"),
			expectedError:       errors.New("database error"),
		},
		{
			name: "success - with external account",
			req: types.AddMemberRequest{
				Email:              email,
				TitleID:            titleID,
				ExternalAccountID:  "ext-account-1",
			},
			member: &types.OrganizationMember{
				Email:          email,
				OrganizationID: orgID,
				Username:       "testuser",
			},
			mockUser:           &usertypes.User{ID: userID, Email: email},
			mockExistingMember: nil,
			mockCreatedMember:  &types.OrganizationMember{ID: memberID, UserID: userID, Email: email, OrganizationID: orgID, TitleID: &titleID},
			expectedMember:     &types.OrganizationMember{ID: memberID, UserID: userID, Email: email, OrganizationID: orgID, TitleID: &titleID},
		},
		{
			name: "success - with manager",
			req: types.AddMemberRequest{
				Email:     email,
				TitleID:   titleID,
				ManagerID: stringPtr("manager-member-1"),
			},
			member: &types.OrganizationMember{
				Email:          email,
				OrganizationID: orgID,
				Username:       "testuser",
			},
			mockUser:           &usertypes.User{ID: userID, Email: email},
			mockExistingMember: nil,
			mockCreatedMember:  &types.OrganizationMember{ID: memberID, UserID: userID, Email: email, OrganizationID: orgID, TitleID: &titleID},
			expectedMember:     &types.OrganizationMember{ID: memberID, UserID: userID, Email: email, OrganizationID: orgID, TitleID: &titleID},
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

			// Setup mocks
			effectiveUserID := userID
			if tt.mockUser != nil {
				effectiveUserID = tt.mockUser.ID
			} else if tt.mockCreateUser != nil {
				effectiveUserID = tt.mockCreateUser.ID
			}

			if tt.mockUserError != nil {
				mockUserAPI.On("FindUser", usertypes.UserSearchParams{Email: &tt.member.Email}).Return(tt.mockUser, tt.mockUserError)
			} else if tt.mockUser != nil {
				mockUserAPI.On("FindUser", usertypes.UserSearchParams{Email: &tt.member.Email}).Return(tt.mockUser, nil)
			} else {
				mockUserAPI.On("FindUser", usertypes.UserSearchParams{Email: &tt.member.Email}).Return(nil, nil)
				if tt.mockCreateUser != nil {
					mockUserAPI.On("CreateUser", mock.AnythingOfType("*types.User")).Return(tt.mockCreateUser, tt.mockCreateUserError)
				}
			}

			if tt.mockExistingMember != nil {
				// Duplicate check returns existing member
				mockDB.On("GetOrganizationMember", orgID, effectiveUserID).Return(tt.mockExistingMember, tt.mockExistingError).Once()
			} else if tt.mockUser != nil || tt.mockCreateUser != nil {
				// First call: duplicate check (returns nil)
				// Second call: get created member (returns createdMember)
				if tt.mockAddError == nil && tt.mockCreatedMember != nil {
					mockDB.On("GetOrganizationMember", orgID, effectiveUserID).Return(nil, nil).Once()
					mockDB.On("AddOrganizationMember", mock.AnythingOfType("*types.OrganizationMember")).Return(nil).Once()
					mockDB.On("GetOrganizationMember", orgID, effectiveUserID).Return(tt.mockCreatedMember, tt.mockCreatedError).Once()
				} else {
					mockDB.On("GetOrganizationMember", orgID, effectiveUserID).Return(nil, nil).Once()
					if tt.mockAddError != nil {
						mockDB.On("AddOrganizationMember", mock.AnythingOfType("*types.OrganizationMember")).Return(tt.mockAddError).Once()
					}
				}
			}

			// Handle external account
			if tt.req.ExternalAccountID != "" && tt.mockCreatedMember != nil && tt.mockAddError == nil {
				externalAccount := &types.ExternalAccount{ID: tt.req.ExternalAccountID}
				mockDB.On("GetExternalAccount", ctx, tt.req.ExternalAccountID).Return(externalAccount, nil)
				mockDB.On("UpdateExternalAccount", ctx, mock.AnythingOfType("*types.ExternalAccount")).Return(nil)
			}

			// Handle manager
			if tt.req.ManagerID != nil && *tt.req.ManagerID != "" && tt.mockCreatedMember != nil && tt.mockAddError == nil {
				managerMember := &types.OrganizationMember{ID: *tt.req.ManagerID}
				mockDB.On("GetOrganizationMemberByID", ctx, *tt.req.ManagerID).Return(managerMember, nil)
				mockDirectsAPI.On("CreateDirectReport", ctx, mock.MatchedBy(func(params directstypes.CreateDirectReportParams) bool {
					return params.ManagerMemberID == managerMember.ID && params.ReportMemberID == tt.mockCreatedMember.ID
				})).Return(&directstypes.DirectReport{}, nil)
			}

			result, err := api.AddOrganizationMember(ctx, tt.req, tt.member)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectedMember != nil {
					assert.NotNil(t, result)
					assert.Equal(t, tt.expectedMember.ID, result.ID)
					assert.Equal(t, tt.expectedMember.UserID, result.UserID)
					assert.Equal(t, tt.expectedMember.Email, result.Email)
				}
			}

			mockDB.AssertExpectations(t)
			mockUserAPI.AssertExpectations(t)
		})
	}
}
