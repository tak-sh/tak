package stepper

import (
	"context"
	"errors"
	"fmt"
	"github.com/eddieowens/opts"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/step"
	"slices"
	"strings"
	"time"
)

type Stepper interface {
	fmt.Stringer
	Next(c *engine.Context) Handle
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

func New(globalSignals []*step.ConditionalSignal, steps []*step.Step, o ...opts.Opt[Opts]) (Stepper, error) {
	out := &stepper{
		Op:      opts.DefaultApply(o...),
		Signals: globalSignals,
		Root:    &node{},
	}

	if len(steps) == 0 {
		return nil, except.NewInvalid("at least 1 step required")
	}

	idx := slices.IndexFunc(globalSignals, func(signal *step.ConditionalSignal) bool {
		return signal.GetSignal() == v1beta1.ConditionalSignal_success
	})
	if idx < 0 {
		return nil, except.NewInvalid("at least 1 success condition required")
	}

	n := out.Root
	for i := range steps {
		v := steps[i]
		child := &node{Val: v, Parent: n}
		n.Children = append(n.Children, child)
		if _, ok := v.CompiledAction.(step.Branches); !ok {
			n = child
		}
	}

	out.Current = out.Root

	return out, nil
}

type Handle interface {
	fmt.Stringer

	// Val returns the Step to be taken from the Stepper.
	Val() *step.Step

	// Signal returns the applicable step.ConditionalSignal. If this is populated, the
	// Stepper is considered terminated.
	Signal() *step.ConditionalSignal

	// Err returns an error if the step.ConditionalSignal holds an error signal.
	Err() error

	// Idx returns the index that the underlying Step held from the original Stepper.
	// If this value isn't relevant, a value < 0 is returned.
	Idx() int
}

type stepper struct {
	Op      Opts
	Signals []*step.ConditionalSignal
	Root    *node
	Current *node
}

func (s *stepper) String() string {
	if s == nil || s.Root == nil {
		return ""
	}

	return s.Root.String()
}

func (s *stepper) Next(c *engine.Context) Handle {
	ctx, cancel := context.WithTimeout(c.Context, s.Op.Timeout)
	defer cancel()
	ticker := time.NewTicker(s.Op.Tick)
	defer ticker.Stop()
	for {
		for _, v := range s.Signals {
			if v.IsReady(c) {
				if v.Signal == v1beta1.ConditionalSignal_error {
					return &handle{Error: errors.New(v.GetMessage()), Sig: v}
				}
				return &handle{Sig: v}
			}
		}

		for _, v := range s.Current.Children {
			if v.Val.IsReady(c) {
				s.Current = v
				return &handle{Node: v}
			}
		}

		// if we're waiting to decide, it may be that the page is stale.
		_ = c.RefreshPageState()

		select {
		case <-ctx.Done():
			return &handle{Error: context.Cause(ctx)}
		case <-ticker.C:
		}
	}
}

var _ fmt.Stringer = &node{}

type node struct {
	Val      *step.Step
	Idx      int
	Children []*node
	Parent   *node
}

func (n *node) String() string {
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

type handle struct {
	Node  *node
	Sig   *step.ConditionalSignal
	Error error
}

func (h *handle) String() string {
	if h.Node != nil && h.Node.Val != nil {
		return h.Node.Val.CompiledAction.String()
	}

	if h.Sig != nil {
		return h.Sig.String()
	}

	if h.Error != nil {
		return h.Error.Error()
	}

	return ""
}

func (h *handle) Val() *step.Step {
	if h.Node == nil {
		return nil
	}

	return h.Node.Val
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

func (h *handle) Idx() int {
	if h.Node != nil {
		return -1
	}
	return h.Node.Idx
}
