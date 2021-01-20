// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package util is a generated GoMock package.
package util

import (
	bytes "bytes"
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockExecInterface is a mock of ExecInterface interface
type MockExecInterface struct {
	ctrl     *gomock.Controller
	recorder *MockExecInterfaceMockRecorder
}

// MockExecInterfaceMockRecorder is the mock recorder for MockExecInterface
type MockExecInterfaceMockRecorder struct {
	mock *MockExecInterface
}

// NewMockExecInterface creates a new mock instance
func NewMockExecInterface(ctrl *gomock.Controller) *MockExecInterface {
	mock := &MockExecInterface{ctrl: ctrl}
	mock.recorder = &MockExecInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExecInterface) EXPECT() *MockExecInterfaceMockRecorder {
	return m.recorder
}

// commandExecute mocks base method
func (m *MockExecInterface) commandExecute(resultStdout, resultStderr *bytes.Buffer, ctx context.Context, name string, arg ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{resultStdout, resultStderr, ctx, name}
	for _, a := range arg {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "commandExecute", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// commandExecute indicates an expected call of commandExecute
func (mr *MockExecInterfaceMockRecorder) commandExecute(resultStdout, resultStderr, ctx, name interface{}, arg ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{resultStdout, resultStderr, ctx, name}, arg...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "commandExecute", reflect.TypeOf((*MockExecInterface)(nil).commandExecute), varargs...)
}