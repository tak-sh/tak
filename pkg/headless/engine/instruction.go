package engine

import "time"

type EventQueue chan Event

func NewEventQueue() EventQueue {
	return make(EventQueue, 10)
}

type Event interface {
	eventSigil()
}

type PathNode interface {
	GetId() string

	// IsReady lets the Decider know that the PathNode is either  ready
	// to be taken, or is not ready to be taken.
	IsReady(st *TemplateData) bool
}

// Instruction is an individual pieces of work to be evaluated at runtime.
// Every Instruction should have their Eval method called via the Evaluator
// rather than calling it directly.
type Instruction interface {
	GetId() string
	Eval(c *Context, to time.Duration) error
}

var _ Event = &NextInstructionEvent{}

type NextInstructionEvent struct {
	// The Instruction that is now running.
	Instruction Instruction
}

func (c *NextInstructionEvent) eventSigil() {}

// Evaluator evaluates and tracks Instruction's that have been
// previously evaluated.
type Evaluator interface {
	Eval(c *Context, i Instruction) error
	Prev() Instruction
}

func NewEvaluator(eq EventQueue, to time.Duration) Evaluator {
	out := &evaluator{
		Q:         eq,
		Evaluated: make([]Instruction, 0),
		Timeout:   to,
	}

	return out
}

type evaluator struct {
	Q         EventQueue
	Evaluated []Instruction
	Timeout   time.Duration
}

func (e *evaluator) Eval(c *Context, i Instruction) error {
	if e.Q != nil {
		e.Q <- &NextInstructionEvent{
			Instruction: i,
		}
	}

	err := i.Eval(c, e.Timeout)
	if err != nil {
		return err
	}

	e.Evaluated = append(e.Evaluated, i)
	return nil
}

func (e *evaluator) Prev() Instruction {
	n := len(e.Evaluated)
	if n > 0 {
		return e.Evaluated[n-1]
	}

	return nil
}
