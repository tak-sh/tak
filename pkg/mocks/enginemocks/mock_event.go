// Code generated by mockery v2.43.2. DO NOT EDIT.

package enginemocks

import mock "github.com/stretchr/testify/mock"

// Event is an autogenerated mock type for the Event type
type Event struct {
	mock.Mock
}

type Event_Expecter struct {
	mock *mock.Mock
}

func (_m *Event) EXPECT() *Event_Expecter {
	return &Event_Expecter{mock: &_m.Mock}
}

// eventSigil provides a mock function with given fields:
func (_m *Event) eventSigil() {
	_m.Called()
}

// Event_eventSigil_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'eventSigil'
type Event_eventSigil_Call struct {
	*mock.Call
}

// eventSigil is a helper method to define mock.On call
func (_e *Event_Expecter) eventSigil() *Event_eventSigil_Call {
	return &Event_eventSigil_Call{Call: _e.mock.On("eventSigil")}
}

func (_c *Event_eventSigil_Call) Run(run func()) *Event_eventSigil_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Event_eventSigil_Call) Return() *Event_eventSigil_Call {
	_c.Call.Return()
	return _c
}

func (_c *Event_eventSigil_Call) RunAndReturn(run func()) *Event_eventSigil_Call {
	_c.Call.Return(run)
	return _c
}

// NewEvent creates a new instance of Event. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEvent(t interface {
	mock.TestingT
	Cleanup(func())
}) *Event {
	mock := &Event{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
