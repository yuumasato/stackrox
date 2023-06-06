// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
)

// MockDataStore is a mock of DataStore interface.
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore.
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance.
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// GetDeclarativeConfig mocks base method.
func (m *MockDataStore) GetDeclarativeConfig(ctx context.Context, id string) (*storage.DeclarativeConfigHealth, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeclarativeConfig", ctx, id)
	ret0, _ := ret[0].(*storage.DeclarativeConfigHealth)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeclarativeConfig indicates an expected call of GetDeclarativeConfig.
func (mr *MockDataStoreMockRecorder) GetDeclarativeConfig(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeclarativeConfig", reflect.TypeOf((*MockDataStore)(nil).GetDeclarativeConfig), ctx, id)
}

// GetDeclarativeConfigs mocks base method.
func (m *MockDataStore) GetDeclarativeConfigs(ctx context.Context) ([]*storage.DeclarativeConfigHealth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeclarativeConfigs", ctx)
	ret0, _ := ret[0].([]*storage.DeclarativeConfigHealth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeclarativeConfigs indicates an expected call of GetDeclarativeConfigs.
func (mr *MockDataStoreMockRecorder) GetDeclarativeConfigs(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeclarativeConfigs", reflect.TypeOf((*MockDataStore)(nil).GetDeclarativeConfigs), ctx)
}

// RemoveDeclarativeConfig mocks base method.
func (m *MockDataStore) RemoveDeclarativeConfig(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDeclarativeConfig", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveDeclarativeConfig indicates an expected call of RemoveDeclarativeConfig.
func (mr *MockDataStoreMockRecorder) RemoveDeclarativeConfig(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDeclarativeConfig", reflect.TypeOf((*MockDataStore)(nil).RemoveDeclarativeConfig), ctx, id)
}

// UpsertDeclarativeConfig mocks base method.
func (m *MockDataStore) UpsertDeclarativeConfig(ctx context.Context, configHealth *storage.DeclarativeConfigHealth) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertDeclarativeConfig", ctx, configHealth)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertDeclarativeConfig indicates an expected call of UpsertDeclarativeConfig.
func (mr *MockDataStoreMockRecorder) UpsertDeclarativeConfig(ctx, configHealth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertDeclarativeConfig", reflect.TypeOf((*MockDataStore)(nil).UpsertDeclarativeConfig), ctx, configHealth)
}
