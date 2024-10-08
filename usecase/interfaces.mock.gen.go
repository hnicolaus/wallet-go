// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/interfaces.go

// Package usecase is a generated GoMock package.
package usecase

import (
	context "context"
	reflect "reflect"

	model "github.com/WalletService/model"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockUsecaseInterface is a mock of UsecaseInterface interface.
type MockUsecaseInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseInterfaceMockRecorder
}

// MockUsecaseInterfaceMockRecorder is the mock recorder for MockUsecaseInterface.
type MockUsecaseInterfaceMockRecorder struct {
	mock *MockUsecaseInterface
}

// NewMockUsecaseInterface creates a new mock instance.
func NewMockUsecaseInterface(ctrl *gomock.Controller) *MockUsecaseInterface {
	mock := &MockUsecaseInterface{ctrl: ctrl}
	mock.recorder = &MockUsecaseInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecaseInterface) EXPECT() *MockUsecaseInterfaceMockRecorder {
	return m.recorder
}

// CreateUserTransaction mocks base method.
func (m *MockUsecaseInterface) CreateUserTransaction(ctx context.Context, transaction model.Transaction) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserTransaction", ctx, transaction)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserTransaction indicates an expected call of CreateUserTransaction.
func (mr *MockUsecaseInterfaceMockRecorder) CreateUserTransaction(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserTransaction", reflect.TypeOf((*MockUsecaseInterface)(nil).CreateUserTransaction), ctx, transaction)
}

// GetUser mocks base method.
func (m *MockUsecaseInterface) GetUser(ctx context.Context, userID int64) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, userID)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUsecaseInterfaceMockRecorder) GetUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUsecaseInterface)(nil).GetUser), ctx, userID)
}

// GetUsers mocks base method.
func (m *MockUsecaseInterface) GetUsers(ctx context.Context, request model.UserFilter) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx, request)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUsecaseInterfaceMockRecorder) GetUsers(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUsecaseInterface)(nil).GetUsers), ctx, request)
}

// RegisterUser mocks base method.
func (m *MockUsecaseInterface) RegisterUser(ctx context.Context, user model.User) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, user)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockUsecaseInterfaceMockRecorder) RegisterUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockUsecaseInterface)(nil).RegisterUser), ctx, user)
}

// UserLogin mocks base method.
func (m *MockUsecaseInterface) UserLogin(ctx context.Context, phoneNumber, password string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserLogin", ctx, phoneNumber, password)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserLogin indicates an expected call of UserLogin.
func (mr *MockUsecaseInterfaceMockRecorder) UserLogin(ctx, phoneNumber, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserLogin", reflect.TypeOf((*MockUsecaseInterface)(nil).UserLogin), ctx, phoneNumber, password)
}
