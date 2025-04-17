package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authTypes "ems.dev/backend/http/types/auth"
	orgtypes "ems.dev/backend/services/organization/types"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrganizationAPI struct {
	mock.Mock
}

func (m *MockOrganizationAPI) CreateOrganization(ctx context.Context, org *orgtypes.Organization, ownerID string) error {
	args := m.Called(ctx, org, ownerID)
	return args.Error(0)
}

func (m *MockOrganizationAPI) GetUserOrganizations(ctx context.Context, userID string) ([]orgtypes.Organization, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]orgtypes.Organization), args.Error(1)
}

func (m *MockOrganizationAPI) UpdateOrganization(ctx context.Context, org *orgtypes.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationAPI) DeleteOrganization(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrganizationAPI) AddOrganizationMember(ctx context.Context, orgID, userID string) error {
	args := m.Called(ctx, orgID, userID)
	return args.Error(0)
}

func (m *MockOrganizationAPI) RemoveOrganizationMember(ctx context.Context, orgID, userID string) error {
	args := m.Called(ctx, orgID, userID)
	return args.Error(0)
}

func (m *MockOrganizationAPI) GetOrganizationMembers(ctx context.Context, orgID string) ([]orgtypes.OrganizationMember, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).([]orgtypes.OrganizationMember), args.Error(1)
}

func (m *MockOrganizationAPI) IsOrganizationOwner(ctx context.Context, orgID, userID string) (bool, error) {
	args := m.Called(ctx, orgID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrganizationAPI) GetOrganizationByID(ctx context.Context, id string) (*orgtypes.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgtypes.Organization), args.Error(1)
}

type MockUserAPI struct {
	mock.Mock
}

func (m *MockUserAPI) FindUser(params usertypes.UserSearchParams) (*usertypes.User, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usertypes.User), args.Error(1)
}

func (m *MockUserAPI) GetOrCreateUserFromAuthProvider(provider string, providerID string, email string, name string) (*usertypes.User, error) {
	args := m.Called(provider, providerID, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usertypes.User), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockOrganizationAPI, *MockUserAPI) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockOrgAPI := new(MockOrganizationAPI)
	mockUserAPI := new(MockUserAPI)
	handler := NewOrganizationHandler(mockOrgAPI, mockUserAPI)

	// Add mock auth middleware
	router.Use(func(c *gin.Context) {
		// Set user claims in context
		userClaims := &authTypes.UserClaims{
			Email: "test@example.com",
		}
		c.Set("user_claims", userClaims)
		c.Next()
	})

	api := router.Group("/api")
	handler.RegisterRoutes(api)

	return router, mockOrgAPI, mockUserAPI
}

