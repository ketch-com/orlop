// Code generated by MockGen. DO NOT EDIT.
// Source: server_options.go

// Package orlop_test is a generated GoMock package.
package orlop_test

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockServerOption is a mock of ServerOption interface
type MockServerOption struct {
	ctrl     *gomock.Controller
	recorder *MockServerOptionMockRecorder
}

// MockServerOptionMockRecorder is the mock recorder for MockServerOption
type MockServerOptionMockRecorder struct {
	mock *MockServerOption
}

// NewMockServerOption creates a new mock instance
func NewMockServerOption(ctrl *gomock.Controller) *MockServerOption {
	mock := &MockServerOption{ctrl: ctrl}
	mock.recorder = &MockServerOptionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServerOption) EXPECT() *MockServerOptionMockRecorder {
	return m.recorder
}

// apply mocks base method
func (m *MockServerOption) apply(ctx context.Context, opts *serverOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "apply", ctx, opts)
	ret0, _ := ret[0].(error)
	return ret0
}

// apply indicates an expected call of apply
func (mr *MockServerOptionMockRecorder) apply(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "apply", reflect.TypeOf((*MockServerOption)(nil).apply), ctx, opts)
}

// addHandler mocks base method
func (m *MockServerOption) addHandler(ctx context.Context, opt *serverOptions, mux mux) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "addHandler", ctx, opt, mux)
	ret0, _ := ret[0].(error)
	return ret0
}

// addHandler indicates an expected call of addHandler
func (mr *MockServerOptionMockRecorder) addHandler(ctx, opt, mux interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "addHandler", reflect.TypeOf((*MockServerOption)(nil).addHandler), ctx, opt, mux)
}

// Mockmux is a mock of mux interface
type Mockmux struct {
	ctrl     *gomock.Controller
	recorder *MockmuxMockRecorder
}

// MockmuxMockRecorder is the mock recorder for Mockmux
type MockmuxMockRecorder struct {
	mock *Mockmux
}

// NewMockmux creates a new mock instance
func NewMockmux(ctrl *gomock.Controller) *Mockmux {
	mock := &Mockmux{ctrl: ctrl}
	mock.recorder = &MockmuxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockmux) EXPECT() *MockmuxMockRecorder {
	return m.recorder
}

// Handle mocks base method
func (m *Mockmux) Handle(pattern string, handler http.Handler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Handle", pattern, handler)
}

// Handle indicates an expected call of Handle
func (mr *MockmuxMockRecorder) Handle(pattern, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*Mockmux)(nil).Handle), pattern, handler)
}
