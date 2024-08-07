// Code generated by mockery v2.43.2. DO NOT EDIT.

package steppermocks

import (
	mock "github.com/stretchr/testify/mock"
	step "github.com/tak-sh/tak/pkg/headless/step"

	stepper "github.com/tak-sh/tak/pkg/headless/step/stepper"
)

// Factory is an autogenerated mock type for the Factory type
type Factory struct {
	mock.Mock
}

type Factory_Expecter struct {
	mock *mock.Mock
}

func (_m *Factory) EXPECT() *Factory_Expecter {
	return &Factory_Expecter{mock: &_m.Mock}
}

// NewStepper provides a mock function with given fields: globalSignals, steps
func (_m *Factory) NewStepper(globalSignals []*step.ConditionalSignal, steps []*step.Step) stepper.Stepper {
	ret := _m.Called(globalSignals, steps)

	if len(ret) == 0 {
		panic("no return value specified for NewStepper")
	}

	var r0 stepper.Stepper
	if rf, ok := ret.Get(0).(func([]*step.ConditionalSignal, []*step.Step) stepper.Stepper); ok {
		r0 = rf(globalSignals, steps)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(stepper.Stepper)
		}
	}

	return r0
}

// Factory_NewStepper_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewStepper'
type Factory_NewStepper_Call struct {
	*mock.Call
}

// NewStepper is a helper method to define mock.On call
//   - globalSignals []*step.ConditionalSignal
//   - steps []*step.Step
func (_e *Factory_Expecter) NewStepper(globalSignals interface{}, steps interface{}) *Factory_NewStepper_Call {
	return &Factory_NewStepper_Call{Call: _e.mock.On("NewStepper", globalSignals, steps)}
}

func (_c *Factory_NewStepper_Call) Run(run func(globalSignals []*step.ConditionalSignal, steps []*step.Step)) *Factory_NewStepper_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]*step.ConditionalSignal), args[1].([]*step.Step))
	})
	return _c
}

func (_c *Factory_NewStepper_Call) Return(_a0 stepper.Stepper) *Factory_NewStepper_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Factory_NewStepper_Call) RunAndReturn(run func([]*step.ConditionalSignal, []*step.Step) stepper.Stepper) *Factory_NewStepper_Call {
	_c.Call.Return(run)
	return _c
}

// NewFactory creates a new instance of Factory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *Factory {
	mock := &Factory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
