// Code generated by MockGen. DO NOT EDIT.
// Source: operation.go

// Package mock_v1 is a generated GoMock package.
package mock_v1

import (
	context "context"
	reflect "reflect"
	entity "usersegments/internal/entity"
	service "usersegments/internal/service"

	gomock "github.com/golang/mock/gomock"
)

// MockOperationService is a mock of OperationService interface.
type MockOperationService struct {
	ctrl     *gomock.Controller
	recorder *MockOperationServiceMockRecorder
}

// MockOperationServiceMockRecorder is the mock recorder for MockOperationService.
type MockOperationServiceMockRecorder struct {
	mock *MockOperationService
}

// NewMockOperationService creates a new mock instance.
func NewMockOperationService(ctrl *gomock.Controller) *MockOperationService {
	mock := &MockOperationService{ctrl: ctrl}
	mock.recorder = &MockOperationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOperationService) EXPECT() *MockOperationServiceMockRecorder {
	return m.recorder
}

// GetOperationHistory mocks base method.
func (m *MockOperationService) GetOperationHistory(ctx context.Context, user_id int, input service.GetOperationHistoryInput) ([]entity.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOperationHistory", ctx, user_id, input)
	ret0, _ := ret[0].([]entity.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOperationHistory indicates an expected call of GetOperationHistory.
func (mr *MockOperationServiceMockRecorder) GetOperationHistory(ctx, user_id, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOperationHistory", reflect.TypeOf((*MockOperationService)(nil).GetOperationHistory), ctx, user_id, input)
}