func TestCreateOrganization(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		user           *usertypes.User
		userError      error
		createError    error
		expectedStatus int
	}{
		{
			name:           "successful creation",
			requestBody:    `{"name": "Test Org", "slug": "test-org"}`,
			user:           &usertypes.User{ID: "user-1", Email: "test@example.com"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid request body",
			requestBody:    `{"name": "Test Org"}`,
			user:           &usertypes.User{ID: "user-1", Email: "test@example.com"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "user not found",
			requestBody:    `{"name": "Test Org", "slug": "test-org"}`,
			userError:      errors.New("user not found"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "organization creation failed",
			requestBody:    `{"name": "Test Org", "slug": "test-org"}`,
			user:           &usertypes.User{ID: "user-1", Email: "test@example.com"},
			createError:    errors.New("creation failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockOrgAPI, mockUserAPI := setupTestRouter()

			// Set up mock expectations
			email := "test@example.com"
			mockUserAPI.On("FindUser", usertypes.UserSearchParams{Email: &email}).Return(tt.user, tt.userError)
			if tt.user != nil {
				mockOrgAPI.On("CreateOrganization", mock.Anything, mock.Anything, tt.user.ID).Return(tt.createError)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/organizations", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestListOrganizations(t *testing.T) {
	tests := []struct {
		name           string
		user           *usertypes.User
		userError      error
		orgs           []orgtypes.Organization
		orgsError      error
		expectedStatus int
	}{
		{
			name: "successful list",
			user: &usertypes.User{ID: "user-1", Email: "test@example.com"},
			orgs: []orgtypes.Organization{
				{ID: "org-1", Name: "Org 1"},
				{ID: "org-2", Name: "Org 2"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user not found",
			userError:      errors.New("user not found"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "failed to get organizations",
			user:           &usertypes.User{ID: "user-1", Email: "test@example.com"},
			orgsError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockOrgAPI, mockUserAPI := setupTestRouter()

			// Set up mock expectations
			email := "test@example.com"
			mockUserAPI.On("FindUser", usertypes.UserSearchParams{Email: &email}).Return(tt.user, tt.userError)
			if tt.user != nil {
				mockOrgAPI.On("GetUserOrganizations", mock.Anything, tt.user.ID).Return(tt.orgs, tt.orgsError)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/organizations", nil)
			req.Header.Set("Authorization", "Bearer test-token")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetOrganization(t *testing.T) {
	tests := []struct {
		name           string
		orgID          string
		user           *usertypes.User
		userError      error
		orgs           []orgtypes.Organization
		orgsError      error
		expectedStatus int
	}{
		{
			name:  "successful get",
			orgID: "org-1",
			user:  &usertypes.User{ID: "user-1"},
			orgs: []orgtypes.Organization{
				{ID: "org-1", Name: "Org 1"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user not found",
			orgID:          "org-1",
			userError:      errors.New("user not found"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "organization not found",
			orgID:          "org-1",
			user:           &usertypes.User{ID: "user-1"},
			orgs:           []orgtypes.Organization{},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "failed to get organizations",
			orgID:          "org-1",
			user:           &usertypes.User{ID: "user-1"},
			orgsError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockOrgAPI, mockUserAPI := setupTestRouter()

			mockUserAPI.On("FindUser", mock.Anything).Return(tt.user, tt.userError)
			if tt.user != nil {
				mockOrgAPI.On("GetUserOrganizations", mock.Anything, tt.user.ID).Return(tt.orgs, tt.orgsError)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/organizations/"+tt.orgID, nil)
			req.Header.Set("Authorization", "Bearer test-token")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	tests := []struct {
		name           string
		orgID          string
		requestBody    string
		user           *usertypes.User
		userError      error
		orgs           []orgtypes.Organization
		orgsError      error
		updateError    error
		expectedStatus int
	}{
		{
			name:        "successful update",
			orgID:       "org-1",
			requestBody: `{"name": "Updated Org", "slug": "updated-org"}`,
			user:        &usertypes.User{ID: "user-1"},
			orgs: []orgtypes.Organization{
				{ID: "org-1", Name: "Org 1", IsOwner: true},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user not found",
			orgID:          "org-1",
			requestBody:    `{"name": "Updated Org", "slug": "updated-org"}`,
			userError:      errors.New("user not found"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:        "not organization owner",
			orgID:       "org-1",
			requestBody: `{"name": "Updated Org", "slug": "updated-org"}`,
			user:        &usertypes.User{ID: "user-1"},
			orgs: []orgtypes.Organization{
				{ID: "org-1", Name: "Org 1", IsOwner: false},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "update failed",
			orgID:       "org-1",
			requestBody: `{"name": "Updated Org", "slug": "updated-org"}`,
			user:        &usertypes.User{ID: "user-1"},
			orgs: []orgtypes.Organization{
				{ID: "org-1", Name: "Org 1", IsOwner: true},
			},
			updateError:    errors.New("update failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockOrgAPI, mockUserAPI := setupTestRouter()

			mockUserAPI.On("FindUser", mock.Anything).Return(tt.user, tt.userError)
			if tt.user != nil {
				mockOrgAPI.On("GetUserOrganizations", mock.Anything, tt.user.ID).Return(tt.orgs, tt.orgsError)
				if len(tt.orgs) > 0 && tt.orgs[0].IsOwner {
					mockOrgAPI.On("UpdateOrganization", mock.Anything, mock.Anything).Return(tt.updateError)
				}
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/organizations/"+tt.orgID, strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestDeleteOrganization(t *testing.T) {
	tests := []struct {
		name           string
		orgID          string
		deleteError    error
		expectedStatus int
	}{
		{
			name:           "successful delete",
			orgID:          "org-1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing organization ID",
			orgID:          "",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "delete failed",
			orgID:          "org-1",
			deleteError:    errors.New("delete failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockOrgAPI, _ := setupTestRouter()

			if tt.orgID != " " {
				mockOrgAPI.On("DeleteOrganization", mock.Anything, tt.orgID).Return(tt.deleteError)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/organizations/"+strings.TrimSpace(tt.orgID), nil)
			req.Header.Set("Authorization", "Bearer test-token")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAddOrganizationMember(t *testing.T) {
	tests := []struct {
		name           string
		orgID          string
		requestBody    string
		user           *usertypes.User
		userError      error
		isOwner        bool
		isOwnerError   error
		addError       error
		expectedStatus int
	}{
		{
			name:           "successful add",
			orgID:          "org-1",
			requestBody:    `{"userId": "user-2"}`,
			user:           &usertypes.User{ID: "user-1"},
			isOwner:        true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "user not found",
			orgID:          "org-1",
			requestBody:    `{"userId": "user-2"}`,
			userError:      errors.New("user not found"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "not organization owner",
			orgID:          "org-1",
			requestBody:    `{"userId": "user-2"}`,
			user:           &usertypes.User{ID: "user-1"},
			isOwner:        false,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "add failed",
			orgID:          "org-1",
			requestBody:    `{"userId": "user-2"}`,
			user:           &usertypes.User{ID: "user-1"},
			isOwner:        true,
			addError:       errors.New("add failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockOrgAPI, mockUserAPI := setupTestRouter()

			mockUserAPI.On("FindUser", mock.Anything).Return(tt.user, tt.userError)
			if tt.user != nil {
				mockOrgAPI.On("IsOrganizationOwner", mock.Anything, tt.orgID, tt.user.ID).Return(tt.isOwner, tt.isOwnerError)
				if tt.isOwner {
					mockOrgAPI.On("AddOrganizationMember", mock.Anything, tt.orgID, "user-2").Return(tt.addError)
				}
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/organizations/"+tt.orgID+"/members", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRemoveOrganizationMember(t *testing.T) {
	tests := []struct {
		name           string
		orgID          string
		userID         string
		user           *usertypes.User
		userError      error
		isOwner        bool
		isOwnerError   error
		removeError    error
		expectedStatus int
	}{
		{
			name:           "successful remove",
			orgID:          "org-1",
			userID:         "user-2",
			user:           &usertypes.User{ID: "user-1"},
			isOwner:        true,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "user not found",
			orgID:          "org-1",
			userID:         "user-2",
			userError:      errors.New("user not found"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "not organization owner",
			orgID:          "org-1",
			userID:         "user-2",
			user:           &usertypes.User{ID: "user-1"},
			isOwner:        false,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "remove failed",
			orgID:          "org-1",
			userID:         "user-2",
			user:           &usertypes.User{ID: "user-1"},
			isOwner:        true,
			removeError:    errors.New("remove failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockOrgAPI, mockUserAPI := setupTestRouter()

			mockUserAPI.On("FindUser", mock.Anything).Return(tt.user, tt.userError)
			if tt.user != nil {
				mockOrgAPI.On("IsOrganizationOwner", mock.Anything, tt.orgID, tt.user.ID).Return(tt.isOwner, tt.isOwnerError)
				if tt.isOwner {
					mockOrgAPI.On("RemoveOrganizationMember", mock.Anything, tt.orgID, tt.userID).Return(tt.removeError)
				}
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/organizations/"+tt.orgID+"/members/"+tt.userID, nil)
			req.Header.Set("Authorization", "Bearer test-token")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestListOrganizationMembers(t *testing.T) {
	tests := []struct {
		name           string
		orgID          string
		members        []orgtypes.OrganizationMember
		membersError   error
		expectedStatus int
	}{
		{
			name:  "successful list",
			orgID: "org-1",
			members: []orgtypes.OrganizationMember{
				{ID: "member-1", UserID: "user-1"},
				{ID: "member-2", UserID: "user-2"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing organization ID",
			orgID:          "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "list failed",
			orgID:          "org-1",
			membersError:   errors.New("list failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockOrgAPI, _ := setupTestRouter()

			if tt.orgID != "" {
				mockOrgAPI.On("GetOrganizationMembers", mock.Anything, tt.orgID).Return(tt.members, tt.membersError)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/organizations/"+tt.orgID+"/members", nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
