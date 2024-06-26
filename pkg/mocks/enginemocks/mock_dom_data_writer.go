// Code generated by mockery v2.43.2. DO NOT EDIT.

package enginemocks

import (
	mock "github.com/stretchr/testify/mock"
	engine "github.com/tak-sh/tak/pkg/headless/engine"
)

// DOMDataWriter is an autogenerated mock type for the DOMDataWriter type
type DOMDataWriter struct {
	mock.Mock
}

type DOMDataWriter_Expecter struct {
	mock *mock.Mock
}

func (_m *DOMDataWriter) EXPECT() *DOMDataWriter_Expecter {
	return &DOMDataWriter_Expecter{mock: &_m.Mock}
}

// GetQueries provides a mock function with given fields:
func (_m *DOMDataWriter) GetQueries() []engine.DOMQuery {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetQueries")
	}

	var r0 []engine.DOMQuery
	if rf, ok := ret.Get(0).(func() []engine.DOMQuery); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]engine.DOMQuery)
		}
	}

	return r0
}

// DOMDataWriter_GetQueries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQueries'
type DOMDataWriter_GetQueries_Call struct {
	*mock.Call
}

// GetQueries is a helper method to define mock.On call
func (_e *DOMDataWriter_Expecter) GetQueries() *DOMDataWriter_GetQueries_Call {
	return &DOMDataWriter_GetQueries_Call{Call: _e.mock.On("GetQueries")}
}

func (_c *DOMDataWriter_GetQueries_Call) Run(run func()) *DOMDataWriter_GetQueries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DOMDataWriter_GetQueries_Call) Return(_a0 []engine.DOMQuery) *DOMDataWriter_GetQueries_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DOMDataWriter_GetQueries_Call) RunAndReturn(run func() []engine.DOMQuery) *DOMDataWriter_GetQueries_Call {
	_c.Call.Return(run)
	return _c
}

// NewDOMDataWriter creates a new instance of DOMDataWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDOMDataWriter(t interface {
	mock.TestingT
	Cleanup(func())
}) *DOMDataWriter {
	mock := &DOMDataWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
