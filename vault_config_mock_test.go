// Code generated by MockGen. DO NOT EDIT.
// Source: vault_config.go

// Package orlop_test is a generated GoMock package.
package orlop_test

import (
	gomock "github.com/golang/mock/gomock"
	orlop "github.com/switch-bit/orlop"
	reflect "reflect"
)

// MockHasVaultConfig is a mock of HasVaultConfig interface
type MockHasVaultConfig struct {
	ctrl     *gomock.Controller
	recorder *MockHasVaultConfigMockRecorder
}

// MockHasVaultConfigMockRecorder is the mock recorder for MockHasVaultConfig
type MockHasVaultConfigMockRecorder struct {
	mock *MockHasVaultConfig
}

// NewMockHasVaultConfig creates a new mock instance
func NewMockHasVaultConfig(ctrl *gomock.Controller) *MockHasVaultConfig {
	mock := &MockHasVaultConfig{ctrl: ctrl}
	mock.recorder = &MockHasVaultConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHasVaultConfig) EXPECT() *MockHasVaultConfigMockRecorder {
	return m.recorder
}

// GetEnabled mocks base method
func (m *MockHasVaultConfig) GetEnabled() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetEnabled indicates an expected call of GetEnabled
func (mr *MockHasVaultConfigMockRecorder) GetEnabled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnabled", reflect.TypeOf((*MockHasVaultConfig)(nil).GetEnabled))
}

// GetAddress mocks base method
func (m *MockHasVaultConfig) GetAddress() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddress")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAddress indicates an expected call of GetAddress
func (mr *MockHasVaultConfigMockRecorder) GetAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddress", reflect.TypeOf((*MockHasVaultConfig)(nil).GetAddress))
}

// GetToken mocks base method
func (m *MockHasVaultConfig) GetToken() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToken")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetToken indicates an expected call of GetToken
func (mr *MockHasVaultConfigMockRecorder) GetToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToken", reflect.TypeOf((*MockHasVaultConfig)(nil).GetToken))
}

// GetPrefix mocks base method
func (m *MockHasVaultConfig) GetPrefix() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrefix")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPrefix indicates an expected call of GetPrefix
func (mr *MockHasVaultConfigMockRecorder) GetPrefix() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrefix", reflect.TypeOf((*MockHasVaultConfig)(nil).GetPrefix))
}

// GetTLS mocks base method
func (m *MockHasVaultConfig) GetTLS() orlop.HasTLSConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTLS")
	ret0, _ := ret[0].(orlop.HasTLSConfig)
	return ret0
}

// GetTLS indicates an expected call of GetTLS
func (mr *MockHasVaultConfigMockRecorder) GetTLS() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTLS", reflect.TypeOf((*MockHasVaultConfig)(nil).GetTLS))
}
