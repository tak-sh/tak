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

type Stepper interface {
	stepper.Stepper
	Replay()
	PreviousStep()

	// Step allows for the next Step to take place. Similar to "step over"
	// in a debugger.
	Step()
}

func NewStepper(signals []*step.ConditionalSignal, steps []*step.Step, o ...opts.Opt[stepper.Opts]) Stepper {
	out := &debugStepper{
		Stepper:   stepper.New(signals, steps, o...),
		History:   linkedliststack.New[*HandleEntry](),
		readyChan: make(chan struct{}, 1),
	}

	return out
}

var _ stepper.Stepper = &debugStepper{}

type debugStepper struct {
	stepper.Stepper
	History stacks.Stack[*HandleEntry]

	// waits for messages on this chan to move to the next step.
	readyChan chan struct{}
}

func (s *debugStepper) Step() {
	if len(s.readyChan) == 0 {
		s.Stepper.Current().Val.Cancel(engine.ErrSkip)
		s.readyChan <- struct{}{}
	}
}

func (s *debugStepper) Replay() {
	s.jump(1)
}

func (s *debugStepper) PreviousStep() {
	s.jump(2)
}

func (s *debugStepper) jump(lvls int) {
	h, _ := s.History.Pop()
	if h != nil && h.Handle.Node() != nil {
		s.Stepper.Jump(stepper.NavUp(h.Handle.Node(), lvls))
	}
	s.Step()
}

func (s *debugStepper) String() string {
	return s.Stepper.String()
}

func (s *debugStepper) Next(c *engine.Context) stepper.Handle {
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
