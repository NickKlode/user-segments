// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_v1 is a generated GoMock package.
package mock_v1

import (
	context "context"
	reflect "reflect"
	service "usersegments/internal/service"

	gomock "github.com/golang/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// AddUserToSegment mocks base method.
func (m *MockUserService) AddUserToSegment(ctx context.Context, input service.AddUserToSegmentsInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToSegment", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToSegment indicates an expected call of AddUserToSegment.
func (mr *MockUserServiceMockRecorder) AddUserToSegment(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToSegment", reflect.TypeOf((*MockUserService)(nil).AddUserToSegment), ctx, input)
}

// CreateUser mocks base method.
func (m *MockUserService) CreateUser(ctx context.Context, input service.UserCreateInput) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, input)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserServiceMockRecorder) CreateUser(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserService)(nil).CreateUser), ctx, input)
}

// DeleteUserFromSegment mocks base method.
func (m *MockUserService) DeleteUserFromSegment(ctx context.Context, input service.DeleteUserFromSegmentsInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserFromSegment", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserFromSegment indicates an expected call of DeleteUserFromSegment.
func (mr *MockUserServiceMockRecorder) DeleteUserFromSegment(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserFromSegment", reflect.TypeOf((*MockUserService)(nil).DeleteUserFromSegment), ctx, input)
}

// GetAllUserSegments mocks base method.
func (m *MockUserService) GetAllUserSegments(ctx context.Context, user_id int) ([]service.AllUserSegmentsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUserSegments", ctx, user_id)
	ret0, _ := ret[0].([]service.AllUserSegmentsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUserSegments indicates an expected call of GetAllUserSegments.
func (mr *MockUserServiceMockRecorder) GetAllUserSegments(ctx, user_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUserSegments", reflect.TypeOf((*MockUserService)(nil).GetAllUserSegments), ctx, user_id)
}
