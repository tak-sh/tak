package step

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/api/script/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"time"
)

var _ Action = &ForEachElementAction{}

func NewForEachAction(id string, fe *v1beta1.Action_ForEachElement) (*ForEachElementAction, error) {
	out := &ForEachElementAction{
		Action_ForEachElement: fe,
		ID:                    id,
		CompiledActions:       make([]Action, 0, len(fe.GetActions())),
	}

	for i, act := range fe.GetActions() {
		a, err := New(id, act)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("action #%d", i), err)
		}
		out.CompiledActions = append(out.CompiledActions, a)
	}

	var err error
	out.CompiledSelector, err = engine.CompileTemplate(fe.GetSelector())
	if err != nil {
		return nil, errors.Join(except.NewInvalid("invalid selector %s", fe.GetSelector()), err)
	}

	return out, nil
}

type ForEachElementAction struct {
	*v1beta1.Action_ForEachElement
	ID               string
	CompiledSelector *engine.TemplateRenderer
	CompiledActions  []Action
}

func (f *ForEachElementAction) GetId() string {
	return f.ID
}

func (f *ForEachElementAction) Eval(c *engine.Context, to time.Duration) error {
	cont, err := c.Browser.Content(c, f.CompiledSelector.Render(c.TemplateData))
	if err != nil {
		return err
	}

	co := c.Copy()
	for _, v := range cont {
		co.TemplateData = c.TemplateData.Merge(&engine.TemplateData{ScriptTemplateData: &v1beta1.ScriptTemplateData{
			Element: v,
		}})

		for _, a := range f.CompiledActions {
			err = a.Eval(co, to)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *ForEachElementAction) String() string {
	return "iterating through " + f.GetSelector()
}

func (f *ForEachElementAction) Validate() error {
	return nil
}
