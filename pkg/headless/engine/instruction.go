package engine

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

var (
	ErrRetryEval = errors.New("retry eval")
	ErrSkip      = errors.New("skip eval")
)

type Operation string

const (
	OperationListAccounts         Operation = "list_accounts"
	OperationLogin                Operation = "login"
	OperationDownloadTransactions Operation = "download_transactions"
)

func (o Operation) ActionString() string {
	switch o {
	case OperationLogin:
		return "logging in"
	case OperationDownloadTransactions:
		return "downloading transactions"
	case OperationListAccounts:
		return "listing accounts"
	}
	return ""
}

type EventQueue chan Event

func NewEventQueue() EventQueue {
	return make(EventQueue, 10)
}

type Event interface {
	fmt.Stringer
	eventSigil()
}

type PathNode interface {
	// IsReady lets the Decider know that the PathNode is either  ready
	// to be taken, or is not ready to be taken.
	IsReady(c *Context) bool
}

// Instruction is an individual pieces of work to be evaluated at runtime.
// Every Instruction should have their Eval method called via the Evaluator
// rather than calling it directly.
type Instruction interface {
	fmt.Stringer
	Message() proto.Message
	GetId() string
	Eval(c *Context, to time.Duration) error
}

var _ Event = &ChangeOperationEvent{}

type ChangeOperationEvent struct {
	To      Operation
	Message string
}

func (c *ChangeOperationEvent) String() string {
	sb := strings.Builder{}
	sb.WriteString(c.To.ActionString())
	if c.Message != "" {
		sb.WriteString(" ")
		sb.WriteString(c.Message)
	}
	return sb.String()
}

func (c *ChangeOperationEvent) eventSigil() {}

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
	Eval(c *Context, i Instruction) EvalHandle
	EventQueue() EventQueue
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

func (e *evaluator) EventQueue() EventQueue {
	return e.Q
}

func (e *evaluator) Eval(c *Context, i Instruction) EvalHandle {
	if e.Q != nil {
		e.Q <- &NextInstructionEvent{
			Instruction: i,
			Context:     c,
		}
	}

	return EvalAsync(c, i, e.Timeout)
}

type EvalHandle interface {
	Cancel(err error)
	Done() <-chan struct{}
	Err() error
	// Cause can be nil if the context is Done but no error occurred.
	Cause() error
	Instruction() Instruction
}

func newHandle(ctx context.Context, i Instruction) EvalHandle {
	ctx, cancel := context.WithCancelCause(ctx)
	return &handle{
		Can: cancel,
		Ctx: ctx,
		Int: i,
	}
}

type handle struct {
	Can context.CancelCauseFunc
	Ctx context.Context
	Int Instruction
}

func (h *handle) Cause() error {
	err := h.Err()
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}

func (h *handle) Instruction() Instruction {
	return h.Int
}

func (h *handle) Err() error {
	return context.Cause(h.Ctx)
}

func (h *handle) Cancel(err error) {
	h.Can(err)
}

func (h *handle) Done() <-chan struct{} {
	return h.Ctx.Done()
}

func EvalAsync(c *Context, act Instruction, to time.Duration) EvalHandle {
	h := newHandle(c.Context, act)
	go func() {
		var err error
		for iter := 0; iter < 1 || errors.Is(err, ErrRetryEval); iter++ {
			err = act.Eval(c, to)
		}

		if errors.Is(err, ErrSkip) {
			err = nil
		}

		h.Cancel(err)
	}()
	return h
}
