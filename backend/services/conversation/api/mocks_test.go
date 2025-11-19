package api

import (
	"context"

	"ems.dev/backend/services/conversation/types"
	templatetypes "ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockConversationDB is a mock implementation of ConversationDB
type MockConversationDB struct {
	mock.Mock
}

func (m *MockConversationDB) CreateConversation(ctx context.Context, conversation *types.Conversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

func (m *MockConversationDB) GetConversation(ctx context.Context, id string) (*types.Conversation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Conversation), args.Error(1)
}

func (m *MockConversationDB) GetConversationWithDetails(ctx context.Context, id string) (*types.ConversationWithDetails, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ConversationWithDetails), args.Error(1)
}

func (m *MockConversationDB) ListConversations(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.Conversation, error) {
	args := m.Called(ctx, organizationID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.Conversation), args.Error(1)
}

func (m *MockConversationDB) ListConversationsWithDetails(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.ConversationWithDetails, error) {
	args := m.Called(ctx, organizationID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*types.ConversationWithDetails), args.Error(1)
}

func (m *MockConversationDB) UpdateConversation(ctx context.Context, id string, req *types.UpdateConversationRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockConversationDB) DeleteConversation(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockConversationDB) GetConversationStats(ctx context.Context, organizationID string, managerMemberID *string) (map[string]int, error) {
	args := m.Called(ctx, organizationID, managerMemberID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

// MockConversationTemplateAPI is a mock implementation of ConversationTemplateAPIInterface
type MockConversationTemplateAPI struct {
	mock.Mock
}

func (m *MockConversationTemplateAPI) CreateConversationTemplate(params templatetypes.CreateConversationTemplateParams) (*templatetypes.ConversationTemplate, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*templatetypes.ConversationTemplate), args.Error(1)
}

func (m *MockConversationTemplateAPI) GetConversationTemplate(id uuid.UUID) (*templatetypes.ConversationTemplate, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*templatetypes.ConversationTemplate), args.Error(1)
}

func (m *MockConversationTemplateAPI) ListConversationTemplates(params templatetypes.ConversationTemplateSearchParams) ([]*templatetypes.ConversationTemplate, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*templatetypes.ConversationTemplate), args.Error(1)
}

func (m *MockConversationTemplateAPI) UpdateConversationTemplate(params templatetypes.UpdateConversationTemplateParams) (*templatetypes.ConversationTemplate, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*templatetypes.ConversationTemplate), args.Error(1)
}

func (m *MockConversationTemplateAPI) DeleteConversationTemplate(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Helper functions for tests
func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func statusPtr(s types.ConversationStatus) *types.ConversationStatus {
	return &s
}
