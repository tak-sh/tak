package step

import (
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/validate"
)

type Action interface {
	validate.Validator
	engine.Instruction
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
	} else if act := a.GetForEachElement(); act != nil {
		return NewForEachAction(id, act)
	}

	return nil, except.NewInvalid("empty actions are not allowed")
}
