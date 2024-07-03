package engine

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrRetryEval = errors.New("retry eval")
	ErrSkip      = errors.New("skip eval")
)

type EventQueue chan Event

func NewEventQueue() EventQueue {
	return make(EventQueue, 10)
}

type Event interface {
	fmt.Stringer
	eventSigil()
}

type PathNode interface {
	GetId() string

	// IsReady lets the Decider know that the PathNode is either  ready
	// to be taken, or is not ready to be taken.
	IsReady(c *Context) bool
}

// Instruction is an individual pieces of work to be evaluated at runtime.
// Every Instruction should have their Eval method called via the Evaluator
// rather than calling it directly.
type Instruction interface {
	fmt.Stringer
	GetId() string
	Eval(c *Context, to time.Duration) error
	Cancel(err error)
}

var _ Event = &NextInstructionEvent{}

type NextInstructionEvent struct {
	// The Instruction that is now running.
	Instruction Instruction
	Context     *Context
}

func (c *NextInstructionEvent) String() string {
	return c.Instruction.String()
}

func (c *NextInstructionEvent) eventSigil() {}

// Evaluator evaluates and tracks Instruction's that have been
// previously evaluated.
type Evaluator interface {
	Eval(c *Context, i Instruction) error
}

func NewEvaluator(eq EventQueue, to time.Duration) Evaluator {
	out := &evaluator{
		Q:       eq,
		Timeout: to,
	}

	return out
}

type evaluator struct {
	Q       EventQueue
	Timeout time.Duration
}

func (e *evaluator) Eval(c *Context, i Instruction) (err error) {
	if e.Q != nil {
		e.Q <- &NextInstructionEvent{
			Instruction: i,
			Context:     c,
		}
	}

	for iter := 0; iter < 1 || errors.Is(err, ErrRetryEval); iter++ {
		err = i.Eval(c, e.Timeout)
	}

	if errors.Is(err, ErrSkip) {
		err = nil
	}

	return
}
