// Code generated by MockGen. DO NOT EDIT.
// Source: srv/srvgql/service.go

// Package mock_srvgql is a generated GoMock package.
package mock_srvgql

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// BuildServiceToken mocks base method
func (m *MockService) BuildServiceToken(idTokenStr string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildServiceToken", idTokenStr)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildServiceToken indicates an expected call of BuildServiceToken
func (mr *MockServiceMockRecorder) BuildServiceToken(idTokenStr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildServiceToken", reflect.TypeOf((*MockService)(nil).BuildServiceToken), idTokenStr)
}
