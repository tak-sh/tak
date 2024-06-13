package step

import (
	"context"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"time"
)

type Decider interface {
	ChoosePath(c *engine.Context, to time.Duration, pe ...PathNode) (node PathNode, stop bool, err error)
}

type PathNode interface {
	// IsReady lets the Decider know that the PathNode is either  ready
	// to be taken, or is not ready to be taken.
	IsReady(st *engine.TemplateData) bool
}

// Signaller signals to the engine that they carry some state that may affect how to continue.
type Signaller interface {
	// CheckSignal returns a non-nil signal if it has been triggered.
	CheckSignal(st *engine.TemplateData) *v1beta1.ConditionalSignal
}

func NewDecider(scriptSigs []*ConditionalSignal) Decider {
	out := &decider{
		ScriptSignals: scriptSigs,
	}

	return out
}

type decider struct {
	// Top-level signals that run every time a fork in a path is encountered.
	ScriptSignals []*ConditionalSignal
}

func (d *decider) ChoosePath(c *engine.Context, to time.Duration, pe ...PathNode) (PathNode, bool, error) {
	ctx, cancel := context.WithTimeout(c.Context, to)
	defer cancel()
	ticker := time.NewTicker(250 * time.Millisecond)
	for {
		for _, v := range d.ScriptSignals {
			sig := v.CheckSignal(c.TemplateData)
			switch sig.GetSignal() {
			case v1beta1.ConditionalSignal_success:
				return nil, true, nil
			case v1beta1.ConditionalSignal_error:
				return nil, false, except.NewInvalid(sig.GetMessage())
			}
		}

		for _, v := range pe {
			if t, ok := v.(Signaller); ok {
				sig := t.CheckSignal(c.TemplateData)
				switch sig.GetSignal() {
				case v1beta1.ConditionalSignal_success:
					continue
				case v1beta1.ConditionalSignal_error:
					return v, false, except.NewInvalid(sig.GetMessage())
				}
			}

			isReady := v.IsReady(c.TemplateData)

			if isReady {
				return v, false, nil
			}
		}

		// if we're waiting to decide, it may be that the page is stale.
		_ = c.RefreshPageState()

		select {
		case <-ctx.Done():
			return nil, true, ctx.Err()
		case <-ticker.C:
		}
	}
}
