// Code generated by mockery v1.0.0
package cache

import mock "github.com/stretchr/testify/mock"

// MockRequest is an autogenerated mock type for the Request type
type MockRequest struct {
	mock.Mock
}

// CacheKey provides a mock function with given fields:
func (_m *MockRequest) CacheKey() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// CacheMinIndex provides a mock function with given fields:
func (_m *MockRequest) CacheMinIndex() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}
