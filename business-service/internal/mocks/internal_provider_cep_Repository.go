// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"

	shared "github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	mock "github.com/stretchr/testify/mock"
)

// Repository_internal_provider_cep is an autogenerated mock type for the Repository type
type Repository_internal_provider_cep struct {
	mock.Mock
}

// GetAddress provides a mock function with given fields: ctx, _a1
func (_m *Repository_internal_provider_cep) GetAddress(ctx context.Context, _a1 string) (*shared.AddressResponse, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetAddress")
	}

	var r0 *shared.AddressResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*shared.AddressResponse, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *shared.AddressResponse); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shared.AddressResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRepository_internal_provider_cep creates a new instance of Repository_internal_provider_cep. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository_internal_provider_cep(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository_internal_provider_cep {
	mock := &Repository_internal_provider_cep{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
