// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/storage.go
//
// Generated by this command:
//
//	mockgen -source=internal/storage/storage.go -destination=internal/mock_storage/mock.go
//

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	storage "github.com/LilLebowski/loyalty-system/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddOrderForUser mocks base method.
func (m *MockStorage) AddOrderForUser(ctx context.Context, externalOrderID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrderForUser", ctx, externalOrderID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrderForUser indicates an expected call of AddOrderForUser.
func (mr *MockStorageMockRecorder) AddOrderForUser(ctx, externalOrderID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrderForUser", reflect.TypeOf((*MockStorage)(nil).AddOrderForUser), ctx, externalOrderID, userID)
}

// AddWithdrawalForUser mocks base method.
func (m *MockStorage) AddWithdrawalForUser(ctx context.Context, userID string, withdrawal storage.Withdrawal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddWithdrawalForUser", ctx, userID, withdrawal)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddWithdrawalForUser indicates an expected call of AddWithdrawalForUser.
func (mr *MockStorageMockRecorder) AddWithdrawalForUser(ctx, userID, withdrawal any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddWithdrawalForUser", reflect.TypeOf((*MockStorage)(nil).AddWithdrawalForUser), ctx, userID, withdrawal)
}

// GetOrdersByUser mocks base method.
func (m *MockStorage) GetOrdersByUser(ctx context.Context, userID string) ([]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByUser", ctx, userID)
	ret0, _ := ret[0].([]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByUser indicates an expected call of GetOrdersByUser.
func (mr *MockStorageMockRecorder) GetOrdersByUser(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByUser", reflect.TypeOf((*MockStorage)(nil).GetOrdersByUser), ctx, userID)
}

// GetOrdersInProgress mocks base method.
func (m *MockStorage) GetOrdersInProgress(ctx context.Context) ([]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersInProgress", ctx)
	ret0, _ := ret[0].([]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersInProgress indicates an expected call of GetOrdersInProgress.
func (mr *MockStorageMockRecorder) GetOrdersInProgress(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersInProgress", reflect.TypeOf((*MockStorage)(nil).GetOrdersInProgress), ctx)
}

// GetUserBalance mocks base method.
func (m *MockStorage) GetUserBalance(ctx context.Context, userID string) (storage.UserBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", ctx, userID)
	ret0, _ := ret[0].(storage.UserBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockStorageMockRecorder) GetUserBalance(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockStorage)(nil).GetUserBalance), ctx, userID)
}

// GetUserByLogin mocks base method.
func (m *MockStorage) GetUserByLogin(ctx context.Context, authData storage.Auth) (storage.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", ctx, authData)
	ret0, _ := ret[0].(storage.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockStorageMockRecorder) GetUserByLogin(ctx, authData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockStorage)(nil).GetUserByLogin), ctx, authData)
}

// GetWithdrawalsForUser mocks base method.
func (m *MockStorage) GetWithdrawalsForUser(ctx context.Context, userID string) ([]storage.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawalsForUser", ctx, userID)
	ret0, _ := ret[0].([]storage.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawalsForUser indicates an expected call of GetWithdrawalsForUser.
func (mr *MockStorageMockRecorder) GetWithdrawalsForUser(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawalsForUser", reflect.TypeOf((*MockStorage)(nil).GetWithdrawalsForUser), ctx, userID)
}

// Register mocks base method.
func (m *MockStorage) Register(ctx context.Context, login, passwordHash string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, login, passwordHash)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockStorageMockRecorder) Register(ctx, login, passwordHash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockStorage)(nil).Register), ctx, login, passwordHash)
}

// UpdateOrder mocks base method.
func (m *MockStorage) UpdateOrder(ctx context.Context, order storage.Accrual) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockStorageMockRecorder) UpdateOrder(ctx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockStorage)(nil).UpdateOrder), ctx, order)
}
