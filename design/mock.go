// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_design is a generated GoMock package.
package design

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockShape is a mock of Shape interface
type MockShape struct {
	ctrl     *gomock.Controller
	recorder *MockShapeMockRecorder
}

// MockShapeMockRecorder is the mock recorder for MockShape
type MockShapeMockRecorder struct {
	mock *MockShape
}

// NewMockShape creates a new mock instance
func NewMockShape(ctrl *gomock.Controller) *MockShape {
	mock := &MockShape{ctrl: ctrl}
	mock.recorder = &MockShapeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockShape) EXPECT() *MockShapeMockRecorder {
	return m.recorder
}

// Area mocks base method
func (m *MockShape) Area() float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Area")
	ret0, _ := ret[0].(float64)
	return ret0
}

// Area indicates an expected call of Area
func (mr *MockShapeMockRecorder) Area() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Area", reflect.TypeOf((*MockShape)(nil).Area))
}