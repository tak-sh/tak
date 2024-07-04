package debug

import (
	"context"
	"github.com/eddieowens/opts"
	"github.com/emirpasic/gods/v2/stacks"
	"github.com/emirpasic/gods/v2/stacks/linkedliststack"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"github.com/tak-sh/tak/pkg/headless/step/stepper"
)

type ControlStepper interface {
	Controller
	stepper.Stepper
}

type Controller interface {
	Replay()
	PreviousStep()

	// Step allows for the next Step to take place. Similar to "step over"
	// in a debugger.
	Step()
}

type ControllerFactory interface {
	Controller
	stepper.Factory
}

func NewStepper(signals []*step.ConditionalSignal, steps []*step.Step, o ...opts.Opt[stepper.Opts]) ControlStepper {
	out := newController()
	out.Stepper = stepper.New(signals, steps, o...)

	return out
}

func newController() *controller {
	out := &controller{
		History:   linkedliststack.New[*HandleEntry](),
		readyChan: make(chan struct{}, 1),
	}

	return out
}

func NewFactory(o ...opts.Opt[stepper.Opts]) ControllerFactory {
	out := &factory{
		Opts:       o,
		Controller: newController(),
	}

	return out
}

type factory struct {
	Opts       []opts.Opt[stepper.Opts]
	Controller *controller
}

func (f *factory) Replay() {
	f.Controller.Replay()
}

func (f *factory) PreviousStep() {
	f.Controller.PreviousStep()
}

func (f *factory) Step() {
	f.Controller.Step()
}

func (f *factory) NewStepper(globalSignals []*step.ConditionalSignal, steps []*step.Step) stepper.Stepper {
	f.Controller.Stepper = stepper.New(globalSignals, steps, f.Opts...)
	return f.Controller
}

var _ stepper.Stepper = &controller{}

type controller struct {
	stepper.Stepper
	History stacks.Stack[*HandleEntry]

	// waits for messages on this chan to move to the next step.
	readyChan chan struct{}
}

func (s *controller) Step() {
	if len(s.readyChan) == 0 && s.Stepper != nil {
		s.Stepper.Current().Val.Cancel(engine.ErrSkip)
		s.readyChan <- struct{}{}
	}
}

func (s *controller) Replay() {
	s.jump(1)
}

func (s *controller) PreviousStep() {
	s.jump(2)
}

func (s *controller) jump(lvls int) {
	if s.Stepper == nil {
		return
	}
	h, _ := s.History.Pop()
	if h != nil && h.Handle.Node() != nil {
		s.Stepper.Jump(stepper.NavUp(h.Handle.Node(), lvls))
	}
	s.Step()
}

func (s *controller) String() string {
	return s.Stepper.String()
}

func (s *controller) Next(c *engine.Context) stepper.Handle {
	for {
		select {
		case <-c.Done():
			return stepper.NewErrHandle(context.Cause(c))
		case _, ok := <-s.readyChan:
			if !ok {
				return stepper.NewErrHandle(context.Canceled)
			}
		}

		handle := s.Stepper.Next(c)
		s.History.Push(&HandleEntry{
			Data:   c.TemplateData.Merge(),
			Handle: handle,
		})
		return handle
	}
}

type HandleEntry struct {
	Data   *engine.TemplateData
	Handle stepper.Handle
}
