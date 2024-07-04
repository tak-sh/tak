package stepper

import (
	"context"
	"errors"
	"fmt"
	"github.com/eddieowens/opts"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"strings"
	"sync"
	"time"
)

type Stepper interface {
	fmt.Stringer
	// Next checks for either a step.ConditionalSignal or step.Step is ready, or a timeout
	// (WithTimeout) occurs. These conditions are checked every tick duration (WithTickDuration).
	Next(c *engine.Context) Handle

	// Current retrieves the current Node.
	Current() *Node

	// Jump sets the current Node to be n.
	Jump(n *Node)
}

type Factory interface {
	NewStepper(globalSignals []*step.ConditionalSignal, steps []*step.Step) Stepper
}

func NewFactory(o ...opts.Opt[Opts]) Factory {
	out := &factory{
		Opts: o,
	}

	return out
}

type factory struct {
	Opts []opts.Opt[Opts]
}

func (f *factory) NewStepper(globalSignals []*step.ConditionalSignal, steps []*step.Step) Stepper {
	return New(globalSignals, steps, f.Opts...)
}

type Opts struct {
	// Timeout to decide a new path.
	Timeout time.Duration
	// Tick duration when deciding a new path.
	Tick time.Duration
}

func WithTimeout(to time.Duration) opts.Opt[Opts] {
	return func(o *Opts) {
		o.Timeout = to
	}
}

func WithTickDuration(dur time.Duration) opts.Opt[Opts] {
	return func(o *Opts) {
		o.Tick = dur
	}
}

func (o Opts) DefaultOptions() Opts {
	return Opts{
		Timeout: 10 * time.Second,
		Tick:    250 * time.Millisecond,
	}
}

func New(globalSignals []*step.ConditionalSignal, steps []*step.Step, o ...opts.Opt[Opts]) Stepper {
	out := &stepper{
		Op:      opts.DefaultApply(o...),
		Signals: globalSignals,
		Root:    NewGraph(steps...),
	}

	out.Ticker = time.NewTicker(out.Op.Tick)
	out.Curr = out.Root

	return out
}

func NewGraph(steps ...*step.Step) *Node {
	root := &Node{}
	n := root
	for i := range steps {
		v := steps[i]
		child := &Node{Val: v, Parent: n}
		n.Children = append(n.Children, child)
		if _, ok := v.CompiledAction.(step.Branches); !ok {
			n = child
		}
	}
	return root
}

type Handle interface {
	fmt.Stringer

	Node() *Node

	// Signal returns the applicable step.ConditionalSignal. If this is populated, the
	// Stepper is considered terminated.
	Signal() *step.ConditionalSignal

	// Err returns an error if the step.ConditionalSignal holds an error signal.
	Err() error
}

type stepper struct {
	Op      Opts
	Signals []*step.ConditionalSignal
	Root    *Node
	Curr    *Node
	Ticker  *time.Ticker

	currLock sync.RWMutex
}

func (s *stepper) Current() *Node {
	s.currLock.RLock()
	defer s.currLock.RUnlock()
	return s.Curr
}

func (s *stepper) Jump(n *Node) {
	s.Curr.Val.Cancel(engine.ErrSkip)
	s.setCurr(n)
	s.Ticker.Reset(s.Op.Tick)
}

func (s *stepper) String() string {
	if s == nil || s.Root == nil {
		return ""
	}

	return s.Root.String()
}

func (s *stepper) Next(c *engine.Context) Handle {
	ctx := c.Context
	if s.Op.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.Op.Timeout)
		defer cancel()
	}
	c = c.WithContext(ctx)
	for {
		h := FirstReady(c, s.Signals, s.getChildren())
		if h != nil {
			s.Curr = h.Node()
			return h
		}

		// if we're waiting to decide, it may be that the page is stale.
		_ = c.RefreshPageState()

		select {
		case <-c.Done():
			return NewErrHandle(context.Cause(ctx))
		case <-s.Ticker.C:
		}
	}
}

func (s *stepper) getChildren() []*Node {
	s.currLock.RLock()
	defer s.currLock.RUnlock()
	return s.Curr.Children
}

func (s *stepper) setCurr(n *Node) {
	s.currLock.Lock()
	defer s.currLock.Unlock()
	s.Curr = n
}

func FirstReady(c *engine.Context, signals []*step.ConditionalSignal, children []*Node) Handle {
	for _, v := range signals {
		if v.IsReady(c) {
			return NewSignalHandle(v)
		}
	}

	for _, v := range children {
		if v.Val.IsReady(c) {
			return NewNodeHandle(v)
		}
	}

	return nil
}

var _ fmt.Stringer = &Node{}

type Node struct {
	Val      *step.Step
	Idx      int
	Children []*Node
	Parent   *Node
}

// NavUp follows the Node.Parent until either the Parent is nil (root)
// or it has navigated upwards n times.
func NavUp(node *Node, n int) *Node {
	for i := 0; i < n && node.Parent != nil; i++ {
		node = node.Parent
	}
	return node
}

func (n *Node) String() string {
	out := make([]string, 0, len(n.Children)+1)
	if n.Val != nil && n.Val.CompiledAction != nil {
		out = append(out, n.Val.CompiledAction.String())
	}

	child := make([]string, 0, len(n.Children))
	for _, v := range n.Children {
		child = append(child, v.String())
	}
	if len(child) > 0 {
		out = append(out, strings.Join(child, ", "))
	}

	return strings.Join(out, " -> ")
}

func NewErrHandle(err error) Handle {
	return &handle{Error: err}
}

func NewNodeHandle(n *Node) Handle {
	return &handle{N: n}
}

func NewSignalHandle(s *step.ConditionalSignal) Handle {
	if s.Signal == v1beta1.ConditionalSignal_error {
		return &handle{Error: errors.New(s.GetMessage()), Sig: s}
	}

	return &handle{Sig: s}
}

type handle struct {
	N     *Node
	Sig   *step.ConditionalSignal
	Error error
}

func (h *handle) Node() *Node {
	return h.N
}

func (h *handle) String() string {
	if h.N != nil && h.N.Val != nil {
		return h.N.Val.CompiledAction.String()
	}

	if h.Sig != nil {
		return h.Sig.String()
	}

	if h.Error != nil {
		return h.Error.Error()
	}

	return ""
}

func (h *handle) Signal() *step.ConditionalSignal {
	return h.Sig
}

func (h *handle) Err() error {
	if h.Error != nil {
		return h.Error
	}

	return nil
}
