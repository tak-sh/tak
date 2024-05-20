package action

import (
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/headless"
	"github.com/tak-sh/tak/pkg/validate"
)

type Action interface {
	fmt.Stringer
	validate.Validator

	// Act performs the browser action. Typically, leverages chromedp.Action.
	Act(ctx *headless.Context) error

	// GetID a unique ID for the Action. Used for things like calling Context.RenderTemplate.
	GetID() string
}

func New(id string, a *v1beta1.Action) Action {
	if act := a.GetInput(); act != nil {
		return NewInput(id, act)
	} else if act := a.GetMouseClick(); act != nil {
		return NewMouseClick(id, act)
	} else if act := a.GetAsk(); act != nil {
		return NewPromptAction(id, act)
	} else if act := a.GetNav(); act != nil {
		return NewNav(id, act)
	}
	return &NoOpAction{
		ID: id,
	}
}
