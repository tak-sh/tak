package step

import (
	"context"
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/validate"
	"time"
)

type Action interface {
	fmt.Stringer
	validate.Validator

	// Act performs the browser action. Typically, leverages chromedp.Action.
	Act(ctx *engine.Context) error

	// GetID a unique ID for the Action. Used for things like calling Context.RenderTemplate.
	GetID() string
}

func New(id string, a *v1beta1.Action) (Action, error) {
	if act := a.GetInput(); act != nil {
		return NewInput(id, act)
	} else if act := a.GetMouseClick(); act != nil {
		return NewMouseClick(id, act)
	} else if act := a.GetAsk(); act != nil {
		return NewPromptAction(id, act)
	} else if act := a.GetNav(); act != nil {
		return NewNav(id, act), nil
	} else if act := a.GetBranch(); act != nil {
		return NewBranch(id, act)
	} else if act := a.GetStore(); act != nil {
		return NewStoreAction(id, act)
	}

	return nil, except.NewInvalid("empty actions are not allowed")
}

func RunAction(c *engine.Context, act Action, to time.Duration) error {
	switch act.(type) {
	case *PromptAction, *BranchAction:
		return act.Act(c)
	}
	ctx := c.Context
	if to > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(c.Context, to)
		defer cancel()
	}
	err := act.Act(c.WithContext(ctx))
	if errors.Is(err, context.DeadlineExceeded) {
		return except.NewTimeout("took too long")
	}
	return err
}

func RunActionAsync(c *engine.Context, act Action, to time.Duration) <-chan error {
	out := make(chan error, 1)
	go func() {
		out <- RunAction(c, act, to)
	}()
	return out
}
