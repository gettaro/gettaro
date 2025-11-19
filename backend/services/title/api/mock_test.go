package api

import (
	"ems.dev/backend/services/title/types"
	"github.com/stretchr/testify/mock"
)

// MockTitleDB is a mock implementation of the TitleDB interface
type MockTitleDB struct {
	mock.Mock
}

func (m *MockTitleDB) CreateTitle(title *types.Title) error {
	args := m.Called(title)
	return args.Error(0)
}

func (m *MockTitleDB) GetTitle(id string) (*types.Title, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Title), args.Error(1)
}

func (m *MockTitleDB) ListTitles(orgID string) ([]types.Title, error) {
	args := m.Called(orgID)
	return args.Get(0).([]types.Title), args.Error(1)
}

func (m *MockTitleDB) UpdateTitle(title types.Title) error {
	args := m.Called(title)
	return args.Error(0)
}

func (m *MockTitleDB) DeleteTitle(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTitleDB) AssignMemberTitle(memberTitle types.MemberTitle) error {
	args := m.Called(memberTitle)
	return args.Error(0)
}

func (m *MockTitleDB) GetMemberTitle(memberID string, orgID string) (*types.MemberTitle, error) {
	args := m.Called(memberID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.MemberTitle), args.Error(1)
}

func (m *MockTitleDB) RemoveMemberTitle(memberID string, orgID string) error {
	args := m.Called(memberID, orgID)
	return args.Error(0)
}
