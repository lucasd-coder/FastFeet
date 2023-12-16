// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Validator_internal_shared is an autogenerated mock type for the Validator type
type Validator_internal_shared struct {
	mock.Mock
}

// ValidateStruct provides a mock function with given fields: s
func (_m *Validator_internal_shared) ValidateStruct(s interface{}) error {
	ret := _m.Called(s)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewValidator_internal_shared creates a new instance of Validator_internal_shared. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewValidator_internal_shared(t interface {
	mock.TestingT
	Cleanup(func())
}) *Validator_internal_shared {
	mock := &Validator_internal_shared{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
